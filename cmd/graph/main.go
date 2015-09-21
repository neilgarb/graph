package main

import (
	"flag"
	"log"

	"github.com/neilgarb/graph"
)

var (
	address *string
)

func init() {
	address = flag.String("address", ":8080", "Address to bind to")
	flag.Parse()
}

func main() {
	server := graph.NewServer()
	log.Fatal(server.Listen(*address))
}
