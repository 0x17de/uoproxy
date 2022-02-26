package packets

type UpdatePlayerPacket struct {
	BasePacket
	Target    uint
	Model     int
	X         int
	Y         int
	Z         int
	Direction byte
	Color     int
	Status    byte
	Highlight byte
}

func (p *UpdatePlayerPacket) Read(in chan byte) {
	p.Target = p.uint(in)
	p.Model = p.short(in)
	p.X = p.short(in)
	p.Y = p.short(in)
	p.Z = int(<-in)
	p.Direction = <-in
	p.Color = p.short(in)
	p.Status = <-in
	p.Highlight = <-in
}

func (p *UpdatePlayerPacket) Bytes() []byte {
	return p.Writer().
		uint(p.Target).
		short(p.Model).
		short(p.X).
		short(p.Y).
		byteInt(p.Z).
		byte(p.Direction).
		short(p.Color).
		byte(p.Status).
		byte(p.Highlight).
		bytes()
}
