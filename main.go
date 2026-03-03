package main

import (
	"fmt"

	"github.com/YUSHACOD/gen_c/gnrtr"
)

func main() {
	fmt.Printf("Generating c code\n")

	// fields := make([]gnrtr.StructField, 3)
	// fields[0].Name = "X"
	// fields[0].Type = "int"
	//
	// fields[1].Name = "Y"
	// fields[1].Type = "int"
	//
	// fields[2].Name = "Z"
	// fields[2].Type = "int"
	//
	// struct_string, err := gnrtr.GenStruct("Point", fields)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// } else {
	// 	fmt.Printf("%s\n", struct_string)
	// }
	//
	// func_t := gnrtr.FuncType{
	// 	Name:   "add",
	// 	Args:   "(int x, int y)",
	// 	Return: "int",
	// }
	//
	// func_type, err := gnrtr.GenFuncType(func_t)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// } else {
	// 	fmt.Printf("%s\n", func_type)
	// }

	msgBoxA_func := gnrtr.FuncType{
		Name:   "MessageBoxA",
		Params: "(HWND hWnd, LPCSTR lpText, LPCSTR lpCaption, UINT uType)",
		Args:   "(hWnd, lpText, lpCaption, uType)",
		Return: "int",
	}

	// Generate hooks
	hook, err := gnrtr.GenHook(msgBoxA_func)
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", hook)
	}

	hooks := make([]gnrtr.FuncType, 1)
	hooks[0] = msgBoxA_func

	// Generate hook table
	hook_table, err := gnrtr.GenHookTable(hooks)
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", hook_table)
	}
}
