package main

import (
	"fmt"

	"github.com/YUSHACOD/gen_c/gnrtr"
)

func main() {
	fmt.Printf("Generating c code\n")

	fields := make([]gnrtr.StructField, 3)
	fields[0].Name = "X"
	fields[0].Type = "int"

	fields[1].Name = "Y"
	fields[1].Type = "int"

	fields[2].Name = "Z"
	fields[2].Type = "int"

	struct_string, err := gnrtr.GenStruct("Point", fields)
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", struct_string)
	}
}
