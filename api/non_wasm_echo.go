//go:build !js

package api

import (
	"fmt"
//	"log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EchoStart (non-JS Edition) registers routes and 
// handlers and then takes some shortcuts offered
// by the Echo framework to start the server.
func EchoStart() {
	fmt.Println("Running server-side EchoStart...")
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ServerDelay) // non-JS-only, to emulate network latencies 

//	if l, ok := e.Logger.(*log.Logger); ok {
//	   l.SetHeader("${time_rfc3339} ${level}")
	l := e.Logger 
	   fmt.Printf("logger: %T \n", l)
//	}
	// Routes
	e.GET ("/", hRenderTodosRoute)
	e.POST("/add",  hAddTodoRoute)
	e.POST("/toggle/:id", hToggleTodoRoute)
	
	// The next two routes are non-JS-only
	e.GET ("/sync",  hGetTodos)  
	e.POST("/sync", hSyncTodos) 

	fmt.Printf("Everything (server-side) is up and running...")
	// These next two commands are non-JS-only
	e.Static("/", "../../public/")   // non-JS-only; JS gets from server!
	e.Logger.Fatal(e.Start(":3000")) // non-JS-only
}
