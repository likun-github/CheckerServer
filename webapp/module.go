/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package webapp

import (
	"CheckerServer/server/dao"
	"CheckerServer/server/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Module = func() *Web {
	web := new(Web)
	return web
}
type UserInfoJson struct {
	Id int64 `json:"userid"`//用户id
	Name string `json:"name"`//用户真实姓名
	WxName string `json:"nickname"`//微信昵称
	WXImg string `json:"pic"`//头像
	Status int8  `json:"status"`//用户状态，0代表仅获取openid,1代表获取基本用户信息
	Score int64 `json:"score"`//分数
	Level int8 `json:"level"`//关卡

}
type ChessJson struct {
	Id int64 `json:"chessid"`//用户id
	White int64 `json:"white"`//用户真实姓名
	Black int64 `json:"black"`//微信昵称
	King int64 `json:"king"`//头像


}

type Web struct {
	basemodule.BaseModule
}

func (self *Web) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Webapp"
}
func (self *Web) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *Web) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		//[26/Oct/2017:19:07:04 +0800]`-`"GET /g/c HTTP/1.1"`"curl/7.51.0"`502`[127.0.0.1]`-`"-"`0.006`166`-`-`127.0.0.1:8030`-`0.000`xd
		log.Info("%s %s %s [%s] in %v", r.Method, r.URL.Path, r.Proto, r.RemoteAddr, time.Since(start))
	})
}
func Statushandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
func (self *Web) Run(closeSig chan bool) {
	//这里如果出现异常请检查8080端口是否已经被占用
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Error("webapp server error", err.Error())
		return
	}
	go func() {
		log.Info("webapp server Listen : %s", ":8081")
		root := mux.NewRouter()
		root.HandleFunc("/", HomeHandler);
		root.HandleFunc("/login", LoginHandler);//


		root.HandleFunc("/getinfo",GetInfoHandler)
		root.HandleFunc("/getuserid",GetUseridHandler)
		root.HandleFunc("/register",RegisterHandler)
		root.HandleFunc("/ChessManual",ChessManualHandler)
		status := root.PathPrefix("/status")
		status.HandlerFunc(Statushandler)

		static := root.PathPrefix("/checkerserver/")
		static.Handler(http.StripPrefix("/checkerserver/", http.FileServer(http.Dir(self.GetModuleSettings().Settings["StaticPath"].(string)))))
		//r.Handle("/static",static)
		ServeMux := http.NewServeMux()
		ServeMux.Handle("/", root)
		http.Serve(l, loggingHandler(ServeMux))
	}()
	<-closeSig
	log.Info("webapp server Shutting down...")
	l.Close()
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	for k,_ := range vars {
		log.Info(k)
	}
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "login")
}
//根据code获取用户id，没有的新建用户
func GetUseridHandler(writer http.ResponseWriter, request *http.Request) {
	appsecret:="5ee539127cb87ad7294f491648bc401c"
	appid:="wxa54181747176608d"
	request.ParseForm() //解析参数，默认是不会解析的
	code:=strings.Join(request.Form["code"], "")//解析code

	url:="https://api.weixin.qq.com/sns/jscode2session?appid="+appid+"&secret="+appsecret+"&js_code="+code+"&grant_type=authorization_code"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	s,err:=ioutil.ReadAll(resp.Body)
	b :=string(s)
	m := make(map[string]string)
	json.Unmarshal([]byte(b), &m)
	openid:=m["openid"]//从微信后台获取openid

	infoDao:=dao.NewUserInfoDao()
	userInfo:=infoDao.SelectByOpenid(openid)//在用户信息表中查询openid
	//如果无此openid，新建用户
	if userInfo==nil{
		user := new(model.UserInfo)
		user.WxOpenId=openid
		user.Status=0
		user.Level=1
		if !infoDao.Insert(user) {
			log.Info("db error")
		}
		u:=&UserInfoJson{Id:user.Id,Name:user.Name,WxName:user.WxName,WXImg:user.WXImg,Status:user.Status,Score:user.Score,Level:user.Level}
		j,_:=json.Marshal(u)
		fmt.Fprintf(writer, string(j))

	}else {
		//转json
		u:=&UserInfoJson{Id:userInfo.Id,Name:userInfo.Name,WxName:userInfo.WxName,WXImg:userInfo.WXImg,Status:userInfo.Status,Score:userInfo.Score,Level:userInfo.Level}
		j,_:=json.Marshal(u)
		fmt.Fprintf(writer, string(j))

	}



}
//根据用户id获取用户信息
func GetInfoHandler(writer http.ResponseWriter, request *http.Request)  {
	query := request.URL.Query()
	userId,_ := strconv.Atoi(query["userid"][0])
	infoDao := dao.NewUserInfoDao()
	userInfo := infoDao.SelectById(int64(userId))
	u:=&UserInfoJson{Id:userInfo.Id,Name:userInfo.Name,WxName:userInfo.WxName,WXImg:userInfo.WXImg,Status:userInfo.Status,Score:userInfo.Score,Level:userInfo.Level}
	j,_:=json.Marshal(u)
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, string(j))



}
//添加用户信息
func RegisterHandler(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	userid,_:=strconv.Atoi(query["userid"][0])
	pic:=query["pic"][0]
	nickname:=query["nickname"][0]
	infoDao:=dao.NewUserInfoDao()
	userInfo := infoDao.SelectById(int64(userid))
	userInfo.WXImg=pic
	userInfo.WxName=nickname
	userInfo.Status=1
	if !infoDao.Update(userInfo){
		log.Info("db error")
	}
	//转json
	u:=&UserInfoJson{Id:userInfo.Id,Name:userInfo.Name,WxName:userInfo.WxName,WXImg:userInfo.WXImg,Status:userInfo.Status,Score:userInfo.Score,Level:userInfo.Level}
	j,_:=json.Marshal(u)
	fmt.Fprintf(writer, string(j))


}


//添加用户信息
func ChessManualHandler(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	passid,_:=strconv.Atoi(query["passid"][0])
	chessdao:=dao.NewChessManualDao()
	chessinfo:=chessdao.SelectById(int64(passid))



	u:=&ChessJson{Id:chessinfo.Id,White:chessinfo.White,Black:chessinfo.Black,King:chessinfo.King}
	j,_:=json.Marshal(u)
	fmt.Fprintf(writer, string(j))


}
func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	for k,_ := range vars {
		log.Info(k)
	}
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "wellcome to home")
}

func (self *Web) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}
