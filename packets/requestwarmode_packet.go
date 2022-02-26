package packets

type RequestWarModePacket struct {
	BasePacket
	Flag byte
	UNK1 []byte
}

func (p *RequestWarModePacket) Read(in chan byte) {
	p.Flag = <-in
	p.UNK1 = p.bstr(in, 3)
}

func (p *RequestWarModePacket) Bytes() []byte {
	return p.Writer().
		byte(p.Flag).
		bstr(p.UNK1).
		bytes()
}
