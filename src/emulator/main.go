package main

import (
	"emulator/go3270"
	"syscall/js"
)

// 🟧 3270 emulator

func main() {
	js.Global().Get("console").Call("log", "%cgo-3270 WASM initialized", "color: pink")
	js.Global().Set("NewGo3270", js.FuncOf(go3270.New))
	select {}
}
