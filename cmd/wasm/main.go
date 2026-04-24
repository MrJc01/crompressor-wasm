//go:build js && wasm

package main

import (
	"syscall/js"
	"fmt"
)

func main() {
	fmt.Println("Crompressor WASM module loaded")

	js.Global().Set("cromPack", js.FuncOf(pack))
	js.Global().Set("cromUnpack", js.FuncOf(unpack))
	js.Global().Set("cromInfo", js.FuncOf(info))

	// Block forever
	select {}
}

func pack(this js.Value, args []js.Value) interface{} {
	// TODO: Wire to crompressor Pack API
	return js.ValueOf("pack: not implemented yet")
}

func unpack(this js.Value, args []js.Value) interface{} {
	// TODO: Wire to crompressor Unpack API
	return js.ValueOf("unpack: not implemented yet")
}

func info(this js.Value, args []js.Value) interface{} {
	return js.ValueOf("crompressor-wasm v0.1.0")
}
