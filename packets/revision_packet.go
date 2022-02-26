package packets

type RevisionPacket struct {
	BasePacket
	Serial   uint
	Revision uint
}

func (p *RevisionPacket) Read(in chan byte) {
	p.Serial = p.uint(in)
	p.Revision = p.uint(in)
}

func (p *RevisionPacket) Bytes() []byte {
	return p.Writer().
		uint(p.Serial).
		uint(p.Revision).
		bytes()
}
