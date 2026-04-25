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
const {{.id}}_table = {
};
`,

	gf.PT_Struct: `
`,

	gf.PT_FuncTypes: `
`,

	gf.PT_FuncGlobals: `
`,

	gf.PT_Custom: `
`,
}

var Templates map[gf.PrimitiveType]*template.Template = map[gf.PrimitiveType]*template.Template{}

//  (section) ---------------------------------------------------------------------- : templates  //

func compTempaltes() {

	for _, s := range []gf.PrimitiveType{
		gf.PT_Enum,
		// gf.PT_Enum2String,
		// gf.PT_Struct,
		// gf.PT_FuncTypes,
		// gf.PT_FuncGlobals,
		// gf.PT_Custom,
		// gf.PT_GenCFile,
		// gf.PT_GenHFile,
		// gf.PT_GenCPPFile,
		// gf.PT_GenHPPFile,
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
