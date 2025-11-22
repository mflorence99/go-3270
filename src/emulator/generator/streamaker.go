package generator

import (
	"emulator/conv"
	"emulator/core"
	"emulator/types"
)

// 🟧 Fabricate an outbound stream for a 12x40 display

// 👇 Caller supplies a screen "image" in the form of an array
//    of 40 character strings, as in the example below
//
//    ■​ indicates an unprotected field
//    ¶​ indicates a protected field

var ExampleImg = []string{
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

// 👇  additional attribute information can be optionally
//     provided for specified row/col positions

type Coord struct {
	Row uint
	Col uint
}

var ExampleAttrs = map[Coord]*types.Attrs{}

// 🟦 Constructor

func MakeStream(cmd types.Command, wcc types.WCC, img []string, attrs map[Coord]*types.Attrs) []byte {
	// 👇 use the convenience of the inbound stream
	out := core.NewInbound()
	out.Put(byte(cmd))
	out.Put(wcc.Bits())
	// 👇 for each row/col
	for row := 1; row <= len(img); row++ {
		for col := 1; col <= len(img[row-1]); col++ {
			addr := uint((row-1)*len(img[row-1]) + col - 1)
			char := rune(img[row-1][col-1])
			switch char {

			case '¶':
				out.Put(byte(types.SBA))
				out.PutSlice(conv.Addr2Bytes(addr))
				out.Put(byte(types.SF))
				out.Put((&types.Attrs{Protected: true}).Bits())

			case '■':
				out.Put(byte(types.SBA))
				out.PutSlice(conv.Addr2Bytes(addr))
				out.Put(byte(types.SF))
				out.Put((&types.Attrs{}).Bits())

			default:
				out.Put(conv.A2E(byte(char)))

			}
		}
	}
	return out.Bytes()
}
