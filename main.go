package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YUSHACOD/gen_c/genc_fmt"
	// "github.com/YUSHACOD/gen_c/gnrtr"
)

func main() {
	fmt.Printf("Generating c code\n")

	genc_file := os.Args[1]
	content, err := os.ReadFile(genc_file)
	if err != nil {
		log.Fatalf("Input file reading error => %v", err)
	}

	// fmt.Printf("%s\n", content)

	genc_fmt.ParseGenc(content)
	// fmt.Printf("%v\n", gen_commands)
}
