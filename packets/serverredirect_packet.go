package packets

type ServerRedirectPacket struct {
	BasePacket
	Address [4]byte
	Port    int
	Key     uint
}

func (p *ServerRedirectPacket) Read(in chan byte) {
	for i := 0; i < 4; i++ {
		p.Address[i] = <-in
	}
	p.Port = p.short(in)
	p.Key = p.uint(in)
}
func (p *ServerRedirectPacket) Bytes() []byte {
	o := p.Writer()
	for i := 0; i < 4; i++ {
		o.byte(p.Address[i])
	}
	o.short(p.Port)
	o.uint(p.Key)
	return o.bytes()
}
