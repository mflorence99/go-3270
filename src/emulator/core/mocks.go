//go:build dev

package core

import (
	"emulator/conv"
	"emulator/fonts"
	"emulator/types"
	"image"
	"math"

	"github.com/golang/freetype/truetype"
)

// 🟧 Generato a mock emulator for testing

func MockEmulator(rows, cols uint) *Emulator {
	bus := NewBus()
	// 👇 constants normally computed by the mediator
	dpi := 96.0
	fontHeight := 16.0
	fontSize := 12.0
	fontWidth := 9.0
	paddedHeight := 1.5
	paddedWidth := 1.1
	// 👇 now we can compute the size of the canvas
	canvasWidth := float64(cols) * math.Round(fontWidth*paddedWidth)
	canvasHeight := float64(rows) * math.Round(fontHeight*paddedHeight)
	// 👇 load the fonts
	normalFont, _ := truetype.Parse(fonts.NormalFontEmbed)
	normalFace := truetype.NewFace(normalFont, &truetype.Options{Size: fontSize, DPI: dpi /* , Hinting: font.HintingFull */})
	boldFont, _ := truetype.Parse(fonts.BoldFontEmbed)
	boldFace := truetype.NewFace(boldFont, &truetype.Options{Size: fontSize, DPI: dpi /* , Hinting: font.HintingFull */})
	// 👇 mock config
	cfg := types.Config{
		BgColor:  "#202020",
		BoldFace: &boldFace,
		CLUT: map[types.Color]string{
			0xf0: "#202020",
			0xf1: "#4169e1",
			0xf2: "#ff0000",
			0xf3: "#ee82ee",
			0xf4: "#04c304",
			0xf5: "#40e0d0",
			0xf6: "#ffff00",
			0xf7: "#ffffff",
			0xf8: "#202020",
			0xf9: "#0000cd",
			0xfa: "#ffa500",
			0xfb: "#800080",
			0xfc: "#90ee90",
			0xfd: "#afeeee",
			0xfe: "#c0c0c0",
			0xff: "#e2e2e9"},
		Cols:         cols,
		FontHeight:   fontHeight,
		FontSize:     fontSize,
		FontWidth:    fontWidth,
		Monochrome:   false,
		NormalFace:   &normalFace,
		PaddedHeight: paddedHeight,
		PaddedWidth:  paddedWidth,
		RGBA:         image.NewRGBA(image.Rect(0, 0, int(canvasWidth), int(canvasHeight))),
		Rows:         rows,
	}
	// 👇 mock emulator!
	return NewEmulator(bus, &cfg)
}

// 🟧 Fabricate an outbound stream for a 12x40 display

// 👇 Caller supplies a screen "image" in the form of an array
//    of 40 character strings, as in the example below
//
//    ■​ indicates an unprotected field
//    ¶​ indicates a protected field

var MockExampleImg = []string{
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

type MockRCCoord struct {
	Row uint
	Col uint
}

var MockExampleAttrs = map[MockRCCoord]*types.Attrs{}

// 🟦 Constructor

func MockStream(cmd types.Command, wcc types.WCC, img []string, attrs map[MockRCCoord]*types.Attrs) []byte {
	// 👇 use the convenience of the inbound stream
	out := NewInbound()
	out.Put(byte(cmd))
	out.Put(wcc.Bits())
	// 👇 for each row/col
	for row := 1; row <= len(img); row++ {
		runes := []rune(img[row-1])
		for col := 1; col <= len(runes); col++ {
			addr := uint((row-1)*len(runes) + col - 1)
			char := runes[col-1]
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
