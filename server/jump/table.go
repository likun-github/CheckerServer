package jump

import (
	"CheckerServer/server/jump/objects"
	"container/list"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/modules/timer"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt64(min, max int64) int64 {
	if min >= max {
		return max
	}
	return rand.Int63n(max-min) + min
}

//游戏逻辑相关的代码
//第一阶段 空档期，等待押注
//第二阶段 押注期  可以押注
//第三阶段 开奖期  开奖
//第四阶段 结算期  结算

type Table struct {
	fsm FSM
	room.BaseTableImp
	room.QueueTable
	room.UnifiedSendMessageTable
	room.TimeOutTable
	module                  module.RPCModule
	seats                   []*objects.Player
	currentPlayer           int//0,1
	viewer                  *list.List //观众
	seatMax                 int        //房间最大座位数
	current_id              int64		//当前房间人数
	current_frame           int64 //当前帧
	sync_frame              int64 //上一次同步数据的帧
	stoped                  bool
	writelock               sync.Mutex

	MatchPeriodHandler         FSMHandler//空挡转匹配完成
	Match2ControlHandler       FSMHandler//匹配完成到控制期
	Control2WithdrawHandler    FSMHandler//控制期转悔棋期
	Withdraw2ControlHandler    FSMHandler//悔棋期转控制期
	Control2DrawHandler        FSMHandler//控制期转求和期
	Draw2ControlHandler        FSMHandler//求和期转控制期
	Control2PlayDraughtHandler FSMHandler//控制期转行棋期
	PlayDraught2ControlHandler FSMHandler//行棋期转控制期
	SettlementPeriodHandler    FSMHandler//控制到结算期
	step   					int64 //悔棋期帧
	//composition 			*stack.Stack
}

func NewTable(module module.RPCModule, tableId int) *Table {
	this := &Table{
		module:        module,
		stoped:        true,
		seatMax:       2,
		current_id:    0,//当前房间人数
		current_frame: 0,//当前帧
		sync_frame:    0,//上一帧
		//composition:stack.NewStack(),//棋局
	}
	this.BaseTableImpInit(tableId, this)
	this.QueueInit()
	this.UnifiedSendMessageTableInit(this)
	this.TimeOutTableInit(this, this, 60)
	//游戏逻辑状态机
	this.InitFsm()
	this.seats = make([]*objects.Player, this.seatMax)
	this.viewer = list.New()

	this.Register("SitDown", this.SitDown)
	this.Register("GetLevel", this.getLevel)
	this.Register("StartGame", this.StartGame)
	this.Register("PauseGame", this.PauseGame)
	this.Register("Stake", this.Stake)
	this.Register("Login", this.login)

	for indexSeat, _ := range this.seats {
		this.seats[indexSeat] = objects.NewPlayer(indexSeat)
	}

	return this
}
func (this *Table) GetModule() module.RPCModule {
	return this.module
}
func (this *Table) GetSeats() []room.BasePlayer {
	m := make([]room.BasePlayer, len(this.seats))
	for i, seat := range this.seats {
		m[i] = seat
	}
	return m
}
func (this *Table) GetViewer() *list.List {
	return this.viewer
}

/**
玩家断线,游戏暂停
*/
func (self *Table) OnNetBroken(player room.BasePlayer) {
	player.OnNetBroken()
}

////访问权限校验
func (this *Table) VerifyAccessAuthority(userId string, bigRoomId string) bool {
	_, tableid, transactionId, err := room.ParseBigRoomId(bigRoomId)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if (tableid != this.TableId()) || (transactionId != this.TransactionId()) {
		log.Error("transactionId!=this.TransactionId()", transactionId, this.TransactionId())
		return false
	}
	return true
}
func (this *Table) AllowJoin() bool {
	this.writelock.Lock()
	ready := true
	if this.current_id > 1 {
		this.writelock.Unlock()
		return false
	}
	this.current_id++
	this.writelock.Unlock()
	return true
	//for _, seat := range this.GetSeats() {
	//	if seat.Bind() == false {
	//		//还没有准备好
	//		ready = false
	//		break
	//	}
	//}

	return !ready
}

func (this *Table) OnCreate() {
	this.BaseTableImp.OnCreate()
	this.ResetTimeOut()
	//log.Debug("Table", "OnCreate")
	if this.stoped {
		this.stoped = false
		timewheel.GetTimeWheel().AddTimer(1000*time.Millisecond, nil, this.Update)
		//go func() {
		//	//这里设置为500ms
		//	tick := time.NewTicker(1000 * time.Millisecond)
		//	defer func() {
		//		tick.Stop()
		//	}()
		//	for !this.stoped {
		//		select {
		//		case <-tick.C:
		//			this.Update(nil)
		//		}
		//	}
		//}()
	}
}
func (this *Table) OnStart() {
	log.Debug("Table", "OnStart")
	for _, player := range this.seats {
		player.Score=1000
		//player.Weight = 0
		//player.Target = 0
		//player.Stake = false
	}
	//将游戏状态设置到空闲期
	this.fsm.Call(MatchPeriodEvent)
	this.step=0
	this.current_frame = 0
	this.sync_frame = 0
	this.BaseTableImp.OnStart()
}
func (this *Table) OnResume() {
	this.BaseTableImp.OnResume()
	log.Debug("Table", "OnResume")
	this.NotifyResume()
}
func (this *Table) OnPause() {
	this.BaseTableImp.OnPause()
	log.Debug("Table", "OnPause")
	this.NotifyPause()
}
func (this *Table) OnStop() {
	this.BaseTableImp.OnStop()
	log.Debug("Table", "OnStop")
	//将游戏状态设置到空档期
	this.fsm.Call(VoidPeriodEvent)
	this.NotifyStop()
	this.ExecuteCallBackMsg() //统一发送数据到客户端
	for _, player := range this.seats {
		player.OnUnBind()
	}

	var nv *list.Element
	for e := this.viewer.Front(); e != nil; e = nv {
		nv = e.Next()
		this.viewer.Remove(e)
	}
}
func (this *Table) OnDestroy() {
	this.BaseTableImp.OnDestroy()
	log.Debug("BaseTableImp", "OnDestroy")
	this.stoped = true
}

func (self *Table) onGameOver() {
	self.Finish()
}

/**
牌桌主循环
定帧计算所有玩家的位置
*/


func (self *Table) Update(arge interface{}) {
	self.ExecuteEvent(arge) //执行这一帧客户端发送过来的消息
	if self.State() == room.Active {
		self.current_frame++

		if self.current_frame-self.sync_frame >= 1 {
			//每帧同步一次
			self.sync_frame = self.current_frame
			self.StateSwitch()
		}

		ready := true
		for _, seat := range self.GetSeats() {
			if seat.Bind() == false {
				//还没有准备好
				ready = false
				break
			}
		}
		if ready == false {
			//有玩家离开了牌桌,牌桌退出
			self.Finish()
		}

	} else if self.State() == room.Initialized {
		ready := true
		for _, seat := range self.GetSeats() {
			if seat.SitDown() == false {
				//还没有准备好
				ready = false
				break
			}
		}
		if ready {
			self.Start() //开始游戏了
		}
	}

	self.ExecuteCallBackMsg() //统一发送数据到客户端
	self.CheckTimeOut()
	if !self.stoped {
		timewheel.GetTimeWheel().AddTimer(100*time.Millisecond, nil, self.Update)
	}
}

func (self *Table) Exit(session gate.Session) error {
	player := self.GetBindPlayer(session)
	if player != nil {
		playerImp := player.(*objects.Player)
		playerImp.OnUnBind()
		return nil
	}
	return nil
}

func (self *Table) getSeatsMap() []map[string]interface{} {
	m := make([]map[string]interface{}, len(self.seats))
	for i, player := range self.seats {
		if player != nil {
			m[i] = player.SerializableMap()
		}
	}
	return m
}

/**
玩家获取关卡信息
*/
func (self *Table) getLevel(session gate.Session) {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
	}
}


