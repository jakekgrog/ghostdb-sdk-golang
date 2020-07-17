package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		t.Logf("TEST PASSED")
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func AssertDeepEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if reflect.DeepEqual(a, b) {
		t.Logf("TEST PASSED")
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func readFileByLine(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}