package gnrtr_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
	gnr "github.com/YUSHACOD/gen_c/gnrtr"
	"github.com/davecgh/go-spew/spew"
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

	var s strings.Builder
	for k, typ := range w.TypeMap {
		switch typ {
		case gf.PT_Enum:
			{
				enum_temp_struct := struct {
					Id    string
					Names []string
				}{
					Id:    k,
					Names: w.Enums[k].Value_names,
				}
				spew.Dump(enum_temp_struct)
				gnr.Templates[gf.PT_Enum].Execute(&s, enum_temp_struct)
			}
		}
	}

	fmt.Println(s.String())
}
