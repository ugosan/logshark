package main

import (
	"fmt"

	"github.com/TylerBrock/colorjson"
	"github.com/ugosan/logshark/v1/server"
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
