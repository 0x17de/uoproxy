package packets

type BytePacket struct {
	BasePacket
	Len  int
	Data []byte
}

func (p *BytePacket) Read(in chan byte) {
	p.Data = make([]byte, p.Len)
	for i := 0; i < p.Len; i++ {
		p.Data[i] = <-in
	}
}
func (p *BytePacket) Bytes() []byte {
	return p.Data
}
