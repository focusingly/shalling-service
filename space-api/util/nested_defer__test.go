package util_test

import (
	"fmt"
	"testing"
)

func nestedDeferTest() {
	defer func() {
		defer func() {
			fmt.Println(2)
		}()
		fmt.Println(1)
	}()

	func() {
		defer func() {
			fmt.Println("a")
		}()
		func() {
			fmt.Println("x")
		}()
		defer func() {
			fmt.Println("b")
		}()
	}()

	fmt.Println("c")
}

func TestNestedDefer(t *testing.T) {
	nestedDeferTest()
}
