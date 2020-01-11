package dao
import (
	"CheckerServer/server/database"
	"CheckerServer/server/model"
	"github.com/liangdas/mqant/log"
)

type ChessMenuDao struct {
	Dao
}
func NewChessDao() (dao *ChessMenuDao) {
	return &ChessMenuDao{Dao{Engine:database.Engine}}
}
func (this *ChessMenuDao)SelectById(id int64) (user *model.ChessMenu) {
	u := new(model.ChessMenu)
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
func (this *ChessMenuDao)SelectAll() (chesses []model.ChessMenu) {
	chesses = make([]model.ChessMenu, 0)
	err := this.Engine.Find(&chesses)
	if err != nil {
		log.Error("select all chesses error, %s", err.Error())
	}
	return
}

