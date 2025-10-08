package consts

var (
	BASIC     byte = 0xC0
	COLOR     byte = 0x42
	HIGHLIGHT byte = 0x41
)

var typecodes = map[byte]string{
	0xC0: "BASIC",
	0x41: "HIGHLIGHT",
	0x42: "COLOR",
}

func TypecodeFor(typecode byte) string {
	return typecodes[typecode]
}
