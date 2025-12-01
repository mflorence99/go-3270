package core

import (
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fldsImg = []string{
	/*                 1         2         3         4 */
	/*        1234567890123456789012345678901234567890 */
	/* 01 */ "         ¶Test screen                   ",
	/* 02 */ "                                        ",
	/* 03 */ "¶What is your name ?■Mark Florence     ¶",
	/* 04 */ "                                        ",
	/* 05 */ "¶Where are you from?■Coos Bay, Or      ¶",
	/* 06 */ "                                        ",
	/* 07 */ "                                        ",
	/* 08 */ "                                        ",
	/* 09 */ "                                        ",
	/* 10 */ "                                        ",
	/* 11 */ "                                        ",
	/* 12 */ "                             ¶Test # 46b",
}

var fldsAttrs = map[MockRCCoord]*types.Attrs{
	{3, 21}: {MDT: true},
	{3, 22}: {Color: types.RED, Highlight: true},
	{3, 23}: {Color: types.GREEN, Reverse: true},
	{3, 24}: {Color: types.BLUE, Blink: true},
	{5, 21}: {MDT: true},
	{5, 22}: {Color: types.RED, Highlight: true},
	{5, 23}: {Color: types.GREEN, Reverse: true},
	{5, 24}: {Color: types.BLUE, Blink: true},
}

func TestNewFlds(t *testing.T) {
	emu := MockEmulator(12, 40).Initialize()
	stream := MockStream(types.EW, types.WCC{}, fldsImg, fldsAttrs)
	emu.Bus.PubOutbound(stream)
	flds := emu.Flds.Flds
	t.Run("check Flds created according to stream", func(t *testing.T) {
		assert.Equal(t, 8, len(flds))

		type xpctd struct {
			numCells int
			addr     uint
			prot     bool
		}

		xpctds := []xpctd{
			{71, 9, true},
			{20, 80, true},
			{19, 100, false},
			{41, 119, true},
			{20, 160, true},
			{19, 180, false},
			{270, 199, true},
			{20, 469, true},
		}

		for ix, expected := range xpctds {
			assert.Equal(t, expected.numCells, len(flds[ix].Cells))
			addr, ok := flds[ix].Cells[0].GetFldAddr()
			assert.True(t, ok)
			assert.Equal(t, expected.addr, addr)
			assert.Equal(t, expected.prot, flds[ix].Cells[0].Attrs.Protected)
		}
	})
}

func TestFldsFindFld(t *testing.T) {
	emu := MockEmulator(12, 40).Initialize()
	stream := MockStream(types.EW, types.WCC{}, fldsImg, fldsAttrs)
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

func TestFldsRM(t *testing.T) {
	emu := MockEmulator(12, 40).Initialize()
	stream := MockStream(types.EW, types.WCC{}, fldsImg, fldsAttrs)
	emu.Bus.PubOutbound(stream)
	t.Run("can RM return modified fields * char attrs", func(t *testing.T) {
		chars := emu.Flds.RM()
		_ = chars
	})
}
