package model

type UserInfo struct {
	Id int64 `xorm:"pk autoincr 'userId'"`
	Name string `xorm:"varchar(40) notnull 'userName'"`
	WxOpenId string `xorm:"varchar(32) notnull 'userWXOpenID'"`
	WxName string `xorm:"varchar(32) notnull 'userWXName'"`
	Phone string `xorm:"varchar(20) notnull 'userPhone'"`
	WXImg string `xorm:"varchar(128) notnull 'userWXIMG'"`
	QQOpenId string `xorm:"varchar(32) notnull 'userQQOpenID'"`
	QQName string `xorm:"varchar(32) notnull 'userQQName'"`
	QQImg string  `xorm:"varchar(128) notnull 'userQQImg'"`
	Type int8 `xorm:"tinyint notnull 'userType'"`
	Pwd string `xorm:"varchar(40) notnull 'userPwd'"`
	Score int64 `xorm:"bigint notnull 'userScore'"`
	Level int8 `xorm:"tinyint notnull 'userLevel'"`
	Money float64 `xorm:"decimal notnull 'userMoney'"`
	//CreateTime time.Time `xorm:"created"`
	//LastModifyTime time.Time `xorm:"updated"`

}

/*
  [userWXImg_GXI] VARCHAR(2048)    NOT NULL,
  [userQQOpenID]  NVARCHAR(32)    NOT NULL,
  [userQQName]    NVARCHAR(32)    NOT NULL,
  [userQQImg]     VARBINARY(MAX)    NOT NULL,
  [userQQImg_GXI] VARCHAR(2048)    NOT NULL,
  [userType]      SMALLINT    NOT NULL,
  [userPwd]       NVARCHAR(40)    NOT NULL,
  [userScore]     DECIMAL(16)    NOT NULL,
  [userLevel]     SMALLINT    NOT NULL,
  [userMoney]     INT    NOT NULL,

*/
