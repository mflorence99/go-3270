package main

import (
	"fmt"
	"syscall/js"
)


func test(this js.Value, args []js.Value) interface{} {
	name := args[0].String()
	Log.Invoke(fmt.Sprintf("%%cHello, %s from Go %%cWebAssembly!", name), "color: pink", "color: skyblue")
	Log.Invoke(len(TerminalFontData))
	return name
}

func render(this js.Value, args []js.Value) interface{} {
	// 👇 simulate response
	data := []byte{193, 194, 195 /* 👈 EBCDIC "ABC" */ } 
	uint8ArrayConstructor := js.Global().Get("Uint8Array")
	result := uint8ArrayConstructor.New(len(data))
	js.CopyBytesToJS(result, data)
	return result
}
