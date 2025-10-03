package go3270

import (
	_ "embed"
	"emulator/device"
	"emulator/utils"
	"fmt"
	"image"
	"slices"
	"syscall/js"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"golang.org/x/image/font/opentype"
)

// 🔥 Hack alert! we must use extension {js, wasm} and we can't use symlinks, so this file is a copy of the font renamed

//go:embed 3270Font.wasm
var go3270Font []uint8

// 🟧 Bridge between Typescript UI and Go-powered emulator

// The design objective is that all Go <-> UI communication goes through this package. No other package must use syscall/js. That way, everything but here can be tested with go test.

// The device package is handed a drawing context into which it renders the 3270 stream and any operator input. Using requestAnimationFrame, this module actually draws the context onto a supplied HTML canvas whenever the context changes

type Go3270 struct {
	bus    EventBus.Bus
	device *device.Device

	// 👇 manage frame rendering
	lastImage     []uint8
	lastTimestamp float64
	renderContext js.Func
	reqID         js.Value
}

// 🔥 main.go places this function name on the DOM's global window object
func NewGo3270(this js.Value, args []js.Value) any {
	go3270 := new(Go3270)
	// 👇 get the bus ready right away
	go3270.bus = EventBus.New()
	// 👇 properties
	canvas := args[0]
	bgColor := args[1].String()
	color := args[2].String()
	fontSize := args[3].Float()
	cols := args[4].Float()
	rows := args[5].Float()
	dpi := args[6].Float()
	// 👇 constants
	maxFPS := 30.0
	paddedHeight := 1.05
	paddedWidth := 1.1
	// 🔥 scaling 2x does produce slightly crisper font rendering, but it takes about 2x as long to render
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
		gg,
		bgColor,
		color,
		cols,
		rows,
		fontHeight,
		fontSize,
		fontWidth,
		paddedHeight,
		paddedWidth,
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
	go3270.bus.Subscribe("go3270", go3270Message)
	// 👇 start the requestAnimationFrame to render "gg" context changes
	go3270.startRenderContextLoop(canvas, rgba, maxFPS)
	// 👇 finally, this is thhe interface through which TypeScript will call us
	return js.ValueOf(tsInterface)
}

// 🟦 Render drawing context when changed via requestAnimationFrame

func (go3270 *Go3270) startRenderContextLoop(canvas js.Value, rgba *image.RGBA, maxFPS float64) {
	go3270.renderContext = js.FuncOf(func(this js.Value, args []js.Value) any {
		timestamp := args[0].Float()
		// 👇 make sure we don't bust the max FPS we were given
		if timestamp-go3270.lastTimestamp >= (1000 / maxFPS) {
			if go3270.lastImage == nil || !slices.Equal(go3270.lastImage, rgba.Pix) {
				timeNow := time.Now()
				// 🔥 I copied this from go-canvas where the author was worried about 3 separate copies -- I haven't figured how to reduce it to 2 even when using Uint8ClampedArray -- but it only takes ~1ms anyway
				u8 := js.Global().Get("Uint8ClampedArray").New(len(rgba.Pix))
				js.CopyBytesToJS(u8, rgba.Pix)
				canvasHeight := canvas.Get("offsetHeight")
				canvasWidth := canvas.Get("offsetWidth")
				ctx := canvas.Call("getContext", "2d")
				pixels := ctx.Call("createImageData", canvasWidth, canvasHeight)
				pixels.Get("data").Call("set", u8)
				ctx.Call("putImageData", pixels, 0, 0)
				// 👇 set up for next time
				go3270.lastImage = make([]uint8, len(rgba.Pix))
				copy(go3270.lastImage, rgba.Pix)
				go3270.lastTimestamp = timestamp
				utils.ElapsedTime(timeNow, "requestAnimationFrame")
			}
		}
		go3270.reqID = js.Global().Call("requestAnimationFrame", go3270.renderContext)
		return nil
	})
	// 👇 kick off the rendering loop
	js.Global().Call("requestAnimationFrame", go3270.renderContext)
}

// 🟦 Go WASM methods callable by Javascript via go3270.ts

func (go3270 *Go3270) Close() {
	js.Global().Get("console").Call("log", "%cGo3270 closing", "color: orange")
	// 👇 perform any cleanup
	js.Global().Call("cancelAnimationFrame", go3270.reqID)
	go3270.device.Close()
	// 🟦 Go WASM functions invoked by go test-able code
	go3270.bus.Unsubscribe("go3270", go3270Message)
}

func (go3270 *Go3270) Keystroke(code string, key string, alt bool, ctrl bool, shift bool) {
	// 🔥 simulate handling of Keystroke
	js.Global().Get("console").Call("log", fmt.Sprintf("%%ccode=%s %%ckey=%s %%calt=%t ctrl=%t shift=%t", code, key, alt, ctrl, shift), "color: coral", "color: skyblue", "color: gray")
}

func (go3270 *Go3270) ReceiveFromApp(u8in js.Value) {
	bytes := make([]uint8, u8in.Get("length").Int())
	js.CopyBytesToGo(bytes, u8in)
	go3270.device.ReceiveFromApp(bytes)
}

// 🟦 Messages from go test-able code sent to the UI for action

func go3270Message(eventType string, bytes []uint8, params map[string]any, args []any) {
	// 👇 params and args may be nil
	if params == nil {
		params = map[string]any{}
	}
	if args != nil && args[0] != nil {
		params["args"] = args
	}
	// 👇 bytes may be nil, but if not convert to JS
	var u8 js.Value
	if bytes != nil {
		u8 = js.Global().Get("Uint8ClampedArray").New(len(bytes))
		js.CopyBytesToJS(u8, bytes)
		params["bytes"] = u8
	}
	// 👇 dispatch event to JS
	params["eventType"] = eventType
	event := js.Global().Get("CustomEvent").New("go3270", map[string]any{
		"detail": params,
	})
	js.Global().Get("window").Call("dispatchEvent", event)
}
