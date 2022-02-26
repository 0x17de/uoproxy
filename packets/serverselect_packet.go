package packets

type ServerSelectPacket struct {
	BasePacket
	Index int
}

func (p *ServerSelectPacket) Read(in chan byte) {
	p.Index = p.short(in)
}
func (p *ServerSelectPacket) Bytes() []byte {
	return p.Writer().
		short(p.Index).
		bytes()
}
