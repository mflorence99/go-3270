package main

import (
	"emulator/go3270"
	"fmt"
	"syscall/js"
)

// 🟧 3270 emulator

func main() {
	fmt.Println("go-3270 WASM initialized")
	js.Global().Set("NewGo3270", js.FuncOf(go3270.NewGo3270))
	select {}
}
