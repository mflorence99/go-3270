package conv

var ASCII = [256]byte{}

func A2E(a []byte) []byte {
	e := make([]byte, len(a))
	for ix := 0; ix < len(a); ix++ {
		if a[ix] == ' ' {
			e[ix] = 0x40
		} else {
			e[ix] = ASCII[a[ix]]
		}
	}
	return e
}

func init() {
	// 👇 the EBCDIC table starts at 0x40 because of the difficulty of eyeballing 64 0x00's
	for ix := 0; ix < len(EBCDIC); ix++ {
		ASCII[EBCDIC[ix]] = byte(ix + 0x40)
	}
}
