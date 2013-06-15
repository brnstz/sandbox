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

	for {
		if n == nil {
			break
		}

		str += fmt.Sprint(n.data)

		n = n.next
	}

	return str
}

// question 2.1
func removeDuplicates(n *node) {
	dupes := map[interface{}]bool{}

	var last *node

	for {

		// If no more nodes, we're done
		if n == nil {
			break
		}

		if dupes[n.data] {
			// Skip nodes that are dupes
			last.next = n.next
		} else {
			// Otherwise record letter for potential dupes later and
			// set last valid node
			dupes[n.data] = true
			last = n
		}

		// Always increment to next node
		n = n.next
	}
}

// question 2.2
func findKthToLast(n *node, k int) interface{} {
	nCount := n

	i := 0
	for {

		if nCount == nil {
			break
		}

		nCount = nCount.next
		i++
	}

	// k is longer than i or k is negative, no element to be found
	if k > i || k < 0 {
		return nil
	}

	for j := 0; j < i-k-1; j++ {
		fmt.Println("j: ", j)
		n = n.next
	}

	return n.data
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

	if node.String() != "FOLW UP" {
		t.Error("remove duplicates did not match FOLW UP")
	}
}

func TestFindKthToLast(t *testing.T) {
	data := "ZHIS IS A STRING"
	node := CreateNodesFromString(data, 0, len(data))

	item := findKthToLast(node, 15)
	if item != "Z" {
		t.Error("15th node from end should be Z")
	}

	item = findKthToLast(node, -1)
	if item != nil {
		t.Error("-1th node from end should be nil")
	}

}
