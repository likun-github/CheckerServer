package jump

import (
	"encoding/json"
)

/**
定期刷新所有玩家的位置
*/
/*
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
 */

/**
通知所有玩家有新玩家加入
*/
/*
func (self *Table) NotifyJoin(player *objects.Player) {
	b, _ := json.Marshal(player.SerializableMap())
	self.NotifyCallBackMsg("PlayerInfo", b)
}
 */

///////////////////////////////////////////////////////////////////////////////////////////////////////大爷我写的……
/**
通知所有玩家匹配完成，游戏开始
返回整桌的玩家信息，包括：userid, name, avatar, score, level
*/
func (self *Table) NotifyMatchFinish() {
	players_info := map[string]interface{} {
		"w_uid": 	self.seats[0].UserId,
		"w_name":	self.seats[0].Username,
		"w_avatar":	self.seats[0].Avatar,
		"w_score":  self.seats[0].Score,
		"w_level":  self.seats[0].Level,

		"b_uid": 	self.seats[1].UserId,
		"b_name":	self.seats[1].Username,
		"b_avatar":	self.seats[1].Avatar,
		"b_score":  self.seats[1].Score,
		"b_level":  self.seats[1].Level,
	}
	players_info_json,_ := json.Marshal(players_info)
	self.NotifyCallBackMsg("MatchFinish", players_info_json)
}

/**
通知非控制方玩家走子完成，轮到其走子
返回控制方玩家的走子信息，包括：composition
*/
func (self *Table) NotifyUpdateComposition() {
	new_composition := map[string]interface{} {
		"W": self.composition.Top().(*Chess).white,
		"B": self.composition.Top().(*Chess).black,
		"K": self.composition.Top().(*Chess).king,
	}
	new_composition_json,_ := json.Marshal(new_composition)
	self.seats[self.currentPlayer].Session().Send("UpdateComposition", new_composition_json)
}

/**
通知非控制方玩家悔棋信息，轮到其判断是否同意
返回控制方玩家的索引，包括：withdraw_requeseted_player(用于验证)
*/
func (self *Table) NotifyWithdrawRequested() {
	withdraw_requeseted_player := -1
	if self.currentPlayer == 0 {
		withdraw_requeseted_player = 1
	} else {
		withdraw_requeseted_player = 0
	}
	/*
	withdraw_requested := map[string]interface{} {
		"withdraw_requeseted_player": withdraw_requeseted_player,
	}
	withdraw_requested_json,_ := json.Marshal(withdraw_requested)*/
	self.seats[withdraw_requeseted_player].Session().Send("WithdrawRequested", nil)
}

/**
通知非控制方玩家悔棋超时，直接认为非控制方玩家同意悔棋
返回空值
*/
func (self *Table) NotifyWithdrawTimeout() {
	withdraw_timeout_player := -1
	if self.currentPlayer == 0 {
		withdraw_timeout_player = 1
	} else {
		withdraw_timeout_player = 0
	}
	withdraw_timeout_json,_ := json.Marshal(nil)
	self.seats[withdraw_timeout_player].Session().Send("WithdrawTimeout", withdraw_timeout_json)
}

/**
通知控制方玩家悔棋结果，不论结果如何，控制方玩家继续走子
返回非控制方玩家的悔棋结果，包括：withdraw_agreed
*/
func (self *Table) NotifyWithdrawDecided() {
	withdraw_decided := map[string]interface{}{
		"withdraw_agreed": self.withdraw_agreed,
	}
	withdraw_decided_json,_ := json.Marshal(withdraw_decided)
	self.seats[self.currentPlayer].Session().Send("WithdrawDecided", withdraw_decided_json)
}

/**
通知非控制方玩家控制方玩家求和棋，轮到其判断是否同意
返回空值
*/
func (self *Table) NotifyDrawRequested() {
	draw_requeseted_player := -1
	if self.currentPlayer == 0 {
		draw_requeseted_player = 1
	} else {
		draw_requeseted_player = 0
	}
	/*
	draw_requested := map[string]interface{} {
		"draw_requeseted_player": draw_requeseted_player,
	}
	draw_requested_json,_ := json.Marshal(draw_requested)
	 */
	self.seats[draw_requeseted_player].Session().Send("DrawRequested", nil)
}

/**
通知非控制方玩家和棋超时，直接认为非控制方玩家同意和棋
返回空值
*/
func (self *Table) NotifyDrawTimeout() {
	draw_timeout_player := -1
	if self.currentPlayer == 0 {
		draw_timeout_player = 1
	} else {
		draw_timeout_player = 0
	}
	draw_timeout_json,_ := json.Marshal(nil)
	self.seats[draw_timeout_player].Session().Send("DrawTimeout", draw_timeout_json)
}

/**
通知控制方玩家非控制方拒绝和棋
返回空值
*/
func (self *Table) NotifyDrawDenied() {
	draw_denied_json,_ := json.Marshal(nil)
	self.seats[self.currentPlayer].Session().Send("DrawDenied", draw_denied_json)
	self.draw_agreed = -1
}

/**
通知控制方玩家非控制方同意和棋
返回空值
*/
func (self *Table) NotifyDrawAgreed() {
	draw_agreed_json,_ := json.Marshal(nil)
	self.seats[self.currentPlayer].Session().Send("DrawAgreed", draw_agreed_json)
	self.draw_agreed = -1
}

/**
通知另一方玩家对方认输辽
返回空值
*/
func (self *Table) NotifyOpponentLost() {
	the_other_side := -1
	if(self.lose_requested == 0) { // 白方认输，通知黑方
		the_other_side = 1
	} else if (self.lose_requested == 1) { // 黑方认输，通知白方
		the_other_side = 0
	}
	the_other_side_lost_json,_ := json.Marshal(nil)
	self.seats[the_other_side].Session().Send("OpponentLost", the_other_side_lost_json)
	self.lose_requested = -1
}

/**
通知所有玩家游戏结果
返回两个玩家的信息，包括：userid, score, level, result:0(输),1(赢),2（平）
*/
func (self *Table) NotifyResult() {
	game_result := map[string]interface{} {
		"w_uid": 	self.seats[0].UserId,
		"w_score":  self.seats[0].Score,
		"w_level":  self.seats[0].Level,
		"w_result":	self.seats[0].Result,

		"b_uid": 	self.seats[1].UserId,
		"b_score":  self.seats[1].Score,
		"b_level":  self.seats[1].Level,
		"b_result":	self.seats[1].Result,
	}
	game_result_json,_ := json.Marshal(game_result)
	self.NotifyCallBackMsg("GameEnd", game_result_json)
}
///////////////////////////////////////////////////////////////////////////////////////////////////////……框架，谁都别动






/**
通知所有玩家进入匹配完成期了
*/
/*
func (self *Table) NotifyMatch() {
	b, _ := json.Marshal(map[string]interface{}{
		"match": true,
	})
	self.NotifyCallBackMsg("Jump/Match", b)
}
 */

/**
通知所有玩家开始游戏了
*/
/*
func (self *Table) NotifyResume() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("Jump/OnResume", b)
}
 */

/**
通知所有玩家开始游戏了
*/
/*
func (self *Table) NotifyPause() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnPause", b)
}
*/
/**
通知所有玩家开始游戏了
*/
/*
func (self *Table) NotifyStop() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnStop", b)
}
 */

/**
通知所有玩家进入空闲期了
*/
/*
func (self *Table) NotifyIdle() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Idle", b)
}
 */

/**
通知所有玩家开始押注了
*/
/*
func (self *Table) NotifyBetting() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Betting", b)
}
 */

/**
通知所有玩家开始开奖了
*/
/*
func (self *Table) NotifyOpening() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("Jump/Opening", b)
}
 */

/**
通知所有玩家开奖结果出来了
*/
/*
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
 */


