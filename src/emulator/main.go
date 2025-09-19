package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Go WebAssembly initialized!")
	js.Global().Set("greetGo", js.FuncOf(greet))
	select {} 
}
