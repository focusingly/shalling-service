package util_test

import (
	"fmt"
	"testing"

	"github.com/yanyiwu/gojieba"
)

func TestWorCut(*testing.T) {
	var s string
	x := gojieba.NewJieba()
	s = "# 现在是凌晨2点多, 窗外下着连绵不断的雨"
	w := x.Tokenize(s, gojieba.SearchMode, true)
	for _, w := range w {
		fmt.Println("|", w.Str, "|")
	}
	x.Free()
}
