package genc_fmt_test

import (
	"fmt"
	"testing"
	"os"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
	// gen "github.com/YUSHACOD/gen_c/gnrtr"
)

func Test(test *testing.T) {
	fmt.Println("Testing")

	input, err := os.ReadFile("../template.genc")
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(string(input))

	t := gf.NewTokenizer(string(input))
	genc := gf.ParseGenc(t)
	for k,v  := range genc.Primitives {
		fmt.Println("Primitive Id: ", k)
		fmt.Println("Primitive Val:")
		v.Print()
	}

	w := gf.GenerateWritables(genc)
	w.Print()
}


