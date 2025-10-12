package attrs

import (
	"go3270/consts"
)

type Attrs struct {
	Blink      bool
	Color      consts.Color
	Hidden     bool
	Highlight  bool
	Modified   bool
	Numeric    bool
	Protected  bool
	Reverse    bool
	Underscore bool
}

func NewBasic(basic byte) *Attrs {
	a := new(Attrs)
	a.fromByte(basic)
	return a
}

func NewExtended(bytes []byte) *Attrs {
	a := new(Attrs)
	for ix := 0; ix < len(bytes)-1; ix += 2 {
		chunk := bytes[ix : ix+2]
		typecode := consts.Typecode(chunk[0])
		basic := chunk[1]
		color := consts.Color(chunk[1])
		highlight := consts.Highlight(chunk[1])
		switch typecode {
		case consts.BASIC:
			a.fromByte(basic)
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
			a.Color = color
		}
	}
	return a
}

func (a *Attrs) fromByte(char byte) {
	a.Hidden = ((char & 0b00001000) != 0) && ((char & 0b00000100) != 0)
	a.Highlight = ((char & 0b00001000) != 0) && ((char & 0b00000100) == 0)
	a.Modified = (char & 0b00000001) != 0
	a.Numeric = (char & 0b00010000) != 0
	a.Protected = (char & 0b00100000) != 0
}

func (a *Attrs) Byte() byte {
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

func (a *Attrs) String() string {
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
