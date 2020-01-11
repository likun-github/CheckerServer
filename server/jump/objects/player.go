package objects
import (
	"encoding/json"
	"github.com/liangdas/mqant-modules/room"
)

type Player struct {
	room.BasePlayerImp  //继承父类，
	SeatIndex  int
	Time       int//计时10min,600s

	Controller bool //是否为控制方
	BackNumber int //悔棋次数
	Score      int //分数



	//Coin       int //金币数量
	//timeToMove int64
	//Target     int64 //押注目标
	//Stake      bool  //是否已押注
	//Weight     int64 //计算后权重
}

func NewPlayer(SeatIndex int) *Player {
	this := new(Player)
	this.SeatIndex = SeatIndex
	this.BackNumber=3
	this.Session().GetUserId()
	//this.Coin = 1000
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
		//"Coin":      this.Coin,
		//"Stake":     this.Stake,
		//"Target":    this.Target,
		//"Weight":    this.Weight,
		"SitDown":   this.SitDown(),
	}
}
