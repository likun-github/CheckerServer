package dao

import (
	"CheckerServer/server/database"
	"CheckerServer/server/model"
	"github.com/liangdas/mqant/log"
)

type CollectionDao struct {
	Dao
}

func NewCollectionDao() (dao *CollectionDao) {
	return &CollectionDao{Dao{Engine:database.Engine}}
}

func (this *CollectionDao)InsertCollection(collections []model.Collection) bool{
	session := this.Engine.NewSession()
	err := session.Begin()
	if err!= nil {
		log.Error("begin session error, %s", err.Error())
		return false
	}
	for _, collection := range collections {
		_, err = session.InsertOne(collection)
		if err != nil {
			_ = session.Rollback()
			log.Error("insert user error, rollback. %s", err.Error())
			return false
		}
	}
	_ = session.Commit()
	return true
}

func (this *CollectionDao)SelectByUserId(userid string)(collections *model.Collection) {
	c:=new(model.Collection)
	has, err := this.Engine.Where("userId=?", userid).Get(c)
	if err!=nil{
		log.Error("select user id=%s ", userid)
		return nil
	}
	if !has {
		log.Info("user id = %s not exist", userid)
		return nil
	}
	return c
}