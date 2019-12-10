package dao

import (
	"github.com/liangdas/mqant/log"
	"github.com/xormplus/xorm"
)

type Dao struct {
	Engine *xorm.Engine
}

func (this *Dao) Insert(m interface{}) bool{
	_, err :=this.Engine.InsertOne(m)
	if err != nil {
		log.Error("insert error: ",err.Error())
		return false
	}
	return true
}

func (this *Dao)Delete(m interface{}) bool {
	_, err :=this.Engine.Delete(m)
	if err != nil {
		log.Error("delete error: ", err.Error())
		return false
	}
	return true
}

func (this *Dao)Update(m interface{}) bool {
	_, err :=this.Engine.Update(m)
	if err != nil {
		log.Error("update error: ", err.Error())
		return false
	}
	return true
}
