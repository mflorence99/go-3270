package main

import (
	"syscall/js"
)

var console = js.Global().Get("console")
var Log = console.Get("log")

func main() {
	Log.Invoke("%cGo WebAssembly initialized!", "color: pink")
	js.Global().Set("NewGo3270", js.FuncOf(NewGo3270))
	select {}
}
