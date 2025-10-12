package mediator

import (
	"go3270/emulator"
	"go3270/emulator/bus"
	"syscall/js"
)

type Mediator struct {
	bus      *bus.Bus
	emulator *emulator.Emulator
}

func NewMediator(this js.Value, args []js.Value) any {
	m := new(Mediator)
	m.bus = bus.NewBus()
	m.emulator = emulator.NewEmulator(m.bus)
	// 🟦 Go WASM methods callable by Javascript
	// 👁️ go3270.d.ts
	tsInterface := map[string]any{
		"close": js.FuncOf(func(this js.Value, args []js.Value) any {
			return nil
		}),
		"focussed": js.FuncOf(func(this js.Value, args []js.Value) any {
			return nil
		}),
		"keystroke": js.FuncOf(func(this js.Value, args []js.Value) any {
			return nil
		}),
		"receiveFromApp": js.FuncOf(func(this js.Value, args []js.Value) any {
			return nil
		}),
	}
	return js.ValueOf(tsInterface)
}
