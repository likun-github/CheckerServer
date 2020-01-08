package dao

import (
	"CheckerServer/server/database"
	"CheckerServer/server/model"
	"github.com/liangdas/mqant/log"
)

type UserInfoDao struct {
	Dao
}
func NewUserInfoDao() (dao *UserInfoDao) {
	return &UserInfoDao{Dao{Engine:database.Engine}}
}
func (this *UserInfoDao)SelectByOpenid(openid string)(user *model.UserInfo)  {
	u:=new(model.UserInfo)
	has, err := this.Engine.Where("userWXOpenID=?", openid).Get(u)
	if err!=nil{
		log.Error("select user openid=%s ", openid)
		return nil
	}
	if !has {
		log.Info("user openid = %s not exist", openid)
		return nil
	}
	return u
}

func (this *UserInfoDao)SelectById(id int64) (user *model.UserInfo) {
	u := new(model.UserInfo)
	has, err := this.Engine.Id(id).Get(u)
	if err!=nil{
		log.Error("select user id=%d error", id)
		return nil
	}
	if !has {
		log.Info("user id = %d not exist", id)
		return nil
	}
	return u
}

func (this *UserInfoDao)SelectAll() (users []model.UserInfo) {
	users = make([]model.UserInfo, 0)
	err := this.Engine.Find(&users)
	if err != nil {
		log.Error("select all user error, %s", err.Error())
	}
	return
}

func (this *UserInfoDao)InsertUserInfo(users []model.UserInfo) bool{
	session := this.Engine.NewSession()
	err := session.Begin()
	if err!= nil {
		log.Error("begin session error, %s", err.Error())
		return false
	}
	for _, user := range users {
		_, err = session.InsertOne(user)
		if err != nil {
			_ = session.Rollback()
			log.Error("insert user error, rollback. %s", err.Error())
			return false
		}
	}
	_ = session.Commit()
	return true
}