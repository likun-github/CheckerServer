package test

import (
	"CheckerServer/server/dao"
	"CheckerServer/server/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"math/rand"
	"time"
)

var Module = func() module.Module {
	gate := new(Test)
	return gate
}

type Test struct {
	basemodule.BaseModule
}

func (m *Test) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Test"
}
func (m *Test) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *Test) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)

	m.GetServer().RegisterGO("HD_Login", m.login) //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().RegisterGO("HD_FindUserByName", m.findUserByName)
	m.GetServer().RegisterGO("HD_AddUser", m.addUser)
	m.GetServer().RegisterGO("HD_RemoveUser", m.removeUser)
	m.GetServer().RegisterGO("HD_ModifyUser", m.modifyUser)
	m.GetServer().RegisterGO("track", m.track)    //演示后台模块间的rpc调用
	m.GetServer().RegisterGO("track2", m.track2)  //演示后台模块间的rpc调用
	m.GetServer().RegisterGO("track3", m.track3)  //演示后台模块间的rpc调用
	m.GetServer().Register("HD_TestRobot", m.robot)
	m.GetServer().RegisterGO("HD_TestRobot_GO", m.robot) //我们约定所有对客户端的请求都以Handler_开头
}

func (m *Test) Run(closeSig chan bool) {
}

func (m *Test) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}
func (m *Test) robot(session gate.Session, msg map[string]interface{}) (result string, err string) {
	//time.Sleep(1)
	return "sss", ""
}
func (m *Test) login(session gate.Session, msg map[string]interface{}) (result string, err string) {
	if msg["userName"] == nil || msg["passWord"] == nil {
		result = "userName or passWord cannot be nil"
		return
	}
	userName := msg["userName"].(string)
	passWord := msg["passWord"].(string)
	var(
		id int64
		name string
	)
	if passWord != "" {
		//stmt,e := database.Db.Prepare("select id,name from User where id=?")
		//
		//if e != nil {
		//	log.Error(e.Error())
		//	result = "pass error"
		//	return
		//}
		//row := stmt.QueryRow(1)
		//defer stmt.Close()
		//row.Scan(&id, &name)
	}
	err = session.Bind(userName)
	if err != "" {
		return
	}
	session.Set("login", "true")
	session.Push() //推送到网关
	return fmt.Sprintf("login success %d, name = %s", id, name), ""
}

func (m *Test) findUserByName(session gate.Session, msg map[string]interface{}) (result string, err string){
	if msg["userName"] == nil {
		result = "userName cannot be nil"
		return
	}
	userName := msg["userName"].(string)
	//passWord := msg["passWord"].(string)
	if userName != session.GetUserId() || session.Get("login") != "true" {
		result = "user not login"
		return
	}
	userDao:=dao.NewUserDao()
	user := userDao.SelectUserByName(userName)
	if user ==nil {
		result = "invalid username"
		return
	}

	log.Info("find user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	result = fmt.Sprintf("find user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	return
}

func (m *Test) addUser(session gate.Session, msg map[string]interface{}) (result string, err string){
	if msg["userName"] == nil || msg["passWord"] == nil {
		result = "userName or passWord cannot be nil"
		return
	}
	userName := msg["userName"].(string)
	passWord := msg["passWord"].(string)
	userDao:=dao.NewUserDao()
	user_ := userDao.SelectUserByName(userName)
	if user_ != nil {
		result = "userName exists"
		return
	}
	user := new(model.User)
	user.Name = userName
	user.Password = passWord
	if msg["age"] != nil {
		user.Age = int8(msg["age"].(float64))
	}
	if msg["tel"]!=nil {
		user.Tel = int64(msg["tel"].(float64))
	}
	if msg["email"] != nil {
		user.Email = msg["email"].(string)
	}

	if !userDao.Insert(user) {
		result = "db error"
		return
	}

	log.Info("add user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	result = fmt.Sprintf("add user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	return
}

func (m *Test) modifyUser(session gate.Session, msg map[string]interface{}) (result string, err string){
	if msg["userName"] == nil {
		result = "userName cannot be nil"
		return
	}

	userName := msg["userName"].(string)
	if userName != session.GetUserId() || session.Get("login") != "true" {
		result = "user not login"
		return
	}
	//passWord := msg["passWord"].(string)
	userDao:=dao.NewUserDao()
	user := userDao.SelectUserByName(userName)
	if user ==nil {
		result = "invalid username"
		return
	}

	if msg["age"] != nil {
		user.Age = int8(msg["age"].(float64))
	}
	if msg["tel"]!=nil {
		user.Tel = int64(msg["tel"].(float64))
	}
	if msg["email"] != nil {
		user.Email = msg["email"].(string)
	}

	if !userDao.Update(user) {
		result = "db error"
		return
	}

	log.Info("modify user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	result = fmt.Sprintf("modify user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	return
}

func (m *Test) removeUser(session gate.Session, msg map[string]interface{}) (result string, err string){
	if msg["userName"] == nil {
		result = "userName cannot be nil"
		return
	}

	userName := msg["userName"].(string)
	if userName != session.GetUserId() || session.Get("login") != "true" {
		result = "user not login"
		return
	}
	userDao:=dao.NewUserDao()
	user := userDao.SelectUserByName(userName)
	if user ==nil {
		result = "invalid username"
		return
	}

	if !userDao.Delete(user) {
		result = "db error"
		return
	}

	log.Info("delete user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	result = fmt.Sprintf("delete user: id=%d, name=%s, age=%d, tel=%d, email=%s", user.Id, user.Name, user.Age, user.Tel, user.Email)
	return
}

func (m *Test) track(session gate.Session) (result string, err string) {
	//演示后台模块间的rpc调用
	time.Sleep(time.Millisecond * 10)
	log.TInfo(session, "Login %v", "track1")
	m.RpcInvoke("Login", "track2", session)
	return fmt.Sprintf("My is Login Module %s"), ""
}

func (m *Test) track2(session gate.Session) (result string, err string) {
	//演示后台模块间的rpc调用
	time.Sleep(time.Millisecond * 10)
	log.TInfo(session, "Login %v", "track2")
	r := rand.Intn(100)
	if r > 30 {
		m.RpcInvoke("Login", "track3", session)
	}

	return fmt.Sprintf("My is Login Module"), ""
}
func (m *Test) track3(session gate.Session) (result string, err string) {
	//演示后台模块间的rpc调用
	time.Sleep(time.Millisecond * 10)
	log.TInfo(session, "Login %v", "track3")
	return fmt.Sprintf("My is Login Module"), ""
}
