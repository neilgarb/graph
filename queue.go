package graph

type Queue struct {
	q []interface{}
}

func NewQueue() *Queue {
	return &Queue{
		q: make([]interface{}, 0),
	}
}

func (self *Queue) Enqueue(i interface{}) {
	if i == nil {
		return
	}
	self.q = append(self.q, i)
}

func (self *Queue) Dequeue() interface{} {
	if len(self.q) == 0 {
		return nil
	}
	i := self.q[0]
	self.q = self.q[1:]
	return i
}

func (self *Queue) IsEmpty() bool {
	return len(self.q) == 0
}
