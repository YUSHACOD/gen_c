package genc_fmt

import (
	// "fmt"
	// "strings"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"unicode"

	"text/template"

	"github.com/olekukonko/tablewriter"
)

//  write data structs : ------------------------------------------------------------- (section)  //

type Table struct {
	Cols []string
	Rows []map[string]string
}

type Enum struct {
	Value_names []string
}

type Struct struct {
	Types []string
	Ids   []string
}

type FuncType struct {
	Name string
	Ret  string
	Args string
}

type FuncGlobal struct {
	Type string
	Name string
}

type Enum2String string

type Custom string

type GencWritables struct {
	curr_req map[string]Table
	req_len  uint32

	Tables      map[string]Table
	Enums       map[string]Enum
	Structs     map[string]Struct
	FuncTypes   map[string][]FuncType
	FuncGlobals map[string][]FuncGlobal
	Customs     map[string]Custom

	Enum2Strings map[string]Enum2String


	TypeMap   map[string]PrimitiveType
	PrimOrder []string
}

type WriteType string

const (
	C_file   WriteType = "c_file"
	H_file   WriteType = "h_file"
	CPP_file WriteType = "cpp_file"
	HPP_file WriteType = "hpp_file"
)

type WriteCommand struct {
	typ         WriteType
	write_order []string
}

//  (section) ------------------------------------------------------------- : write data structs  //

// writables print helpers : -------------------------------------------------------- (section)  //
func (t Table) Print() {

	colKeys := t.Cols

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(colKeys)

	for i := range len(t.Rows) {

		row := make([]string, len(colKeys))

		for j, key := range colKeys {
			if i < len(t.Rows) {
				row[j] = t.Rows[i][key]
			}
		}

		table.Append(row)
	}

	table.Render()
}

func (e Enum) Print() {
	fmt.Println("Value: ")
	for _, s := range e.Value_names {
		fmt.Println(s)
	}
}

func (s Struct) Print() {
	fmt.Println("Struct Fields: ")
	for idx := range s.Types {
		fmt.Printf("type: %s, identifier: %s\n", s.Types[idx], s.Ids[idx])
	}
}

func (f FuncType) Print() {
	fmt.Printf("Return: %s,  Identifier: %s, Args: %s\n", f.Ret, f.Name, f.Args)
}

func (f FuncGlobal) Print() {
	fmt.Printf("Type: %s,  Identifier: %s\n", f.Type, f.Name)
}

func (w GencWritables) Print() {

	fmt.Println(Blue)
	fmt.Println("Writables")

	for id, t := range w.Tables {
		fmt.Printf("Table %s ->\n", id)
		t.Print()
	}

	fmt.Println()
	for id, e := range w.Enums {
		fmt.Printf("Enum %s ->\n", id)
		e.Print()
	}

	fmt.Println()
	for id, s := range w.Structs {
		fmt.Printf("Struct %s ->\n", id)
		s.Print()
	}

	fmt.Println()
	for id, fs := range w.FuncTypes {
		fmt.Printf("FuncTypes %s ->\n", id)
		for _, f := range fs {
			f.Print()
		}
	}

	fmt.Println()
	for id, fs := range w.FuncGlobals {
		fmt.Printf("FuncGlobals %s ->\n", id)
		for _, f := range fs {
			f.Print()
		}
	}

	fmt.Println()
	fmt.Println("Enum2String tables or func don't know for now")
	for k, v := range w.Enum2Strings {
		fmt.Printf("%s : %s\n", k, v)
	}

	fmt.Println()
	fmt.Println("Custom Template Expansions ->")
	for id, t := range w.Customs {
		fmt.Println("Id : ", id)
		fmt.Println(t)
	}

	fmt.Println()

	fmt.Println(Yellow)
	fmt.Println("Type Map ->")
	for k, v := range w.TypeMap {
		fmt.Printf(" %s : %s\n", k, v)
	}

	fmt.Print(Reset)

}

//  (section) -------------------------------------------------------- : writables print helpers  //

//  expression evaluation : ---------------------------------------------------------- (section)  //

func uppercase(s string) string {
	return strings.ToUpper(s)
}

func lowercase(s string) string {
	return strings.ToLower(s)
}

// hello_world -> HelloWorld
func snake2pascal(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// hello_world -> helloWorld
func snake2camel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(p)
		} else if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// HelloWorld -> hello_world
func pascal2snake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// HelloWorld -> helloWorld
func pascal2camel(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// helloWorld -> hello_world
func camel2snake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// helloWorld -> HelloWorld
func camel2pascal(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

var OpMap template.FuncMap

func CompFuncMap() {
	OpMap = template.FuncMap{
		"uppercase":     uppercase,
		"lowercaseorld": lowercase,
		"snake2pascal":  snake2pascal,
		"snake2camel":   snake2camel,
		"pascal2snake":  pascal2snake,
		"pascal2camel":  pascal2camel,
		"camel2snake":   camel2snake,
		"camel2pascal":  camel2pascal,
	}
}

func (e *Expression) evaluate(idx uint32, w *GencWritables) string {
	var res string

	switch e.typ {

	case ET_Value:
		return e.value

	case ET_Array:
		log.Fatalf("This expresssion shouldn't have been a Array type %v", e)

	case ET_ColId:
		if t, ok := w.curr_req[e.value]; ok {
			return t.Rows[idx][e.arr[0].value]
		} else {
			log.Fatalf("This table alias is not present in the current requirement cache %s, %v",
				e.value, w.curr_req)
		}

	case ET_PrimIdAlias:

	case ET_OP_Concat:
		res := ""
		for _, exp := range e.arr {
			res = res + exp.evaluate(idx, w)
		}
		return res

	case ET_OP_Uppercase:
		return uppercase(e.arr[0].evaluate(idx, w))

	case ET_OP_Lowercase:
		return lowercase(e.arr[0].evaluate(idx, w))

	case ET_OP_Snake2Pascal:
		return snake2pascal(e.arr[0].evaluate(idx, w))

	case ET_OP_Snake2Camel:
		return snake2camel(e.arr[0].evaluate(idx, w))

	case ET_OP_Pascal2Snake:
		return pascal2snake(e.arr[0].evaluate(idx, w))

	case ET_OP_Pascal2Camel:
		return pascal2camel(e.arr[0].evaluate(idx, w))

	case ET_OP_Camel2Snake:
		return camel2snake(e.arr[0].evaluate(idx, w))

	case ET_OP_Camel2Pascal:
		return camel2pascal(e.arr[0].evaluate(idx, w))

	}

	return res
}

func (e *Expression) evaluateArray(w *GencWritables) []string {

	if e.typ != ET_Array {
		log.Panicf("This is not a Array Expression %s", e.typ)
	}

	res := make([]string, 0)

	for _, exp := range e.arr {
		res = append(res, exp.evaluate(0, w))
	}

	return res
}

//  (section) ---------------------------------------------------------- : expression evaluation  //

// gen writables : ------------------------------------------------------------------ (section)  //

func (w *GencWritables) genTable(p Primitive) Table {
	table := Table{
		Rows: make([]map[string]string, 0),
	}

	for i := range 2 {
		field := p.fields[i]
		switch field.typ {
		case FT_Table_Cols:
			{
				for _, exp := range field.val.arr {
					table.Cols = append(table.Cols, exp.evaluate(0, w))
				}
			}

		case FT_Table_Rows:
			{
				for _, exp := range field.val.arr {
					row := exp.evaluateArray(w)

					t_row := make(map[string]string)
					for i, row_elem := range row {
						t_row[table.Cols[i]] = row_elem
					}

					table.Rows = append(
						table.Rows,
						t_row,
					)
				}
			}

		default:
			log.Panicf("This is invalid field type for table %s", field.typ)
		}
	}

	return table
}

func (w *GencWritables) generateRequiresTable(req SubPrimitive) {

	tables := make(map[string]Table)
	var ln uint32 = math.MaxUint32

	for _, exp := range req.args {
		for _, exp := range exp.arr {
			table_id := exp.arr[0].evaluate(0, w)
			table_alias := exp.arr[1].evaluate(0, w)

			if t, ok := w.Tables[table_id]; ok {
				tables[table_alias] = t
				ln = min(ln, uint32(len(t.Rows)))
			} else {
				log.Fatalf("This table id is not currently present in the list of tables %s => %v",
					table_id, w.Tables)
			}
		}
	}

	w.curr_req = tables
	w.req_len = ln
}

func (w *GencWritables) genEnum(p Primitive) (Enum, Table) {
	var enum Enum

	var req SubPrimitive
	for _, s_prim := range p.sub_prims {
		if s_prim.typ == ST_Requires {
			req = s_prim
		}
	}
	w.generateRequiresTable(req)

	for _, field := range p.fields {
		switch field.typ {

		case FT_Enum_ValueName:
			for idx := range w.req_len {
				enum.Value_names = append(enum.Value_names, field.val.evaluate(idx, w))
			}

		default:
			log.Fatalf("Unkown Field found")
		}
	}

	enum_t_rows := make([]map[string]string, 0)
	for _, value_name := range enum.Value_names {
		enum_t_rows = append(enum_t_rows, map[string]string{
			"value_name": value_name,
		})
	}
	enum_table := Table{
		Cols: []string{"value_name"},
		Rows: enum_t_rows,
	}

	return enum, enum_table
}

func (w *GencWritables) genEnum2String(p Primitive) Enum2String {

	res := ""
	for _, field := range p.fields {
		switch field.typ {

		case FT_Enum2String_Enum:
			for idx := range w.req_len {
				res = field.val.evaluate(idx, w)
			}

		default:
			log.Fatalf("Unkown Field found")
		}
	}

	if t, ok := w.TypeMap[res]; ok {
		if t != PT_Enum {
			log.Fatalf("This enum prim id provided doesnt point to a enum primitive %s", res)
		}
	}

	return Enum2String(res)
}

func (w *GencWritables) genStruct(p Primitive) (Struct, Table) {

	var req SubPrimitive
	for _, s_prim := range p.sub_prims {
		if s_prim.typ == ST_Requires {
			req = s_prim
		}
	}
	w.generateRequiresTable(req)

	struc := Struct{}

	for _, field := range p.fields {
		switch field.typ {

		case FT_Struct_FieldTypes:
			for idx := range w.req_len {
				struc.Types = append(struc.Types, field.val.evaluate(idx, w))
			}

		case FT_Struct_FieldIds:
			for idx := range w.req_len {
				struc.Ids = append(struc.Ids, field.val.evaluate(idx, w))
			}

		default:
			log.Fatalf("Unkown Field found")
		}
	}

	struct_table := Table{
		Cols: []string{"field_types", "field_ids"},
	}
	for idx := range len(struc.Types) {
		struct_table.Rows = append(struct_table.Rows, map[string]string{
			"field_types": struc.Types[idx],
			"field_ids":   struc.Ids[idx],
		})
	}
	return struc, struct_table
}

func (w *GencWritables) genFuncTypes(p Primitive) ([]FuncType, Table) {

	var req SubPrimitive
	for _, s_prim := range p.sub_prims {
		if s_prim.typ == ST_Requires {
			req = s_prim
		}
	}
	w.generateRequiresTable(req)

	func_types := make([]FuncType, w.req_len)

	for _, field := range p.fields {
		switch field.typ {

		case FT_FuncTypes_Args:
			for idx := range w.req_len {
				func_types[idx].Args = field.val.evaluate(idx, w)
			}

		case FT_FuncTypes_Identifier:
			for idx := range w.req_len {
				func_types[idx].Name = field.val.evaluate(idx, w)
			}

		case FT_FuncTypes_Ret:
			for idx := range w.req_len {
				func_types[idx].Ret = field.val.evaluate(idx, w)
			}
		}
	}

	func_types_table := Table{
		Cols: []string{"args", "identifier", "ret"},
	}
	for _, ft := range func_types {
		func_types_table.Rows = append(func_types_table.Rows, map[string]string{
			"args":       ft.Args,
			"identifier": ft.Name,
			"ret":        ft.Ret,
		})
	}

	return func_types, func_types_table
}

func (w *GencWritables) genFuncGlobals(p Primitive) ([]FuncGlobal, Table) {

	var req SubPrimitive
	for _, s_prim := range p.sub_prims {
		if s_prim.typ == ST_Requires {
			req = s_prim
		}
	}
	w.generateRequiresTable(req)

	func_globals := make([]FuncGlobal, w.req_len)
	for _, field := range p.fields {
		switch field.typ {

		case FT_FuncGlobals_Typ:
			for idx := range w.req_len {
				func_globals[idx].Type = field.val.evaluate(idx, w)
			}

		case FT_FuncGlobals_Identifier:
			for idx := range w.req_len {
				func_globals[idx].Name = field.val.evaluate(idx, w)
			}
		}
	}

	func_globals_table := Table{
		Cols: []string{"type", "identifier"},
	}
	for _, ft := range func_globals {
		func_globals_table.Rows = append(func_globals_table.Rows, map[string]string{
			"type":       ft.Type,
			"identifier": ft.Name,
		})
	}

	return func_globals, func_globals_table
}

func (w *GencWritables) expandCustom(p Primitive) Custom {

	var req SubPrimitive
	for _, s_prim := range p.sub_prims {
		if s_prim.typ == ST_Requires {
			req = s_prim
		}
	}
	w.generateRequiresTable(req)

	res := strings.Builder{}
	field := p.fields[0]

	if field.typ == FT_Custom_Template {

		temp := field.val.evaluate(0, w)
		tmpl := template.Must(template.New("custom").Funcs(OpMap).Parse(temp))

		data := make([]map[string]map[string]string, w.req_len)
		for idx := range data {
			dat := make(map[string]map[string]string)
			for k, v := range w.curr_req {
				dat[k] = v.Rows[idx]
			}
			data[idx] = dat
		}

		err := tmpl.Execute(&res, data)

		if err != nil {
			fmt.Print(
				Red,
				"Some error occured executing custom template \n",
				temp,
				err,
				"\n",
				Reset)
		}
	}

	return Custom(strings.ReplaceAll(res.String(), "\r", ""))
}


func GenerateWritables(genc *GenC) GencWritables {

	CompFuncMap()

	//  gen writables core : --------------------------------------------------------- (section)  //
	wrtb := GencWritables{
		Tables:       make(map[string]Table),
		Enums:        make(map[string]Enum),
		Enum2Strings: make(map[string]Enum2String),
		Structs:      make(map[string]Struct),
		FuncTypes:    make(map[string][]FuncType),
		FuncGlobals:  make(map[string][]FuncGlobal),
		Customs:      make(map[string]Custom),

		TypeMap: make(map[string]PrimitiveType),
	}

	for _, id := range genc.Ids {
		prim := genc.Primitives[id]

		wrtb.PrimOrder = append(wrtb.PrimOrder, id)
		switch prim.Typ {

		case PT_Table:
			{
				wrtb.Tables[id] = wrtb.genTable(prim)
				wrtb.TypeMap[id] = PT_Table
			}

		case PT_Enum:
			{
				wrtb.Enums[id], wrtb.Tables[id] = wrtb.genEnum(prim)
				wrtb.TypeMap[id] = PT_Enum
			}

		case PT_Enum2String:
			{
				wrtb.Enum2Strings[id] = wrtb.genEnum2String(prim)
				wrtb.TypeMap[id] = PT_Enum2String
			}

		case PT_Struct:
			{
				wrtb.Structs[id], wrtb.Tables[id] = wrtb.genStruct(prim)
				wrtb.TypeMap[id] = PT_Struct
			}

		case PT_FuncTypes:
			{
				wrtb.FuncTypes[id], wrtb.Tables[id] = wrtb.genFuncTypes(prim)
				wrtb.TypeMap[id] = PT_FuncTypes
			}

		case PT_FuncGlobals:
			{
				wrtb.FuncGlobals[id], wrtb.Tables[id] = wrtb.genFuncGlobals(prim)
				wrtb.TypeMap[id] = PT_FuncGlobals
			}

		case PT_Custom:
			{
				wrtb.Customs[id] = wrtb.expandCustom(prim)
				wrtb.TypeMap[id] = PT_Custom
			}


		default:
			{
				log.Fatalf("This ain't no primitive %s", string(genc.Primitives[id].Typ))
			}

		}
	}

	return wrtb
}

//  (section) ------------------------------------------------------------------ : gen writables  //
