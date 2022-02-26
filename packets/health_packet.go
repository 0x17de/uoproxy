package packets

type HealthPacket struct {
	BasePacket
	Target  uint
	Maximum int
	Current int
}

func (p *HealthPacket) Read(in chan byte) {
	p.Target = p.uint(in)
	p.Maximum = p.short(in)
	p.Current = p.short(in)
}

func (p *HealthPacket) Bytes() []byte {
	return p.Writer().
		uint(p.Target).
		short(p.Maximum).
		short(p.Current).
		bytes()
}
