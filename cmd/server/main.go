package main

import (
        "fmt"
	// "github.com/elijahmorg/lhmtx/api" 
	"github.com/fbaube/go-browserver/api"
)

func main() {
     	fmt.Printf("Running main...")
	api.EchoStart() // This	will start the func's non-JS version
}
