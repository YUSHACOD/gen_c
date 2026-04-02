package genc_fmt_test

import (
	"fmt"
	"testing"
	"os"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
)

func Test(test *testing.T) {
	fmt.Println("Testing")

	input, err := os.ReadFile("../template.genc")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(input))

	t := gf.NewTokenizer(string(input))
	for token := t.NextToken(); token.Typ != gf.Eof; token = t.NextToken() {
		token.Print()
		fmt.Println()
	}
}


