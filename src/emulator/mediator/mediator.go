package mediator

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"math"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"

	"emulator/core"
	"emulator/samples"
	"emulator/types"
	"syscall/js"
)

// 🔥 Hack alert! we must use extension {js, wasm}
//    and we can't use symlinks, so this file is a copy of the font renamed

var (
	//go:embed JuliaMono-Regular.ttf.wasm
	normalFontEmbed []byte
	//go:embed JuliaMono-Bold.ttf.wasm
	boldFontEmbed []byte
)

// 🟧 Bridge between Typescript UI and Go-powered emulator

// 🟦 The design objective is that all Go <-> UI communication
//    goes through this package. No other package must use
//    syscall/js. That way, everything but here can be tested with go test.

// 🟦 The emulator package is handed an image into which it renders
//    the 3270 stream and any operator input. Using requestAnimationFrame,
//    this module actually draws the context onto a supplied HTML
//    canvas whenever the context changes

type Mediator struct {
	bus *core.Bus
	emu *core.Emulator
}

// 👁️ go3270.ts

// args[0] canvas
// args[1] bgColor
// args[2] monochrome
// args[3] clut [map attr -> [color, name]]
// args[4] fontSize
// args[5] rows
// args[6] cols
// args[7] dpi
// 👇 for testing
// args[8] testpage

func NewGo3270(this js.Value, args []js.Value) any {
	m := new(Mediator)
	m.bus = core.NewBus()
	// 🔥 must subscribe BEFORE we create the emulator
	m.bus.SubInbound(m.inbound)
	m.bus.SubPanic(m.panic)
	m.bus.SubStatus(m.status)
	// 👇 create and configure the emulator and its children
	cfg := m.configure(args)
	m.emu = core.NewEmulator(m.bus, cfg)
	m.emu.Initialize()
	// 👇 if debugging, show screenshot
	if cfg.Testpage != "" {
		m.bus.PubOutbound(samples.Index[cfg.Testpage])
	}
	return m.jsInterface()
}

func (m *Mediator) close() {
	m.bus.PubClose()
	m.bus.UnsubscribeAll()
}

func (m *Mediator) configure(args []js.Value) *types.Config {
	// 👇 from the args
	canvas := args[0]
	bgColor := args[1].String()
	monochrome := args[2].Bool()
	obj := args[3]
	clut := make(map[types.Color]string)
	keys := js.Global().Get("Object").Call("keys", obj)
	for ix := 0; ix < keys.Length(); ix++ {
		k := keys.Index(ix).String()
		v := obj.Get(k).Index(0).String()
		c, _ := strconv.ParseUint(k, 10, 8)
		clut[types.Color(c)] = v
	}
	fontSize := args[4].Float()
	rows := uint(args[5].Int())
	cols := uint(args[6].Int())
	dpi := args[7].Float()
	testpage := args[8].String()
	// 👇 constants
	maxFPS := 30.0
	paddedHeight := 1.5
	paddedWidth := 1.1
	tickMs := 333
	// 👇 load the fonts
	normalFont, _ := truetype.Parse(normalFontEmbed)
	normalFace := truetype.NewFace(normalFont, &truetype.Options{Size: fontSize, DPI: dpi /* , Hinting: font.HintingFull */})
	boldFont, _ := truetype.Parse(boldFontEmbed)
	boldFace := truetype.NewFace(boldFont, &truetype.Options{Size: fontSize, DPI: dpi /* , Hinting: font.HintingFull */})
	// 👇 measure the cell size
	temp := gg.NewContext(100, 100)
	temp.SetFontFace(boldFace)
	fontWidth, fontHeight := temp.MeasureString("M")
	// 👇 resize the canvas
	canvasWidth := float64(cols) * math.Round(fontWidth*paddedWidth)
	canvasHeight := float64(rows) * math.Round(fontHeight*paddedHeight)
	wrapper := canvas.Get("parentNode")
	wrapper.Get("style").Set("width", fmt.Sprintf("%fpx", canvasWidth))
	wrapper.Get("style").Set("height", fmt.Sprintf("%fpx", canvasHeight))
	canvas.Set("width", canvasWidth)
	canvas.Set("height", canvasHeight)
	// 👇 prepare the rendering surface
	rgba := image.NewRGBA(image.Rect(0, 0, int(canvasWidth), int(canvasHeight)))
	// 👇 kick off loops
	m.rcLoop(canvas, rgba, maxFPS)
	m.tickLoop(tickMs)
	// 👇 finally!
	cfg := types.Config{
		BgColor:      bgColor,
		BoldFace:     &boldFace,
		CLUT:         clut,
		Cols:         cols,
		FontHeight:   fontHeight,
		FontSize:     fontSize,
		FontWidth:    fontWidth,
		Monochrome:   monochrome,
		NormalFace:   &normalFace,
		PaddedHeight: paddedHeight,
		PaddedWidth:  paddedWidth,
		RGBA:         rgba,
		Rows:         rows,
		Testpage:     testpage,
	}
	return &cfg

}

// 🟦 Create the Javascript interface through which the UI calls the Go code

func (m *Mediator) jsInterface() js.Value {
	functions := map[string]any{
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			m.close()
			return nil
		}),
		"focus": js.FuncOf(func(this js.Value, args []js.Value) any {
			state := args[0].Bool()
			m.bus.PubFocus(state)
			return nil
		}),
		"keystroke": js.FuncOf(func(this js.Value, args []js.Value) any {
			key := types.Keystroke{
				Code:  args[0].String(),
				Key:   args[1].String(),
				ALT:   args[2].Bool(),
				CTRL:  args[3].Bool(),
				SHIFT: args[4].Bool(),
			}
			m.bus.PubKeystroke(key)
			return nil
		}),
		"outbound": js.FuncOf(func(this js.Value, args []js.Value) any {
			chars := make([]byte, args[0].Get("length").Int())
			js.CopyBytesToGo(chars, args[0])
			// 👇 data can be split into multiple frames
			slices := bytes.Split(chars, types.LT)
			for _, slice := range slices {
				if len(slice) > 0 {
					m.bus.PubOutbound(slice)
				}
			}
			return nil
		}),
	}
	return js.ValueOf(functions)
}

// 🟦 Forward messages from subscriptions to UI via dispatchEvent

func (m *Mediator) dispatchEvent(params map[string]any) {
	event := js.Global().Get("CustomEvent").New("go3270", map[string]any{
		"detail": params,
	})
	js.Global().Get("window").Call("dispatchEvent", event)
}

func (m *Mediator) inbound(chars []byte, _ core.PubInboundHints) {
	u8s := js.Global().Get("Uint8ClampedArray").New(len(chars))
	js.CopyBytesToJS(u8s, chars)
	params := map[string]any{
		"eventType": "inbound",
		"chars":     u8s,
	}
	m.dispatchEvent(params)
}

func (m *Mediator) panic(msg string) {
	params := map[string]any{
		"eventType": "panic",
		"args":      msg,
	}
	m.dispatchEvent(params)
	// 🔥 shut 'er down!
	m.bus.UnsubscribeAll()
}

func (m *Mediator) status(stat *types.Status) {
	params := map[string]any{
		"eventType": "status",
		"alarm":     stat.Alarm,
		"cursorAt":  stat.CursorAt,
		"error":     stat.Error,
		"insert":    stat.Insert,
		"locked":    stat.Locked,
		"message":   stat.Message,
		"numeric":   stat.Numeric,
		"protected": stat.Protected,
		"waiting":   stat.Waiting,
	}
	m.dispatchEvent(params)
}

// 🟦 Render drawing context when changed via requestAnimationFrame

func (m *Mediator) rcLoop(canvas js.Value, rgba *image.RGBA, maxFPS float64) {
	var (
		lastImage     []byte
		lastTimestamp float64
		rc            js.Func
	)
	rc = js.FuncOf(func(this js.Value, args []js.Value) any {
		timestamp := args[0].Float()
		// 👇 make sure we don't bust the max FPS we were given
		if timestamp-lastTimestamp >= (1000 / maxFPS) {
			if lastImage == nil || !bytes.Equal(lastImage, rgba.Pix) {
				// 🔥 I copied this from go-canvas where the author was worried
				// about 3 separate copies -- I haven't figured how to reduce
				// it to 2 even when using Uint8ClampedArray --
				// but it only takes ~2ms anyway
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
		js.Global().Call("requestAnimationFrame", rc)
		return nil
	})
	// 👇 kick off the rendering loop
	js.Global().Call("requestAnimationFrame", rc)
}

// 🟦 Inject ticks into the system eg: to support blinking

func (m *Mediator) tickLoop(interval int) {
	counter := 0
	ticker := js.FuncOf(func(this js.Value, args []js.Value) any {
		m.bus.PubTick(counter)
		counter++
		return nil
	})
	js.Global().Call("setInterval", ticker, interval)
}
