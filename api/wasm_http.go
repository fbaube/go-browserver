//go:build js

package api

import (
        "fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

// EchoStart (JS Edition) registers routes and handlers 
// and then takes the echo.Server.Handler which is of 
// type http.Handler and passes that into wasmhttp.Serve() 
func EchoStart() {
     	fmt.Printf("Running in-browser EchoStart...")
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(SyncToServer) // JS-only

	// Routes
	e.GET("/", renderTodosRoute)
	e.POST("/add", addTodoRoute)
	e.POST("/toggle/:id", toggleTodoRoute)

	// Start server (JS-only)
	wasmhttp.Serve(e.Server.Handler) 

	fmt.Printf("Everything (in-browser) is up and running...")
	select {} // JS-only, to keep in memory in browser 
}
