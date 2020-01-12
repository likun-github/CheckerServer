package jump




import (
	"encoding/json"
	"fmt"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/server"
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

// 新建桌子
func (self *Jump) newTable(module module.RPCModule, tableId int) (room.BaseTable, error) {
	table := NewTable(module, tableId);
	return table, nil;
}

// 確定table是否可加入
func (self *Jump) usableTable(table room.BaseTable) bool {
	return table.AllowJoin();
}

func (self *Jump) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings, server.Metadata(map[string]string{
		"type": "helloworld",
	}))
	//房间号
	self.gameId = 7;
	self.room = room.NewRoom(self, self.gameId, self.newTable, self.usableTable);
	//我们约定所有对客户端的请求都以Handler_开头
	self.GetServer().RegisterGO("HD_Match", self.match);
	//房间号


}

func (self *Jump) Run(closeSig chan bool) {

}

func (self *Jump) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

// 開始匹配
func (self *Jump) match(session gate.Session, msg map[string]interface{}) (result string, err string) {
	fmt.Println(msg["userid"]);
	k,_ :=self.getUsableTable(session);
	fmt.Println(k);


	m2 := make(map[string]string)
	// 然后赋值
	m2["a"] = "就是简单的尝试"
	m2["b"] = "bb"
	j,_:=json.Marshal(m2)
	session.Send("try",j)
	return

}

func (self *Jump) getUsableTable(session gate.Session) (map[string]interface{}, string) {
	fmt.Println("运行了吗")
	//这个桌子分配逻辑还是不智能，如果空闲的桌子多了,人数少了不容易将他们分配到相同的桌子里面快速组局
	table, err := self.room.GetUsableTable()

	fmt.Println("看看桌子")

	if err == nil {
		table.Create()
		tableInfo := map[string]interface{}{
			"BigRoomId": room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId()),
		}
		fmt.Println(tableInfo)
		return tableInfo, ""
	} else {
		return nil, "There is no available table"
	}
}














