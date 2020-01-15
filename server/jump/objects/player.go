package objects
import (
	"encoding/json"
	"github.com/liangdas/mqant-modules/room"
)

type Player struct {
	room.BasePlayerImp  //继承父类，
	SeatIndex  int	//id,0为白子方，1为黑子方
	Time       int//计时10min,600s
	BackNumber int //悔棋次数
	Score      int64 //分数
	Username   string//用户昵称
	Avatar     string//用户头像

}
//新建用户，基本属性
func NewPlayer(SeatIndex int) *Player {
	this := new(Player)
	this.SeatIndex = SeatIndex
	this.Time=600
	this.BackNumber=3
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
		//"Coin":      this.Coin,
		//"Stake":     this.Stake,
		//"Target":    this.Target,
		//"Weight":    this.Weight,
		"SitDown":   this.SitDown(),
	}
}
