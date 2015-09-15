package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
)

type Node struct {
	mu    sync.Mutex
	Links []*Node
}

type Server struct {
	mu    sync.Mutex
	Nodes []*Node
}

var server Server

func (self *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Trim(line, "\r\n")
		switch line {
		default:
			conn.Write([]byte("ERROR\n"))
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	server = Server{
		Nodes: make([]*Node, 0),
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go server.handleConnection(conn)
	}
}
