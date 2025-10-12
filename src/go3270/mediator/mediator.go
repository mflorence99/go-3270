package mediator

import (
	"go3270/emulator"
	"go3270/emulator/bus"
	"syscall/js"
)

type Mediator struct {
	bus *bus.Bus
	emu *emulator.Emulator
}

// 👁️ go3270.ts

// args[0] canvas
// args[1] bgColor
// args[2] color [normal, highlight]
// args[3] clut [map color -> [normal, highlight]]
// args[4] fontSise
// args[5] cols
// args[6] rows
// args[7] dpi

func NewMediator(this js.Value, args []js.Value) any {
	m := new(Mediator)
	m.bus = bus.NewBus()
	m.bus.Subscribe("close", m.close)
	m.emu = emulator.NewEmulator(m.bus)
	cfg := emulator.Config{
		Cols: 80,
		Rows: 24,
	}
	m.bus.Publish("config", &cfg)
	return m.jsInterface()
}

func (m *Mediator) close() {
	m.bus.UnsubscribeAll()
}

func (m *Mediator) configure() {
	m.bus.UnsubscribeAll()
}

func (m *Mediator) jsInterface() js.Value {
	functions := map[string]any{
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			m.bus.Publish("close")
			return nil
		}),
		"focussed": js.FuncOf(func(this js.Value, args []js.Value) any {
			return nil
		}),
		"keystroke": js.FuncOf(func(this js.Value, args []js.Value) any {
			m.bus.Publish("keystroke", args[0].String(), args[1].String(), args[2].Bool(), args[3].Bool(), args[4].Bool())
			return nil
		}),
		"outbound": js.FuncOf(func(this js.Value, args []js.Value) any {
			return nil
		}),
	}
	return js.ValueOf(functions)
}
