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
	for y := range matrix {
		for x := range matrix[y] {
			fmt.Printf("%5d", matrix[y][x])
		}
		fmt.Println()
	}
}

// question 1.6 using a copy
func rotateMatrix(matrix [][]int) [][]int {
	m := len(matrix)

	newMatrix := make([][]int, m)
	for y := range newMatrix {
		newMatrix[y] = make([]int, m)
	}

	for y := range matrix {
		for x := range matrix[y] {
			new_y := x
			new_x := m - y - 1
			newMatrix[new_y][new_x] = matrix[y][x]
		}
	}

	return newMatrix
}

// question 1.6 in place
func rotateMatrixInPlace(matrix [][]int) {
	// First we transpose the matrix, algorithm from:
	// http://en.wikipedia.org/wiki/In-place_matrix_transposition#Square_matrices
	m := len(matrix)
	for x := 0; x <= m-2; x++ {
		for y := x + 1; y <= m-1; y++ {
			matrix[y][x], matrix[x][y] = matrix[x][y], matrix[y][x]
		}
	}

	// Then we swap lower half of y's with upper half of y's
	// For example with an m=3 matrix, for the first col, we would just do:
	// swap matrix[0][0] with matrix[2][0]
	// matrix[1][0] remains (middle value)
	for x := range matrix {
		for y := 0; y < m/2; y++ {
			yp := m - y - 1
			matrix[x][y], matrix[x][yp] = matrix[x][yp], matrix[x][y]
		}
	}
}

// question 1.7
func zeroRowCol(matrix [][]int) {
	rowZs := map[int]bool{}
	colZs := map[int]bool{}

	// First create of map of all cols and rows that have a zero in them
	for col := range matrix {
		for row := range matrix[col] {
			if matrix[col][row] == 0 {
				rowZs[row] = true
				colZs[col] = true
			}
		}
	}

	// Then go back through matrix and set all values in those rows
	// and cols to 0
	for col := range matrix {
		for row := range matrix[col] {
			if rowZs[row] || colZs[col] {
				matrix[col][row] = 0
			}
		}
	}

}

// from https://groups.google.com/group/golang-nuts/browse_thread/thread/571811b0ea0da610
func funcName(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
}

// hacky way to see if slices are equal
func matrixEquals(a, b [][]int) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func dummy() {
	fmt.Fprintln(os.Stderr, "let me keep fmt and os in the packages")
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

	// indexed by y, x
	matrixArr := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	// for comparing results
	compareMatrixArr := [][]int{
		{7, 4, 1},
		{8, 5, 2},
		{9, 6, 3},
	}

	matrix := matrixArr[:][:]
	compareMatrix := compareMatrixArr[:][:]

	// Try algo that makes a copy
	rotatedCopyMatrix := rotateMatrix(matrix)

	if !matrixEquals(rotatedCopyMatrix, compareMatrix) {
		t.Error("rotatedCopyMatrix does not match")
	}

	// Try algo that rotates in place
	rotateMatrixInPlace(matrix)

	if !matrixEquals(matrix, compareMatrix) {
		t.Error("matrixMatrixInPlace() does not match")
	}
}

func TestZeroMatris(t *testing.T) {
	matrixArr := [][]int{
		{2, 3, 3, 5, 3, 5, 0},
		{5, 6, 62, 1, 76, 464, 1},
		{0, 445, 666, 66, 66, 76, 4},
		{22, 90, 32, 0, 8, 0, 123},
		{32, 3, 5, 44, 2, 13, 99},
		{1, 1, 1, 1, 1, 1, 8},
	}

	resultMatrixArr := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 6, 62, 0, 76, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 3, 5, 0, 2, 0, 0},
		{0, 1, 1, 0, 1, 0, 0},
	}

	matrix := matrixArr[:][:]
	resultMatrix := resultMatrixArr[:][:]

	zeroRowCol(matrix)

	if !matrixEquals(matrix, resultMatrix) {
		t.Error("zeroRowCol() did not succeed")
	}

}
