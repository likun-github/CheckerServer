package objects

type Chess struct {
	//number int
	White string
	Black string
	King  string
}
func NewChess(white string,black string,king string)  *Chess{
	p:=&Chess{}
	p.White =white
	p.Black=black
	p.King=king
	return p
}