package jump

type Chess struct {
	//number int
	White string
	black string
	king  string
}
func NewChess(white string,black string,king string)  *Chess{
	p:=&Chess{}
	p.White =white
	p.black=black
	p.king=king
	return p
}