package chap2_test

import (
	"fmt"
	"testing"
)

type node struct {
	next *node
	data interface{}
}

func NewNode(inData interface{}) *node {
	n := &node{data: inData}
	return n
}

func (n *node) appendNode(inNode *node) {
	for n.next != nil {
		n = n.next
	}

	n.next = inNode
}

func (n *node) String() string {
	str := ""

	for n.next != nil {
		str += fmt.Sprintf("%v", n.data)
		n = n.next
	}
	// print last one
	str += fmt.Sprintf("%v", n.data)

	return str
}

// quesiton 2.1
func removeDuplicates(n *node) {
	dupes := map[interface{}]bool{}

	// Note it's impossible for first node to be a dupe
	// FIXME: not getting O, maybe recursive better?
	var last *node
	for n.next != nil {
		if dupes[n.data] {
			last.next = n.next
		}
		dupes[n.data] = true
		last = n
		n = n.next
	}
}

func CreateNodesFromString(data string, i, max int) *node {
	if i < max {
		curN := NewNode(data[i : i+1])
		curN.next = CreateNodesFromString(data, i+1, max)
		return curN
	} else {
		return nil
	}
}

func TestRemoveDuplicates(t *testing.T) {
	data := "FOLLOW UP"
	node := CreateNodesFromString(data, 0, len(data))
	removeDuplicates(node)
	fmt.Println(node.String())

}
