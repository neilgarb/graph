package graph

import (
	"fmt"
	"sync"
)

type Node struct {
	ID int

	linksMu sync.Mutex
	Links   map[*Node]bool

	propsMu    sync.Mutex
	Properties map[string]string
}

func NewNode() *Node {
	return &Node{
		Links:      make(map[*Node]bool, 0),
		Properties: make(map[string]string, 0),
	}
}

func (self *Node) String() string {
	return fmt.Sprintf("%d", self.ID)
}

func (self *Node) SetProperty(name, val string) {
	self.propsMu.Lock()
	self.Properties[name] = val
	self.propsMu.Unlock()
}

func (self *Node) GetProperty(name string) string {
	self.propsMu.Lock()
	defer self.propsMu.Unlock()
	return self.Properties[name]
}

func (self *Node) GetProperties() map[string]string {
	self.propsMu.Lock()
	defer self.propsMu.Unlock()
	return self.Properties
}

func (self *Node) LinkTo(to *Node) {
	self.linksMu.Lock()
	self.Links[to] = true
	self.linksMu.Unlock()
}

func (self *Node) UnlinkTo(to *Node) {
	self.linksMu.Lock()
	delete(self.Links, to)
	self.linksMu.Unlock()
}

func (self *Node) Walk(fn func(*Node), visited map[*Node]bool) {
	_, ok := visited[self]
	if ok {
		return
	}

	self.linksMu.Lock()
	defer self.linksMu.Unlock()

	visited[self] = true
	fn(self)
	for l, _ := range self.Links {
		l.Walk(fn, visited)
	}
}

func (self *Node) PathTo(dest *Node, progress []*Node) []*Node {
	self.linksMu.Lock()
	defer self.linksMu.Unlock()

	parents := make(map[*Node]*Node, 0)

	q := NewQueue()
	q.Enqueue(self)

	for !q.IsEmpty() {
		u := q.Dequeue().(*Node)
		for l, _ := range u.Links {
			if _, ok := parents[l]; !ok {
				parents[l] = u
				q.Enqueue(l)
			}
		}
	}

	if _, ok := parents[dest]; !ok {
		return nil
	}

	path := make([]*Node, 0)
	cur := dest
	for cur != nil && cur != self {
		path = append([]*Node{cur}, path...)
		cur = parents[cur]
	}
	path = append([]*Node{self}, path...)

	return path
}
