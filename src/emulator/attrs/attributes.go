package attrs

import (
	"emulator/consts"
)

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
			attrs.hidden = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) != 0)
			attrs.highlight = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) == 0)
			attrs.modified = (chunk[1] & 0b00000001) != 0
			attrs.numeric = (chunk[1] & 0b00010000) != 0
			attrs.protected = (chunk[1] & 0b00100000) != 0
		case consts.HIGHLIGHT:
			switch highlight {
			case consts.BLINK:
				attrs.blink = true
			case consts.REVERSE:
				attrs.reverse = true
			case consts.UNDERSCORE:
				attrs.underscore = true
			}
		case consts.COLOR:
			attrs.color = chunk[1]
		}
	}
	return attrs
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

func (attrs *Attributes) Byte() byte {
	var char byte = 0b00000000
	if attrs.Hidden() {
		char |= 0b00001100
	}
	if attrs.Highlight() {
		char |= 0b00001000
	}
	if attrs.Modified() {
		char |= 0b00000001
	}
	if attrs.Numeric() {
		char |= 0b00010000
	}
	if attrs.Protected() {
		char |= 0b00100000
	}
	return char
}

func (attrs *Attributes) String() string {
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
