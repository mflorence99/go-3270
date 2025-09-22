package main

import (
	_ "embed"
	"syscall/js"
)

var console = js.Global().Get("console")
var Log = console.Get("log")

//go:embed 3270Medium.wasm
var TerminalFontData []byte

func main() {
	Log.Invoke("%cGo WebAssembly initialized!", "color: pink")
	js.Global().Set("fontGo", js.FuncOf(font))
	js.Global().Set("renderGo", js.FuncOf(render))
	js.Global().Set("testGo", js.FuncOf(test))
	select {}
}
