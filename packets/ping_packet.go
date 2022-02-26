package packets

type PingPacket struct {
	BasePacket
	Value byte
}

func (p *PingPacket) Read(in chan byte) {
	p.Value = <-in
}

func (p *PingPacket) Bytes() []byte {
	return p.Writer().
		byte(p.Value).
		bytes()
}
