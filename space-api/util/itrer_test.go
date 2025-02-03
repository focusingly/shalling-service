package util_test

import (
	"fmt"
	"iter"
	"slices"
	"space-api/util"
	"testing"
)

func TestGetSerialResume(*testing.T) {
	var i int32
	it := util.GetSerial(func() int32 {
		tmp := i
		i++
		return tmp
	}, 10)

	for v := range it {
		fmt.Println(v)
		if v > 4 {
			break
		}
	}
	println()

	for v := range it {
		fmt.Println(v)
	}
}

func TestGetSerialNextFn(t *testing.T) {
	var i int32
	it := util.GetSerial(func() int32 {
		tmp := i
		i++
		return tmp
	}, 10)

	next, stop := iter.Pull(it)
	defer stop()
	cols := []int32{}
	for i := 0; i < 4; i++ {
		val, _ := next()
		cols = append(cols, val)
	}
	stop()

	for {
		val, ok := next()
		if !ok {
			break
		}
		cols = append(cols, val)
	}

	if len(cols) != 4 {
		t.Fatalf("excepted got 4 nums, but got: %d --> %v", len(cols), cols)
	} else {
		fmt.Println(cols)
	}
}

func TestGetSerialRange(t *testing.T) {
	var i int32
	it := util.GetSerial(func() int32 {
		tmp := i
		i++
		return tmp
	}, 10)

	bf := slices.Collect(it)
	if len(bf) != 10 {
		t.Fatalf("expected %d nums, but got %d --> %v", 10, len(bf), bf)
	}

	fmt.Println(bf)
}
