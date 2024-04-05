package main

import (
	"space-go/internal/server"
)

func main() {
	s := server.FromDump("dumps/planet_dump_2.json", "dumps/graph_dump_2.json")
	err := server.Run(s)
	if err != nil {
		panic(err)
	}
}
