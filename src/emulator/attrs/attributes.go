package attrs

import (
	"emulator/consts"
)

type Attributes struct {
	Blink      bool
	Color      byte
	Hidden     bool
	Highlight  bool
	Modified   bool
	Numeric    bool
	Protected  bool
	Reverse    bool
	Underscore bool
}

func New(bytes []byte) *Attributes {
	// 👇 treat a single-byte attribute as BASIC
	if len(bytes) == 1 {
		bytes = []byte{byte(consts.BASIC), bytes[0]}
	}
	// 👇 now just look at extended attributes in pairs
	a := new(Attributes)
	for ix := 0; ix < len(bytes)-1; ix += 2 {
		chunk := bytes[ix : ix+2]
		typecode := consts.Typecode(chunk[0])
		highlight := consts.Highlight(chunk[1])
		switch typecode {
		case consts.BASIC:
			a.Hidden = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) != 0)
			a.Highlight = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) == 0)
			a.Modified = (chunk[1] & 0b00000001) != 0
			a.Numeric = (chunk[1] & 0b00010000) != 0
			a.Protected = (chunk[1] & 0b00100000) != 0
		case consts.HIGHLIGHT:
			switch highlight {
			case consts.BLINK:
				a.Blink = true
			case consts.REVERSE:
				a.Reverse = true
			case consts.UNDERSCORE:
				a.Underscore = true
			}
		case consts.COLOR:
			a.Color = chunk[1]
		}
	}
	return a
}

func (a *Attributes) Byte() byte {
	var char byte = 0b00000000
	if a.Hidden {
		char |= 0b00001100
	}
	if a.Highlight {
		char |= 0b00001000
	}
	if a.Modified {
		char |= 0b00000001
	}
	if a.Numeric {
		char |= 0b00010000
	}
	if a.Protected {
		char |= 0b00100000
	}
	return char
}

func (a *Attributes) String() string {
	str := "ATTR=[ "
	if a.Blink {
		str += "BLINK "
	}
	if a.Hidden {
		str += "HIDDEN "
	}
	if a.Highlight {
		str += "HILITE "
	}
	if a.Modified {
		str += "MDT "
	}
	if a.Numeric {
		str += "NUM "
	}
	if a.Protected {
		str += "PROT "
	}
	if a.Reverse {
		str += "REV "
	}
	if a.Underscore {
		str += "USCORE "
	}
	str += "]"
	return str
}
