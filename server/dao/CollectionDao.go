package dao

import (
	"CheckerServer/server/common/stack"
	"CheckerServer/server/database"
	"CheckerServer/server/model"
	"github.com/liangdas/mqant/log"
	"CheckerServer/server/jump"
	"strconv"
)

type CollectionDao struct {
	Dao
}

func NewCollectionDao() (dao *CollectionDao) {
	return &CollectionDao{Dao{Engine:database.Engine}}
}

// 新增一条collection记录
func (this *CollectionDao)InsertCollection(userid int64, step int8, composition *stack.Stack) bool{
	collection := new(model.Collection)
	collection.UserId = userid
	collection.Step = step
	var white []int64
	var black []int64
	var king []int64
	for i := 0; i<composition.Size(); i++ {
		comp := composition.Get(i)
		w,_ := strconv.ParseInt(comp.(*jump.Chess).White, 2, 64)
		b,_ := strconv.ParseInt(comp.(*jump.Chess).White, 2, 64)
		k,_ := strconv.ParseInt(comp.(*jump.Chess).White, 2, 64)
		white=append(white,w)
		black=append(white,b)
		king=append(white,k)
	}
	collection.White = white
	collection.Black = black
	collection.King = king
	// 插入新数据
	if !this.Insert(collection) {
		log.Info("Collection插入失败")
		return false
	} else {
		return true
	}
}

// 选取userid对应用户所收藏的棋局
func (this *CollectionDao)SelectByUserId(UserId int64)(collections *model.Collection) {
	c:=new(model.Collection)
	has, err := this.Engine.Where("userId=?", UserId).Get(c)
	if err!=nil{
		log.Error("select user id=%s ", UserId)
		return nil
	}
	if !has {
		log.Info("user id = %s not exist", UserId)
		return nil
	}
	return c
}