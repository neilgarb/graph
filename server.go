package graph

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type Server struct {
	Seq int

	mu    sync.Mutex
	Nodes map[*Node]bool
}

func NewServer() *Server {
	return &Server{
		Nodes: make(map[*Node]bool, 0),
	}
}

func (self *Server) Flush() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.Seq = 0
	self.Nodes = make(map[*Node]bool, 0)
}

func (self *Server) AddNode(node *Node) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.Seq += 1
	node.ID = self.Seq
	self.Nodes[node] = true
}

func (self *Server) GetNodeByID(ID int) *Node {
	self.mu.Lock()
	defer self.mu.Unlock()

	for n, _ := range self.Nodes {
		if n.ID == ID {
			return n
		}
	}
	return nil
}

func (self *Server) GetNodesByProperty(prop string, val string) []*Node {
	self.mu.Lock()
	defer self.mu.Unlock()

	nodes := make([]*Node, 0)
	for n, _ := range self.Nodes {
		if n.GetProperty(prop) == val {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (self *Server) GetNodeByProperty(prop string, val string) *Node {
	self.mu.Lock()
	defer self.mu.Unlock()

	for n, _ := range self.Nodes {
		if n.GetProperty(prop) == val {
			return n
		}
	}
	return nil
}

func (self *Server) RemoveNode(node *Node) {
	if node == nil {
		return
	}
	self.mu.Lock()
	delete(self.Nodes, node)
	self.mu.Unlock()
}

func (self *Server) Listen(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		} else {
			go self.HandleConnection(conn)
		}
	}
	panic("Unreachable code")
}

func (self *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connection from %s", conn.RemoteAddr())
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		line = strings.Trim(line, "\r\n")
		command := MakeCommand(line)
		if command == nil {
			conn.Write([]byte("ERROR Bad command.\n"))
		} else {
			err := command.Execute(self, conn)
			if err != nil {
				conn.Write([]byte("ERROR "))
				conn.Write([]byte(err.Error()))
				conn.Write([]byte("\n"))
			} else {
				conn.Write([]byte("END\n"))
			}
		}
	}
}
