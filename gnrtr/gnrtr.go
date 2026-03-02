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

const struct_field_template = "    {{.Type}} {{.Name}};\n"

func GenStruct(name string, fields []StructField) (string, error) {
	field_t := template.Must(template.New("field_t").Parse(struct_field_template))


	var res strings.Builder

	res.WriteString("typedef struct {\n")

	for _, field := range fields {
		err := field_t.Execute(&res, field)
		if err != nil {
			return "", fmt.Errorf("Error executing field %v template: %v", field, err)
		}
	}

	fmt.Fprintf(&res, "} %s;\n", name)
	return res.String(), nil
}
