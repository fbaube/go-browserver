//go:build js

package api

import (
        "fmt"
//	"log"
	"syscall/js"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	wasmhttp "github.com/nlepage/go-wasm-http-server"
	WU "github.com/fbaube/wasmutils"
)

// EchoStart (JS Edition) registers routes and handlers 
// and then takes the echo.Server.Handler which is of 
// type http.Handler and passes that into wasmhttp.Serve() 
func EchoStart() {
     	println("println: Running in-browser EchoStart...")
     	fmt.Printf("fmt.Printf: Running in-browser EchoStart...")
	js.Global().Get("console").Call("log", "Hello1 from Go WebAssembly!")
	// WU.G.Get("console").Call("log", "Hello1 from Go WebAssembly!")
       	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(SyncToServer) // JS-only
	/*
	if l, ok := e.Logger.(*log.Logger); ok {
	   l.SetHeader("${time_rfc3339} ${level}")
	} */
	// Routes
	e.GET("/", hRenderTodosRoute)
	e.POST("/add", hAddTodoRoute)
	e.POST("/toggle/:id", hToggleTodoRoute)

	// Start server (JS-only)
	wasmhttp.Serve(e.Server.Handler) 

	fmt.Printf("Everything (in-browser) is up and running...")
	select {} // JS-only, to keep in memory in browser 
}
