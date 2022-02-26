package packets

import "bytes"

type BOut struct {
	info PacketInfo
	o    bytes.Buffer
}

func NewBOut(info PacketInfo) *BOut {
	o := &BOut{info: info}
	o.byteInt(info.Id)
	if info.Size == -1 {
		o.short(0)
	}
	return o
}

func (o *BOut) byteInt(b int) *BOut {
	o.o.WriteByte(byte(b & 0xff))
	return o
}
func (o *BOut) zstrFixed(s string, l int) *BOut {
	for i := 0; i < l; i++ {
		if i < len(s) {
			o.o.WriteByte(s[i])
		} else {
			o.o.WriteByte(0)
		}
	}
	return o
}
func (o *BOut) bstr(s []byte) *BOut {
	for _, c := range s {
		o.o.WriteByte(c)
	}
	return o
}
func (o *BOut) byte(b byte) *BOut {
	o.o.WriteByte(b)
	return o
}
func (o *BOut) uint(v uint) *BOut {
	for i := 3; i >= 0; i-- {
		o.o.WriteByte(byte((v >> (i * 8)) & 0xff))
	}
	return o
}
func (o *BOut) short(i int) *BOut {
	o.o.WriteByte(byte(i >> 8))
	o.o.WriteByte(byte(i & 0xff))
	return o
}
func (o *BOut) byteArr(b []byte) *BOut {
	for _, c := range b {
		o.o.WriteByte(c)
	}
	return o
}
func (o *BOut) bytes() []byte {
	b := o.o.Bytes()
	if o.info.Size == -1 {
		l := len(b)
		b[1] = byte(l >> 8)
		b[2] = byte(l & 0xff)
	}
	return b
}
