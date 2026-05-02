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
	w.Print()

	var s strings.Builder
	for _, prim_id := range w.PrimOrder {
		switch w.TypeMap[prim_id] {
		case gf.PT_Enum:
			{
				enum_temp_struct := struct {
					Id    string
					Names []string
				}{
					Id:    prim_id,
					Names: w.Enums[prim_id].Value_names,
				}
				spew.Dump(enum_temp_struct)
				gnr.Templates[gf.PT_Enum].Execute(&s, enum_temp_struct)
			}

		case gf.PT_Enum2String:
			{
				values := make([]string, 0)
				enum_id := string(w.Enum2Strings[prim_id])
				enum_table := w.Tables[enum_id]
				for _, v := range enum_table.Rows {
					values = append(values, v[enum_table.Cols[0]])
				}
				enum_to_string_struct := struct {
					Id         string
					ValueNames []string
				}{
					Id:         enum_id,
					ValueNames: values,
				}

				spew.Dump(enum_to_string_struct)
				gnr.Templates[gf.PT_Enum2String].Execute(&s, enum_to_string_struct)
			}

		case gf.PT_Struct:
			{
				type Field struct {
					Type string
					Id string
				}
				fields := make([]Field, 0) 
				struct_prim := w.Structs[prim_id]
				for i := range struct_prim.Ids {
					fields = append(fields, Field{
						Type: struct_prim.Types[i],
						Id: struct_prim.Ids[i],
					})
				}
				struct_temp_struct := struct {
					Id string
					Fields []Field
				}{
					Id: prim_id,
					Fields: fields,
				}

				spew.Dump(struct_temp_struct)
				gnr.Templates[gf.PT_Struct].Execute(&s, struct_temp_struct)
			}

		case gf.PT_FuncTypes:
			{
				spew.Dump(w.FuncTypes[prim_id])
				gnr.Templates[gf.PT_FuncTypes].Execute(&s, w.FuncTypes[prim_id])
			}
		case gf.PT_FuncGlobals:
			{
				spew.Dump(w.FuncGlobals[prim_id])
				gnr.Templates[gf.PT_FuncGlobals].Execute(&s, w.FuncGlobals[prim_id])
			}

		case gf.PT_Custom:
			{
				spew.Dump(w.Customs[prim_id])
				s.WriteString(string(w.Customs[prim_id]))
			}
		}
	}

	fmt.Println(s.String())
}
