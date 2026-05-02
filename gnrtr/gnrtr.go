package gnrtr

import (
	"log"
	"text/template"

	gf "github.com/YUSHACOD/gen_c/genc_fmt"
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

func Gen(w gf.GencWritables) {
}
