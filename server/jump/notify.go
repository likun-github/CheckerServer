package jump



import (
	"CheckerServer/server/jump/objects"
	"encoding/json"
)

/**
定期刷新所有玩家的位置
*/
func (self *Table) NotifyAxes() {
	seats := []map[string]interface{}{}
	for _, player := range self.seats {
		if player.Bind() {
			seats = append(seats, player.SerializableMap())
		}
	}
	b, _ := json.Marshal(map[string]interface{}{
		"State":     self.State(),
		"StateGame": self.fsm.getState(),
		"Seats":     seats,
	})
	self.NotifyCallBackMsg("XaXb/OnSync", b)
}

/**
通知所有玩家有新玩家加入
*/
func (self *Table) NotifyJoin(player *objects.Player) {
	b, _ := json.Marshal(player.SerializableMap())
	self.NotifyCallBackMsg("XaXb/OnEnter", b)
}

/**
通知所有玩家开始游戏了
*/
func (self *Table) NotifyResume() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnResume", b)
}

/**
通知所有玩家开始游戏了
*/
func (self *Table) NotifyPause() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnPause", b)
}

/**
通知所有玩家开始游戏了
*/
func (self *Table) NotifyStop() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnStop", b)
}
/**
通知所有玩家进入匹配完成期了
*/
func (self *Table) NotifyMatch() {
	b, _ := json.Marshal(map[string]interface{}{
		"match": true,
	})
	self.NotifyCallBackMsg("Jump/Match", b)
}
/**
通知所有玩家进入空闲期了
*/
func (self *Table) NotifyIdle() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Idle", b)
}

/**
通知所有玩家开始押注了
*/
func (self *Table) NotifyBetting() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Betting", b)
}

/**
通知所有玩家开始开奖了
*/
func (self *Table) NotifyOpening() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Opening", b)
}

/**
通知所有玩家开奖结果出来了
*/
func (self *Table) NotifySettlement(Result int64) {
	seats := []map[string]interface{}{}
	for _, player := range self.seats {
		if player.Bind() {
			seats = append(seats, player.SerializableMap())
		}
	}
	b, _ := json.Marshal(map[string]interface{}{
		"Result": Result,
		"Seats":  seats,
	})
	self.NotifyCallBackMsg("XaXb/Settlement", b)
}

