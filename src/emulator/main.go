package main

import (
	"emulator/go3270"
	"syscall/js"
)

// 🟧 3270 emulator

func main() {
	println("🐞 Go3270 WASM initialized")
	js.Global().Set("NewGo3270", js.FuncOf(go3270.New))
	select {}
}
