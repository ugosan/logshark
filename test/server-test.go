package main

import (
	
	"fmt"
	"github.com/ugosan/logshark/cmd/server"
	"github.com/TylerBrock/colorjson"
)

func main() {
	c := make(chan map[string]interface{})

	go server.Start(c)
	
	f := colorjson.NewFormatter()
	f.Indent = 4
		
	for {
		obj := <-c

    s, _ := f.Marshal(obj)
		fmt.Println(string(s))

	}
}