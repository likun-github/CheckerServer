package jump

import (
	"CheckerServer/server/xaxb/objects"
	"fmt"
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
