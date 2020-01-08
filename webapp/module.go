/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package webapp

import (
	"CheckerServer/server/dao"
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
		root.HandleFunc("/login", LoginHandler);
		root.HandleFunc("/getinfo",GetInfoHandler)
		root.HandleFunc("/getuserid",GetUseridHandler)
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

func GetUseridHandler(writer http.ResponseWriter, request *http.Request) {
	appsecret:="5ee539127cb87ad7294f491648bc401c"
	appid:="wxa54181747176608d"
	request.ParseForm() //解析参数，默认是不会解析的
	code:=strings.Join(request.Form["code"], "")
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
	openid:=m["openid"]
	infoDao:=dao.NewUserInfoDao()
	userinfo:=infoDao.SelectByOpenid(openid)
	fmt.Println("我就是喜欢乱写")
	fmt.Println(userinfo.Id)
	fmt.Fprintf(writer, " %s",openid)

}
func GetInfoHandler(writer http.ResponseWriter, request *http.Request)  {
	query := request.URL.Query()
	userId,_ := strconv.Atoi(query["userid"][0])
	infoDao := dao.NewUserInfoDao()
	userInfo := infoDao.SelectById(int64(userId))
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "%s, %d",userInfo.Name, userInfo.Level)

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
