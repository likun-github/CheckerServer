package model

import (
	"time"
)


type Collection struct {
	Id int64 `xorm:"pk AUTO_INCREMENT 'collectionId'"`
	UserId int64 `xorm:"bigint notnull 'userId'"`
	CreateTime time.Time `xorm:"created"`
	LastModifyTime time.Time `xorm:"updated"`
	//DeleteTime time.Time `xorm:"deleted"`
	Step int8 `xorm:"int notnull 'step'"`
	White []int64 `xorm:"TEXT notnull White"`
	Black []int64 `xorm:"TEXT notnull Black"`
	King []int64 `xorm:"TEXT notnull King"`
}