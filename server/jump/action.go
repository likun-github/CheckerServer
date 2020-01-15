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
"CheckerServer/server/jump/objects"
"encoding/json"
"fmt"
"github.com/liangdas/mqant-modules/room"
"github.com/liangdas/mqant/gate"
)
//坐下
func (self *Table) SitDown(session gate.Session) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()


		return nil
	}
	return nil
}

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

/**
玩家加入场景
*/
func (self *Table) Join(session gate.Session) error {
	self.writelock.Lock()
	defer self.writelock.Unlock()
	player := self.GetBindPlayer(session)
	if player != nil {
		playerImp := player.(*objects.Player)
		playerImp.OnRequest(session)

		//回复当前状态
		result := map[string]interface{}{
			"State":     self.State(),
			"Rid":       player.Session().GetUserId(),
			"SeatIndex": playerImp.SeatIndex,
		}
		b, _ := json.Marshal(result)
		session.Send("Jump/OnEnter", b)

		return nil
	}
	var indexSeat int = -1
	for i, player := range self.seats {
		if !player.Bind() {
			indexSeat = i
			player.OnBind(session)
			self.NotifyJoin(player) //广播给所有其他玩家

			result := map[string]interface{}{
				"State":     self.State(),
				"Rid":       player.Session().GetUserId(),
				"SeatIndex": indexSeat,
			}
			b, _ := json.Marshal(result)
			session.Send("Jump/OnEnter", b)
			break
		}
	}

	if indexSeat == -1 {

		return fmt.Errorf("房间已满,无法加入游戏")
	}
	return nil
}

/**
玩家押注
*/
func (self *Table) Stake(session gate.Session, target int64) error {
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.OnRequest(session)
		player.OnSitDown()
		return nil
	}
	return nil
}
//角色绑定
func (self *Table) login(session gate.Session, Score int64,username string,avatar string) error {
	fmt.Println("函数运行了吗")
	playerImp := self.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*objects.Player)
		player.Score=Score
		player.Username=username
		player.Avatar=avatar
		player.OnRequest(session)
		player.OnSitDown()
		fmt.Println("绑定成功")
		return nil
	}
	return nil
}

