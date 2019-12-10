package model

import "time"

type User struct {
	Id int64
	CreateTime time.Time `xorm:"created"`
	LastModifyTime time.Time `xorm:"updated"`
	DeleteTime time.Time `xorm:"deleted"`
	Name string `xorm:"varchar(25) unique notnull 'name'"`
	Password string `xorm:"varchar(16) not null 'password'"`
	Age int8 `xorm:"tinyint(3) 'age'"`
	//Sex int8 `xorm:"tinyint(1) 'sex'"`
	Email string `xorm:"varchar(25) 'email'"`
	//Address string //`xorm:"varchar(255) 'address'"`
	Tel int64 `xorm:"bigint(11) 'tel'"`
}
