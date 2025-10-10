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

func TypecodeFor(t Typecode) string {
	return typecodes[t]
}

func (t Typecode) String() string {
	return TypecodeFor(t)
}
