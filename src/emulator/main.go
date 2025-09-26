package main

import (
	"emulator/go3270"
	"syscall/js"
)

var console = js.Global().Get("console")
var Log = console.Get("log")

// 🟧 3270 emulator

func main() {
	Log.Invoke("%cGo WebAssembly initialized!", "color: pink")
	js.Global().Set("NewGo3270", js.FuncOf(go3270.NewGo3270))
	select {}
}
