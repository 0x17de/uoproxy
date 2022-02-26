package packets

import (
	"bytes"
)

type UnknownPacket struct {
	BasePacket
	Len     int
	DataLen int
	Data    []byte
}

func (p *UnknownPacket) Read(in chan byte) {
	var o bytes.Buffer
	if p.Info.Size == -1 {
		p.Len = p.short(in)
		p.DataLen = p.Len - 3
	} else {
		p.Len = p.Info.Size
		p.DataLen = p.Len - 1
	}
	for i := 0; i < p.DataLen; i++ {
		o.WriteByte(<-in)
	}
	p.Data = o.Bytes()
}

func (p *UnknownPacket) Bytes() []byte {
	return p.Writer().byteArr(p.Data).bytes()
}
