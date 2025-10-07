package device

type Attributes struct {
	blink      bool
	color      byte
	hidden     bool
	highlight  bool
	modified   bool
	numeric    bool
	protected  bool
	reverse    bool
	underscore bool
}

func NewAttribute(u8 byte) *Attributes {
	return NewAttributes([]byte{TypeCodeLookup["BASIC"], u8})
}

func NewAttributes(u8s []byte) *Attributes {
	// 👇 now pretend we have an extended attribute and analze bytes in pairs
	attrs := new(Attributes)
	for ix := 0; ix < len(u8s)-1; ix += 2 {
		chunk := u8s[ix : ix+2]
		switch chunk[0] {
		case TypeCodeLookup["BASIC"]:
			attrs.hidden = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) != 0)
			attrs.highlight = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) == 0)
			attrs.modified = (chunk[1] & 0b00000001) != 0
			attrs.numeric = (chunk[1] & 0b00010000) != 0
			attrs.protected = (chunk[1] & 0b00100000) != 0
		case TypeCodeLookup["HIGHLIGHT"]:
			switch chunk[1] {
			case HighlightLookup["BLINK"]:
				attrs.blink = true
			case HighlightLookup["REVERSE"]:
				attrs.reverse = true
			case HighlightLookup["UNDERSCORE"]:
				attrs.underscore = true
			}
		case TypeCodeLookup["COLOR"]:
			attrs.color = chunk[1]
		}
	}
	return attrs
}

func NewProtectedAttribute() *Attributes {
	return NewAttributes([]byte{TypeCodeLookup["BASIC"], 0b00100000})
}

func (attrs *Attributes) Color(dflt [2]string) string {
	colors, ok := CLUT[attrs.color]
	if !ok {
		colors = dflt
	}
	if attrs.Highlight() {
		return colors[0]
	} else {
		return colors[1]
	}
}

func (attrs *Attributes) Blink() bool {
	return attrs.blink
}

func (attrs *Attributes) Hidden() bool {
	return attrs.hidden
}

func (attrs *Attributes) Highlight() bool {
	return attrs.highlight
}

func (attrs *Attributes) Modified() bool {
	return attrs.modified
}

func (attrs *Attributes) Numeric() bool {
	return attrs.numeric
}

func (attrs *Attributes) Protected() bool {
	return attrs.protected
}

func (attrs *Attributes) Reverse() bool {
	return attrs.reverse
}

func (attrs *Attributes) Underscore() bool {
	return attrs.underscore
}

func (attrs *Attributes) ToByte() byte {
	var u8 byte = 0b00000000
	if attrs.Hidden() {
		u8 &= 0b00001100
	}
	if attrs.Highlight() {
		u8 &= 0b00001000
	}
	if attrs.Modified() {
		u8 &= 0b00000001
	}
	if attrs.Numeric() {
		u8 &= 0b00010000
	}
	if attrs.Protected() {
		u8 &= 0b00100000
	}
	return Six2E[u8]
}

func (attrs *Attributes) ToString() string {
	str := "ATTR=[ "
	if attrs.Blink() {
		str += "BLINK "
	}
	if attrs.Hidden() {
		str += "HIDDEN "
	}
	if attrs.Highlight() {
		str += "HILITE "
	}
	if attrs.Modified() {
		str += "MDT "
	}
	if attrs.Numeric() {
		str += "NUM "
	}
	if attrs.Protected() {
		str += "PROT "
	}
	if attrs.Reverse() {
		str += "REV "
	}
	if attrs.Underscore() {
		str += "USCORE "
	}
	str += "]"
	return str
}
