package jump

import (
	"CheckerServer/server/dao"
	"CheckerServer/server/jump/objects"
	"fmt"
	"math"
	//"github.com/liangdas/mqant/utils"
	//"math"
)

var (
	VoidPeriod 			  = FSMState("空档期")
	MatchPeriod			  = FSMState("匹配期")
	ControlPeriod 		  = FSMState("控制期")
	PlayFinishPeriod      = FSMState("行棋完成期")
	WithdrawPeriod		  = FSMState("悔棋期")
	DrawPeriod			  = FSMState("求和期")
	SettlementPeriod	  = FSMState("结算期")

	VoidPeriodEvent 	  = FSMEvent("进入空档期")
	MatchPeriodEvent	  = FSMEvent("进入匹配期")
	ControlPeriodEvent 	  = FSMEvent("进入控制期")
	PlayFinishPeriodEvent = FSMEvent("进入行棋完成期")
	WithdrawPeriodEvent	  = FSMEvent("进入悔棋期")
	DrawPeriodEvent		  = FSMEvent("进入求和期")
	SettlementPeriodEvent = FSMEvent("进入结算期")
)

// 这里的都是瞬时动作
func (this *Table) InitFsm() {
	this.fsm = *NewFSM(VoidPeriod)

	//（空档期）進入匹配完成期
	this.MatchPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("空档期转匹配期")
		this.start_match_step = this.current_frame

		return MatchPeriod
	})

	this.Match2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("匹配期转控制期")
		this.start_control_step = this.current_frame
		this.NotifyMatchFinish()	//通知所有玩家进入控制期，游戏开始

		return ControlPeriod
	})

	this.Control2WithdrawHandler = FSMHandler(func() FSMState {
		fmt.Println("控制期转悔棋期")
		this.start_withdraw_step = this.current_frame

		this.NotifyWithdrawRequested() //通知非控制方玩家决定是否同意悔棋
		return WithdrawPeriod
	})

	this.Withdraw2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("悔棋期转控制期")

		if this.withdraw_agreed == 1 { // 非控制方同意悔棋
			// 悔两步棋：非控制方玩家走的和控制方玩家走的
			this.composition.Pop()
			this.composition.Pop()
			this.composition_num = this.composition.Size()
			this.seats[this.currentPlayer].WithdrawNumber -= 1
		}

		this.NotifyWithdrawDecided() // 通知控制方玩家悔棋结果
		this.withdraw_agreed = -1
		return ControlPeriod
	})

	this.Control2DrawHandler = FSMHandler(func() FSMState {
		fmt.Println("控制期转求和期")
		this.start_draw_step = this.current_frame

		this.NotifyDrawRequested() // 通知非控制方玩家决定是否同意和棋
		return DrawPeriod
	})

	this.Draw2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("求和期转控制期")

		this.NotifyDrawDenied() //通知控制方玩家非控制方玩家拒绝和棋
		return ControlPeriod
	})

	this.Control2PlayFinishHandler = FSMHandler(func() FSMState {
		fmt.Println("控制期转行棋完成期")
		return PlayFinishPeriod
	})

	this.PlayFinish2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("行棋完成期转控制期")
		this.composition_num = this.composition.Size()
		// 改变走子方
		if this.currentPlayer == 0 {
			this.currentPlayer = 1
		} else {
			this.currentPlayer = 0
		}

		this.NotifyUpdateComposition()	// 通知非控制方玩家更新棋局，并开始走子
		return ControlPeriod
	})

	this.SettlementPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("进入结算期")

		// 先确定胜负方是谁，如果winner和loser都是-1则为平局
		winner := -1
		loser := -1
		if this.lose_requested != -1 { // 某方认输，通知另一方这个消息，再结算
			loser = this.lose_requested
			if loser == 0 {
				winner = 1
			} else if loser == 1 {
				winner = 0
			}
			this.NotifyTheOtherSideLost()
		} else if this.draw_agreed == 1 { // 非控制方同意和棋请求，通知另一方这个消息，再结算
			this.NotifyDrawAgreed()
		}

		// 计算分数-平局
		NewScoreWhite := int64(0)
		NewScoreBlack := int64(0)
		ExpWhite := 1/(math.Pow(10, float64((this.seats[1].Score-this.seats[0].Score)/400))+1)
		ExpBlack := 1/(math.Pow(10, float64((this.seats[0].Score-this.seats[1].Score)/400))+1)
		if winner == -1 || loser == -1 {
			// 计算白方的分数
			NewScoreWhite = this.seats[0].Score /*OldScoreWhite*/ + int64((float64(this.seats[0].K())) * (0.5-ExpWhite))
			NewScoreBlack = this.seats[1].Score /*OldScoreBlack*/ + int64((float64(this.seats[1].K())) * (0.5-ExpBlack))
		} else {
			if winner == 0 { // 白方赢
				NewScoreWhite = this.seats[0].Score /*OldScoreWhite*/ + int64((float64(this.seats[0].K())) * (1.0-ExpWhite))
				NewScoreBlack = this.seats[1].Score /*OldScoreBlack*/ + int64((float64(this.seats[1].K())) * (0.0-ExpBlack))
			} else { // 黑方赢
				NewScoreWhite = this.seats[0].Score /*OldScoreWhite*/ + int64((float64(this.seats[0].K())) * (0.0-ExpWhite))
				NewScoreBlack = this.seats[1].Score /*OldScoreBlack*/ + int64((float64(this.seats[1].K())) * (1.0-ExpBlack))
			}
		}

		// 修改table里黑白玩家的积分以及等级数据，一会儿发的消息里的数据是从这里拿的
		this.seats[0].Score = NewScoreWhite
		this.seats[0].Level = this.seats[0].GetLevel()
		this.seats[1].Score = NewScoreBlack
		this.seats[1].Level = this.seats[1].GetLevel()

		// 修改数据库里玩家分数以及等级
		infoDao := dao.NewUserInfoDao()
		resultWhite := infoDao.ModifyScoreNLevel(this.seats[0].UserId, this.seats[0].Score, this.seats[0].Level)
		if resultWhite != nil {
			fmt.Print("白方分数与等级数据修改失败")
		}
		resultBlack := infoDao.ModifyScoreNLevel(this.seats[1].UserId, this.seats[1].Score, this.seats[1].Level)
		if resultBlack != nil {
			fmt.Print("黑方分数与等级数据修改失败")
		}

		this.NotifyResult() // 通知所有玩家游戏结果
		return SettlementPeriod
	})

	this.fsm.AddHandler(VoidPeriod, MatchPeriodEvent,this.MatchPeriodHandler)
	this.fsm.AddHandler(MatchPeriod, ControlPeriodEvent,this.Match2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, WithdrawPeriodEvent,this.Control2WithdrawHandler)
	this.fsm.AddHandler(WithdrawPeriod, ControlPeriodEvent,this.Withdraw2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, DrawPeriodEvent,this.Control2DrawHandler)
	this.fsm.AddHandler(DrawPeriod, ControlPeriodEvent,this.Draw2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, PlayFinishPeriodEvent,this.Control2PlayFinishHandler)
	this.fsm.AddHandler(PlayFinishPeriod, ControlPeriodEvent,this.PlayFinish2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, SettlementPeriodEvent,this.SettlementPeriodHandler)
	this.fsm.AddHandler(DrawPeriod, SettlementPeriodEvent,this.SettlementPeriodHandler)
}

// 这里的是循环动作，在不停地检测
func (this *Table) StateSwitch() {
	switch this.fsm.getState() {
	case VoidPeriod:
		if this.seats[0].Bind() {
			this.fsm.Call(MatchPeriodEvent)
		}

	case MatchPeriod: 	// 匹配时长为10s，如果匹配成功，则直接进入ControlPeriod，超时则丢ai
		if (this.current_frame - this.start_match_step) > 1000 {
			this.fsm.Call(ControlPeriodEvent)
		} else {
			ready := true
			// 遍历所有座位，如果有座位没有跟session绑定，那么就没有匹配成功
			for _, seat := range this.GetSeats() {
				player := seat.(*objects.Player)
				if !player.Bind() || player.UserId == -1 || player.SitDown() == false || this.seats[0].UserId == this.seats[1].UserId{
					ready = false
				}
			}
			if ready {
				//匹配完成，直接进入控制期
				this.fsm.Call(ControlPeriodEvent)
			}
		}

	case ControlPeriod: // 控制方玩家需在行棋时间内完成行棋，否则判输
		if (this.current_frame - this.start_control_step) > this.control_time { // 玩家走子超时，判输
			this.fsm.Call(SettlementPeriodEvent)
		} else { // 控制方玩家在规定时间内完成了某些操作：走子、悔棋、和棋、认输
			if this.composition.Size() - 1 == this.composition_num { // 控制方玩家在规定时间内完成走子，通过composition的长度来判断
				this.fsm.Call(PlayFinishPeriodEvent)
			} else if this.withdraw_requested == 1 { // 控制方玩家要求悔棋
				this.withdraw_requested = 0
				this.fsm.Call(WithdrawPeriodEvent)
			} else if this.draw_requested == 1 { // 控制方玩家要求和棋
				this.draw_requested = 0
				this.fsm.Call(DrawPeriodEvent)
			} else if this.lose_requested != -1 { // 某方玩家认输
				this.fsm.Call(SettlementPeriodEvent)
			}
		}

	case PlayFinishPeriod: // 行棋完成期直接转入控制期
		this.fsm.Call(ControlPeriodEvent)

	case WithdrawPeriod: // 非控制方玩家需在规定时间内决定是否同意悔棋，否则直接认为同意
		if (this.current_frame - this.start_withdraw_step) > 1000 { // 非控制方玩家决定时间超时，直接认为同意悔棋
			this.withdraw_agreed = 1
			this.fsm.Call(ControlPeriodEvent)
		} else { // 非控制方在规定时间内决定悔棋结果
			if this.withdraw_agreed != -1 {
				this.fsm.Call(ControlPeriodEvent)
			}
		}

	case DrawPeriod: // 非控制方玩家需在规定时间内决定是否同意和棋，否则直接认为同意
		if (this.current_frame - this.start_draw_step) > 10 { // 非控制方玩家决定时间超时，直接认为同意和棋
			// 直接进入结算期
			this.fsm.Call(ControlPeriodEvent)
		} else {
			if this.draw_agreed == 0 { // 非控制方不同意和棋
				this.fsm.Call(ControlPeriodEvent)
			} else if this.draw_agreed == 1 { // 非控制方同意和棋
				this.fsm.Call(SettlementPeriodEvent)
			}
		}

	case SettlementPeriod:
	}
}
