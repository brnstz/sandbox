package chap1_test

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sort"
	"testing"
)

// question 1.1
func isUniq(input string) bool {
	var count = map[rune]int{}
	for _, c := range input {
		count[c]++
		if count[c] > 1 {
			return false
		}
	}

	return true
}

type SortableRunes []rune

func (s SortableRunes) Len() int {
	return len(s)
}
func (s SortableRunes) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s SortableRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func isUniqNoDataStruct(input string) bool {
	// Convert to rune array just for sorting
	runes := SortableRunes(input)
	sort.Sort(runes)

	// Start with 0
	var last rune = 0

	for _, c := range input {
		// If last character is same as current, then we have a repeat
		// in our current string
		if last == c {
			return false
		}

		// For next iteration, save current to last
		last = c
	}

	// If we get this far without returning false, we have no dupes
	return true
}

// question 2.1
func reverse(input []rune) []rune {
	max_index := len(input) - 1

	for i := 0; i <= max_index; i++ {
		// Find index we should swap with
		swap_index := max_index - i

		// If i has passed swap index, we are finished
		if i >= swap_index {
			break
		}

		// Otherwise, swap it
		input[i], input[swap_index] = input[swap_index], input[i]

	}

	return input

}

// from https://groups.google.com/group/golang-nuts/browse_thread/thread/571811b0ea0da610
func funcName(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
}

func dummy() {
	fmt.Fprintln(os.Stderr, "let me keep fmt and sys in the packages")
}

func TestIsUniq(t *testing.T) {
	altFuncs := []func(string) bool{isUniqNoDataStruct, isUniq}
	for _, myfunc := range altFuncs {

		if myfunc("abcdefghi") == false {
			t.Error("abcdefghi should return true, using", funcName(myfunc))
		}
		if myfunc("aaa") == true {
			t.Error("aaa should return false, using", funcName(myfunc))
		}
	}
}

func TestReverse(t *testing.T) {
	input := []rune("everything is amazing")
	output := reverse(input)
	output_string := string(output)
	if output_string != "gnizama si gnihtyreve" {
		t.Error("Reversed string does not match")
	}

}
