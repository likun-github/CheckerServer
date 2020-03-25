package jump



// Copyright 2014 loolgame Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"CheckerServer/server/common/moveGeneration"
	"CheckerServer/server/jump/objects"
	"errors"
	"fmt"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/gate"
	"strconv"
)
//坐下
/*
func (self *Table) SitDown(session gate.Session) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		if player.SitDown() {
			fmt.Println("玩家成功坐下")
		}
		return nil
	}
	return nil
}

 */
/*
func (self *Table) StartGame(session gate.Session) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		if self.State() == room.Active {

		} else if self.State() == room.Initialized {
			self.Start()
		} else if self.State() == room.Paused {
			self.Resume()
		}
		return nil
	}
	return nil
}

 */
/*
func (self *Table) PauseGame(session gate.Session) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		self.Pause()
		return nil
	}
	return nil
}

 */

/**
玩家加入场景
*/
func (self *Table) Join(session gate.Session/*, userinfo *model.UserInfo*/) error {
	self.writelock.Lock()
	defer self.writelock.Unlock()

	// 更改table状态
	if self.State() == room.Active {
	} else if self.State() == room.Initialized {
		self.Start()
	} else if self.State() == room.Paused {
		self.Resume()
	}

	player := self.GetBindPlayer(session)
	if player != nil {
		fmt.Println("player已经跟session绑定")
		playerImp := player.(*objects.Player)
		playerImp.OnRequest(session)


		return nil
	}
	fmt.Println("player还没有跟session绑定")
	var indexSeat int = -1
	for i, player := range self.seats {
		if !player.Bind() {
			indexSeat = i

			if i == 1 {
				userid,_ := strconv.ParseInt(session.GetUserId(), 10, 64)
				if self.seats[0].UserId == userid {
					return fmt.Errorf("此userid已加入过桌子")
				}
			}

			player.OnBind(session)
			fmt.Println("player完成跟session的绑定")
			player.OnRequest(session)
			//self.NotifyJoin(player) //广播给所有其他玩家
			break
		}
	}
	if indexSeat == -1 {
		return fmt.Errorf("房间已满,无法加入游戏")
	}
	return nil
}


//角色绑定

func (self *Table) Login(session gate.Session, userid int64, level int8, score int64,username string,avatar string) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)

		player.UserId = userid
		player.Username=username
		player.Avatar=avatar
		player.Score=score
		player.Level = level
		fmt.Println("player与数据库里的信息绑定成功")

		player.OnRequest(session)
		player.OnSitDown()
		if player.SitDown() {
			fmt.Println("玩家成功坐下")
		}

		return nil
	} else {
		return errors.New("session没有和player绑定")
	}

}


// 控制方走子完成
func (self *Table) PlayOneTurn(session gate.Session, composition *objects.Chess) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)

		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	self.composition.Push(composition)
	fmt.Println("控制方走子加入棋局栈")

	// 判断游戏是否结束
	// 初始化bitboard
	W,_ := strconv.ParseInt(composition.White, 2, 64)
	B,_ := strconv.ParseInt(composition.Black, 2, 64)
	K,_ := strconv.ParseInt(composition.King, 2, 64)
	var cb = &moveGeneration.CheckerBitboard{W:uint64(W),B:uint64(B),K:uint64(K)}
	// 初始化padded array board
	moveGeneration.BitboardToPaddedArrayBoard(cb)

	if cb.GetMoversWhite() == 0 && cb.GetJumpersWhite() == 0 { // 白子没法走，黑子胜
		self.winner = 1
		self.loser = 0
	} else if cb.GetMoversBlack() == 0 && cb.GetJumpersBlack() == 0 { // 黑子没法走，白子胜
		self.winner = 0
		self.loser = 1
	}

	return nil
}

// 控制方悔棋
func (self *Table) Withdraw(session gate.Session) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
		// 判断控制方是否还有悔棋次数
		if player.WithdrawNumber > 0 { // 玩家还有悔棋次数，可以悔棋
			self.withdraw_requested = 1
		} else { // 玩家没有悔棋次数，不能悔棋
			session.Send("CannotWithdraw",nil)
		}
	}

	return nil
}

// 非控制方确定悔棋结果
func (self *Table) WithdrawDecided(session gate.Session, withdrawAgreed int) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	self.withdraw_agreed = withdrawAgreed
	return nil
}

// 控制方请求和棋
func (self *Table) Draw(session gate.Session) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	self.draw_requested = 1
	return nil
}

// 非控制方确定和棋结果
func (self *Table) DrawDecided(session gate.Session, drawAgreed int) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	self.draw_agreed = drawAgreed
	return nil
}

// 控制方认输
func (self *Table) Lose(session gate.Session,  lose_checker_color int) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	self.lose_requested = lose_checker_color

	return nil
}

// 某方收藏本局
func (self *Table) Collect(session gate.Session,  collect_checker_color int) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	if collect_checker_color == 0 { // 白方收藏本局
		self.collect_requested_white = self.seats[0].UserId
	} else if collect_checker_color == 1 { // 黑方收藏本局
		self.collect_requested_white = self.seats[0].UserId
	}

	return nil
}

// 某方再来一局
func (self *Table) Again(session gate.Session,  collect_checker_color int) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	if collect_checker_color == 0 { // 白方再来一局
		self.game_finished_action_white = 1
	} else if collect_checker_color == 1 { // 黑方再来一局
		self.game_finished_action_white = 1
	}

	return nil
}

// 某方返回大厅
func (self *Table) Exit_(session gate.Session,  collect_checker_color int) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
	}

	if collect_checker_color == 0 { // 白方返回大厅
		self.game_finished_action_white = 2
	} else if collect_checker_color == 1 { // 黑方返回大厅
		self.game_finished_action_white = 2
	}

	return nil
}
