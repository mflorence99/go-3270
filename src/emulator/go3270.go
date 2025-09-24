package main

// 🟧 3270 data stream protocol

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

import (
	_ "embed"
	"fmt"
	"image"
	"math/rand"
	"syscall/js"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// 🔥 Hack alert! we must use extension {js, wasm} and we can't use symlinks, so this file is a copy of the font renamed

//go:embed 3270Font.wasm
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
	font         *opentype.Font
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
	c.paddedHeight = 1.05
	c.paddedWidth = 1.1
	c.scaleFactor = 2
	// 👇 load the 3270 font
	c.font, _ = opentype.Parse(go3270Font)
	c.face, _ = opentype.NewFace(c.font, &opentype.FaceOptions{Size: c.fontSize * c.scaleFactor, DPI: c.dpi, Hinting: font.HintingFull})
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
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			return c.Close()
		}),
		"datastream": js.FuncOf(func(this js.Value, args []js.Value) any {
			return c.Datastream(args[0])
		}),
		"restore": js.FuncOf(func(this js.Value, args []js.Value) any {
			c.Restore(args[0])
			return nil
		}),
		"testPattern": js.FuncOf(func(this js.Value, args []js.Value) any {
			c.TestPattern()
			return nil
		}),
	}
	return js.ValueOf(obj)
}

func (c *Go3270) Close() js.Value {
	// 🔥 simulate the state of the device
	data := []byte{193, 194, 195 /* 👈 EBCDIC "ABC" */}
	uint8ArrayConstructor := js.Global().Get("Uint8Array")
	result := uint8ArrayConstructor.New(len(data))
	js.CopyBytesToJS(result, data)
	return result
}

func (c *Go3270) Coords(col, row float64) (float64, float64, float64, float64, float64) {
	w := c.fontWidth * c.paddedWidth
	h := c.fontHeight * c.paddedHeight
	x := col * w
	y := row * h
	// 🔥 we could do better calculating the baseline - this is just a WAG, because an em is drawn with a significantly different height than that returned by MeasureString()
	baseline := y + h - (c.fontSize / 3 * c.scaleFactor)
	return x, y, w, h, baseline
}

func (c *Go3270) Datastream(bytes js.Value) js.Value {
	// 🔥 do something with stream
	_ = bytes
	c.TestPattern()
	// 🔥 simulate response
	data := []byte{193, 194, 195 /* 👈 EBCDIC "ABC" */}
	uint8ArrayConstructor := js.Global().Get("Uint8Array")
	result := uint8ArrayConstructor.New(len(data))
	js.CopyBytesToJS(result, data)
	return result
}

func (c *Go3270) DispatchEvent(eventType string, data any) {
	event := js.Global().Get("CustomEvent").New(eventType, map[string]any{
		"detail": data,
	})
	js.Global().Get("document").Call("dispatchEvent", event)
}

func (c *Go3270) Restore(bytes js.Value) {
	// 🔥 simulate restoration of state of device
	_ = bytes
	c.TestPattern()
	c.DispatchEvent("go3270-alarm", true)
}

func (c *Go3270) TestPattern() {
	c.dc.SetRGBA(0, 0, 0, 0)
	c.dc.Clear()
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/?"
	chs := []rune(str)
	for col := 0.0; col < c.cols; col++ {
		for row := 0.0; row < c.rows; row++ {
			x, y, w, h, baseline := c.Coords(col, row)
			if int(col)%2 == 0 && int(row)%2 != 0 {
				c.dc.SetHexColor(c.color)
				c.dc.DrawRectangle(x, y, w, h)
				c.dc.Fill()
				c.dc.SetHexColor("#000000")
			} else if col < 10 {
				c.dc.SetHexColor("#42a5f5")
			} else if col < 20 {
				c.dc.SetHexColor("#f44336")
			} else if col < 30 {
				c.dc.SetHexColor("#66be6a")
			} else if col < 40 {
				c.dc.SetHexColor("#ec407a")
			} else if col < 50 {
				c.dc.SetHexColor("#26c6da")
			} else if col < 60 {
				c.dc.SetHexColor("#ffee58")
			} else if col < 70 {
				c.dc.SetHexColor("#f5f5f5")
			} else if col < 80 {
				c.dc.SetHexColor(c.color)
			}
			ich := rand.Intn(len(chs))
			ch := string(chs[ich])
			c.dc.DrawString(ch, x, baseline)
		}
	}
	c.imgCopy()
}

func (c *Go3270) imgCopy() {
	js.CopyBytesToJS(c.copybuff, c.image.Pix)
	c.imgData.Get("data").Call("set", c.copybuff)
	c.ctx.Call("putImageData", c.imgData, 0, 0)
}
