package jump

import (
	"fmt"
	//"github.com/liangdas/mqant/utils"
	//"math"
)

var (
	VoidPeriod 			  = FSMState("等待匹配期")
	MatchPeriod			  = FSMState("匹配完成期")
	ControlPeriod 		  = FSMState("控制期")
	PlayDraughtPeriod     = FSMState("行棋期")
	WithdrawPeriod		  = FSMState("悔棋期")
	DrawPeriod			  = FSMState("求和期")
	SettlementPeriod	  = FSMState("结算期")

	VoidPeriodEvent 	  = FSMEvent("进入等待匹配期")
	MatchPeriodEvent	  = FSMEvent("进入匹配完成期")
	ControlPeriodEvent 	  = FSMEvent("进入控制期")
	PlayDraughtPeriodEvent = FSMEvent("进入行棋期")
	WithdrawPeriodEvent	  = FSMEvent("进入悔棋期")
	DrawPeriodEvent		  = FSMEvent("进入求和期")
	SettlementPeriodEvent = FSMEvent("进入结算期")



)

func (this *Table) InitFsm() {
	this.fsm = *NewFSM(VoidPeriod)

	//進入匹配完成期
	this.MatchPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("进入匹配完成期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return MatchPeriod
	})
	this.Match2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("匹配转控制期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return ControlPeriod
	})
	this.Control2WithdrawHandler = FSMHandler(func() FSMState {
		fmt.Println("控制转匹配期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return WithdrawPeriod
	})
	this.Withdraw2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("悔棋转控制期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return ControlPeriod
	})
	this.Control2DrawHandler = FSMHandler(func() FSMState {
		fmt.Println("控制转求和期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return DrawPeriod
	})
	this.Draw2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("求和转控制期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return ControlPeriod
	})
	this.Control2PlayDraughtHandler = FSMHandler(func() FSMState {
		fmt.Println("控制转行棋期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return PlayDraughtPeriod
	})
	this.PlayDraught2ControlHandler = FSMHandler(func() FSMState {
		fmt.Println("行棋转控制期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return ControlPeriod
	})
	this.SettlementPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("结算期")
		this.NotifyMatch()//通知所有玩家进入匹配完成期
		return SettlementPeriod
	})




	this.fsm.AddHandler(VoidPeriod, MatchPeriodEvent,this.MatchPeriodHandler)
	this.fsm.AddHandler(MatchPeriod, ControlPeriodEvent,this.Match2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, WithdrawPeriodEvent,this.Control2WithdrawHandler)
	this.fsm.AddHandler(WithdrawPeriod, ControlPeriodEvent,this.Withdraw2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, DrawPeriodEvent,this.Control2DrawHandler)
	this.fsm.AddHandler(DrawPeriod, ControlPeriodEvent,this.Draw2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, PlayDraughtPeriodEvent,this.Control2PlayDraughtHandler)
	this.fsm.AddHandler(PlayDraughtPeriod, ControlPeriodEvent,this.PlayDraught2ControlHandler)
	this.fsm.AddHandler(ControlPeriod, SettlementPeriodEvent,this.SettlementPeriodHandler)


}

/**
进入空闲期
*/
/**
进入空闲期
*/
func (this *Table) StateSwitch() {
	switch this.fsm.getState() {
	case VoidPeriod:

	case MatchPeriod:

	case ControlPeriod:

	case PlayDraughtPeriod:

	case WithdrawPeriod:

	case DrawPeriod:

	case SettlementPeriod:
	}
}
