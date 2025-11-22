package generator

import (
	_ "embed"
	"emulator/core"
	"emulator/types"
	"image"

	"github.com/golang/freetype/truetype"
)

var (
	//go:embed JuliaMono-Regular.ttf
	normalFontEmbed []byte
	//go:embed JuliaMono-Bold.ttf
	boldFontEmbed []byte
)

// 🟧 Generato a mock 12x40 emulator for testing

func MockEmulator() *core.Emulator {
	bus := core.NewBus()
	// 👇 load the fonts
	fontSize := 12.0
	normalFont, _ := truetype.Parse(normalFontEmbed)
	normalFace := truetype.NewFace(normalFont, &truetype.Options{Size: fontSize, DPI: 96})
	boldFont, _ := truetype.Parse(boldFontEmbed)
	boldFace := truetype.NewFace(boldFont, &truetype.Options{Size: fontSize, DPI: 96})
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
		Cols:         uint(40),
		FontHeight:   16,
		FontSize:     fontSize,
		FontWidth:    9,
		Monochrome:   false,
		NormalFace:   &normalFace,
		PaddedHeight: 1.5,
		PaddedWidth:  1.1,
		RGBA:         image.NewRGBA(image.Rect(0, 0, 400, 300)),
		Rows:         uint(12),
	}
	// 👇 mock emulator!
	return core.NewEmulator(bus, &cfg)
}
