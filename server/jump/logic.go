package jump

var (
	VoidPeriod 			  = FSMState("等待匹配期")
	MatchPeriod			  = FSMState("匹配完成期")
	ControlPeriod 		  = FSMState("控制期")
	WithdrawPeriod		  = FSMState("悔棋期")
	DrawPeriod			  = FSMState("求和期")
	LosePeriod		      = FSMState("认输期")
	PlayDraughtPeriod     = FSMState("认输期")
	SettlementPeriod	  = FSMState("结算期")

	VoidPeriodEvent 	  = FSMEvent("进入等待匹配期")
	MatchPeriodEvent	  = FSMEvent("进入匹配完成期")
	ControlPeriodEvent 	  = FSMEvent("进入控制期")
	WithdrawPeriodEvent	  = FSMEvent("进入悔棋期")
	DrawPeriodEvent		  = FSMEvent("进入求和期")
	LosePeriodEvent		  = FSMEvent("进入认输期")
	PlayDraughtPeriodEvent = FSMEvent("进入行棋期")
	SettlementPeriodEvent = FSMEvent("进入结算期")



)

func (this *Table) InitFsm() {
	this.fsm = *NewFSM(VoidPeriod)
	//this.PlayDraughtHandler = FSMHandler(func() FSMState {
	//	fmt.Println("已进入等待匹配期")
	//	return VoidPeriod
	//})
	////進入匹配期
	//this.MatchPeriodHandler = FSMHandler(func() FSMState {
	//	fmt.Println("进入匹配完成期")
	//	this.NotifyMatch()//通知所有玩家进入匹配完成期
	//	return MatchPeriod
	//})
	//this.ControlPeriodHandler = FSMHandler(func() FSMState {
	//	fmt.Println("已进入控制期")
	//	this.NotifyBetting()
	//	return ControlPeriod
	//})
	//this.SettlementPeriodHandler = FSMHandler(func() FSMState {
	//	fmt.Println("已进入结算期")
	//	var mixWeight int64 = math.MaxInt64
	//	var winer *objects.Player = nil
	//	Result := utils.RandInt64(0, 10)
	//	for _, seat := range this.GetSeats() {
	//		player := seat.(*objects.Player)
	//		if player.Stake {
	//			player.Weight = int64(math.Abs(float64(player.Target - Result)))
	//			if mixWeight > player.Weight {
	//				mixWeight = player.Weight
	//				winer = player
	//			}
	//		}
	//	}
	//	if winer != nil {
	//		winer.Coin += 800
	//	}
	//
	//	//this.step4 = this.current_frame
	//	this.NotifySettlement(Result)
	//	return SettlementPeriod
	//})

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

}
