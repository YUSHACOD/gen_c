package genc_fmt_test

import (
	"fmt"
	"testing"

	"github.com/YUSHACOD/gen_c/genc_fmt"
)

func Test(test *testing.T) {
	fmt.Println("Testing")

	input := `
@table(Funcs) {
	@fields(name type_name args ret)
	{ Add                add_op_ft                  ` + "`" + `int x, int y` + "`" + `         int }
	{ Sub                sub_op_ft                  ` + "`" + `int x, int y` + "`" + `         int }
	{ Mul                mul_op_ft                  ` + "`" + `int x, int y` + "`" + `         int }
}
`
	fmt.Println("Input :", input)
	tokens := genc_fmt.Tokenize(input)
	for _, t := range tokens {
		t.Print()
	}
}
