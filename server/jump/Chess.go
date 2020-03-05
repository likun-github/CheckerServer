package jump

type Chess struct {
	//number int
	white string
	black string
	king string
}
func NewChess(white string,black string,king string)  *Chess{
	p:=&Chess{}
	p.white=white
	p.black=black
	p.king=king
	return p
}