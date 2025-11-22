package core

import (
	"emulator/conv"
	"emulator/types"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
)

// 🟧 Generato a mock 12x40 emulator for testing

func MockEmulator() *Emulator {
	bus := NewBus()
	// 🔥 ignore any errors while loading font -- this code is just
	//    for testing -- plus, we use the same font for normal and bold
	//    as the mock emulator will never actually be rendered
	fontSize := 12.0
	fontPath, _ := findfont.Find("UbuntuMono-R.ttf")
	fontData, _ := os.ReadFile(fontPath)
	ttfFont, _ := truetype.Parse(fontData)
	ttfFace := truetype.NewFace(ttfFont, &truetype.Options{Size: fontSize, DPI: 96})
	// 👇 mock config
	cfg := types.Config{
		BgColor:  "#202020",
		BoldFace: &ttfFace,
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
		Cols:         uint(40),
		FontHeight:   16,
		FontSize:     fontSize,
		FontWidth:    9,
		Monochrome:   false,
		NormalFace:   &ttfFace,
		PaddedHeight: 1.5,
		PaddedWidth:  1.1,
		RGBA:         image.NewRGBA(image.Rect(0, 0, 400, 300)),
		Rows:         uint(12),
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

// 🟧 Smoke test

func TestStreammaker(t *testing.T) {
	emu := MockEmulator()
	emu.Bus.SubRender(func() {
		assert.True(t, true, "smoke test for mock render")
	})
	emu.Init()
	stream := MockStream(types.EW, types.WCC{}, MockExampleImg, MockExampleAttrs)
	emu.Bus.PubOutbound(stream)
}
