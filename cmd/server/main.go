package main

import (
	"space-go/internal/server"
)

func main() {
	s := server.FromDump("dumps/planet_dump_3.json", "dumps/graph_dump_3.json")
	err := server.Run(s)
	if err != nil {
		panic(err)
	}
}
