package statistics

import "fmt"

var	EMPTYEQUEUEERROR = fmt.Errorf("empty queue")

type Node struct {
	value interface{}
	next *Node
	pre *Node
}

func NewNode(v interface{}) *Node {
	return & Node {
		value: v,
		next: nil,
		pre: nil,
	}
}

type Queue struct {
	size int
	capacity int
	head *Node
	tail *Node
}

func NewQueue(capacity int) *Queue {
	head := NewNode(nil)
	tail := NewNode(nil)
	head.next, tail.pre = tail, head

	return & Queue {
		size: 0,
		capacity: capacity,
		head: head,
		tail: tail,
	}
}

func (queue *Queue) Push(v interface{}) {
	node := NewNode(v)
	queue.tail.pre.next = node
	node.pre = queue.tail.pre
	node.next = queue.tail
	queue.tail.pre = node
	queue.size++

	if queue.size > queue.capacity {
		queue.Pop()
	}
}

func (queue *Queue) Pop() (interface{}, error) {
	if queue.size <= 0 {
		return nil, EMPTYEQUEUEERROR
	}

	node := queue.head.next;
	queue.head.next = node.next
	node.next.pre = queue.head
	queue.size--;

	return node.value, nil
}

func (queue *Queue) Front() (interface{}, error) {
	if queue.size <= 0 {
		return nil, EMPTYEQUEUEERROR
	}

	return queue.head.next.value, nil
}

func (queue *Queue) Back() (interface{}, error) {
	if queue.size <= 0 {
		return nil, EMPTYEQUEUEERROR
	}

	return queue.tail.pre, nil
}

func (queue *Queue) Len() int {
	return queue.size
}

func (queue *Queue) All() []interface{} {
	node := queue.head.next
	res := []interface{}{}
	for node != queue.tail {
		res = append(res, node.value)
		node = node.next
	}

	return res
}


func (queue *Queue) String() string {
	node := queue.head.next
	res := "[ "
	for node != queue.tail {
		res = fmt.Sprintf("%v -> %v", res, node.value)
		node = node.next
	}
	res += " ]"

	return res
}
