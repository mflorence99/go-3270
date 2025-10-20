package main

import (
	"go3270/mediator"
	"syscall/js"
)

func main() {
	println("🐞 WASM initialized")
	js.Global().Set("NewGo3270", js.FuncOf(mediator.NewMediator))
	select {}
}
