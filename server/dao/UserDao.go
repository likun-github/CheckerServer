package dao

import (
	"CheckerServer/server/database"
	"CheckerServer/server/model"
	"github.com/liangdas/mqant/log"
)

type UserDao struct {
	Dao
}
func NewUserDao() (dao *UserDao) {
	return &UserDao{Dao{Engine:database.Engine}}
}
func (this *UserDao)SelectById(id int64) (user *model.User) {
	u := new(model.User)
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

func (this *UserDao)SelectAll() (users []model.User) {
	users = make([]model.User, 0)
	err := this.Engine.Find(&users)
	if err != nil {
		log.Error("select all user error, %s", err.Error())
	}
	return
}

func (this *UserDao)InsertUsers(users []model.User) bool{
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

func (this *UserDao) SelectUserByName(name string) *model.User {
	user := new(model.User)
	has, err := this.Engine.Where("name=?", name).Get(user)
	if err!=nil{
		log.Error("select user %s error", name)
	}
	if !has {
		log.Info("user name = %s not exist", name)
		return nil
	}
	return user
}