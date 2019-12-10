package model

import "time"

type Record struct {
	Id int64
	CreateTime time.Time `xorm:"created"`
	LastModifyTime time.Time `xorm:"updated"`
	DeleteTime time.Time `xorm:"deleted"`
	UserId1 int64 `xorm:"bigint notnull user_id1"`
	UserId2 int64 `xorm:"bigint notnull user_id2"`
	Winner int8
}
