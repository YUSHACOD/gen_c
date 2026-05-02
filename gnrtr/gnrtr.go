package gnrtr

import (
	"log"
	"strings"
	"text/template"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
	"github.com/davecgh/go-spew/spew"
)

//  templates : ---------------------------------------------------------------------- (section)  //

var TemplateString map[gf.PrimitiveType]string = map[gf.PrimitiveType]string{

	gf.PT_Enum: `
typedef enum {
    {{range .Names}}{{.}},
    {{end -}}
    EnumCount({{.Id}}),
} {{.Id}};
`,

	gf.PT_Enum2String: `
const char* {{.Id}}_to_string[] = {
    {{range .ValueNames}} "{{.}}",
    {{end -}}
};
`,

	gf.PT_Struct: `
typedef struct {
    {{range .Fields}}{{.Type}} {{.Id}};
    {{end -}}
} {{.Id}};
`,

	gf.PT_FuncTypes: `
{{range .}}typedef {{.Ret}} {{.Name}}({{.Args}});
{{end -}}
`,

	gf.PT_FuncGlobals: `
{{range .}}global {{.Type}}* {{.Name}};
{{end -}}
`,

}

var Templates map[gf.PrimitiveType]*template.Template = map[gf.PrimitiveType]*template.Template{}

//  (section) ---------------------------------------------------------------------- : templates  //

func compTempaltes() {

	for _, s := range []gf.PrimitiveType{
		gf.PT_Enum,
		gf.PT_Enum2String,
		gf.PT_Struct,
		gf.PT_FuncTypes,
		gf.PT_FuncGlobals,
	} {
		t, err := template.New(string(s)).Parse(TemplateString[s])
		if err != nil {
			log.Fatalf("Template CompError %v, for template prim %s", err, s)
		}
		Templates[s] = t
	}

}

func InitGen() {
	compTempaltes()
}

func Gen(w gf.GencWritables) string {

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
				Templates[gf.PT_Enum].Execute(&s, enum_temp_struct)
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
				Templates[gf.PT_Enum2String].Execute(&s, enum_to_string_struct)
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
				Templates[gf.PT_Struct].Execute(&s, struct_temp_struct)
			}

		case gf.PT_FuncTypes:
			{
				spew.Dump(w.FuncTypes[prim_id])
				Templates[gf.PT_FuncTypes].Execute(&s, w.FuncTypes[prim_id])
			}
		case gf.PT_FuncGlobals:
			{
				spew.Dump(w.FuncGlobals[prim_id])
				Templates[gf.PT_FuncGlobals].Execute(&s, w.FuncGlobals[prim_id])
			}

		case gf.PT_Custom:
			{
				spew.Dump(w.Customs[prim_id])
				s.WriteString(string(w.Customs[prim_id]))
			}
		}
	}
	
	return s.String()
}
