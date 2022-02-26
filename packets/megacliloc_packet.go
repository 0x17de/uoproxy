package packets

type ClilocEntry struct {
	Id   uint
	Text []byte
}
type ClientMegaCliloc struct {
	Serials []uint
}
type ServerMegaCliloc struct {
	UNK1    int
	Serial1 uint
	UNK2    int
	Serial2 uint

	Clilocs []ClilocEntry
}
type MegaClilocPacket struct {
	BasePacket
	Server *ServerMegaCliloc
	Client *ClientMegaCliloc
}

//go:generate stringer -type=ClilocId
type ClilocId uint

const (
	NamedItem                 ClilocId = 0x1011ab
	ItemStack                 ClilocId = 0x1005b7 // {COUNT}  #{CLILOCID}  e.g.: Bandage:1023617 Arrow:1023903
	Name                      ClilocId = 0x1036d5
	Insured                   ClilocId = 0x103332
	SkillBonus                ClilocId = 0x102e63 // {CLILOCID}  {COUNT}
	Weight1                   ClilocId = 0x105e94 // rune
	Weight2                   ClilocId = 0x105e95 // usually
	RuneLabel                 ClilocId = 0x10395b
	Contents                  ClilocId = 0x105c71 // bags
	Exceptional               ClilocId = 0x102f1c
	ElementalSlayer           ClilocId = 0x102e70
	DamageIncrease            ClilocId = 0x102e31
	FasterCastRecovery        ClilocId = 0x102e3c
	FasterCasting             ClilocId = 0x102e3d
	HitChanceIncrease         ClilocId = 0x102e3f
	DexterityBonus            ClilocId = 0x102e39
	HitFireArea               ClilocId = 0x102e43
	HitFireball               ClilocId = 0x102e44
	HitLifeLeech              ClilocId = 0x102e46
	HitManaLeech              ClilocId = 0x102e4b
	HitPhysicalArea           ClilocId = 0x102e4c
	SwingSpeedIncrease        ClilocId = 0x102e86
	HitPointIncrease          ClilocId = 0x102e4f
	LowerManaCost             ClilocId = 0x102e51
	LowerReagentCost          ClilocId = 0x102e52
	Luck                      ClilocId = 0x102e54
	FireResist                ClilocId = 0x102e5f
	ManaIncrease              ClilocId = 0x102e57
	ManaRegeneration          ClilocId = 0x102e58
	NightSight                ClilocId = 0x102e59
	StaminaRegeneration       ClilocId = 0x102e5b
	HitPointRegeneration      ClilocId = 0x102e5c
	SpellChanneling           ClilocId = 0x102e82
	FireDamage                ClilocId = 0x102e35
	EnergyDamage              ClilocId = 0x102e37
	WeaponDamage              ClilocId = 0x103130
	WeaponSpeed               ClilocId = 0x10312f
	StrengthRequirement       ClilocId = 0x103132
	PhysicalDamage            ClilocId = 0x102e33
	OneHandedWeapon           ClilocId = 0x1033c0
	TwoHandedWeapon           ClilocId = 0x103133
	SkillRequiredSwordmanship ClilocId = 0x103134
	Durability                ClilocId = 0x102f1f
	WeaponLevel               ClilocId = 0xfe472
	Resistances               ClilocId = 0x102f33
)

func (p *MegaClilocPacket) Read(in chan byte) {
	if p.ClientToServer {
		l := p.short(in)
		p.Client = &ClientMegaCliloc{}
		c := p.Client
		c.Serials = make([]uint, 0)
		for i := 3; i < l; i += 4 {
			c.Serials = append(c.Serials, p.uint(in))
		}
	} else {
		p.Server = &ServerMegaCliloc{}
		s := p.Server
		p.short(in) // len
		s.UNK1 = p.short(in)
		s.Serial1 = p.uint(in)
		s.UNK2 = p.short(in)
		s.Serial2 = p.uint(in)
		s.Clilocs = make([]ClilocEntry, 0)
		for {
			id := p.uint(in)
			if id == 0 {
				break
			}
			cliloc := ClilocEntry{Id: id}
			l := p.short(in)
			if l > 0 {
				cliloc.Text = p.bstr(in, l)
			}
			s.Clilocs = append(s.Clilocs, cliloc)
		}
	}
}

func (p *MegaClilocPacket) Bytes() []byte {
	o := p.Writer()
	if p.ClientToServer {
		for _, s := range p.Client.Serials {
			o.uint(s)
		}
	} else {
		s := p.Server
		o.short(s.UNK1)
		o.uint(s.Serial1)
		o.short(s.UNK2)
		o.uint(s.Serial2)
		for _, cliloc := range s.Clilocs {
			o.uint(cliloc.Id)
			l := len(cliloc.Text)
			o.short(l)
			if l > 0 {
				o.bstr(cliloc.Text)
			}
		}
		o.uint(0)
	}
	return o.bytes()
}
