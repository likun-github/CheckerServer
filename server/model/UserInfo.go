package model

type UserInfo struct {
	Id int64 `xorm:"pk AUTO_INCREMENT 'userId'"`//用户id
	Name string `xorm:"varchar(40)  'userName'"`//用户真实姓名
	WxOpenId string `xorm:"varchar(32) notnull 'userWXOpenID'"`//微信openid
	WxName string `xorm:"varchar(32)  'userWXName'"`//微信昵称
	Phone string `xorm:"varchar(20)  'userPhone'"`//电话
	WXImg string `xorm:"varchar(128)  'userWXIMG'"`//头像
	QQOpenId string `xorm:"varchar(32)  'userQQOpenID'"`//QQopenid
	QQName string `xorm:"varchar(32)  'userQQName'"`//QQ昵称
	QQImg string  `xorm:"varchar(128)  'userQQImg'"`//QQ头像
	Type int8 `xorm:"tinyint  'userType'"`//认证类型
	//Pwd string `xorm:"varchar(40)  'userPwd'"`
	Status int8  `xorm:"tinyint  'userStatus'"`//用户状态，0代表仅获取openid,1代表获取基本用户信息
	Score int64 `xorm:"bigint notnull 'userScore'"`//分数
	Level int8 `xorm:"tinyint notnull 'userLevel'"`//关卡
	Money float64 `xorm:"decimal  'userMoney'"`//金币
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
