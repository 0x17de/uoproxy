package packets

import (
	"log"
)

type ServerListEntry struct {
	Name     string
	Full     byte
	Timezone byte
	Address  [4]byte
}
type ServerListPacket struct {
	BasePacket
	Flags   int
	Servers []ServerListEntry
}

func (p *ServerListPacket) Read(in chan byte) {
	p.short(in) // Len
	p.Flags = int(<-in)
	p.Servers = make([]ServerListEntry, p.short(in))
	for i := 0; i < len(p.Servers); i++ {
		s := &p.Servers[p.short(in)]
		s.Name = p.zstrFixed(in, 32)
		s.Full = <-in
		s.Timezone = <-in
		for i := 0; i < 4; i++ {
			s.Address[i] = <-in
		}
		log.Printf("Found server: %s", s.Name)
	}
}
func (p *ServerListPacket) Bytes() []byte {
	o := p.Writer().
		byteInt(p.Flags).
		short(len(p.Servers))
	for i := 0; i < len(p.Servers); i++ {
		s := &p.Servers[i]
		o.short(i).
			zstrFixed(s.Name, 32).
			byte(s.Full).
			byte(s.Timezone)
		for i := 3; i >= 0; i-- {
			o.byte(s.Address[i])
		}
	}
	return o.bytes()
}
