package jump

import (
	"CheckerServer/server/dao"
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

	self.GetServer().RegisterGO("HD_Match", self.match) //我们约定所有对客户端的请求都以Handler_开头
	self.GetServer().RegisterGO("HD_GetUsableTable", self.HDGetUsableTable)
	self.GetServer().RegisterGO("HD_Enter", self.enter)
	self.GetServer().RegisterGO("HD_Exit", self.exit)
	self.GetServer().RegisterGO("HD_SitDown", self.sitdown)
	self.GetServer().RegisterGO("HD_StartGame", self.startGame)
	self.GetServer().RegisterGO("HD_PauseGame", self.pauseGame)
	self.GetServer().RegisterGO("HD_Stake", self.stake)
	self.GetServer().RegisterGO("HD_Hello", func(session gate.Session, msg map[string]interface{}) (string, string) {
		//log.Info("HD_Hello")
		return "success", ""
	})
	self.GetServer().RegisterGO("HD_Login", self.login)//行棋
	self.GetServer().RegisterGO("HD_Control", self.control)//行棋
	self.GetServer().RegisterGO("HD_Withdraw", self.withdraw)//悔棋
	self.GetServer().RegisterGO("HD_WithdrawDecided", self.withdrawdecided)//同意悔棋
	self.GetServer().RegisterGO("HD_Draw", self.draw)//和棋
	self.GetServer().RegisterGO("HD_Lose", self.lose)//认输



	//房间号


}

func (self *Jump) Run(closeSig chan bool) {

}

func (self *Jump) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}
func (self *Jump) match(session gate.Session, msg map[string]interface{}) (result string, err string) {
	fmt.Println(msg["userid"])
	fmt.Println("match login")
	m2 := make(map[string]string)
	// 然后赋值
	m2["a"] = "就是简单的尝试"
	m2["b"] = "bb"
	j,_:=json.Marshal(m2)
	session.Send("try",j)
	return

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
	fmt.Println("看看hd")
	return self.getUsableTable(session)
}

/**
查找可用座位
*/

func (self *Jump) getUsableTable(session gate.Session) (map[string]interface{}, string) {
	//这个桌子分配逻辑还是不智能，如果空闲的桌子多了,人数少了不容易将他们分配到相同的桌子里面快速组局
	table, err := self.room.GetUsableTable()


	if err == nil {
		table.Create()
		tableInfo := map[string]interface{}{
			"BigRoomId": room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId()),
		}
		b, _ := json.Marshal(tableInfo)
		session.Send("table",b)
		fmt.Println(tableInfo)
		return tableInfo, ""
	} else {
		return nil, "There is no available table"
	}
}
//进入桌子
func (self *Jump) enter(session gate.Session, msg map[string]interface{}) (string, string) {
	//fmt.Println("看一看enter会不会跑")

	if BigRoomId, ok := msg["BigRoomId"]; !ok {
		return "", "No BigRoomId found"
	} else {
		bigRoomId := BigRoomId.(string)

		moduleId, tableid, _, err := room.ParseBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		if session.Get("BigRoomId") != "" {
			//用户当前已经加入过一个BigRoomId
			if session.Get("BigRoomId") != bigRoomId {
				//先从上一个桌子退出
				_, e := self.RpcInvoke(moduleId, "HD_Exit", session, map[string]interface{}{
					"BigRoomId": session.Get("BigRoomId"),
				})
				if e != "" {
					return "", e
				}
			}
		}
		table := self.room.GetTable(tableid)
		if table != nil {
			tableimp := table.(*Table)
			if table.VerifyAccessAuthority(session.GetUserId(), bigRoomId) == false {
				return "", "Access rights validation failed"
			}
			erro := tableimp.Join(session)
			if erro == nil {
				bigRoomId = room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId())
				session.Set("BigRoomId", bigRoomId) //设置到session
				session.Push()
				return bigRoomId, ""
			}
			return "", erro.Error()
		} else {
			return "", "No room found"
		}
	}

}

func (self *Jump) exit(session gate.Session, msg map[string]interface{}) (string, string) {
	if BigRoomId, ok := msg["BigRoomId"]; !ok {
		return "", "No BigRoomId found"
	} else {
		bigRoomId := BigRoomId.(string)
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
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

func (self *Jump) sitdown(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("SitDown", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}
func (self *Jump) startGame(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("StartGame", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}
func (self *Jump) pauseGame(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("PauseGame", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}

func (self *Jump) stake(session gate.Session, msg map[string]interface{}) (string, string) {
	if Target, ok := msg["Target"]; !ok {
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
		err = table.PutQueue("Stake", session, int64(Target.(float64)))
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//player绑定
func (self *Jump) login(session gate.Session, msg map[string]interface{}) (string, string) {
	if msg["userid"] == nil  {
		return "","没有用户id"
	}
	userid := msg["userid"].(string)
	infoDao := dao.NewUserInfoDao()

	userId, _ := strconv.ParseInt(userid, 10, 64)
	userInfo := infoDao.SelectById(userId)
	session.Bind(userInfo.WxName)
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("Login", session, userInfo.Score,userInfo.WxName,userInfo.WXImg)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""

}
//行棋
func (self *Jump) control(session gate.Session, msg map[string]interface{}) (string, string) {
	if Target, ok := msg["Target"]; !ok {
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
		err = table.PutQueue("Stake", session, int64(Target.(float64)))
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//悔棋
func (self *Jump) withdraw(session gate.Session, msg map[string]interface{}) (string, string) {
	if Target, ok := msg["Target"]; !ok {
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
		err = table.PutQueue("Stake", session, int64(Target.(float64)))
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//同意悔棋
func (self *Jump) withdrawdecided(session gate.Session, msg map[string]interface{}) (string, string) {
	if Target, ok := msg["Target"]; !ok {
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
		err = table.PutQueue("Stake", session, int64(Target.(float64)))
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//求和
func (self *Jump) draw(session gate.Session, msg map[string]interface{}) (string, string) {
	if Target, ok := msg["Target"]; !ok {
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
		err = table.PutQueue("Stake", session, int64(Target.(float64)))
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
//认输
func (self *Jump) lose(session gate.Session, msg map[string]interface{}) (string, string) {
	if Target, ok := msg["Target"]; !ok {
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
		err = table.PutQueue("Stake", session, int64(Target.(float64)))
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}








