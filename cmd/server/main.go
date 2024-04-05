package main

import (
	"space-go/internal/server"
)

func main() {
	s := server.FromDump("dumps/planet_dump_1.json", "dumps/graph_dump_1.json")
	err := server.Run(s)
	if err != nil {
		panic(err)
	}
}
