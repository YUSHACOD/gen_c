package gnrtr

import (
// "fmt"
// "strings"
// "text/template"
)

//  templates : ---------------------------------------------------------------------- (section)  //

const struct_field_template = "    {{.Type}} {{.Name}};\n"

//  (section) ---------------------------------------------------------------------- : templates  //

//  write data structs : ------------------------------------------------------------- (section)  //

type Struct struct {
	typ []string
	ids []string
}

type Table struct {
	cols []string
	rows map[string][]string
}

type Enum struct {
	Id string
	vals []string
}

type GencWritables struct {
	tables []Table
	enums  []Enum
}

//  (section) ------------------------------------------------------------- : write data structs  //
