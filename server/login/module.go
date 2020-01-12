/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package login

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	gate := new(Login)
	return gate
}

type Login struct {
	basemodule.BaseModule
}

func (m *Login) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Login"
}
func (m *Login) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *Login) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)
	m.GetServer().RegisterGO("HD_Login", m.login) //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().Register("HD_Robot", m.robot)
	m.GetServer().RegisterGO("HD_Robot_GO", m.robot) //我们约定所有对客户端的请求都以Handler_开头
}

func (m *Login) Run(closeSig chan bool) {
}

func (m *Login) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}
func (m *Login) robot(session gate.Session, msg map[string]interface{}) (result string, err string) {
	//time.Sleep(1)

	fmt.Println("真的吗，就是简单的试一下")
	return "sss", ""
}
func (m *Login) login(session gate.Session, msg map[string]interface{}) (result string, err string) {
	if msg["userid"] == nil  {
		result = "userid cannot be nil"
		return
	}
	fmt.Println("尝试绑定")
	userid := msg["userid"].(string)
	fmt.Println(userid)
	//string := strconv.FormatFloat(userid,'E',-1,64)
	err = session.Bind(userid)
	// 直接创建1
	//m2 := make(map[string]string)
	//// 然后赋值
	//m2["a"] = "就是简单的尝试"
	//m2["b"] = "bb"
	//j,_:=json.Marshal(m2)
	//session.Send("try",j)
	if err != "" {
		return
	}
	session.Set("login", "true")
	session.Push() //推送到网关
	return fmt.Sprintf("login success %s", userid), ""
}

