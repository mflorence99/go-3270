package consts

// 🟧 3270 field and extended attributes

type Attrs struct {
	Autoskip   bool
	Blink      bool
	Color      Color
	Hidden     bool
	Highlight  bool
	LCID       LCID
	MDT        bool
	Numeric    bool
	Outline    Outline
	Protected  bool
	Reverse    bool
	Underscore bool

	// 🔥 character attributes are distinguished from field attributes
	CharAttr bool

	// 🔥 initial setting for cells has Default: true to indicate that cell attributes can be overridden by the containing field
	Default bool
}

// 🟦 Constructors

func NewBasicAttrs(char byte) *Attrs {
	a := new(Attrs)
	a.fromByte(char)
	return a
}

func NewExtendedAttrs(chars []byte) *Attrs {
	a := new(Attrs)
	a.fromBytes(chars)
	return a
}

// 🔥 note that we are taking a copy and overwriting deltas
func NewModifiedAttrs(attrs *Attrs, chars []byte) *Attrs {
	a := *attrs
	a.fromBytes(chars)
	a.CharAttr = true
	return &a
}

func (a *Attrs) fromByte(char byte) {
	a.Hidden = ((char & 0b00001000) != 0) && ((char & 0b00000100) != 0)
	a.Highlight = ((char & 0b00001000) != 0) && ((char & 0b00000100) == 0)
	a.MDT = (char & 0b00000001) != 0
	a.Numeric = (char & 0b00010000) != 0
	a.Protected = (char & 0b00100000) != 0
	a.Autoskip = a.Protected && a.Numeric
	// 🔥 set the default color attributes - ignored if monochrome -- checking for "hidden" is a more accurate reading of the spec, but only affects the cursor color
	switch {
	case !a.Protected && (a.Highlight || a.Hidden):
		a.Color = 0xF2
	case !a.Protected && !a.Highlight:
		a.Color = 0xF4
	case a.Protected && (a.Highlight || a.Hidden):
		a.Color = 0xF7
	case a.Protected && !a.Highlight:
		a.Color = 0xF1
	}
}

func (a *Attrs) fromBytes(chars []byte) {
	for ix := 0; ix < len(chars)-1; ix += 2 {
		chunk := chars[ix : ix+2]
		typecode := Typecode(chunk[0])
		switch typecode {

		case BASIC:
			basic := chunk[1]
			a.fromByte(basic)

		case HIGHLIGHT:
			a.Blink = false
			a.Reverse = false
			a.Underscore = false
			a.Highlight = false
			highlight := Highlight(chunk[1])
			switch highlight {
			case BLINK:
				a.Blink = true
			case REVERSE:
				a.Reverse = true
			case UNDERSCORE:
				a.Underscore = true
			case INTENSIFY:
				a.Highlight = true
			}

		case COLOR:
			color := Color(chunk[1])
			a.Color = color

		case CHARSET:
			lcid := LCID(chunk[1])
			a.LCID = lcid

		case OUTLINE:
			outline := Outline(chunk[1])
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
	if a.MDT {
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
