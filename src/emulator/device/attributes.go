package device

import (
	"emulator/types"
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

func NewAttributes(bytes []byte) *Attributes {
	// 👇 quick exit for one-byte attribute
	if len(bytes) == 1 {
		return NewAttributes([]byte{types.TypeCodeLookup["BASIC"], bytes[0]})
	}
	// 👇 now pretend we have an extended attribute and analze bytes in pairs
	attrs := new(Attributes)
	for ix := 0; ix < len(bytes)-1; ix += 2 {
		chunk := bytes[ix : ix+2]
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
