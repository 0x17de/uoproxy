package packets

type PostLoginPacket struct {
	BasePacket
	Key  uint
	User string
	Pass string
}

func (p *PostLoginPacket) Read(in chan byte) {
	p.Key = p.uint(in)
	p.User = p.zstrFixed(in, 30)
	p.Pass = p.zstrFixed(in, 30)
}
func (p *PostLoginPacket) Bytes() []byte {
	return p.Writer().
		uint(p.Key).
		zstrFixed(p.User, 30).
		zstrFixed(p.Pass, 30).
		bytes()
}
