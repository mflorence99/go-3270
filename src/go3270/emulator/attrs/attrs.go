package attrs

import "go3270/emulator/consts"

// 🟧 3270 field and extended attributes

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
	Outline    consts.Outline
	LCID       consts.LCID
}

// 🟦 Constructors

func NewBasic(char byte) *Attrs {
	a := new(Attrs)
	a.fromByte(char)
	return a
}

func NewExtended(chars []byte) *Attrs {
	a := new(Attrs)
	a.fromBytes(chars)
	return a
}

// 🔥 note that we are taking a copy and overwriting deltas
func NewModified(attrs *Attrs, chars []byte) *Attrs {
	a := *attrs
	a.fromBytes(chars)
	return &a
}

func (a *Attrs) fromByte(char byte) {
	a.Hidden = ((char & 0b00001000) != 0) && ((char & 0b00000100) != 0)
	a.Highlight = ((char & 0b00001000) != 0) && ((char & 0b00000100) == 0)
	a.Modified = (char & 0b00000001) != 0
	a.Numeric = (char & 0b00010000) != 0
	a.Protected = (char & 0b00100000) != 0
	// 👇 set the default color attributes - ignored if monochrome
	switch {
	case !a.Protected && !a.Highlight:
		a.Color = 0xF4
	case !a.Protected && a.Highlight:
		a.Color = 0xF2
	case a.Protected && !a.Highlight:
		a.Color = 0xF1
	case a.Protected && a.Highlight:
		a.Color = 0xF7
	}
}

func (a *Attrs) fromBytes(chars []byte) {
	for ix := 0; ix < len(chars)-1; ix += 2 {
		chunk := chars[ix : ix+2]
		typecode := consts.Typecode(chunk[0])
		switch typecode {

		case consts.BASIC:
			basic := chunk[1]
			a.fromByte(basic)

		case consts.HIGHLIGHT:
			a.Blink = false
			a.Reverse = false
			a.Underscore = false
			a.Highlight = false
			highlight := consts.Highlight(chunk[1])
			switch highlight {
			case consts.BLINK:
				a.Blink = true
			case consts.REVERSE:
				a.Reverse = true
			case consts.UNDERSCORE:
				a.Underscore = true
			case consts.INTENSIFY:
				a.Highlight = true
			}

		case consts.COLOR:
			color := consts.Color(chunk[1])
			a.Color = color

		case consts.CHARSET:
			lcid := consts.LCID(chunk[1])
			a.LCID = lcid

		case consts.OUTLINE:
			outline := consts.Outline(chunk[1])
			a.Outline = outline
		}
	}
}

// 🟦 Public functions

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
