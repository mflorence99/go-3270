package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBasicAttrs(t *testing.T) {
	a := &Attrs{
		Autoskip:  true,
		Highlight: true,
		MDT:       true,
		Numeric:   true,
		Protected: true,
	}
	assert.Equal(t, *a, *NewBasicAttrs(0b00111001), "create attrs from bits")
}

func TestNewExtendedAttrs(t *testing.T) {
	a := &Attrs{
		Autoskip:   true,
		Blink:      false,
		Color:      Color(0xf4),
		Hidden:     true,
		Highlight:  false,
		LCID:       LCID(0xf1),
		MDT:        true,
		Numeric:    true,
		Outline:    Outline(0b00001111),
		Protected:  true,
		Reverse:    false,
		Underscore: true,
	}
	b := []byte{
		byte(BASIC),
		0b00111101,
		byte(HIGHLIGHT),
		byte(UNDERSCORE),
		byte(COLOR),
		0xf4,
		byte(CHARSET),
		0xf1,
		byte(OUTLINE),
		0b00001111,
	}
	assert.Equal(t, *a, *NewExtendedAttrs(b), "create attrs from bytes")
}

func TestBits(t *testing.T) {
	a := &Attrs{
		Hidden:    true,
		Highlight: true,
		MDT:       true,
		Numeric:   true,
		Protected: true,
	}
	assert.Equal(t, a.Bits(), byte(0b00111101), "decode attrs to bit settings")
}

func TestBytes(t *testing.T) {
	a := &Attrs{
		Autoskip:   true,
		Blink:      false,
		Color:      Color(0xf4),
		Hidden:     true,
		Highlight:  false,
		LCID:       LCID(0xf1),
		MDT:        true,
		Numeric:    true,
		Outline:    Outline(0b00001111),
		Protected:  true,
		Reverse:    false,
		Underscore: true,
	}
	b := []byte{
		byte(BASIC),
		0b00111101,
		byte(HIGHLIGHT),
		byte(UNDERSCORE),
		byte(COLOR),
		0xf4,
		byte(CHARSET),
		0xf1,
		byte(OUTLINE),
		0b00001111,
	}
	assert.Equal(t, b, a.Bytes(), "decode attrs to bytes")
}

func TestDiff(t *testing.T) {
	a := &Attrs{
		Autoskip: true,
		Blink:    true,
		Color:    Color(0xf4),
	}
	b := &Attrs{
		Autoskip:  true,
		Color:     Color(0xf7),
		Hidden:    true,
		Highlight: true,
	}
	d := &Attrs{
		Autoskip: false,
		Blink:    true,
		Color:    Color(0xf4),
	}
	assert.Equal(t, d, a.Diff(b), "diff of two attrs")
}

func TestAttrsStringer(t *testing.T) {
	a := &Attrs{
		Autoskip:   true,
		Blink:      true,
		Color:      Color(0xf4),
		Hidden:     true,
		Highlight:  true,
		LCID:       LCID(0xf1),
		MDT:        true,
		Numeric:    true,
		Outline:    Outline(0b00001111),
		Protected:  true,
		Reverse:    true,
		Underscore: true,
	}
	str := "SKIP BLINK GREEN HIDDEN HILITE MDT NUM PROT REV USCORE BRTL f1 "
	assert.Equal(t, str, a.String(), "attrs stringified")
	assert.Equal(t, str, AttrsFor(a), "attrs stringified")
}
