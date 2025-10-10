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
	attrs := new(Attributes)
	for ix := 0; ix < len(bytes)-1; ix += 2 {
		chunk := bytes[ix : ix+2]
		typecode := consts.Typecode(chunk[0])
		highlight := consts.Highlight(chunk[1])
		switch typecode {
		case consts.BASIC:
			attrs.Hidden = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) != 0)
			attrs.Highlight = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) == 0)
			attrs.Modified = (chunk[1] & 0b00000001) != 0
			attrs.Numeric = (chunk[1] & 0b00010000) != 0
			attrs.Protected = (chunk[1] & 0b00100000) != 0
		case consts.HIGHLIGHT:
			switch highlight {
			case consts.BLINK:
				attrs.Blink = true
			case consts.REVERSE:
				attrs.Reverse = true
			case consts.UNDERSCORE:
				attrs.Underscore = true
			}
		case consts.COLOR:
			attrs.Color = chunk[1]
		}
	}
	return attrs
}

func (attrs *Attributes) Byte() byte {
	var char byte = 0b00000000
	if attrs.Hidden {
		char |= 0b00001100
	}
	if attrs.Highlight {
		char |= 0b00001000
	}
	if attrs.Modified {
		char |= 0b00000001
	}
	if attrs.Numeric {
		char |= 0b00010000
	}
	if attrs.Protected {
		char |= 0b00100000
	}
	return char
}

func (attrs *Attributes) String() string {
	str := "ATTR=[ "
	if attrs.Blink {
		str += "BLINK "
	}
	if attrs.Hidden {
		str += "HIDDEN "
	}
	if attrs.Highlight {
		str += "HILITE "
	}
	if attrs.Modified {
		str += "MDT "
	}
	if attrs.Numeric {
		str += "NUM "
	}
	if attrs.Protected {
		str += "PROT "
	}
	if attrs.Reverse {
		str += "REV "
	}
	if attrs.Underscore {
		str += "USCORE "
	}
	str += "]"
	return str
}
