package main

import (
	"fmt"
	"log"
	"os"
	fp "path/filepath"
	"testing"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
	gnr "github.com/YUSHACOD/gen_c/gnrtr"
)

func Test(_ *testing.T) {
	fmt.Println("Testing")

	gnr.InitGen()

	root := "." // or any directory

	err := fp.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if len(path) > 5 {

			last5 := path[len(path)-5:]
			dir := fp.Dir(path)

			switch last5 {
			case ".genc":
				fallthrough
			case ".genh":
				fallthrough
			case ".gcpp":
				fallthrough
			case ".ghpp":
				{
					dir = fp.Join(dir, "generated")
					fmt.Println(dir)
					os.Mkdir(dir, 0755)
					ext := "c"
					if path[len(path)-1] == 'p' {
						ext = path[len(path)-3:]
					} else {
						ext = path[len(path)-1:]
					}
					fmt.Println(ext)

					input, err := os.ReadFile(path)
					if err != nil {
						fmt.Println(err)
					}

					// fmt.Println(string(input))

					t := gf.NewTokenizer(string(input))
					genc := gf.ParseGenc(t)

					w := gf.GenerateWritables(genc)

					// w.Print()

					gen := gnr.Gen(w)
					fmt.Println(gen)

					os.WriteFile(fp.Join(dir, path[:len(path)-5]+ "." + ext), []byte(gen), 0644)

				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}
