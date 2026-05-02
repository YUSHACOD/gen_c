package gnrtr_test

import (
	"fmt"
	"os"
	"testing"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
	gnr "github.com/YUSHACOD/gen_c/gnrtr"
)

func Test(_ *testing.T) {
	gnr.InitGen()

	fmt.Println("Testing")

	input, err := os.ReadFile("../template.genc")
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(string(input))

	t := gf.NewTokenizer(string(input))
	genc := gf.ParseGenc(t)
	// for k,v  := range genc.Primitives {
	// 	fmt.Println("Primitive Id: ", k)
	// 	fmt.Println("Primitive Val:")
	// 	v.Print()
	// }

	w := gf.GenerateWritables(genc)
	fmt.Println(gnr.Gen(w))
}
