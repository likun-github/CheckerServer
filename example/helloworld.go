// Copyright 2014 mqantserver Author. All Rights Reserved.
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
package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/liangdas/armyant/work"
	"time"
)

func main() {
	this := new(work.MqttWork)
	opts := this.GetDefaultOptions("ws://47.107.157.238:3654")
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		fmt.Println("连接断开", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		fmt.Println("连接成功")
	})
	err := this.Connect(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	//访问HelloWorld001模块的HD_Say函数
	//msg, err := this.Request("Test/HD_AddUser", []byte(`{"userName":"test","passWord":"xxx","age":26, "email:":"xxxx.com", "tel":151111111}`))
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(fmt.Sprintf("topic :%s  body :%s", msg.Topic(), string(msg.Payload())))
	Request(this, "Test/HD_AddUser", []byte(`{"userName":"test","passWord":"xxx","age":26, "email":"xxxx.com", "tel":151111111}`))
	fmt.Println("------------------------------")
	Request(this, "Test/HD_AddUser", []byte(`{"userName":"test","passWord":"xxx","age":26, "email":"xxxx.com", "tel":151111111}`))
	fmt.Println("------------------------------")
	//Request(this, "Test/HD_AddUser", []byte(`{"passWord":"xxx","age":26, "email:":"xxxx.com", "tel":151111111}`))
	//fmt.Println("------------------------------")
	//Request(this, "Test/HD_AddUser", []byte(`{"userName":"test","age":26, "email:":"xxxx.com", "tel":151111111}`))
	//fmt.Println("------------------------------")
	//Request(this, "Test/HD_AddUser", []byte(`{"age":26, "email:":"xxxx.com", "tel":151111111}`))
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")
	Request(this, "Test/HD_Login", []byte(`{"userName":"test","passWord":"xxx"}`))
	fmt.Println("------------------------------")
	Request(this, "Test/HD_Login", []byte(`{"userName":"test","passWord":"xxxx"}`))
	//fmt.Println("------------------------------")
	//Request(this, "Test/HD_Login", []byte(`{"passWord":"xxx"}`))
	//fmt.Println("------------------------------")
	//Request(this, "Test/HD_Login", []byte(`{"userName":"test"`))
	//fmt.Println("------------------------------")
	//Request(this, "Test/HD_Login", []byte(`{}`))
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")
	Request(this, "Test/HD_ModifyUser", []byte(`{"userName":"test","age":30}`))
	fmt.Println("------------------------------")
	Request(this, "Test/HD_ModifyUser", []byte(`{"userName":"test1","age":80}`))
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")

	Request(this, "Test/HD_FindUserByName", []byte(`{"userName":"test"}`))
	fmt.Println("------------------------------")
	Request(this, "Test/HD_FindUserByName", []byte(`{"userName":"test1"}`))
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")

	Request(this, "Test/HD_RemoveUser", []byte(`{"userName":"test"}`))
	fmt.Println("------------------------------")
	Request(this, "Test/HD_RemoveUser", []byte(`{"userName":"test1"}`))
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")


	time.Sleep(2*time.Minute)

}

func Request(this*work.MqttWork, topic string, body[]byte) {
	msg, err := this.Request(topic, body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(fmt.Sprintf("topic :%s  body :%s", msg.Topic(), string(msg.Payload())))
}
