package jump

import (
	"CheckerServer/server/dao"
	"CheckerServer/server/jump/objects"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/server"
	"strconv"
)

var Module = func() module.Module {
	this := new(Jump)
	return this
}

type Jump struct {
	basemodule.BaseModule
	room    *room.Room
	proTime int64
	gameId  int
}

func (self *Jump) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Jump"
}
func (self *Jump) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *Jump) GetFullServerId() string {
	return self.GetServerId()
}

func (self *Jump) usableTable(table room.BaseTable) bool {
	return table.AllowJoin()
}

func (self *Jump) newTable(module module.RPCModule, tableId int) (room.BaseTable, error) {
	table := NewTable(module, tableId)
	return table, nil
}

func (self *Jump) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings, server.Metadata(map[string]string{
		"type": "helloworld",
	}))
	//房间号
	self.gameId = 11
	self.room = room.NewRoom(self, self.gameId, self.newTable, self.usableTable)

	self.GetServer().RegisterGO("HD_GetUsableTable", self.HDGetUsableTable)		//我们约定所有对客户端的请求都以Handler_开头
	//self.GetServer().RegisterGO("HD_Enter", self.enter)
	self.GetServer().RegisterGO("HD_Exit", self.exit)
	//self.GetServer().RegisterGO("HD_SitDown", self.sitdown)
	//self.GetServer().RegisterGO("HD_StartGame", self.startGame)
	//self.GetServer().RegisterGO("HD_PauseGame", self.pauseGame)
	//self.GetServer().RegisterGO("HD_Stake", self.stake)

	//self.GetServer().RegisterGO("HD_Login", self.login)//行棋
	self.GetServer().RegisterGO("HD_Control", self.control)//行棋
	self.GetServer().RegisterGO("HD_Withdraw", self.withdraw)//悔棋
	self.GetServer().RegisterGO("HD_WithdrawDecided", self.withdrawdecided) //悔棋结果确定
	self.GetServer().RegisterGO("HD_Draw", self.draw)//和棋
	self.GetServer().RegisterGO("HD_DrawDecided", self.drawdecided)//和棋结果确定
	self.GetServer().RegisterGO("HD_Lose", self.lose)//认输
	self.GetServer().RegisterGO("HD_Collect", self.collect)//收藏
	self.GetServer().RegisterGO("HD_Again", self.again)//再来一局
	self.GetServer().RegisterGO("HD_Exit", self.exit)//再来一局
}

func (self *Jump) Run(closeSig chan bool) {
}

func (self *Jump) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

/**
检查参数是否存在
*/
func (self *Jump) ParameterCheck(msg map[string]interface{}, paras ...string) error {
	for _, v := range paras {
		if _, ok := msg[v]; !ok {
			return fmt.Errorf("No %s found", v)
		}
	}
	return nil
}

/**
检查参数是否存在
*/
func (self *Jump) GetTableByBigRoomId(bigRoomId string) (*Table, error) {
	_, tableid, _, err := room.ParseBigRoomId(bigRoomId)
	if err != nil {
		return nil, err
	}
	table := self.room.GetTable(tableid)
	if table != nil {
		tableimp := table.(*Table)
		return tableimp, nil
	} else {
		return nil, errors.New("No table found")
	}
}

/**
查找可用座位
*/
func (self *Jump) HDGetUsableTable(session gate.Session, msg map[string]interface{}) (map[string]interface{}, string) {
	// 客户端消息验证：是否发了userid
	if msg["userid"] == nil  {
		fmt.Println("客户端没有发送userid")
		session.Send("Match_Fail/No_userid",nil)
		return nil,"没有用户id"
	} else {
		fmt.Println("客户端发送的userid为", msg["userid"].(string))
	}
	session.Bind(msg["userid"].(string))

	return self.getUsableTable(session,msg)
}

/**
查找可用座位
*/
func (self *Jump) getUsableTable(session gate.Session,msg map[string]interface{}) (map[string]interface{}, string) {
	table, err := self.room.GetUsableTable()
	if err == nil {
		table.Create()
		tableInfo := map[string]interface{}{
			"BigRoomId": room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId()),
		}
		// 玩家加入桌子
		_, error := self.enter(session,tableInfo)	// 先将session跟table里的player绑定，再将session跟BigRoomId绑定
		if error != "" {
			fmt.Println(error)
			session.Send("Match_Fail/Enter_Table_Fail",nil)
			return nil, error
		}

		// 获取玩家数据库信息并绑定
		_, error = self.login(session,msg)
		if error != "" {
			fmt.Println(error)
			session.Send("Match_Fail/Player_Info_Fail",nil)
			return nil,error
		}

		// 给玩家发送桌子信息
		table_info_json,_ := json.Marshal(tableInfo)
		session.Send("TableInfo",table_info_json)

		fmt.Println(tableInfo)
		return tableInfo, ""
	} else {
		fmt.Println("没有可用的桌子")
		session.Send("Match_Fail/No_Available_Table",nil)
		return nil, "There is no available table"
	}
}

// player与session绑定
func (self *Jump) enter(session gate.Session, msg map[string]interface{}) (string, string) {
	if BigRoomId, ok := msg["BigRoomId"]; !ok { // 找不到BigRoomId对应的桌子
		fmt.Println("没有找到BigRoomId对应的桌子")
		session.Send("Match_Fail/No_BigRoomId_Found",nil)
		return "", "No BigRoomId found"
	} else { // 找到了BigRoomId对应的桌子
		bigRoomId := BigRoomId.(string)
		moduleId, tableid, _, err := room.ParseBigRoomId(bigRoomId)
		if err != nil { // 解析BigRoomId失败
			fmt.Println(err.Error())
			return "", err.Error()
		}
		if session.Get("BigRoomId") != "" { // 如果用户已经加入了别的桌子，那么先从上一个桌子退出
			if session.Get("BigRoomId") != bigRoomId {
				_, e := self.RpcInvoke(moduleId, "HD_Exit", session, map[string]interface{}{
					"BigRoomId": session.Get("BigRoomId"),
				})
				if e != "" {
					fmt.Println(e)
					return "", e
				}
			}
		}
		// 获取相应的桌子实例并在验证后加入，并将BigRoomId与session绑定
		table := self.room.GetTable(tableid)
		if table != nil {
			tableimp := table.(*Table)
			if table.VerifyAccessAuthority(session.GetUserId(), bigRoomId) == false {
				fmt.Println("Access rights validation failed")
				session.Send("Match_Fail/Access_Rights_Validation_Failed",nil)
				return "", "Access rights validation failed"
			}
			erro := tableimp.Join(session)
			if erro == nil {
				bigRoomId = room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId())
				session.Set("BigRoomId", bigRoomId)
				session.Push()
				return bigRoomId, ""
			} else {
				fmt.Println(erro.Error())
				return "", erro.Error()
			}
		} else {
			fmt.Println("No table found")
			return "", "No table found"
		}
	}
	return "", ""
}



//player绑定
func (self *Jump) login(session gate.Session, msg map[string]interface{}) (string, string) {
	// 查询用户数据
	userid := msg["userid"].(string)
	infoDao := dao.NewUserInfoDao()
	userId, _ := strconv.ParseInt(userid, 10, 64)
	userInfo := infoDao.SelectById(userId)

	//session.Bind(userInfo.WxName)
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.Login(session, userInfo.Id, userInfo.Level, userInfo.Score,userInfo.WxName,userInfo.WXImg)
	//err = table.PutQueue("Login", session, userInfo.Id, userInfo.Level, userInfo.Score,userInfo.WxName,userInfo.WXImg)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}

//行棋
func (self *Jump) control(session gate.Session, msg map[string]interface{}) (string, string) {
	if _, ok := msg["W"]; !ok {
		return "", "No W Found!"
	} else if _, ok := msg["B"]; !ok {
		return "", "No B Found!"
	} else if _, ok := msg["K"]; !ok {
		return "", "No K Found!"
	} else {
		// 把composition转成Chess格式
//		W,_ :=strconv.ParseInt(msg["W"].(string),10, 64)
//		B,_ :=strconv.ParseInt(msg["B"].(string),10, 64)
//		K,_ :=strconv.ParseInt(msg["K"].(string),10, 64)
		composition := objects.NewChess(msg["W"].(string), msg["B"].(string), msg["K"].(string))
		bigRoomId := session.Get("BigRoomId")
		if bigRoomId == "" {
			return "", "fail"
		}
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		if table.fsm.getState() != ControlPeriod { // 玩家在非控制期不知怎么地给服务器发了走子信息
			fmt.Print("非控制期发送的走子信息不予处理")
			return "fail",""
		}
		err = table.PlayOneTurn(session, composition)
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//悔棋
func (self *Jump) withdraw(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("Withdraw", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}
//悔棋结果确定
func (self *Jump) withdrawdecided(session gate.Session, msg map[string]interface{}) (string, string) {
	if withdrawAgreed, ok := msg["withdrawAgreed"]; !ok {
		return "", "No Target found"
	} else {
		bigRoomId := session.Get("BigRoomId")
		if bigRoomId == "" {
			return "", "fail"
		}
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		withdraw_agreed := int( withdrawAgreed.(float64))
		err = table.PutQueue("WithdrawDecided", session, withdraw_agreed)
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//求和
func (self *Jump) draw(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("Draw", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}
//悔棋结果确定
func (self *Jump) drawdecided(session gate.Session, msg map[string]interface{}) (string, string) {
	if drawAgreed, ok := msg["drawAgreed"]; !ok {
		return "", "No Target found"
	} else {
		bigRoomId := session.Get("BigRoomId")
		if bigRoomId == "" {
			return "", "fail"
		}
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		draw_agreed := int(drawAgreed.(float64))
		err = table.PutQueue("DrawDecided", session, draw_agreed)
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//认输
func (self *Jump) lose(session gate.Session, msg map[string]interface{}) (string, string) {
	if loseCheckerColor, ok := msg["checkercolor"]; !ok {
		return "", "Who sent the lose request?! Tell me your checker color!"
	} else {
		bigRoomId := session.Get("BigRoomId")
		if bigRoomId == "" {
			return "", "fail"
		}
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		lose_checker_color := int(loseCheckerColor.(float64))
		err = table.PutQueue("Lose", session,lose_checker_color)
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
// 收藏本局
func (self *Jump) collect(session gate.Session, msg map[string]interface{}) (string, string) {
	if collectCheckerColor, ok := msg["checker_color"]; !ok {
		return "", "Who sent the collect request?! Tell me your checker color!"
	} else {
		bigRoomId := session.Get("BigRoomId")
		if bigRoomId == "" {
			return "", "fail"
		}
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		collect_checker_color := int(collectCheckerColor.(float64))
		err = table.PutQueue("Collect", session,collect_checker_color)
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
// 再来一局
func (self *Jump) again(session gate.Session, msg map[string]interface{}) (string, string) {
	if collectCheckerColor, ok := msg["checker_color"]; !ok {
		return "", "Who sent the again request?! Tell me your checker color!"
	} else {
		bigRoomId := session.Get("BigRoomId")
		if bigRoomId == "" {
			return "", "fail"
		}
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		collect_checker_color := int(collectCheckerColor.(float64))
		err = table.PutQueue("Again", session,collect_checker_color)
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
// 返回大厅
func (self *Jump) exit(session gate.Session, msg map[string]interface{}) (string, string) {
	if BigRoomId, ok := msg["BigRoomId"]; !ok {
		return "", "No BigRoomId found"
	} else if collectCheckerColor, ok := msg["checker_color"]; !ok {
		return "", "Who sent the exit request?! Tell me your checker color!"
	} else {
		// 先修改用户状态
		bigRoomId := BigRoomId.(string)
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		collect_checker_color := int(collectCheckerColor.(float64))
		err = table.PutQueue("Exit_", session,collect_checker_color)
		if err != nil {
			return "", err.Error()
		}
		// 在将session与table解绑
		err = table.Exit(session)
		if err == nil {
			bigRoomId = room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId())
			session.Set("BigRoomId", "") //设置到session
			session.Push()
			return bigRoomId, ""
		}
		return "", err.Error()
	}
}





