package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"database/sql"

	"github.com/YUSHACOD/gen_c/gnrtr"

	"github.com/hashicorp/go-set/v3"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(file_path string) (*sql.DB, func() error) {
	db, err := sql.Open("sqlite3", file_path)
	if err != nil {
		log.Panicln("Cannot open ntdocs.db :", err)
	}
	return db, db.Close
}

func filterFuncSig(sig *gnrtr.FuncType) {

	firstWord := func(s string) string {
		fields := strings.Fields(s)
		if len(fields) == 0 {
			return ""
		}
		return fields[0]
	}

	removeBracketed := func(s string) string {
		re := regexp.MustCompile(`\[[^\]]*\]`)
		return strings.Trim(re.ReplaceAllString(s, ""), " ")
	}

	paramsToArgs := func(params string) string {

		parts := strings.Split(params, ",")
		args := make([]string, 0, len(parts))

		for _, p := range parts {
			fields := strings.Fields(strings.TrimSpace(p))
			if len(fields) == 0 {
				continue
			}

			name := fields[len(fields)-1]
			name = strings.TrimLeft(name, "*") // handle "*c"
			args = append(args, name)
		}

		return strings.Join(args, ", ")
	}

	sig.Params = removeBracketed(sig.Params)
	sig.Args = paramsToArgs(sig.Params)
	sig.Header = firstWord(sig.Header)
}

func GetFuncSigsInDll(db *sql.DB, dll_name string) ([]gnrtr.FuncType, error) {

	res := []gnrtr.FuncType{}

	query_string := fmt.Sprintf(
		"SELECT FunctionSignatures.name, FunctionSignatures.parameters, FunctionSignatures.ret, FunctionSignatures.header FROM FunctionSignatures WHERE FunctionSignatures.dll like %s ;",
		fmt.Sprintf("'%%%s%%'", dll_name),
	)
	func_names_stmnt, err := db.Prepare(query_string)
	if err != nil {
		return res, fmt.Errorf("Failed to prepare the function names query, due to: %v", err)
	}
	defer func_names_stmnt.Close()

	{ // Executing the query and appending the signatures to the result list
		func_name_query, er := func_names_stmnt.Query()
		if er != nil {
			return res, fmt.Errorf("query of func names failed for %s : %v", dll_name, err)
		}
		defer func_name_query.Close()

		for func_name_query.Next() {
			err = func_name_query.Err()
			if err != nil {
				fmt.Printf("iteration error %v\n", err)
			}

			var func_sig gnrtr.FuncType
			err := func_name_query.Scan(
				&func_sig.Name,
				&func_sig.Params,
				&func_sig.Return,
				&func_sig.Header,
			)
			filterFuncSig(&func_sig)
			if err != nil {

				fmt.Printf("error scaning row : %v\n", err)
			} else {
				res = append(res, func_sig)
			}
		}
		err = func_name_query.Err()
		if err != nil {
			fmt.Printf("iteration error %v\n", err)
		}
	}

	return res, nil
}

func main() {

	db, dbclose := OpenDB("ntdocs.sqlite3")
	defer dbclose()

	dll_name := "Kernel32.dll"
	func_sigs, err := GetFuncSigsInDll(db, dll_name)
	headers := []string{}
	if err != nil {
		fmt.Printf("Error retreiving all funcs in dll %s: %v", dll_name, err)
	} else {
		for _, sig := range func_sigs {
			// fmt.Printf("(%d) ", i+1)
			// sig.Print()
			// fmt.Println()
			headers = append(headers, sig.Header)
		}
	}



	headers = set.From(headers).Slice()

	var sb  strings.Builder

	sb.WriteString(`#pragma comment(lib, "onecore.lib")`)
	sb.WriteString("\n")

	sb.WriteString("#include <windows.h>\n")
	for _, h := range headers {
		fmt.Fprintf(&sb, "#include <%s>\n", h)
	}

	sb.WriteString(`#include "hook_utils.cpp"`)

	hook, err := gnrtr.GenHook(func_sigs)
	if err != nil {
		log.Fatalf("%s\n", err)
	} else {
		fmt.Fprintf(&sb, "%s\n", hook)
	}

	// Generate hook table
	hook_table, err := gnrtr.GenHookTable(func_sigs)
	if err != nil {
		log.Fatalf("%s\n", err)
	} else {
		fmt.Fprintf(&sb, "%s\n", hook_table)
	}

	f, err := os.OpenFile("generated/hooks.cpp", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(f, sb.String())
}
