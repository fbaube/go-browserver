//go:build js

package api

import (
        "fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func EchoStart() {
     	fmt.Printf("Executing...")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(SyncToServer)

	// Routes
	e.GET("/", renderTodosRoute)
	e.POST("/toggle/:id", toggleTodoRoute)
	e.POST("/add", addTodoRoute)

	// Start server
	wasmhttp.Serve(e.Server.Handler)

	fmt.Printf("Everything is up and running...")
	select {}
}
