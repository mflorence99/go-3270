package go3270

import (
	_ "embed"
	"emulator/device"
	"fmt"
	"image"
	"syscall/js"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"golang.org/x/image/font/opentype"
)

// 🔥 Hack alert! we must use extension {js, wasm} and we can't use symlinks, so this file is a copy of the font renamed

//go:embed 3270Font.wasm
var go3270Font []uint8

// 🟧 Bridge between Typescript UI and Go-powered emulator

// The design objective is that all Go <-> UI communication goes through this module. No other modulw must use syscall/js. That way, everything but here can be tested with go test.

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

type Go3270 struct {
	bus      EventBus.Bus
	device   *device.Device
	renderer func()
}

// 🔥 main.go places this function name on the DOM's global window object
func NewGo3270(this js.Value, args []js.Value) any {
	go3270 := &Go3270{}
	// 👇 get the bus ready right away
	go3270.bus = EventBus.New()
	// 👇 properties
	canvas := args[0]
	color := args[1].String()
	fontSize := args[2].Float()
	cols := args[3].Float()
	rows := args[4].Float()
	dpi := args[5].Float()
	// 👇 constants
	paddedHeight := 1.05
	paddedWidth := 1.1
	// 🔥 scaling 2x does produce slightly crisper font rendering, but it takes about 2x as long to render (see function TestPattern)
	scaleFactor := 1.0
	// 👇 load the 3270 font
	font, _ := opentype.Parse(go3270Font)
	face, _ := opentype.NewFace(font, &opentype.FaceOptions{Size: fontSize * scaleFactor, DPI: dpi /* , Hinting: font.HintingFull */})
	// 👇 resize canvas to fit font, using temporary context
	rgba := image.NewRGBA(image.Rect(0, 0, 100, 100))
	dc := gg.NewContextForRGBA(rgba)
	dc.SetFontFace(face)
	fontWidth, fontHeight := dc.MeasureString("M")
	canvasWidth := cols * fontWidth * paddedWidth
	canvasHeight := rows * fontHeight * paddedHeight
	wrapper := canvas.Get("parentNode")
	wrapper.Get("style").Set("width", fmt.Sprintf("%fpx", canvasWidth/scaleFactor))
	wrapper.Get("style").Set("height", fmt.Sprintf("%fpx", canvasHeight/scaleFactor))
	canvas.Set("width", canvasWidth)
	canvas.Set("height", canvasHeight)
	// 👇 prepare the rendering surface
	rgba = image.NewRGBA(image.Rect(0, 0, int(canvasWidth), int(canvasHeight)))
	gg := gg.NewContextForRGBA(rgba)
	gg.SetFontFace(face)
	gg.Scale(1/scaleFactor, 1/scaleFactor)
	// 👇 delegate all device handling to go test-able handler
	go3270.device = device.NewDevice(
		go3270.bus,
		color,
		cols,
		gg,
		fontHeight,
		fontSize,
		fontWidth,
		paddedHeight,
		paddedWidth,
		rows,
		scaleFactor)
	// 🟦 Go WASM methods callable by Javascript
	// 👁️ go3270.d.ts
	tsInterface := map[string]any{
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.Close()
			return nil
		}),
		"keystroke": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.Keystroke(args[0].String(), args[1].String(), args[2].Bool(), args[3].Bool(), args[4].Bool())
			return nil
		}),
		"receiveFromApp": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.ReceiveFromApp(args[0])
			return nil
		}),
	}
	// 🟦 Go WASM functions invoked by go test-able code
	go3270.bus.Subscribe("go3270-alarm", alarm)
	go3270.bus.Subscribe("go3270-dispatchEvent", dispatchEvent)
	go3270.bus.Subscribe("go3270-dumpBytes", dumpBytes)
	go3270.bus.Subscribe("go3270-log", log)
	go3270.renderer = render(canvas, rgba)
	go3270.bus.Subscribe("go3270-render", go3270.renderer)
	go3270.bus.Subscribe("go3270-sendToApp", sendToApp)
	// 👇 finally, this is thhe interface through which TypeScript will call us
	return js.ValueOf(tsInterface)
}

// 🟦 Go WASM methods callable by Javascript via window.xxx

func (go3270 *Go3270) Close() {
	log("%cGo3270 closing", "color: orange")
	// 👇 perform any cleanup
	go3270.device.Close()
	// 🟦 Go WASM functions invoked by go test-able code
	go3270.bus.Unsubscribe("go3270-alarm", alarm)
	go3270.bus.Unsubscribe("go3270-dispatchEvent", dispatchEvent)
	go3270.bus.Unsubscribe("go3270-dumpBytes", dumpBytes)
	go3270.bus.Unsubscribe("go3270-log", log)
	go3270.bus.Unsubscribe("go3270-render", go3270.renderer)
	go3270.bus.Unsubscribe("go3270-sendToApp", sendToApp)
}

func (go3270 *Go3270) Keystroke(code string, key string, alt bool, ctrl bool, shift bool) {
	// 🔥 simulate handling of Keystroke
	log(fmt.Sprintf("%%ccode=%s %%ckey=%s %%calt=%t ctrl=%t shift=%t", code, key, alt, ctrl, shift), "color: coral", "color: skyblue", "color: gray")
}

func (go3270 *Go3270) ReceiveFromApp(u8in js.Value) {
	request := make([]uint8, u8in.Get("length").Int())
	js.CopyBytesToGo(request, u8in)
	// 🔥 do something with stream
	_ = request
	go3270.device.TestPattern()
	// 🔥 simulate response
	sendToApp([]uint8{193, 194, 195 /* 👈 EBCDIC "ABC" */})
}

// 🟦 Go WASM functions invoked by go test-able code via EventBus

func alarm() {
	dispatchEvent("go3270-alarm", true)
}

func dispatchEvent(eventType string, data any) {
	event := js.Global().Get("CustomEvent").New(eventType, map[string]any{
		"detail": data,
	})
	js.Global().Get("document").Call("dispatchEvent", event)
}

func dumpBytes(data []uint8, title string, ebcdic bool, color string) {
	u8 := js.Global().Get("Uint8ClampedArray").New(len(data))
	js.CopyBytesToJS(u8, data)
	dispatchEvent("go3270-dumpBytes", map[string]any{
		"bytes":  u8,
		"title":  title,
		"ebcdic": ebcdic,
		"color":  color,
	})
}

func log(args ...any) {
	dispatchEvent("go3270-log", map[string]any{
		"args": args,
	})
}

// 🔥 I copied this from go-canvas and the author was worried about 3 separate copies -- I haven't figured how to reduce it to 2 even when using Uint8ClampedArray -- but it only takes ~1ms anyway
func render(canvas js.Value, rgba *image.RGBA) func() {
	return func() {
		u8 := js.Global().Get("Uint8ClampedArray").New(len(rgba.Pix))
		js.CopyBytesToJS(u8, rgba.Pix)
		canvasHeight := canvas.Get("offsetHeight")
		canvasWidth := canvas.Get("offsetWidth")
		ctx := canvas.Call("getContext", "2d")
		pixels := ctx.Call("createImageData", canvasWidth, canvasHeight)
		pixels.Get("data").Call("set", u8)
		ctx.Call("putImageData", pixels, 0, 0)
	}
}

func sendToApp(data []uint8) {
	u8 := js.Global().Get("Uint8ClampedArray").New(len(data))
	js.CopyBytesToJS(u8, data)
	dispatchEvent("go3270-sendToApp", map[string]any{
		"bytes": u8,
	})
}
