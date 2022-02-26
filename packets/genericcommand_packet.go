package packets

type GenericCommandPacket struct {
	BasePacket
	SubCmd int
	Data   []byte
}

func (p *GenericCommandPacket) Read(in chan byte) {
	l := p.short(in)
	p.SubCmd = p.short(in)
	p.Data = p.bstr(in, l-5)
}

func (p *GenericCommandPacket) Bytes() []byte {
	return p.Writer().
		short(p.SubCmd).
		bstr(p.Data).
		bytes()
}
