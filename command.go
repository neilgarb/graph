package graph

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

var (
	flushCommand    *regexp.Regexp
	getNodeCommand  *regexp.Regexp
	getNodesCommand *regexp.Regexp
	addNodeCommand  *regexp.Regexp
	delNodeCommand  *regexp.Regexp
	setPropCommand  *regexp.Regexp
	getPropCommand  *regexp.Regexp
	getPropsCommand *regexp.Regexp
	linkCommand     *regexp.Regexp
	unlinkCommand   *regexp.Regexp
	walkCommand     *regexp.Regexp
	pathCommand     *regexp.Regexp
)

func init() {
	flushCommand = regexp.MustCompile(`FLUSH`)
	getNodeCommand = regexp.MustCompile(`^GETNODE "([^"]+)" "([^"]+)"$`)
	getNodesCommand = regexp.MustCompile(`^GETNODES "([^"]+)" "([^"]+)"$`)
	addNodeCommand = regexp.MustCompile(`^ADDNODE$`)
	delNodeCommand = regexp.MustCompile(`^DELNODE (\d+)$`)
	setPropCommand = regexp.MustCompile(`^SETPROP (\d+) "([^"]+)" "([^"]+)"$`)
	getPropCommand = regexp.MustCompile(`^GETPROP (\d+) "([^"]+)"$`)
	getPropsCommand = regexp.MustCompile(`^GETPROPS (\d+)$`)
	linkCommand = regexp.MustCompile(`^LINK (\d+) TO (\d+)$`)
	unlinkCommand = regexp.MustCompile(`^UNLINK (\d+) TO (\d+)$`)
	walkCommand = regexp.MustCompile(`^WALK (\d+)( RETURN "([^"]+)")?$`)
	pathCommand = regexp.MustCompile(`^PATH (\d+) TO (\d+)$`)
}

type Command interface {
	Execute(*Server, net.Conn) error
}

func MakeCommand(line string) Command {
	if flushCommand.MatchString(line) {
		return &FlushCommand{}
	} else if getNodeCommand.MatchString(line) {
		return &GetNodeCommand{
			Prop: getNodeCommand.FindStringSubmatch(line)[1],
			Val:  getNodeCommand.FindStringSubmatch(line)[2],
		}
	} else if getNodesCommand.MatchString(line) {
		return &GetNodesCommand{
			Prop: getNodesCommand.FindStringSubmatch(line)[1],
			Val:  getNodesCommand.FindStringSubmatch(line)[2],
		}
	} else if addNodeCommand.MatchString(line) {
		return &AddNodeCommand{}
	} else if delNodeCommand.MatchString(line) {
		return &DelNodeCommand{
			ID: delNodeCommand.FindStringSubmatch(line)[1],
		}
	} else if setPropCommand.MatchString(line) {
		return &SetPropCommand{
			ID:   setPropCommand.FindStringSubmatch(line)[1],
			Prop: setPropCommand.FindStringSubmatch(line)[2],
			Val:  setPropCommand.FindStringSubmatch(line)[3],
		}
	} else if getPropCommand.MatchString(line) {
		return &GetPropCommand{
			ID:   getPropCommand.FindStringSubmatch(line)[1],
			Prop: getPropCommand.FindStringSubmatch(line)[2],
		}
	} else if getPropsCommand.MatchString(line) {
		return &GetPropsCommand{
			ID: getPropsCommand.FindStringSubmatch(line)[1],
		}
	} else if linkCommand.MatchString(line) {
		return &LinkCommand{
			From: linkCommand.FindStringSubmatch(line)[1],
			To:   linkCommand.FindStringSubmatch(line)[2],
		}
	} else if unlinkCommand.MatchString(line) {
		return &UnlinkCommand{
			From: unlinkCommand.FindStringSubmatch(line)[1],
			To:   unlinkCommand.FindStringSubmatch(line)[2],
		}
	} else if walkCommand.MatchString(line) {
		return &WalkCommand{
			Start:  walkCommand.FindStringSubmatch(line)[1],
			Return: walkCommand.FindStringSubmatch(line)[3],
		}
	} else if pathCommand.MatchString(line) {
		return &PathCommand{
			From: pathCommand.FindStringSubmatch(line)[1],
			To:   pathCommand.FindStringSubmatch(line)[2],
		}
	}
	return nil
}

type FlushCommand struct {
}

func (self *FlushCommand) Execute(server *Server, conn net.Conn) error {
	server.Flush()
	return nil
}

type GetNodeCommand struct {
	Prop string
	Val  string
}

func (self *GetNodeCommand) Execute(server *Server, conn net.Conn) error {
	node := server.GetNodeByProperty(self.Prop, self.Val)
	if node == nil {
		return errors.New("No such node.")
	}
	conn.Write([]byte(fmt.Sprintf("%d", node.ID)))
	conn.Write([]byte("\n"))
	return nil
}

type GetNodesCommand struct {
	Prop string
	Val  string
}

func (self *GetNodesCommand) Execute(server *Server, conn net.Conn) error {
	nodes := server.GetNodesByProperty(self.Prop, self.Val)
	for _, n := range nodes {
		conn.Write([]byte(fmt.Sprintf("%d", n.ID)))
		conn.Write([]byte("\n"))
	}
	return nil
}

type AddNodeCommand struct {
}

func (self *AddNodeCommand) Execute(server *Server, conn net.Conn) error {
	node := NewNode()
	server.AddNode(node)
	conn.Write([]byte(fmt.Sprintf("%d\n", node.ID)))
	return nil
}

type DelNodeCommand struct {
	ID string
}

func (self *DelNodeCommand) Execute(server *Server, conn net.Conn) error {
	ID, _ := strconv.Atoi(self.ID)
	node := server.GetNodeByID(ID)
	server.RemoveNode(node)
	return nil
}

type SetPropCommand struct {
	ID   string
	Prop string
	Val  string
}

func (self *SetPropCommand) Execute(server *Server, conn net.Conn) error {
	ID, _ := strconv.Atoi(self.ID)
	node := server.GetNodeByID(ID)
	if node == nil {
		return errors.New("No such node.")
	}
	node.SetProperty(self.Prop, self.Val)
	return nil
}

type GetPropCommand struct {
	ID   string
	Prop string
}

func (self *GetPropCommand) Execute(server *Server, conn net.Conn) error {
	ID, _ := strconv.Atoi(self.ID)
	node := server.GetNodeByID(ID)
	if node == nil {
		return errors.New("No such node.")
	}
	conn.Write([]byte("\""))
	conn.Write([]byte(node.GetProperty(self.Prop)))
	conn.Write([]byte("\"\n"))
	return nil
}

type GetPropsCommand struct {
	ID string
}

func (self *GetPropsCommand) Execute(server *Server, conn net.Conn) error {
	ID, _ := strconv.Atoi(self.ID)
	node := server.GetNodeByID(ID)
	if node == nil {
		return errors.New("No such node.")
	}
	properties := node.GetProperties()
	for k, v := range properties {
		conn.Write([]byte("\""))
		conn.Write([]byte(k))
		conn.Write([]byte("\""))
		conn.Write([]byte("\n"))
		conn.Write([]byte("\""))
		conn.Write([]byte(v))
		conn.Write([]byte("\""))
		conn.Write([]byte("\n"))
	}
	return nil
}

type LinkCommand struct {
	From string
	To   string
}

func (self *LinkCommand) Execute(server *Server, conn net.Conn) error {
	fromID, _ := strconv.Atoi(self.From)
	from := server.GetNodeByID(fromID)
	if from == nil {
		return errors.New("No such node.")
	}
	toID, _ := strconv.Atoi(self.To)
	to := server.GetNodeByID(toID)
	if to == nil {
		return errors.New("No such node.")
	}
	from.LinkTo(to)
	return nil
}

type UnlinkCommand struct {
	From string
	To   string
}

func (self *UnlinkCommand) Execute(server *Server, conn net.Conn) error {
	fromID, _ := strconv.Atoi(self.From)
	from := server.GetNodeByID(fromID)
	if from == nil {
		return errors.New("No such node.")
	}
	toID, _ := strconv.Atoi(self.To)
	to := server.GetNodeByID(toID)
	if to == nil {
		return errors.New("No such node.")
	}
	from.UnlinkTo(to)
	return nil
}

type WalkCommand struct {
	Start  string
	Return string
}

func (self *WalkCommand) Execute(server *Server, conn net.Conn) error {
	id, _ := strconv.Atoi(self.Start)
	start := server.GetNodeByID(id)
	if start == nil {
		return errors.New("No such node.")
	}
	start.Walk(func(n *Node) {
		if self.Return != "" {
			conn.Write([]byte("\""))
			conn.Write([]byte(n.GetProperty(self.Return)))
			conn.Write([]byte("\""))
		} else {
			conn.Write([]byte(fmt.Sprintf("%d", n.ID)))
		}
		conn.Write([]byte("\n"))
	}, make(map[*Node]bool, 0))
	return nil
}

type PathCommand struct {
	From string
	To   string
}

func (self *PathCommand) Execute(server *Server, conn net.Conn) error {
	fromID, _ := strconv.Atoi(self.From)
	from := server.GetNodeByID(fromID)
	if from == nil {
		return errors.New("No such node.")
	}
	toID, _ := strconv.Atoi(self.To)
	to := server.GetNodeByID(toID)
	if to == nil {
		return errors.New("No such node.")
	}
	path := from.PathTo(to, make([]*Node, 0))
	if path == nil {
		return errors.New("No path.")
	}
	for _, n := range path {
		conn.Write([]byte(fmt.Sprintf("%d", n.ID)))
		conn.Write([]byte("\n"))
	}
	return nil
}
