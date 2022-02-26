package packets

type MovePacket struct {
	BasePacket
	SeqKey    byte
	Notoriety byte
}

func (p *MovePacket) Read(in chan byte) {
	p.SeqKey = <-in
	p.Notoriety = <-in
}

func (p *MovePacket) Bytes() []byte {
	return p.Writer().
		byte(p.SeqKey).
		byte(p.Notoriety).
		bytes()
}
