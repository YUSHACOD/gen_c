package gnrtr

import (
	"fmt"
	"strings"
	"text/template"
)

type StructField struct {
	Name string
	Type string
}

func GenStruct(name string, fields []StructField) (string, error) {

	const struct_field_template = "    {{.Type}} {{.Name}};\n"
	field_t := template.Must(template.New("field_t").Parse(struct_field_template))

	var res strings.Builder

	res.WriteString("typedef struct {\n")

	for _, field := range fields {
		err := field_t.Execute(&res, field)
		if err != nil {
			return "", fmt.Errorf("Error executing field %v : %v", field, err)
		}
	}

	fmt.Fprintf(&res, "} %s;\n", name)
	return res.String(), nil
}

type FuncType struct {
	Name   string
	Params string
	Args   string
	Return string
}

func GenFuncType(funcType FuncType) (string, error) {

	const type_template = "typedef {{.Return}} {{.Name}}_F_TYPE{{.Args}};\n"
	type_template_t := template.Must(template.New("type_template").Parse(type_template))

	var res strings.Builder
	err := type_template_t.Execute(&res, funcType)
	if err != nil {
		return "", fmt.Errorf("Error executing func type %v : %v", funcType, err)
	}

	return res.String(), nil
}

func GenHook(funcType FuncType) (string, error) {
	const hook_template = `
static {{.Return}} (WINAPI *og_{{.Name}}){{.Params}} = {{.Name}};
static {{.Return}} WINAPI hooked_{{.Name}}{{.Params}} {

    if(IsDebuggerPresent()) {
    	__debugbreak();
    }
	{{if eq .Return "void"}}
	TIME({ og_{{.Name}}{{.Args}}; });
	{{else}}
    {{.Return}} result;
    TIME({ result = og_{{.Name}}{{.Args}}; });

    if(IsDebuggerPresent()) {
    	__debugbreak();
    }

    return result;
	{{end}}
}
`
	hook_template_t := template.Must(template.New("hook_template").Parse(hook_template))

	var res strings.Builder
	err := hook_template_t.Execute(&res, funcType)
	if err != nil {
		return "", fmt.Errorf("Error executing func type %v : %v", funcType, err)
	}

	return res.String(), nil
}

func GenHookTable(funcs []FuncType) (string, error) {
	const table_element_template = "{&(void *&)og_{{.Name}}, (void *)hooked_{{.Name}}},\n"
	table_element_template_t := template.Must(template.New("element_template").Parse(table_element_template))

	var res strings.Builder
	res.WriteString("static Hook GLBL_hooks[] = {\n")

	for _, f := range funcs {
		err := table_element_template_t.Execute(&res, f)
		if err != nil {
			return "", fmt.Errorf("Error executing func element %v : %v", f, err)
		}
	}

	res.WriteString("};\n")
	return res.String(), nil
}
