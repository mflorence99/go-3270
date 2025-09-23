package main

import (
	_ "embed"
	"fmt"
	"image"
	"math"
	"math/rand"
	"syscall/js"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// 🔥 Hack alert! we must use extension {js, wasm} and we can't use symlinks, so this file is a copy of the font renamed

//go:embed 3270Medium.wasm
var go3270Font []byte

type Go3270 struct {
	canvas       js.Value
	canvasHeight float64
	canvasWidth  float64
	color        string
	cols         float64
	copybuff     js.Value
	ctx          js.Value
	dc           *gg.Context
	dpi          float64
	face         font.Face
	font         *truetype.Font
	fontHeight   float64
	fontSize     float64
	fontWidth    float64
	image        *image.RGBA
	imgData      js.Value
	paddedHeight float64
	paddedWidth  float64
	rows         float64
	scaleFactor  float64
}

func NewGo3270(this js.Value, args []js.Value) any {
	c := &Go3270{}
	// 👇 properties
	c.canvas = args[0]
	c.color = args[1].String()
	c.fontSize = args[2].Float()
	c.cols = args[3].Float()
	c.rows = args[4].Float()
	c.dpi = args[5].Float()
	// 👇 constants
	c.paddedHeight = 1.1
	c.paddedWidth = 1.1
	c.scaleFactor = 2
	// 👇 load the 3270 font
	c.font, _ = truetype.Parse(go3270Font)
	c.face = truetype.NewFace(c.font, &truetype.Options{Size: float64(c.fontSize * c.scaleFactor), DPI: c.dpi, Hinting: font.HintingFull})
	// 👇 resize canvas to fit font, using temporary context
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	dc := gg.NewContextForRGBA(img)
	dc.SetFontFace(c.face)
	c.fontWidth, c.fontHeight = dc.MeasureString("M")
	c.canvasWidth = c.cols * c.fontWidth * c.paddedWidth
	c.canvasHeight = c.rows * c.fontHeight * c.paddedHeight
	wrapper := c.canvas.Get("parentNode")
	wrapper.Get("style").Set("width", fmt.Sprintf("%fpx", c.canvasWidth/c.scaleFactor))
	wrapper.Get("style").Set("height", fmt.Sprintf("%fpx", c.canvasHeight/c.scaleFactor))
	c.canvas.Set("width", c.canvasWidth)
	c.canvas.Set("height", c.canvasHeight)
	c.canvas.Get("style").Set("scale", 1/c.scaleFactor)
	// 👇 derivatives
	c.ctx = c.canvas.Call("getContext", "2d")
	c.imgData = c.ctx.Call("createImageData", c.canvasWidth, c.canvasHeight)
	c.image = image.NewRGBA(image.Rect(0, 0, int(c.canvasWidth), int(c.canvasHeight)))
	c.copybuff = js.Global().Get("Uint8Array").New(len(c.image.Pix))
	c.dc = gg.NewContextForRGBA(c.image)
	c.dc.SetFontFace(c.face)
	// 👇 methods callable by Javascript
	obj := map[string]any{
		"inbound": js.FuncOf(func(this js.Value, args []js.Value) any {
			return c.Inbound()
		}),
		"testPattern": js.FuncOf(func(this js.Value, args []js.Value) any {
			return c.TestPattern()
		}),
	}
	return js.ValueOf(obj)
}

func (c *Go3270) Inbound() any {
	// 👇 simulate response
	data := []byte{193, 194, 195 /* 👈 EBCDIC "ABC" */}
	uint8ArrayConstructor := js.Global().Get("Uint8Array")
	result := uint8ArrayConstructor.New(len(data))
	js.CopyBytesToJS(result, data)
	return result
}

func (c *Go3270) TestPattern() any {
	c.dc.SetRGBA(0, 0, 0, 0)
	c.dc.Clear()
	c.dc.SetHexColor(c.color)
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/?"
	chars := []rune(str)
	for x := 0.0; x < c.cols; x++ {
		for y := 0.0; y < c.rows; y++ {
			ix := rand.Intn(len(chars))
			c.dc.DrawString(string(chars[ix]), math.Round(x*c.fontWidth*c.paddedWidth), math.Round((y+1)*c.fontHeight*c.paddedHeight))
		}
	}
	c.imgCopy()
	return nil
}

func (c *Go3270) imgCopy() {
	js.CopyBytesToJS(c.copybuff, c.image.Pix)
	c.imgData.Get("data").Call("set", c.copybuff)
	c.ctx.Call("putImageData", c.imgData, 0, 0)
}
