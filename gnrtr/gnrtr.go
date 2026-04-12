package gnrtr

import (
	// "fmt"
	"sort"
	// "strings"
	"os"
	// "text/template"
	"github.com/olekukonko/tablewriter"
)

//  templates : ---------------------------------------------------------------------- (section)  //

//  (section) ---------------------------------------------------------------------- : templates  //

//  write data structs : ------------------------------------------------------------- (section)  //

type Table struct {
	Cols []string
	Rows map[string][]string
}

type Enum struct {
	Id   string
	Vals []string
}

type Struct struct {
	Typ []string
	Ids []string
}

type FuncType struct {
	Id   string
	Ret  string
	Args string
}

type FuncGlobal struct {
	Typ string
	Id  string
}

type Custom string

type GencWritables struct {
	Tables      map[string]Table
	Enums       map[string]Enum
	FuncTypes   map[string]FuncType
	FuncGlobals map[string]FuncGlobal
	Customs     map[string]Custom
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

func (t Table) Print() {
	colKeys := make([]string, 0, len(t.Rows))
	for k := range t.Rows {
		colKeys = append(colKeys, k)
	}
	sort.Strings(colKeys)

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(colKeys)

	for i := range t.Cols {
		row := make([]string, len(colKeys))
		for j, key := range colKeys {
			if i < len(t.Rows[key]) {
				row[j] = t.Rows[key][i]
			}
		}
		table.Append(row)
	}

	table.Render()
}

func (w GencWritables) Print() {
	for _, t := range w.Tables {
		t.Print()
	}
}
