package core

import (
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

var img = []string{
	/*                 1         2         3         4 */
	/*        1234567890123456789012345678901234567890 */
	/* 01 */ "         ¶Test screen                   ",
	/* 02 */ "                                        ",
	/* 03 */ "¶What is your name ?■                  ¶",
	/* 04 */ "                                        ",
	/* 05 */ "¶Where are you from?■                  ¶",
	/* 06 */ "                                        ",
	/* 07 */ "                                        ",
	/* 08 */ "                                        ",
	/* 09 */ "                                        ",
	/* 10 */ "                                        ",
	/* 11 */ "                                        ",
	/* 12 */ "                             ¶Test # 46b",
}

var attrs = map[MockRCCoord]*types.Attrs{}

func TestNewFlds(t *testing.T) {
	emu := MockEmulator().Initialize()
	stream := MockStream(types.EW, types.WCC{}, img, attrs)
	emu.Bus.PubOutbound(stream)
	flds := emu.Flds.Flds
	t.Run("check Flds created according to stream", func(t *testing.T) {
		assert.Equal(t, 8, len(flds))

		expected := [][]any{
			// #cells addr  prot
			{71, uint(9), true},
			{20, uint(80), true},
			{19, uint(100), false},
			{41, uint(119), true},
			{20, uint(160), true},
			{19, uint(180), false},
			{270, uint(199), true},
			{20, uint(469), true},
		}

		for ix := 0; ix < len(expected); ix++ {
			assert.Equal(t, expected[ix][0], len(flds[ix].Cells))
			addr, ok := flds[ix].Cells[0].GetFldAddr()
			assert.True(t, ok)
			assert.Equal(t, expected[ix][1], addr)
			assert.Equal(t, expected[ix][2], flds[ix].Cells[0].Attrs.Protected)
		}
	})
}

func TestFldsFindFld(t *testing.T) {
	emu := MockEmulator().Initialize()
	stream := MockStream(types.EW, types.WCC{}, img, attrs)
	emu.Bus.PubOutbound(stream)
	var fld *Fld
	var ok bool
	t.Run("can FindFld locate the correct Fld", func(t *testing.T) {
		fld, ok = emu.Flds.FindFld(119)
		assert.Equal(t, 41, len(fld.Cells))
		assert.True(t, ok)

		fld, ok = emu.Flds.FindFld(120)
		assert.False(t, ok)
	})
}
