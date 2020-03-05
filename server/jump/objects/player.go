package objects

import (
	"encoding/json"
	"github.com/liangdas/mqant-modules/room"
)

type Player struct {
	room.BasePlayerImp  		//继承父类，
	SeatIndex  int				//玩家在桌子里的id,0先手，1后手
	Time       int				//计时10min,600s
	WithdrawNumber int 			//悔棋次数
	UserId		int64			//数据库里的id
	Score      int64 			//分数
	Level      int8 			//段位
	Username   string			//用户昵称
	Avatar     string			//用户头像
	Result  	int				//游戏结果:-1(还没出),0(输),1(赢)

}
//新建用户，基本属性
func NewPlayer(SeatIndex int) *Player {
	this := new(Player)
	this.SeatIndex = SeatIndex

	this.Result=-1
	this.UserId = -1
	this.Score = -1
	this.Level = -1
	this.Username = ""
	this.Avatar = ""

	this.Time=600
	this.WithdrawNumber=3

	return this
}
//map转json

func (this *Player) Serializable() ([]byte, error) {

	return json.Marshal(this.SerializableMap())
}


//转化为map

func (this *Player) SerializableMap() map[string]interface{} {
	rid := ""
	if this.Session() != nil {
		rid = this.Session().GetUserId()
	}
	return map[string]interface{}{
		"SeatIndex": this.SeatIndex,
		"Rid":       rid,
		"Username":    this.Username,
		"SitDown":   this.SitDown(),

	}
}

