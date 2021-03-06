package jump

import (
	"CheckerServer/server/common/stack"
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
	currentPlayer           int			//0：白子,1：黑子
	viewer                  *list.List 	//观众
	seatMax                 int        	//房间最大座位数
	current_id              int64		//当前房间人数
	current_frame           int64 		//当前帧
	sync_frame              int64 		//上一次同步数据的帧
	stoped                  bool
	control_time 			int64		// 玩家的控制时间，是个变量。开始时为 s，走了 步后为 s。
	writelock               sync.Mutex

	MatchPeriodHandler         		FSMHandler//空挡转匹配完成
	Match2ControlHandler       		FSMHandler//匹配完成到控制期
	Control2WithdrawHandler    		FSMHandler//控制期转悔棋期
	Withdraw2ControlHandler    		FSMHandler//悔棋期转控制期
	Control2DrawHandler        		FSMHandler//控制期转求和期
	Draw2ControlHandler        		FSMHandler//求和期转控制期
	Control2PlayFinishHandler  		FSMHandler//控制期转行棋完成期
	PlayFinish2ControlHandler  		FSMHandler//行棋完成期转控制期
	SettlementPeriodHandler    		FSMHandler//控制期转结算期
	SettlementPeriod2ControlHandler FSMHandler//结算期转控制期
	SettlementPeriod2MatchHandler 	FSMHandler//结算期转匹配期
	SettlementPeriod2VoidlHandler 	FSMHandler//结算期转空档期

	// 计时相关
	start_match_step   					int64 // 匹配开始时的帧
	start_control_step					int64 // 控制方玩家开始控制时的帧
	start_withdraw_step					int64 // 非控制方玩家开始决定是否同意悔棋时的帧
	start_draw_step						int64 // 非控制方玩家开始决定是否同意和棋时的帧

	// 悔棋相关
	withdraw_requested  int // 是否有玩家要求悔棋：0(没有),1(有)
	withdraw_agreed		int // 悔棋结果：-1(未决定),0(不同意),1(同意)
	withdraw_timeout    int // 悔棋超时：0(未超时),1(超时)

	// 和棋相关
	draw_requested  int // 是否有玩家要求和棋：0(没有),1(有)
	draw_agreed		int // 和棋结果：-1(未决定),0(不同意),1(同意)
	draw_timeout    int // 和棋超时：0(未超时),1(超时)

	// 认输相关
	lose_requested  int // 哪个玩家认输了：-1(没有),0(白),1(黑)

	// 游戏结果
	winner	int		// 胜方
	loser   int		// 负方

	// 收藏相关
	collect_requested_white	int64	// 白方是否收藏本局：若收藏，则为白方userid；若不收藏，则为-1
	collect_requested_black	int64	// 黑方是否收藏本局：若收藏，则为黑方userid；若不收藏，则为-1

	// 游戏结束后的动作：-1:未选择; 1:再来一局;2:离开大厅
	game_finished_action_white	int
	game_finished_action_black	int

	// 棋局记录相关
	composition 			*stack.Stack	// 棋局栈
	composition_num			int				// composition的个数，即总共走了多少步。这个在行棋完成期进行更新
}

// 数据初始化
func (this *Table)initializeData(mode int) {
	this.current_id = 0
	this.withdraw_requested = 0
	this.withdraw_agreed = -1
	this.withdraw_timeout = 0
	this.draw_requested = 0
	this.draw_agreed = -1
	this.draw_timeout = 0
	this.lose_requested = -1

	this.start_withdraw_step = -1
	this.start_draw_step = -1
	this.winner = -1
	this.loser = -1
	this.collect_requested_white = -1
	this.collect_requested_black = -1
	this.game_finished_action_white = -1
	this.game_finished_action_black = -1
	this.composition = stack.NewStack()
	this.composition.Push(objects.NewChess("00000000000000000000000000000011111111111111111111",
		"11111111111111111111000000000000000000000000000000",
		"00000000000000000000000000000000000000000000000000"))
	this.composition_num = 1
	if mode == 0 { // 结算期转控制期
		this.start_match_step = -1
		this.start_control_step = this.current_frame
	} else if mode == 1 { // 结算期转匹配期
		this.start_match_step = this.current_frame
		this.start_control_step = -1
	} else if mode == 2 { // 结算期转空挡期
		this.start_match_step = -1
		this.start_control_step = -1
	}
}



func NewTable(module module.RPCModule, tableId int) *Table {
	this := &Table{
		module:        		module,
		stoped:        		true,
		seatMax:       		2,
		current_id:    		0,					//当前房间人数
		current_frame: 		0,					//当前帧
		sync_frame:    		0,					//上一帧
		withdraw_requested: 0,
		withdraw_agreed:	-1,
		withdraw_timeout: 	0,
		draw_requested: 	0,
		draw_agreed:		-1,
		draw_timeout:		0,
		lose_requested: 	-1,
		start_match_step:   -1,
		start_control_step: -1,
		start_withdraw_step:-1,
		start_draw_step:    -1,
		control_time:       10000,
		winner:				-1,
		loser:				-1,
		collect_requested_white:	-1,
		collect_requested_black:	-1,
		game_finished_action_white:	-1,
		game_finished_action_black:	-1,
	}
	this.BaseTableImpInit(tableId, this)
	this.QueueInit()
	this.UnifiedSendMessageTableInit(this)
	this.TimeOutTableInit(this, this, 60)

	//游戏逻辑状态机
	this.InitFsm()
	this.seats = make([]*objects.Player, this.seatMax)
	this.viewer = list.New()
	this.composition = stack.NewStack()
	this.composition.Push(objects.NewChess("00000000000000000000000000000011111111111111111111",
		                           "11111111111111111111000000000000000000000000000000",
		                            "00000000000000000000000000000000000000000000000000"))
	this.composition_num = 1

	//this.Register("SitDown", this.SitDown)
	//this.Register("GetLevel", this.getLevel)
	//this.Register("StartGame", this.StartGame)
	//this.Register("PauseGame", this.PauseGame)
	//this.Register("Stake", this.Stake)

	this.Register("Login", this.Login)

	this.Register("PlayOneTurn", this.PlayOneTurn)
	this.Register("Withdraw", this.Withdraw)
	this.Register("WithdrawDecided", this.WithdrawDecided)
	this.Register("Draw", this.Draw)
	this.Register("DrawDecided", this.DrawDecided)
	this.Register("Lose", this.Lose)
	this.Register("Collect", this.Collect)
	this.Register("Again", this.Again)
	this.Register("Exit_", this.Exit_)


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
	log.Debug("Table", "OnCreate")
	if this.stoped {
		this.stoped = false
		timewheel.GetTimeWheel().AddTimer(1000*time.Millisecond, nil, this.Update)
	}
}

func (this *Table) OnStart() {
	log.Debug("Table", "OnStart")
	for _, player := range this.seats {
		player.Score=1000
	}
	// 将游戏状态设置到控制器
	this.fsm.Call(ControlPeriodEvent)
	// this.start_match_step=0
	this.current_frame = 0
	this.sync_frame = 0
	this.BaseTableImp.OnStart()
}


func (this *Table) OnResume() {
	this.BaseTableImp.OnResume()
	log.Debug("Table", "OnResume")
	//this.NotifyResume()
}
func (this *Table) OnPause() {
	this.BaseTableImp.OnPause()
	log.Debug("Table", "OnPause")
	//this.NotifyPause()
}



func (this *Table) OnStop() {
	this.BaseTableImp.OnStop()
	log.Debug("Table", "OnStop")
	//将游戏状态设置到空档期
	this.fsm.Call(VoidPeriodEvent)
	//this.NotifyStop()
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

/*
func (self *Table) onGameOver() {
	self.Finish()
}

 */

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

		/*
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
		*/
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

/*
func (self *Table) getSeatsMap() []map[string]interface{} {
	m := make([]map[string]interface{}, len(self.seats))
	for i, player := range self.seats {
		if player != nil {
			m[i] = player.SerializableMap()
		}
	}
	return m
}

 */


