package chap1_test

import (
	"reflect"
	"runtime"
	"sort"
	"testing"
)

// question 1
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
	runes := SortableRunes(input)
	sort.Sort(runes)
	var last rune = 0

	for _, c := range input {
		if last == c {
			return false
		}
		last = c
	}

	return true
}

// from https://groups.google.com/group/golang-nuts/browse_thread/thread/571811b0ea0da610
func funcName(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
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
