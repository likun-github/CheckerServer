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




func (self *Jump) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings, server.Metadata(map[string]string{
		"type": "helloworld",
	}))
	//房间号
	self.gameId = 13

	self.GetServer().RegisterGO("HD_Match", self.match) //我们约定所有对客户端的请求都以Handler_开头
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

func (self *Jump) login(session gate.Session, msg map[string]interface{}) (result string, err string) {
	m2 := make(map[string]string)
	// 然后赋值
	m2["a"] = "就是简单的尝试"
	m2["b"] = "bb"
	j,_:=json.Marshal(m2)
	session.Send("try",j)
	fmt.Println("试一试是否可以运行")
	session.GetUserId()

	return
}














