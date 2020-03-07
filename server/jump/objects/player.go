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

// 计算积分的放大系数K
func (this *Player) K() int {
	switch {
	case this.Score < 1000:
		return 	120
	case this.Score < 1399 && this.Score >= 1000:
		return 60
	case this.Score < 1799 && this.Score >= 1400:
		return 30
	case this.Score < 1999 && this.Score >= 1800:
		return 25
	case this.Score < 2199 && this.Score >= 2000:
		return 20
	case this.Score < 2399 && this.Score >= 2200:
		return 15
	case this.Score >= 2400:
		return 10
	}
	return 0
}

// 计算玩家的等级
func (this *Player) GetLevel() int8 {
	switch {
	case this.Score < 1100:
		return 	0
	case this.Score < 1200 && this.Score >= 1100:
		return 1
	case this.Score < 1300 && this.Score >= 1200:
		return 2
	case this.Score < 1400 && this.Score >= 1300:
		return 3
	case this.Score < 1500 && this.Score >= 1400:
		return 4
	case this.Score < 1600 && this.Score >= 1500:
		return 5
	case this.Score < 1700 && this.Score >= 1600:
		return 6
	case this.Score < 1800 && this.Score >= 1700:
		return 7
	case this.Score < 2000 && this.Score >= 1800:
		return 8
	case this.Score < 2200 && this.Score >= 2000:
		return 9
	case this.Score < 2400 && this.Score >= 2200:
		return 10
	case this.Score < 2600 && this.Score >= 2400:
		return 11
	case this.Score >= 2600:
		return 12
	}
	return 0
}
