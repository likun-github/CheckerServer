package dao
import (
	"CheckerServer/server/database"
	"CheckerServer/server/model"
	"github.com/liangdas/mqant/log"
)

type ChessManualDao struct {
	Dao
}
func NewChessManualDao() (dao *ChessManualDao) {
	return &ChessManualDao{Dao{Engine:database.Engine}}
}
func (this *ChessManualDao)SelectById(id int64) (user *model.ChessManual) {
	u := new(model.ChessManual)
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
func (this *ChessManualDao)SelectAll() (chesses []model.ChessManual) {
	chesses = make([]model.ChessManual, 0)
	err := this.Engine.Find(&chesses)
	if err != nil {
		log.Error("select all chesses error, %s", err.Error())
	}
	return
}

