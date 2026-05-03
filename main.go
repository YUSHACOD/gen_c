package main

import (
	"fmt"
	"log"
	"os"
	fp "path/filepath"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
	gnr "github.com/YUSHACOD/gen_c/gnrtr"
)

func run() {
	gnr.InitGen()

	root := "."

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

					os.Mkdir(dir, 0755)
					ext := "c"

					if path[len(path)-1] == 'p' {
						ext = path[len(path)-3:]
					} else {
						ext = path[len(path)-1:]
					}

					input, err := os.ReadFile(path)
					if err != nil {
						fmt.Println(err)
					}

					t := gf.NewTokenizer(string(input))
					genc := gf.ParseGenc(t)

					w := gf.GenerateWritables(genc)

					gen := gnr.Gen(w)
					// fmt.Println(gen)
					name := fp.Base(path)
					fmt.Println("file: ", fp.Join(dir, name[:len(name)-5]+"."+ext))

					os.WriteFile(fp.Join(dir, fp.Base(name)[:len(name)-5]+"."+ext), []byte(gen), 0644)

				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Gencing")
	run()
}
