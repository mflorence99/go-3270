package main

import (
	"fmt"
	"syscall/js"
)

var console = js.Global().Get("console")
var log = console.Get("log")

func greet(this js.Value, args []js.Value) interface{} {
	name := args[0].String()
	log.Invoke(fmt.Sprintf("%%cHello, %s from Go %%cWebAssembly!", name), "color: pink", "color: skyblue")
	return name
}
