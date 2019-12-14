package database

import (
	"CheckerServer/server/model"
	_ "database/sql"
	_ "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liangdas/mqant/log"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"os"
	"time"
)

//var Db *sql.DB
var Engine *xorm.Engine

func DbInit() bool {

	engine, err := xorm.NewEngine("mysql", "root:test666@tcp(127.0.0.1:3306)/Test?charset=utf8")
	if err != nil {
		log.Error("can not open db")
		return false
	}
	Engine =engine

	//Db = db
	Engine.SetMaxOpenConns(1000)
	Engine.SetMaxIdleConns(100)
	Engine.SetConnMaxLifetime(10 * time.Minute)
	Engine.ShowSQL(true)
	err = Engine.Ping()
	//rows, err := Db.Query("select * from User")
	if err != nil {
		log.Error("db error")
		return false
	}
	engine.Logger().SetLevel(core.LOG_DEBUG)

	f, err := os.Create("bin/logs/sql.log")
	if err != nil {
		log.Error(err.Error())
		return false
	}
	engine.SetLogger(xorm.NewSimpleLogger(f))

	log.Info("open db success")
	err = engine.Sync2(new(model.User), new(model.Record), new(model.UserInfo))
	if err!= nil {
		log.Error(err.Error())
		return false
	}
	return true
}

