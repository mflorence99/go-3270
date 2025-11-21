package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAidOf(t *testing.T) {
	assert.Equal(t, AIDOf("enter", false, false, false), ENTER, "Enter key")
	assert.Equal(t, AIDOf("escape", false, false, false), CLEAR, "Esc key")
	assert.Equal(t, AIDOf("f1", false, false, false), PF1, "F1 key")
	assert.Equal(t, AIDOf("f1", false, false, true), PF13, "F1+shift key")
	assert.Equal(t, AIDOf("f1", true, false, false), PA1, "F1+alt key")
}

func TestPAx(t *testing.T) {
	pa1 := AIDOf("f1", true, false, false)
	assert.True(t, pa1.PAx(), "F1+alt indicates attn key")

	pf1 := AIDOf("f1", false, false, false)
	assert.False(t, pf1.PAx(), "F1 does not indicate attn key")
}

func TestPFx(t *testing.T) {
	pa1 := AIDOf("f1", true, false, false)
	assert.False(t, pa1.PFx(), "F1+alt is not a PFx key")

	pf1 := AIDOf("f1", false, false, false)
	assert.True(t, pf1.PFx(), "F1 is a PFx key")
}

func TestShortRead(t *testing.T) {
	assert.True(t, CLEAR.ShortRead(), "CLEAR triggers a short read")
	assert.True(t, PA3.ShortRead(), "PA3 triggers a short read")
	assert.False(t, ENTER.ShortRead(), "ENTER does not trigger a short read")
}

func TestAIDStringer(t *testing.T) {
	assert.Equal(t, "ENTER", ENTER.String(), "ENTER stringified")
	assert.Equal(t, "ENTER", AIDFor(ENTER), "ENTER stringified")
}
