package consts_test

import (
	"fmt"
	"go3270/emulator/consts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AIDFor(t *testing.T) {
	assert.True(t, consts.AIDFor(0x6C) == "PA1")
	assert.True(t, consts.AIDFor(0x6D) == "CLEAR")
	assert.True(t, consts.AIDFor(0xFF) == "")
}

func Test_AIDFOf(t *testing.T) {
	assert.True(t, consts.AIDOf("Enter", false, false, false) == consts.ENTER)
	assert.True(t, consts.AIDOf("Escape", false, false, false) == consts.CLEAR)
	assert.True(t, consts.AIDOf("F1", false, false, false) == consts.PF1)
	assert.True(t, consts.AIDOf("F3", false, false, false) == consts.PF3)
	assert.True(t, consts.AIDOf("F1", false, false, true) == consts.PF13)
	assert.True(t, consts.AIDOf("F3", false, false, true) == consts.PF15)
	assert.True(t, consts.AIDOf("F1", true, false, false) == consts.PA1)
	assert.True(t, consts.AIDOf("F3", true, false, false) == consts.PA3)
	assert.True(t, consts.AIDOf("F4", true, false, false) == 0)
	assert.True(t, consts.AIDOf("Backspace", true, false, false) == 0)
}

func Test_PAx(t *testing.T) {
	assert.True(t, consts.AID(0x6C).PAx())
	assert.False(t, consts.AID(0x00).PAx())
	assert.False(t, consts.AID(0x88).PAx())
}

func Test_PFx(t *testing.T) {
	assert.True(t, consts.AID(0x4C).PFx())
	assert.True(t, consts.AID(0xC5).PFx())
	assert.False(t, consts.AID(0x88).PFx())
}

func Test_ColorFor(t *testing.T) {
	assert.True(t, consts.ColorFor(0xF2) == "RED")
	assert.True(t, consts.ColorFor(0xF7) == "WHITE")
	assert.True(t, consts.ColorFor(0xFF) == "")
}

func Test_ColorOf(t *testing.T) {
	assert.True(t, consts.ColorOf("RED") == 0xF2)
	assert.True(t, consts.ColorOf("WHITE") == 0xF7)
}

func Test_CommandFor(t *testing.T) {
	assert.True(t, consts.CommandFor(0x6F) == "EAU")
	assert.True(t, consts.CommandFor(0xF5) == "EW")
	assert.True(t, consts.CommandFor(0xFF) == "")
}

func Test_HighlightFor(t *testing.T) {
	assert.True(t, consts.HighlightFor(0xF1) == "BLINK")
	assert.True(t, consts.HighlightFor(0xF2) == "REVERSE")
	assert.True(t, consts.HighlightFor(0xF3) == "")
}

func Test_OrderFor(t *testing.T) {
	assert.True(t, consts.OrderFor(0x05) == "PT")
	assert.True(t, consts.OrderFor(0x29) == "SFE")
	assert.True(t, consts.OrderFor(0xF3) == "")
}

func Test_TypecodeFor(t *testing.T) {
	assert.True(t, consts.TypecodeFor(0xC0) == "BASIC")
	assert.True(t, consts.TypecodeFor(0x42) == "COLOR")
	assert.True(t, consts.TypecodeFor(0xF3) == "")
}

func Test_Strings(t *testing.T) {
	assert.True(t, fmt.Sprintf("AID=%s", consts.AID(0x88)) == "AID=INBOUND")
	assert.True(t, fmt.Sprintf("Color=%s", consts.Color(0xF6)) == "Color=YELLOW")
	assert.True(t, fmt.Sprintf("Command=%s", consts.Command(0x6F)) == "Command=EAU")
	assert.True(t, fmt.Sprintf("Highlight=%s", consts.Highlight(0xF4)) == "Highlight=UNDERSCORE")
	assert.True(t, fmt.Sprintf("Order=%s", consts.Order(0x29)) == "Order=SFE")
	assert.True(t, fmt.Sprintf("Typecode=%s", consts.Typecode(0x41)) == "Typecode=HIGHLIGHT")
}
