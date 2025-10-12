package go3270

import (
	"go3270/mediator"
	"syscall/js"
)

// 🟧 3270 emulator

func main() {
	println("🐞 Go3270 WASM initialized")
	js.Global().Set("NewGo3270", js.FuncOf(mediator.NewMediator))
	select {}
}
