package model

import "time"

type ChessManual struct {
	Id int64
	CreateTime time.Time `xorm:"created"`
	LastModifyTime time.Time `xorm:"updated"`
	//DeleteTime time.Time `xorm:"deleted"`
	White int64 `xorm:"bigint notnull White"`
	Black int64 `xorm:"bigint notnull Black"`
	King int64 `xorm:"bigint notnull King"`
}
