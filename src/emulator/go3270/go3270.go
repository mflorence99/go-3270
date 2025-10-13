package go3270

import (
	_ "embed"
	"emulator/device"
	"fmt"
	"image"
	"math"
	"slices"
	"syscall/js"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

// 🔥 Hack alert! we must use extension {js, wasm} and we can't use symlinks, so this file is a copy of the font renamed

//go:embed IBMPlexMono-Regular.ttf.wasm
var regularFontEmbed []byte

//go:embed IBMPlexMono-Bold.ttf.wasm
var boldFontEmbed []byte

// 🟧 Bridge between Typescript UI and Go-powered emulator

// The design objective is that all Go <-> UI communication goes through this package. No other package must use syscall/js. That way, everything but here can be tested with go test.

// The device package is handed a drawing context into which it renders the 3270 stream and any operator input. Using requestAnimationFrame, this module actually draws the context onto a supplied HTML canvas whenever the context changes

type Go3270 struct {
	bus    EventBus.Bus
	device *device.Device
}

// 🔥 main.go places this function name on the DOM's global window object
func New(this js.Value, args []js.Value) any {
	go3270 := new(Go3270)
	// 👇 get the bus ready right away
	go3270.bus = EventBus.New()
	// 👇 properties
	canvas := args[0]
	bgColor := args[1].String()
	color := [2]string{args[2].Index(0).String(), args[2].Index(1).String()}
	// clut := map[byte][2]string{}
	obj := args[3]
	keys := js.Global().Get("Object").Call("keys", obj)
	for i := 0; i < keys.Length(); i++ {
		k := keys.Index(i).String()
		v := [2]string{obj.Get(k).Index(0).String(), obj.Get(k).Index(1).String()}
		fmt.Printf("🎨 %s = %v\n", k, v)
	}
	fontSize := args[4].Float()
	cols := args[5].Int()
	rows := args[6].Int()
	dpi := args[7].Float()
	// 👇 constants
	maxFPS := 30.0
	paddedHeight := 1.5
	paddedWidth := 1.1
	// 👇 load the 3270 font
	regularFont, _ := truetype.Parse(regularFontEmbed)
	regularFace := truetype.NewFace(regularFont, &truetype.Options{Size: fontSize, DPI: dpi /* , Hinting: font.HintingFull */})
	boldFont, _ := truetype.Parse(boldFontEmbed)
	boldFace := truetype.NewFace(boldFont, &truetype.Options{Size: fontSize, DPI: dpi /* , Hinting: font.HintingFull */})
	temp := gg.NewContext(100, 100)
	temp.SetFontFace(boldFace)
	fontWidth, fontHeight := temp.MeasureString("M")
	canvasWidth := float64(cols) * math.Round(fontWidth*paddedWidth)
	canvasHeight := float64(rows) * math.Round(fontHeight*paddedHeight)
	wrapper := canvas.Get("parentNode")
	wrapper.Get("style").Set("width", fmt.Sprintf("%fpx", canvasWidth))
	wrapper.Get("style").Set("height", fmt.Sprintf("%fpx", canvasHeight))
	canvas.Set("width", canvasWidth)
	canvas.Set("height", canvasHeight)
	// 👇 prepare the rendering surface
	rgba := image.NewRGBA(image.Rect(0, 0, int(canvasWidth), int(canvasHeight)))
	// 👇 delegate all device handling to go test-able handler
	go3270.device = device.NewDevice(
		go3270.bus,
		rgba,
		boldFace,
		regularFace,
		bgColor,
		color,
		cols,
		rows,
		fontHeight,
		fontSize,
		fontWidth,
		paddedHeight,
		paddedWidth)
	// 🟦 Go WASM methods callable by Javascript
	// 👁️ go3270.d.ts
	tsInterface := map[string]any{
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.Close()
			return nil
		}),
		"focus": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.Focus(args[0].Bool())
			return nil
		}),
		"keystroke": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.Keystroke(args[0].String(), args[1].String(), args[2].Bool(), args[3].Bool(), args[4].Bool())
			return nil
		}),
		"outbound": js.FuncOf(func(this js.Value, args []js.Value) any {
			go3270.Outbound(args[0])
			return nil
		}),
	}
	// 👇 Go WASM functions invoked by go test-able code
	go3270.bus.Subscribe("go3270", go3270Message)
	// 👇 start the requestAnimationFrame to render "gg" context changes
	go3270.startRenderContextLoop(canvas, rgba, maxFPS)
	// 👇 finally, this is thhe interface through which TypeScript will call us
	return js.ValueOf(tsInterface)
}

// 🟦 Render drawing context when changed via requestAnimationFrame

func (go3270 *Go3270) startRenderContextLoop(canvas js.Value, rgba *image.RGBA, maxFPS float64) {
	var (
		lastImage     []byte
		lastTimestamp float64
		renderContext js.Func
	)
	renderContext = js.FuncOf(func(this js.Value, args []js.Value) any {
		timestamp := args[0].Float()
		// 👇 make sure we don't bust the max FPS we were given
		if timestamp-lastTimestamp >= (1000 / maxFPS) {
			if lastImage == nil || !slices.Equal(lastImage, rgba.Pix) {
				// 🔥 I copied this from go-canvas where the author was worried about 3 separate copies -- I haven't figured how to reduce it to 2 even when using Uint8ClampedArray -- but it only takes ~1ms anyway
				pixels := js.Global().Get("Uint8ClampedArray").New(len(rgba.Pix))
				js.CopyBytesToJS(pixels, rgba.Pix)
				canvasHeight := canvas.Get("offsetHeight")
				canvasWidth := canvas.Get("offsetWidth")
				ctx := canvas.Call("getContext", "2d")
				img := ctx.Call("createImageData", canvasWidth, canvasHeight)
				img.Get("data").Call("set", pixels)
				ctx.Call("putImageData", img, 0, 0)
				// 👇 set up for next time
				lastImage = make([]byte, len(rgba.Pix))
				copy(lastImage, rgba.Pix)
				lastTimestamp = timestamp
			}
		}
		js.Global().Call("requestAnimationFrame", renderContext)
		return nil
	})
	// 👇 kick off the rendering loop
	js.Global().Call("requestAnimationFrame", renderContext)
}

// 🟦 Go WASM methods callable by Javascript via go3270.ts

func (go3270 *Go3270) Close() {
	js.Global().Get("console").Call("log", "🐞 Go3270 closing")
	// 👇 perform any cleanup
	// js.Global().Call("cancelAnimationFrame", go3270.reqID)
	go3270.device.Close()
	// 🟦 Go WASM functions invoked by go test-able code
	go3270.bus.Unsubscribe("go3270", go3270Message)
}

func (go3270 *Go3270) Focus(focus bool) {
	js.Global().Get("console").Call("log", device.Ternary(focus, "⌨️ Go3270 gains focus", "⌨️ Go3270 loses focus"))
	// 👇 just forward to device
	go3270.device.Focus(focus)
}

func (go3270 *Go3270) Keystroke(code string, key string, alt bool, ctrl bool, shift bool) {
	str := "⌨️ "
	if ctrl {
		str += "CTRL+"
	}
	if shift {
		str += "SHIFT+"
	}
	if alt {
		str += "ALT+"
	}
	js.Global().Get("console").Call("log", fmt.Sprintf("%s%s %s", str, key, device.Ternary(code != key && len(key) > 1, code, "")))
	// 👇 just forward to device
	go3270.device.Keystroke(code, key, alt, ctrl, shift)
}

func (go3270 *Go3270) Outbound(u8in js.Value) {
	// 🔥 do this the hard way, just so we don't have to re-convert bytes
	params := map[string]any{
		"bytes":     u8in,
		"color":     "yellow",
		"ebcdic":    true,
		"eventType": "dumpBytes",
		"title":     "Outbound",
	}
	event := js.Global().Get("CustomEvent").New("go3270", map[string]any{
		"detail": params,
	})
	js.Global().Get("window").Call("dispatchEvent", event)
	// 👇 just forward to device
	u8s := make([]byte, u8in.Get("length").Int())
	js.CopyBytesToGo(u8s, u8in)
	go3270.device.Outbound(u8s)
}

// 🟦 Messages from go test-able code sent to the UI for action

func go3270Message(msg device.Go3270Message) {
	// 👇 params and args may be nil
	if msg.Params == nil {
		msg.Params = map[string]any{}
	}
	if msg.Args != nil && msg.Args[0] != nil {
		msg.Params["args"] = msg.Args
	}
	// 👇 bytes may be nil, but if not convert to JS
	var u8out js.Value
	if msg.U8s != nil {
		u8out = js.Global().Get("Uint8ClampedArray").New(len(msg.U8s))
		js.CopyBytesToJS(u8out, msg.U8s)
		msg.Params["bytes"] = u8out
	}
	// 👇 dispatch event to JS
	msg.Params["eventType"] = msg.EventType
	event := js.Global().Get("CustomEvent").New("go3270", map[string]any{
		"detail": msg.Params,
	})
	js.Global().Get("window").Call("dispatchEvent", event)
	// 👇 special case: dump what we send
	if msg.EventType == "sendToApp" {
		params := map[string]any{
			"bytes":     u8out,
			"color":     "palegreen",
			"ebcdic":    true,
			"eventType": "dumpBytes",
			"title":     "SendToApp",
		}
		event := js.Global().Get("CustomEvent").New("go3270", map[string]any{
			"detail": params,
		})
		js.Global().Get("window").Call("dispatchEvent", event)
	}
}
