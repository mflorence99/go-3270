package main

import (
	"syscall/js"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func font(this js.Value, args []js.Value) interface{} {
	font, err := truetype.Parse(TerminalFontData)
	if err != nil {
		Log.Invoke(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	dc := gg.NewContext(1024, 1024)
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Hello, world!", 512, 512, 0.5, 0.5)
	return 0
}
