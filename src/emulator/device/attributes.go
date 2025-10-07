package device

import (
	"emulator/types"
	"emulator/utils"
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

func NewAttribute(u8 byte) *Attributes {
	return NewAttributes([]byte{types.TypeCodeLookup["BASIC"], u8})
}

func NewAttributes(u8s []byte) *Attributes {
	// 👇 now pretend we have an extended attribute and analze bytes in pairs
	attrs := new(Attributes)
	for ix := 0; ix < len(u8s)-1; ix += 2 {
		chunk := u8s[ix : ix+2]
		switch chunk[0] {
		case types.TypeCodeLookup["BASIC"]:
			attrs.hidden = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) != 0)
			attrs.highlight = ((chunk[1] & 0b00001000) != 0) && ((chunk[1] & 0b00000100) == 0)
			attrs.modified = (chunk[1] & 0b00000001) != 0
			attrs.numeric = (chunk[1] & 0b00010000) != 0
			attrs.protected = (chunk[1] & 0b00100000) != 0
		case types.TypeCodeLookup["HIGHLIGHT"]:
			switch chunk[1] {
			case types.HighlightLookup["BLINK"]:
				attrs.blink = true
			case types.HighlightLookup["REVERSE"]:
				attrs.reverse = true
			case types.HighlightLookup["UNDERSCORE"]:
				attrs.underscore = true
			}
		case types.TypeCodeLookup["COLOR"]:
			attrs.color = chunk[1]
		}
	}
	return attrs
}

func (attrs *Attributes) GetColor(dflt string) string {
	colors := types.CLUT[attrs.color]
	if colors == nil {
		return dflt
	} else if attrs.IsHighlight() {
		return colors[0]
	} else {
		return colors[1]
	}
}

func (attrs *Attributes) IsBlink() bool {
	return attrs.blink
}

func (attrs *Attributes) IsHidden() bool {
	return attrs.hidden
}

func (attrs *Attributes) IsHighlight() bool {
	return attrs.highlight
}

func (attrs *Attributes) IsModified() bool {
	return attrs.modified
}

func (attrs *Attributes) IsNumeric() bool {
	return attrs.numeric
}

func (attrs *Attributes) IsProtected() bool {
	return attrs.protected
}

func (attrs *Attributes) IsReverse() bool {
	return attrs.reverse
}

func (attrs *Attributes) IsUnderscore() bool {
	return attrs.underscore
}

func (attrs *Attributes) ToByte() byte {
	var u8 byte = 0b00000000
	if attrs.IsProtected() {
		u8 &= 0b00100000
	}
	if attrs.IsNumeric() {
		u8 &= 0b00010000
	}
	if attrs.IsHighlight() {
		u8 &= 0b00001000
	}
	if attrs.IsHidden() {
		u8 &= 0b00001100
	}
	if attrs.IsModified() {
		u8 &= 0b00000001
	}
	return utils.Six2E[u8]
}
