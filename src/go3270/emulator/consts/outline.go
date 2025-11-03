package consts

// 🟧 3270 outline extended attribute

type Outline byte

// 🟦 Lookup tables

const (
	OUTLINE_BOTTOM Outline = 0b00000001
	OUTLINE_RIGHT  Outline = 0b00000010
	OUTLINE_TOP    Outline = 0b00000100
	OUTLINE_LEFT   Outline = 0b00001000
)

// 🟦 Stringer implementation

func OutlineFor(o Outline) string {
	if o == 0b00000000 {
		return "NONE"
	} else {
		str := ""
		if (o & 0b00000001) != 0 {
			str += "B"
		}
		if (o & 0b00000010) != 0 {
			str += "R"
		}
		if (o & 0b00000100) != 0 {
			str += "T"
		}
		if (o & 0b00001000) != 0 {
			str += "L"
		}
		return str
	}
}

func (o Outline) String() string {
	return OutlineFor(o)
}
