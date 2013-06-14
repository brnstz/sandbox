package chap1_test

import (
	"bytes"
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

// question 1.2
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

// question 1.3

func createRuneCountMap(input []rune) (map[rune]int, int) {
	count := map[rune]int{}
	totalCount := 0
	for _, c := range input {
		count[c]++
		totalCount++
	}

	return count, totalCount
}

func isPermutation(input1, input2 []rune) bool {
	map1, count1 := createRuneCountMap(input1)
	map2, count2 := createRuneCountMap(input2)

	// if counts don't match, they use different set of characters, so can't
	// possibly be a permutation. This also checks the case where map2
	// has more keys than map1 and a mismatch wouldn't be caught below
	if count1 != count2 {
		return false
	}

	// We have the same amount of keys, so let's check the count of each one
	for k, _ := range map1 {
		if map1[k] != map2[k] {
			return false
		}
	}

	// We made it!
	return true

}

// question 1.4
func urlencodeSpaces(input []rune, trueLen int) []rune {

	dataLen := len(input)
	toLoc := dataLen - 1
	bufferRemaining := dataLen - trueLen

	for toLoc >= 0 {
		fromLoc := toLoc - bufferRemaining

		if input[fromLoc] == ' ' {
			input[toLoc], input[toLoc-1], input[toLoc-2] = '0', '2', '%'
			toLoc -= 3
			bufferRemaining -= 2
		} else {
			input[toLoc] = input[fromLoc]
			toLoc -= 1
		}
	}

	return input
}

// question 1.5

func compressString(input string) string {
	var lastChar rune
	lastChar = 0

	charRunCount := 1
	oldStrLen := 0
	newStrLen := 0

	var buffer bytes.Buffer

	for _, curChar := range input {
		oldStrLen += 1

		if curChar == lastChar {
			charRunCount++

		} else {
			if lastChar != 0 {
				buffer.WriteString(fmt.Sprintf("%c%d", lastChar, charRunCount))
				charRunCount = 1
				newStrLen += 2
			}
		}

		lastChar = curChar
	}

	if lastChar != 0 {
		buffer.WriteString(fmt.Sprintf("%c%d", lastChar, charRunCount))
		newStrLen += 2
	}

	if newStrLen < oldStrLen {
		return buffer.String()
	} else {
		return input
	}
}

func printMatrix(matrix [][]int) {
	for x := range matrix {
		for y := range matrix[x] {
			fmt.Printf("%3d", matrix[x][y])
		}
		fmt.Println()
	}
}

// question 1.6 using a copy
func rotateMatrix(matrix [][]int) [][]int {
	m := len(matrix)

	newMatrix := make([][]int, m)
	for x := range newMatrix {
		newMatrix[x] = make([]int, m)
	}

	for x := range matrix {
		for y := range matrix[x] {
			new_y := x
			new_x := m - y - 1
			newMatrix[new_x][new_y] = matrix[x][y]
		}
	}

	return newMatrix
}

/*

func rotateMatrixInPlace(matrix []int, m int) {
	// First we transpose the matrix, algorithm from:
	// http://en.wikipedia.org/wiki/In-place_matrix_transposition#Square_matrices

	printMatrix(matrix, m)
	for x := 0; x <= m-2; x++ {
		for y := x + 1; y <= m-1; y++ {
			origIndex := mIndex(x, y, m)
			newIndex := mIndex(y, x, m)
			matrix[origIndex], matrix[newIndex] = matrix[newIndex], matrix[origIndex]
		}
	}
	printMatrix(matrix, m)

	for x := 0; x < m; x++ {
		for y := 0; y <= m/2; y++ {
			origIndex := mIndex(x, y, m)
			fmt.Println(x, y, matrix[origIndex])
			newIndex := mIndex(x, m-y-1, m)
			fmt.Println(x, m-y-1, matrix[newIndex])
			fmt.Println()
			matrix[origIndex], matrix[newIndex] = matrix[newIndex], matrix[origIndex]
		}
	}
	//printMatrix(matrix, m)

	// Next we fully rotate by swapping elements in the column. Only need
	// to do half of the ys, because we are swapping top to bottom and vice
	// versa
	// FIXME
		for x := 0; x < m; x++ {
			for y := 0; y <= m/2; y++ {
				origIndex := mIndex(x, y, m)
				newIndex := mIndex(x, m-y-1, m)
				matrix[origIndex], matrix[newIndex] = matrix[newIndex], matrix[origIndex]
			}
		}
}
*/

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

func TestPermutation(t *testing.T) {
	input1 := []rune("this is a permutation")
	input2 := []rune("thsi si a permuttaion")

	input3 := []rune("this is a permutationzzzzz")
	input4 := []rune("thsi si a permuttaion nope")

	if !isPermutation(input1, input2) {
		t.Error("input1 and input2 should return true for isPermutation")
	}

	if isPermutation(input3, input4) {
		t.Error("input3 and input4 should return false for isPermutation")
	}
}

func TestUrlEncodeSpaces(t *testing.T) {
	input1 := []rune("abcd ef hiaaa a      ")
	len1 := 15

	output1 := urlencodeSpaces(input1, len1)

	if string(output1) != "abcd%20ef%20hiaaa%20a" {
		t.Error("output incorrect")
	}
}

func TestCompress(t *testing.T) {
	if compressString("hello") != "hello" {
		t.Error("should return original string")
	}

	if compressString("hhhhhuuuuuuuklmmmmm") != "h5u7k1l1m5" {
		t.Error("should return compressed string")
	}
}

func TestMatrix(t *testing.T) {
	matrixArr := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	matrix := matrixArr[:][:]
	printMatrix(matrix)
	fmt.Println()
	printMatrix(rotateMatrix(matrix))
	/*
		newMatrix := rotateMatrix(origMatrix, m)

		newMatrixStr := fmt.Sprintf("%v", newMatrix)
		compareMatrixStr := fmt.Sprintf("%v", []int{3, 6, 9, 2, 5, 8, 1, 4, 7})
		if newMatrixStr != compareMatrixStr {
			t.Error("matrix rotate with copy failed")
		}

		rotateMatrixInPlace(origMatrix, m)
	*/
}
