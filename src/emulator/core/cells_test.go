package core

import (
	"emulator/conv"
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cellsImg = []string{
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

var cellsAttrs = map[MockRCCoord]*types.Attrs{
	{12, 31}: {Color: types.RED},
	{12, 32}: {Color: types.GREEN},
	{12, 33}: {Color: types.BLUE},
}

func TestNewCells(t *testing.T) {
	emu := MockEmulator(12, 40).Initialize()
	stream := MockStream(types.EW, types.WCC{}, cellsImg, cellsAttrs)
	emu.Bus.PubOutbound(stream)
	t.Run("check Cells created according to stream", func(t *testing.T) {

		type xpctd struct {
			row   uint
			col   uint
			sf    bool
			home  bool
			char  byte
			color types.Color
		}

		xpctds := []xpctd{
			{12, 30, true, false, 0x1d, 0x00},
			{12, 31, false, true, conv.A2E('T'), types.RED},
			{12, 32, false, false, conv.A2E('e'), types.GREEN},
			{12, 33, false, false, conv.A2E('s'), types.BLUE},
			{12, 34, false, false, conv.A2E('t'), types.BLUE},
		}

		for _, expected := range xpctds {

			addr := emu.Cfg.RC2Addr(expected.row, expected.col)
			cell, ok := emu.Buf.Peek(addr)
			assert.True(t, ok)

			assert.Equal(t, expected.sf, cell.IsFldStart())
			assert.Equal(t, expected.home, cell.IsFldHome())
			assert.Equal(t, expected.char, cell.Char)
			assert.Equal(t, expected.color, cell.Attrs.Color)

			if cell.IsFldStart() {
				sf, _ := cell.GetFldStart()
				assert.Equal(t, cell, sf)
			}

			if cell.IsFldHome() {
				home, _ := cell.GetFldHome()
				assert.Equal(t, cell, home)
			}
		}
	})
}
