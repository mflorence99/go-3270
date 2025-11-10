package types

// 🟧 Attribute type

type Typecode byte

// 🟦 Lookup tables

const (
	BASIC     Typecode = 0xC0
	HIGHLIGHT Typecode = 0x41
	COLOR     Typecode = 0x42
	CHARSET   Typecode = 0x43
	OUTLINE   Typecode = 0xC2
)

var typecodes = map[Typecode]string{
	0xC0: "BASIC",
	0x41: "HIGHLIGHT",
	0x42: "COLOR",
	0x43: "CHARSET",
	0xC2: "OUTLINE",
}

// 🟦 Stringer implementation

func TypecodeFor(t Typecode) string {
	return typecodes[t]
}

func (t Typecode) String() string {
	return TypecodeFor(t)
}
