package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2RuneCP037(t *testing.T) {
	assert.Equal(t, E2Rune(0x00, 0x21), rune('\u0020'), "everything below 0x40 is blank")
	assert.Equal(t, E2Rune(0x00, 0x40), rune('\u0020'), "0x40 is blank")
	assert.Equal(t, E2Rune(0x00, 0xF0), rune('\u0030'), "0xF0 is 0")
}

func TestE2RuneCP310(t *testing.T) {
	assert.Equal(t, E2Rune(0xF1, 0x21), rune('\u0020'), "everything below 0x40 is blank")
	assert.Equal(t, E2Rune(0xF1, 0x40), rune('\u0020'), "0x40 is blank")
	assert.Equal(t, E2Rune(0xF1, 0x80), rune(0x00223c), "0x80 is ~")
}

func TestE2RunesCP037(t *testing.T) {
	hello := E2Runes(0x00, string([]byte{200, 133, 147, 147, 150}))
	assert.Equal(t, hello, "Hello", "convert EBCDIC string to runes")
}

func TestE2RunesCP310(t *testing.T) {
	abcde := E2Runes(0xf1, string([]byte{65, 66, 67, 68, 69}))
	assert.Equal(t, abcde, "ABCDE", "convert EBCDIC string to runes")
}
