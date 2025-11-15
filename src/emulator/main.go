package main

import (
	"emulator/mediator"
	"syscall/js"
)

func main() {
	println("🐞 WASM initialized")
	js.Global().Set("NewGo3270", js.FuncOf(mediator.NewMediator))
	select {}
}
