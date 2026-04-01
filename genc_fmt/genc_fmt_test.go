package genc_fmt_test

import (
	"fmt"
	"testing"
	"os"

	"github.com/YUSHACOD/gen_c/genc_fmt"
)

func Test(test *testing.T) {
	fmt.Println("Testing")

	input, err := os.ReadFile("../template.genc")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	tokens := genc_fmt.Tokenize(string(input))
	for _, t := range tokens {
		t.Print()
	}
}
