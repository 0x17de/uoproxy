package packets

type AccountLoginRequestPacket struct {
	BasePacket
	User         string
	Pass         string
	NextLoginKey byte
}

func (p *AccountLoginRequestPacket) Read(in chan byte) {
	p.User = p.zstrFixed(in, 30)
	p.Pass = p.zstrFixed(in, 30)
	p.NextLoginKey = <-in
}

func (p *AccountLoginRequestPacket) Bytes() []byte {
	return p.Writer().
		zstrFixed(p.User, 30).
		zstrFixed(p.Pass, 30).
		byte(p.NextLoginKey).
		bytes()
}
