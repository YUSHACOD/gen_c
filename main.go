package main

import (
	"fmt"
	"os"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
	gnr "github.com/YUSHACOD/gen_c/gnrtr"
)

func main() {
	gnr.InitGen()

	fmt.Println("Testing")

	input, err := os.ReadFile("./template.genc")
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
	gen := gnr.Gen(w)
	fmt.Println(gen)
	os.WriteFile("generated/generated.c", []byte(gen), 0644)
}
