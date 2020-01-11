package Jump

import (
	"math/rand"
	"time"
)

type CMenu struct {
	number int//棋局帧数
	white int64
	black int64
	king int64
}
//整个棋局
type CMenus []*CMenu

var MenuMap map[string]interface{}

//新建一帧
func NewMenu(number int,white int64,black int64,king int64)  *CMenu{
	p:=&CMenu{}
	p.white=white
	p.black=black
	p.king=king
	return p
}

func init()  {
	rand.Seed(time.Now().UnixNano())
	//建立棋局//menus:=new(CMenus)



}


