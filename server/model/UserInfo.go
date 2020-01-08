package model

type UserInfo struct {
	Id int64 `xorm:"pk autoincr 'userId'"`
	Name string `xorm:"varchar(40)  'userName'"`
	WxOpenId string `xorm:"varchar(32) notnull 'userWXOpenID'"`
	WxName string `xorm:"varchar(32)  'userWXName'"`
	Phone string `xorm:"varchar(20)  'userPhone'"`
	WXImg string `xorm:"varchar(128)  'userWXIMG'"`
	QQOpenId string `xorm:"varchar(32)  'userQQOpenID'"`
	QQName string `xorm:"varchar(32)  'userQQName'"`
	QQImg string  `xorm:"varchar(128)  'userQQImg'"`
	Type int8 `xorm:"tinyint  'userType'"`
	Pwd string `xorm:"varchar(40)  'userPwd'"`
	Score int64 `xorm:"bigint notnull 'userScore'"`//分数
	Level int8 `xorm:"tinyint notnull 'userLevel'"`//关卡
	Money float64 `xorm:"decimal  'userMoney'"`
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
