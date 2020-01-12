package main

import (
	"CheckerServer/server/database"
	mgate "CheckerServer/server/gate"
	"CheckerServer/server/helloworld"
	"CheckerServer/server/jump"
	"CheckerServer/server/login"
	"CheckerServer/server/test"
	"CheckerServer/server/xaxb"
	"CheckerServer/webapp"
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/module/modules"
	"github.com/nats-io/nats.go"
	//"github.com/liangdas/mqant-modules/tracing"
	"fmt"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/registry"
	"github.com/liangdas/mqant/selector"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
)


func main() {
	//监听地址，不大懂
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	//相应配置
	rs := registry.DefaultRegistry //etcdv3.NewRegistry()
	nc, err := nats.Connect(nats.DefaultURL, nats.MaxReconnects(10000))
	if err != nil {

	}
	//数据库初始化
	if !database.DbInit() {
		log.Info("Exiting")
		return
	}
	//启动mqant
	app := mqant.CreateApp(
		module.Debug(true),//只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		module.Nats(nc),
		module.Registry(rs),
	)
	app.Options().Selector.Init(selector.SetStrategy(func(services []*registry.Service) selector.Next {
		var nodes []*registry.Node

		// Filter the nodes for datacenter
		for _, service := range services {
			for _, node := range service.Nodes {
				nodes = append(nodes, node)
				//if node.Metadata["type"] == "helloworld" {
				//	nodes = append(nodes, node)
				//}
			}
		}

		var mtx sync.Mutex
		//log.Info("services[0] $v",services[0].Nodes[0])
		return func() (*registry.Node, error) {
			mtx.Lock()
			defer mtx.Unlock()
			if len(nodes) == 0 {
				return nil, fmt.Errorf("no node")
			}
			index := rand.Intn(int(len(nodes)))
			return nodes[index], nil
		}
	}))
	//app.Route("Chat",ChatRoute)
	//启动程序，接口命名
	app.Run( //只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		modules.MasterModule(),
		//hitball.Module(),
		mgate.Module(), //这是默认网关模块,是必须的支持 TCP,websocket,MQTT协议
		helloworld.Module(),
		login.Module(), //这是用户登录验证模块

		jump.Module(),//跳棋模块
		//chat.Module(),
		//user.Module(),
		webapp.Module(),
		xaxb.Module(),
		test.Module(),
		//tracing.Module(), //很多初学者不会改文件路径，先移除了
	) //这是聊天模块

}
