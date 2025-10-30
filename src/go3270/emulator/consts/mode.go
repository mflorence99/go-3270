package consts

type Mode byte

const (
	FIELD_MODE          Mode = 0x00
	EXTENDED_FIELD_MODE Mode = 0x01
	CHARACTER_MODE      Mode = 0x02
)

var modes = map[Mode]string{
	0x00: "FIELD_MODE",
	0x01: "EXTENDED_FIELD_MODE",
	0x02: "CHARACTER_MODE",
}

func ModeFor(m Mode) string {
	return modes[m]
}

func (m Mode) String() string {
	return ModeFor(m)
}
