package packets

type SendSpeechPacket struct {
	BasePacket
	Target    uint
	Model     int
	TextType  byte
	TextColor int
	Font      int
	Name      string
	Text      string
}

func (p *SendSpeechPacket) Read(in chan byte) {
	l := p.short(in)
	p.Target = p.uint(in)
	p.Model = p.short(in)
	p.TextType = <-in
	p.TextColor = p.short(in)
	p.Font = p.short(in)
	p.Name = p.zstrFixed(in, 30)
	p.Text = p.zstrFixed(in, l-44)
}

func (p *SendSpeechPacket) Bytes() []byte {
	return p.Writer().
		uint(p.Target).
		short(p.Model).
		byte(p.TextType).
		short(p.TextColor).
		short(p.Font).
		zstrFixed(p.Name, 30).
		zstrFixed(p.Text, len(p.Text)).
		bytes()
}
