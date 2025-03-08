package main

import (
	"fmt"
	"syscall/js"

	"baby-blog/wasm/pages/test/buttons"
)

type WasmApplication struct {
	Path string
}

func (application *WasmApplication) go_setpath(this js.Value, p []js.Value) interface{} {
	application.Path = p[0].String()
	fmt.Println("[WASM]: Path set to", application.Path)
	application.updateDOMContent()
	return nil
}

func (application *WasmApplication) updateDOMContent() {
	switch application.Path {
	case "test/buttons":
		buttons.Begin_Interactivity()
		break
	default:
		fmt.Println("[WASM]: Path is not implemented")
		return
	}
}

func (application *WasmApplication) init() {
	js.Global().Set("go_setpath", js.FuncOf(application.go_setpath))
}
func main() {
	application := &WasmApplication{}
	application.init()
	ch := make(chan string)
	fmt.Println("[WASM]: Channel created")
	<-ch // Prevent the program from exiting immediately
}
