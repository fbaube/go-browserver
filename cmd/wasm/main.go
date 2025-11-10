package main

import (
	"fmt"
	"github.com/elijahmorg/lhmtx/api" // FIXME
)

func main() {
     	fmt.Println("Executing...")
	err := api.GetData()
	if err != nil {
		fmt.Println("error syncing data with server")
	}
	fmt.Println("Running...")
	go api.GetData()
	api.EchoStart()
}
