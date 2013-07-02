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

		switch nDataTyped := n.data.(type) {
		case int:
			str += fmt.Sprintf("%v ", nDataTyped)
		default:
			str += fmt.Sprint(nDataTyped)
		}

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

// question 2.3
func deleteMiddle(n *node) {

	// Assume n.next != nil, as it's supposed to be a "middle" node

	n.data = n.next.data
	n.next = n.next.next

}

func placeOnNode(prevN *node, newN *node) *node {
	if prevN == nil {
		return newN
	}

	prevN.next = newN
	prevN = newN

	return prevN
}

// question 2.4
func partitionAtX(n *node, x int) *node {
	var leftHeadN, rightHeadN, leftN, rightN *node

	// Don't process if no nodes, makes post-loop processing simpler
	if n == nil {
		return nil
	}

	for {

		if n == nil {
			break
		}

		// assume int for this question
		intVal := n.data.(int)

		if intVal < x {
			leftN = placeOnNode(leftN, n)

			// Save first left node as head
			if leftHeadN == nil {
				leftHeadN = leftN
			}

		} else {
			rightN = placeOnNode(rightN, n)

			// Save first right node as head
			if rightHeadN == nil {
				rightHeadN = rightN
			}
		}

		n = n.next
	}

	if leftHeadN == nil {
		// if no left head, just return right side
		return rightHeadN
	} else {
		// otherwise, set last left node next to righthead (possibly nil)
		leftN.next = rightHeadN
		return leftHeadN
	}
}

func curVal(n *node) int {
	if n != nil {
		switch nDataTyped := n.data.(type) {
		case int:
			return nDataTyped
		default:
			panic("Expected int")
		}
	}

	return 0
}

// q 2.5
func addTwoNodes(n1 *node, n2 *node) *node {
	var resultHead *node
	var curResult *node
	var lastResult *node
	carry := 0

	for {

		// Nothing else to process
		if n1 == nil && n2 == nil {
			break
		}

		curVal := curVal(n1) + curVal(n2) + carry

		if curVal > 9 {
			curResult = NewNode(curVal - 10)
			carry = 1
		} else {
			curResult = NewNode(curVal)
			carry = 0
		}

		if resultHead == nil {
			// Save head if first iteration
			resultHead = curResult
		} else {
			// Otherwise, link to last result
			lastResult.next = curResult
		}

		if n1 != nil {
			n1 = n1.next
		}

		if n2 != nil {
			n2 = n2.next
		}

		lastResult = curResult
	}

	// If there is a final carry, add it
	if carry > 0 {
		curResult = NewNode(1)
		lastResult.next = curResult
	}

	return resultHead
}

// need to pre-count nodes if doing it the other way where the 1's digit
// is last
func countNodes(n *node) int {
	count := 0
	for {
		if n == nil {
			break
		}
		count++
		n = n.next
	}

	return count
}

func addTwoNodesOtherWayRecurse(n1 *node, n2 *node, n1Diff int) (*node, int) {
	var curResult, nextResult *node
	var carry int

	if n1 == nil {
		return nil, 0
	}

	if n1Diff > 0 {
		nextResult, carry = addTwoNodesOtherWayRecurse(n1.next, n2, n1Diff-1)
	} else {
		nextResult, carry = addTwoNodesOtherWayRecurse(n1.next, n2.next, 0)
	}

	curVal := curVal(n1) + curVal(n2) + carry

	if curVal > 9 {
		curResult = NewNode(curVal - 10)
		carry = 1
	} else {
		curResult = NewNode(curVal)
		carry = 0
	}

	curResult.next = nextResult

	return curResult, carry
}

// q2.5 follow up
func addTwoNodesOtherWay(n1, n2 *node) *node {
	n1Count := countNodes(n1)
	n2Count := countNodes(n2)

	// To simplify algo, let's make it so we can assume n1 has >= digits than
	// n2
	if n1Count < n2Count {
		n1, n2 = n2, n1
		n1Count, n2Count = n2Count, n1Count
	}

	resultHead, carry := addTwoNodesOtherWayRecurse(n1, n2, n1Count-n2Count)

	if carry > 0 {
		carryNode := NewNode(carry)
		carryNode.next = resultHead
		resultHead = carryNode
	}

	return resultHead
}

func CreateNodesFromArr(arr []int) *node {
	var headN *node
	var lastN *node

	for _, v := range arr {
		newN := NewNode(v)

		if headN == nil {
			// set head if first
			headN = newN
		} else {
			// otherwise set lastN's next value
			lastN.next = newN
		}

		// set lastN for next iteration
		lastN = newN
	}

	return headN
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

func TestDeleteMiddle(t *testing.T) {
	data := "NODES NODES NODES SO MANY NODES!!!"
	parent := CreateNodesFromString(data, 0, len(data))
	middle := parent

	for i := 0; i < 25; i++ {
		middle = middle.next
	}
	deleteMiddle(middle)

	if parent.String() != "NODES NODES NODES SO MANYNODES!!!" {
		t.Error("Deleting middle node didn't work")
	}
}

func TestPartition(t *testing.T) {
	data := []int{234, 54, 546, 456, 756, 567, 9199, 1, 8, 4, 3, 2349}
	headN := CreateNodesFromArr(data)
	partN := partitionAtX(headN, 50)

	if partN.String() != "1 8 4 3 234 54 546 456 756 567 9199 2349 " {
		t.Error("Partition with left and right sides failed")
	}

	data = []int{234, 54, 546, 456, 756, 567, 9199, 1, 8, 4, 3, 2349}
	headN = CreateNodesFromArr(data)
	allLeftN := partitionAtX(headN, 99999)
	if allLeftN.String() != "234 54 546 456 756 567 9199 1 8 4 3 2349 " {
		t.Error("Partition with all left failed")
	}

	data = []int{234, 54, 546, 456, 756, 567, 9199, 1, 8, 4, 3, 2349}
	headN = CreateNodesFromArr(data)
	allRightN := partitionAtX(headN, 0)
	if allRightN.String() != "234 54 546 456 756 567 9199 1 8 4 3 2349 " {
		t.Error("Partition with all right failed")
	}
}

func TestAddTwoNodes(t *testing.T) {
	data1 := []int{7, 1, 6}
	data2 := []int{5, 9, 2}

	headN1 := CreateNodesFromArr(data1)
	headN2 := CreateNodesFromArr(data2)

	resultHead := addTwoNodesOtherWay(headN1, headN2)

	fmt.Println(resultHead)
}
