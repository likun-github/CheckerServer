package jump

type Chess struct {
	//number int
	white int64
	black int64
	king int64
}
func NewChess(white int64,black int64,king int64)  *Chess{
	p:=&Chess{}
	p.white=white
	p.black=black
	p.king=king
	return p
}