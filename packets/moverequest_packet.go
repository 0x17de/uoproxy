package packets

const (
	DirN = iota
	DirNE
	DirE
	DirSE
	DirS
	DirSW
	DirW
	DirNW
)

type MoveRequestPacket struct {
	BasePacket
	Direction          byte // 0-7: N,NE,E,SE,S,SW,W,NW
	Sequence           byte
	FastwalkPrevention uint
}

func (p *MoveRequestPacket) Read(in chan byte) {
	p.Direction = <-in
	p.Sequence = <-in
	p.FastwalkPrevention = p.uint(in)
}

func (p *MoveRequestPacket) Bytes() []byte {
	return p.Writer().
		byte(p.Direction).
		byte(p.Sequence).
		uint(p.FastwalkPrevention).
		bytes()
}
