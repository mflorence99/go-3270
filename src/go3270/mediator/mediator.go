package mediator

import (
	_ "embed"
	"fmt"
	"image"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"

	"go3270/emulator"
	"go3270/emulator/bus"
	"go3270/emulator/consts"
	"go3270/emulator/keyboard"
	"syscall/js"
)

// 🔥 Hack alert! we must use extension {js, wasm} and we can't use symlinks, so this file is a copy of the font renamed

var (
	//go:embed IBMPlexMono-Regular.ttf.wasm
	normalFontEmbed []byte
	//go:embed IBMPlexMono-Bold.ttf.wasm
	boldFontEmbed []byte
)

// 🟧 Bridge between Typescript UI and Go-powered emulator

// The design objective is that all Go <-> UI communication goes through this package. No other package must use syscall/js. That way, everything but here can be tested with go test.

// The emulator package is handed an image into which it renders the 3270 stream and any operator input. Using requestAnimationFrame, this module actually draws the context onto a supplied HTML canvas whenever the context changes

type Mediator struct {
	bus *bus.Bus
	emu *emulator.Emulator
}

// 👁️ go3270.ts

// args[0] canvas
// args[1] bgColor
// args[2] color [normal, highlight]
// args[3] clut [map color -> [normal, highlight]]
// args[4] fontSize
// args[5] cols
// args[6] rows
// args[7] dpi

func NewMediator(this js.Value, args []js.Value) any {
	m := new(Mediator)
	m.bus = bus.NewBus()
	// 🔥 must subscribe BEFORE we create the emulator
	m.bus.Subscribe("close", m.close)
	// 👇 create and configure the emulator and its childreen
	m.emu = emulator.NewEmulator(m.bus)
	cfg := m.configure(args)
	m.bus.Publish("config", cfg)
	return m.jsInterface()
}

func (m *Mediator) close() {
	m.bus.UnsubscribeAll()
}

func (m *Mediator) configure(args []js.Value) *emulator.Config {
	// 👇 from the args
	canvas := args[0]
	bgColor := args[1].String()
	color := [2]string{args[2].Index(0).String(), args[2].Index(1).String()}
	obj := args[3]
	clut := make(map[consts.Color][2]string)
	keys := js.Global().Get("Object").Call("keys", obj)
	for i := 0; i < keys.Length(); i++ {
		k := keys.Index(i).String()
		v := [2]string{obj.Get(k).Index(0).String(), obj.Get(k).Index(1).String()}
		fmt.Printf("🎨 %s = %v\n", k, v)
		clut[consts.ColorOf(k)] = v
	}
	fontSize := args[4].Float()
	cols := args[5].Int()
	rows := args[6].Int()
	dpi := args[7].Float()
	// 👇 constants
	paddedHeight := 1.5
	paddedWidth := 1.1
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
	// 👇 finally!
	cfg := emulator.Config{
		BgColor:      bgColor,
		BoldFace:     &boldFace,
		CLUT:         clut,
		Color:        color,
		Cols:         80,
		FontHeight:   fontHeight,
		FontSize:     fontSize,
		FontWidth:    fontWidth,
		NormalFace:   &normalFace,
		PaddedHeight: paddedHeight,
		PaddedWidth:  paddedWidth,
		RGBA:         rgba,
		Rows:         24,
	}
	return &cfg

}

func (m *Mediator) jsInterface() js.Value {
	functions := map[string]any{
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			m.bus.Publish("close")
			return nil
		}),
		"focus": js.FuncOf(func(this js.Value, args []js.Value) any {
			state := args[0].Bool()
			m.bus.Publish("focus", state)
			return nil
		}),
		"keystroke": js.FuncOf(func(this js.Value, args []js.Value) any {
			key := keyboard.Keystroke{
				Code:  args[0].String(),
				Key:   args[1].String(),
				ALT:   args[3].Bool(),
				CTRL:  args[4].Bool(),
				SHIFT: args[5].Bool(),
			}
			m.bus.Publish("keystroke", key)
			return nil
		}),
		"outbound": js.FuncOf(func(this js.Value, args []js.Value) any {
			bytes := make([]byte, args[0].Get("length").Int())
			js.CopyBytesToGo(bytes, args[0])
			m.bus.Publish("outbound", bytes)
			return nil
		}),
	}
	return js.ValueOf(functions)
}
