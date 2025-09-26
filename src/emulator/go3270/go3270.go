package go3270

import (
	_ "embed"
	"emulator/types"
	"fmt"
	"image"
	"math"
	"math/rand"
	"syscall/js"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// 🔥 Hack alert! we must use extension {js, wasm} and we can't use symlinks, so this file is a copy of the font renamed

//go:embed 3270Font.wasm
var go3270Font []byte

// 🟧 Bridge between Typescript UI and Go-powered emulator

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

type Go3270 struct {
	canvas       js.Value
	canvasHeight float64
	canvasWidth  float64
	color        string
	cols         float64
	copybuff     js.Value
	ctx          js.Value
	gg           *gg.Context
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

// 🔥 main.go places this function name on the DOM's global window object
func NewGo3270(this js.Value, args []js.Value) any {
	go3270 := &Go3270{}
	// 👇 properties
	go3270.canvas = args[0]
	go3270.color = args[1].String()
	go3270.fontSize = args[2].Float()
	go3270.cols = args[3].Float()
	go3270.rows = args[4].Float()
	go3270.dpi = args[5].Float()
	// 👇 constants
	go3270.paddedHeight = 1.05
	go3270.paddedWidth = 1.1
	go3270.scaleFactor = 2
	// 👇 load the 3270 font
	go3270.font, _ = opentype.Parse(go3270Font)
	go3270.face, _ = opentype.NewFace(go3270.font, &opentype.FaceOptions{Size: go3270.fontSize * go3270.scaleFactor, DPI: go3270.dpi, Hinting: font.HintingFull})
	// 👇 resize canvas to fit font, using temporary context
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	dc := gg.NewContextForRGBA(img)
	dc.SetFontFace(go3270.face)
	go3270.fontWidth, go3270.fontHeight = dc.MeasureString("M")
	go3270.canvasWidth = go3270.cols * go3270.fontWidth * go3270.paddedWidth
	go3270.canvasHeight = go3270.rows * go3270.fontHeight * go3270.paddedHeight
	wrapper := go3270.canvas.Get("parentNode")
	wrapper.Get("style").Set("width", fmt.Sprintf("%fpx", go3270.canvasWidth/go3270.scaleFactor))
	wrapper.Get("style").Set("height", fmt.Sprintf("%fpx", go3270.canvasHeight/go3270.scaleFactor))
	go3270.canvas.Set("width", go3270.canvasWidth)
	go3270.canvas.Set("height", go3270.canvasHeight)
	// 👇 derivatives
	go3270.ctx = go3270.canvas.Call("getContext", "2d")
	go3270.imgData = go3270.ctx.Call("createImageData", go3270.canvasWidth, go3270.canvasHeight)
	go3270.image = image.NewRGBA(image.Rect(0, 0, int(go3270.canvasWidth), int(go3270.canvasHeight)))
	go3270.copybuff = js.Global().Get("Uint8Array").New(len(go3270.image.Pix))
	go3270.gg = gg.NewContextForRGBA(go3270.image)
	go3270.gg.SetFontFace(go3270.face)
	go3270.gg.Scale(1/go3270.scaleFactor, 1/go3270.scaleFactor)
	// 👇 methods callable by Javascript
	// 👁️ go3270.d.ts
	obj := map[string]any{
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			return go3270.Close()
		}),
		"datastream": js.FuncOf(func(this js.Value, args []js.Value) any {
			return go3270.Datastream(args[0])
		}),
		"restore": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.Restore(args[0])
			return nil
		}),
		"testPattern": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.TestPattern()
			return nil
		}),
	}
	return js.ValueOf(obj)
}

func (go3270 *Go3270) Close() js.Value {
	// 🔥 simulate the state of the device
	data := []byte{193, 194, 195 /* 👈 EBCDIC "ABC" */}
	u8 := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(u8, data)
	return u8
}

func (go3270 *Go3270) Coords(col, row float64) (float64, float64, float64, float64, float64) {
	w := go3270.fontWidth * go3270.paddedWidth
	h := go3270.fontHeight * go3270.paddedHeight
	x := col * w
	y := row * h
	// 🔥 we could do better calculating the baseline - this is just a WAG, because an em is drawn with a significantly different height than that returned by MeasureString()
	baseline := y + h - (go3270.fontSize / 3 * go3270.scaleFactor)
	return x, y, w, h, baseline
}

func (go3270 *Go3270) Datastream(u8in js.Value) js.Value {
	bytes := make([]byte, u8in.Get("length").Int())
	js.CopyBytesToGo(bytes, u8in)
	// 🔥 do something with stream
	_ = bytes
	go3270.TestPattern()
	// 🔥 simulate response
	data := []byte{193, 194, 195 /* 👈 EBCDIC "ABC" */}
	u8out := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(u8out, data)
	return u8out
}

func (go3270 *Go3270) DispatchEvent(eventType string, data any) {
	event := js.Global().Get("CustomEvent").New(eventType, map[string]any{
		"detail": data,
	})
	js.Global().Get("document").Call("dispatchEvent", event)
}

func (go3270 *Go3270) Restore(u8 js.Value) {
	// 🔥 simulate restoration of state of device
	_ = u8
	go3270.TestPattern()
}

func (go3270 *Go3270) TestPattern() {
	go3270.gg.SetHexColor(types.CLUT[0xf0][0]) /* 👈 ragged fonts if draw on transparent! */
	go3270.gg.Clear()
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/?"
	chs := []rune(str)
	for col := 0.0; col < go3270.cols; col++ {
		for row := 0.0; row < go3270.rows; row++ {
			x, _, _, _, baseline := go3270.Coords(col, row)
			// 👇 choose colors from the CLUT, using the base color if out of range
			ix := uint8(math.Floor(col/10) + 0xf1)
			bright := go3270.color
			color := go3270.color
			if ix <= 0xf7 {
				bright = types.CLUT[ix][0]
				color = types.CLUT[ix][1]
			}
			// 👇 alternate high intensity, normal
			if int(row)%2 == 0 {
				go3270.gg.SetHexColor(bright)
			} else {
				go3270.gg.SetHexColor(color)
			}
			ich := rand.Intn(len(chs))
			ch := string(chs[ich])
			go3270.gg.DrawString(ch, x, baseline)
		}
	}
	go3270.imgCopy()
}

func (go3270 *Go3270) imgCopy() {
	js.CopyBytesToJS(go3270.copybuff, go3270.image.Pix)
	go3270.imgData.Get("data").Call("set", go3270.copybuff)
	go3270.ctx.Call("putImageData", go3270.imgData, 0, 0)
}
