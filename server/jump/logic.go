package jump

import (
"CheckerServer/server/xaxb/objects"
"fmt"
"github.com/liangdas/mqant/utils"
"math"
)

var (
	MatchPeriod			  = FSMState("匹配等待期")
	ControlPeriod 		  = FSMState("控制期")
	WithdrawPeriod		  = FSMState("悔棋期")
	DrawPeriod			  = FSMState("求和期")
	LosePeriod		      = FSMState("认输期")
	PlayDraughtPeriod     = FSMState("认输期")
	SettlementPeriod	  = FSMState("结算期")


	MatchPeriodEvent	  = FSMEvent("进入匹配等待期")
	ControlPeriodEvent 	  = FSMEvent("进入控制期")
	WithdrawPeriodEvent	  = FSMEvent("进入悔棋期")
	DrawPeriodEvent		  = FSMEvent("进入求和期")
	LosePeriodEvent		  = FSMEvent("进入认输期")
	PlayDraughtPeriodEvent = FSMEvent("进入行棋期")
	SettlementPeriodEvent = FSMEvent("进入结算期")



)

func (this *Table) InitFsm() {
	this.fsm = *NewFSM(MatchPeriod)
	this.MatchPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入匹配等待期")
		return MatchPeriod
	})
	this.ControlPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入控制期")
		//this.step1 = this.current_frame
		this.NotifyIdle()

		for _, seat := range this.GetSeats() {
			player := seat.(*objects.Player)
			if player.Bind() {
				if player.Coin <= 0 {
					player.Session().Send("XaXb/Exit", []byte(`{"Info":"金币不足你被强制离开房间"}`))
					player.OnUnBind() //踢下线
				}
			}
		}

		return ControlPeriod
	})
	this.BettingPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入押注期")
		this.step2 = this.current_frame
		this.NotifyBetting()
		return BettingPeriod
	})
	this.OpeningPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入开奖期")
		this.step3 = this.current_frame
		this.NotifyOpening()
		return OpeningPeriod
	})
	this.SettlementPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入结算期")
		var mixWeight int64 = math.MaxInt64
		var winer *objects.Player = nil
		Result := utils.RandInt64(0, 10)
		for _, seat := range this.GetSeats() {
			player := seat.(*objects.Player)
			if player.Stake {
				player.Weight = int64(math.Abs(float64(player.Target - Result)))
				if mixWeight > player.Weight {
					mixWeight = player.Weight
					winer = player
				}
			}
		}
		if winer != nil {
			winer.Coin += 800
		}

		this.step4 = this.current_frame
		this.NotifySettlement(Result)
		return SettlementPeriod
	})

	//this.fsm.AddHandler(IdlePeriod, VoidPeriodEvent, this.VoidPeriodHandler)
	//this.fsm.AddHandler(SettlementPeriod, VoidPeriodEvent, this.VoidPeriodHandler)
	//this.fsm.AddHandler(BettingPeriod, VoidPeriodEvent, this.VoidPeriodHandler)
	//this.fsm.AddHandler(OpeningPeriod, VoidPeriodEvent, this.VoidPeriodHandler)
	//
	//this.fsm.AddHandler(VoidPeriod, IdlePeriodEvent, this.IdlePeriodHandler)
	//this.fsm.AddHandler(SettlementPeriod, IdlePeriodEvent, this.IdlePeriodHandler)
	//
	//this.fsm.AddHandler(IdlePeriod, BettingPeriodEvent, this.BettingPeriodHandler)
	//this.fsm.AddHandler(BettingPeriod, OpeningPeriodEvent, this.OpeningPeriodHandler)
	//this.fsm.AddHandler(OpeningPeriod, SettlementPeriodEvent, this.SettlementPeriodHandler)
}

/**
进入空闲期
*/
func (this *Table) StateSwitch() {
	//if this.fsm.getState()==WithdrawPeriod:
	//
	//switch this.fsm.getState() {
	//case VoidPeriod:
	//
	//case IdlePeriod:
	//	if (this.current_frame - this.step1) > 5 {
	//		this.fsm.Call(BettingPeriodEvent)
	//	} else {
	//		//this.NotifyAxes()
	//	}
	//case BettingPeriod:
	//	if (this.current_frame - this.step2) > 20 {
	//		this.fsm.Call(OpeningPeriodEvent)
	//	} else {
	//		ready := true
	//		for _, seat := range this.GetSeats() {
	//			player := seat.(*objects.Player)
	//			if player.SitDown() && !player.Stake {
	//				ready = false
	//			}
	//		}
	//		if ready {
	//			//都押注了直接开奖
	//			this.fsm.Call(OpeningPeriodEvent)
	//		}
	//	}
	//case OpeningPeriod:
	//	if (this.current_frame - this.step3) > 5 {
	//		this.fsm.Call(SettlementPeriodEvent)
	//	} else {
	//		//this.NotifyAxes()
	//	}
	//case SettlementPeriod:
	//	if (this.current_frame - this.step4) > 5 {
	//		this.fsm.Call(IdlePeriodEvent)
	//	} else {
	//		//this.NotifyAxes()
	//	}
	//}
}
