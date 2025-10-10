package consts

type Typecode byte

const (
	BASIC     Typecode = 0xC0
	COLOR     Typecode = 0x42
	HIGHLIGHT Typecode = 0x41
)

var typecodes = map[Typecode]string{
	0xC0: "BASIC",
	0x41: "HIGHLIGHT",
	0x42: "COLOR",
}

func TypecodeFor(typecode Typecode) string {
	return typecodes[typecode]
}

func (typecode Typecode) String() string {
	return TypecodeFor(typecode)
}
