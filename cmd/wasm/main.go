package main

// https://elijahm.com/posts/local_first_htmx_part2/
//
// A service worker that runs a Go program compiled to WASM. 
// This service worker is responsible for proxying the fetch 
// requests and returning rendered HTML to the main thread.
// The server runs the same code. In a real world example
// the server would have some additional code and be authori-
// tative, but I’m bypassing that for the purpose of this POC.
//
// If for some reason the service worker is not installed when 
// a fetch request is made, that request will go to the server, 
// be handled by the server and rendered HTML will be returned
// just as if it was a SSR app.

import (
	"fmt"
	"syscall/js"
	"github.com/fbaube/go-browserver/api"
//	WU "github.com/fbaube/wasmutils"
)

func main() {
        js.Global().Get("console").Call("log", "Hello from Go WebAssembly!")
     	fmt.Println("Executing server...")
	err := api.GetData()
	if err != nil {
		fmt.Println("error syncing data with server")
	}
	fmt.Println("Running server...")
	go api.GetData()
	api.EchoStart() // This should start the func's JS version
/*
	// Now some JS+Dom stuff. // but DOES NOTHING
//	p := WU.Doc.Call("createElement", "h1")
	p := js.Global().Get("document").Call("createElement", "h1")
	p.Set("innerHTML", "Hello from Golang!")
	// WU.DocBody.Call("appendChild", p)
	js.Global().Get("body").Call("appendChild", p)
*/
}
