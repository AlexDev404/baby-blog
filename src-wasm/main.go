// This comment block was generated by running the following command:
// go generate
// DO NOT EDIT.
//go:build wasm
// +build wasm

// Attr: https://www.sprkweb.dev/en/posts/go-wasm-with-typescript/
// Attr: https://permify.co/post/wasm-go/
// Attr: https://pkg.go.dev/syscall/js@go1.24.1
package main

import (
	"fmt"
	"syscall/js"
	"net/http"
	"baby-blog/wasm/pages/test/buttons"
)

type WasmApplication struct {
	Path string
}

func (application *WasmApplication) GoSetPath(this js.Value, p []js.Value) interface{} {
	application.Path = p[0].String()
	fmt.Println("[WASM]: Path set to", application.Path)
	application.updateDOMContent()
	return nil
}

func (application *WasmApplication) updateDOMContent() {
	switch application.Path {
	case "test/buttons":
		buttons.BeginInteractivity()
		break
	default:
		fmt.Println("[WASM]: No implementation for this path")
		return
	}
}

func (application *WasmApplication) init() {
	js.Global().Set("go_setpath", js.FuncOf(application.GoSetPath))
}

func combinedHandler(path string) {
	switch path {
	case "home":
		http.ServeFile(w, r, "./templates/home.html")
	case "week1":
		http.ServeFile(w, r, "./templates/week1.html")
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func main() {
	application := &WasmApplication{}
	application.init()
	ch := make(chan string)
	fmt.Println("[WASM]: Channel created")
	<-ch // Prevent the program from exiting immediately
}
