package genc_fmt

import (
	"fmt"
	"log"
)

// tokenizer : ---------------------------------------------------------------------- (section)  //
type TokenType string

const (
	Primitive           TokenType = "Primitive"
	BraceOpen           TokenType = "BraceOpen"
	BraceClose          TokenType = "BraceClose"
	ParanOpen           TokenType = "ParanOpen"
	ParanClose          TokenType = "ParanClose"
	FieldPoint          TokenType = "FieldPoint"
	BacktickStringValue TokenType = "BacktickStringValue"
	NormalStringValue   TokenType = "NormalStringValue"
	CommaSeperator      TokenType = "CommaSeperator"
	FieldToken          TokenType = "FieldToken"
	Operation           TokenType = "Operation"
	Equals              TokenType = "Equals"
	Eof                 TokenType = "Eof"
)

type Token struct {
	Typ TokenType
	Str string
}

type Tokenizer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func (t Token) Print() {
	fmt.Println("Type: ", t.Typ)
	fmt.Println("Value: ", t.Str)
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.readChar()
	return t
}

func (t *Tokenizer) readChar() {
	if t.readPosition >= len(t.input) {
		t.ch = 0
	} else {
		t.ch = t.input[t.readPosition]
		t.position = t.readPosition
		t.readPosition += 1
	}
}

func isValidIdentifierChar(ch byte) bool {
	return (ch == '_') ||
		('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		('0' <= ch && ch <= '9')
}

func (t *Tokenizer) NextToken() Token {

	token := Token{}
	token.Typ = "EMPTY"

start:
	switch {

	case t.ch == '@': // Primitives processor
		{
			primitive := Token{
				Typ: Primitive,
			}

			start := t.readPosition
			for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
			}
			end := t.position

			primitive.Str = t.input[start:end]
			token = primitive
		}

	case t.ch == '$': // Field Identifier pprocessor
		{
			field := Token{
				Typ: FieldToken,
			}

			start := t.readPosition
			for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
			}
			end := t.position

			field.Str = t.input[start:end]
			token = field
		}

	case t.ch == '#': // Operation Identifier processor
		{
			operation := Token{
				Typ: Operation,
			}

			start := t.readPosition
			for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
			}
			end := t.position

			operation.Str = t.input[start:end]
			token = operation
		}

	case t.ch == '`': // BacktickStringValue processor
		{
			bt_string := Token{
				Typ: BacktickStringValue,
			}

			start := t.readPosition
			for t.readChar(); t.ch != '`'; t.readChar() {
			}
			end := t.position

			bt_string.Str = t.input[start:end]
			token = bt_string

			t.readChar()
		}

	case t.ch == '(': // ParanOpen processor
		{
			token = Token{
				Typ: ParanOpen,
			}
			t.readChar()
		}

	case t.ch == ')': // ParanClose processor
		{
			token = Token{
				Typ: ParanClose,
			}
			t.readChar()
		}

	case t.ch == '{': // BraceOpen processor
		{
			token = Token{
				Typ: BraceOpen,
			}
			t.readChar()
		}

	case t.ch == '}': // BraceClose processor
		{
			token = Token{
				Typ: BraceClose,
			}
			t.readChar()
		}

	case t.ch == '.': // FieldPoint processor
		{
			token = Token{
				Typ: FieldPoint,
			}
			t.readChar()
		}

	case t.ch == ',': // CommaSeperator processor
		{
			token = Token{
				Typ: CommaSeperator,
			}
			t.readChar()
		}

	case t.ch == '=': // Equals processor
		{
			token = Token{
				Typ: Equals,
			}
			t.readChar()
		}

	case isValidIdentifierChar(t.ch): // NormalStringValue Processor
		{
			normal_value := Token{
				Typ: NormalStringValue,
			}

			start := t.position
			for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
			}
			end := t.position

			normal_value.Str = t.input[start:end]
			token = normal_value
		}

	case t.ch == 0:
		{
			token = Token{
				Typ: Eof,
			}
		}

	}

	if token.Typ == "EMPTY" {
		t.readChar()
		goto start
	}

	return token
}

//  (section) ---------------------------------------------------------------------- : tokenizer  //

// parser : ------------------------------------------------------------------------- (section)  //
type PrimitiveType string

const (
	Table       PrimitiveType = "table"
	Enum        PrimitiveType = "enum"
	Enum2String PrimitiveType = "enum_to_string"
	FuncTypes   PrimitiveType = "func_types"
	FuncGlobals PrimitiveType = "func_global"
	Custom      PrimitiveType = "custom"
)

type SubPrimType string

const (
	Requires SubPrimType = "requires"
)

type FieldType string

const (
	// Table Fields
	Table_Cols FieldType = FieldType((Table) + "_" + "cols")

	// Enum Fields
	Enum_ValueName FieldType = FieldType((Enum) + "_" + "value_name")

	// Enum2String Fields
	Enum2String_Enum FieldType = FieldType((Enum2String) + "_" + "enum")

	// FunTypes Fields
	FuncTypes_Identifier FieldType = FieldType((FuncTypes) + "_" + "identifier")
	FuncTypes_Args       FieldType = FieldType((FuncTypes) + "_" + "args")
	FuncTypes_Ret        FieldType = FieldType((FuncTypes) + "_" + "ret")

	// FuncGlobals Fields
	FuncGlobals_Identifier FieldType = FieldType((FuncGlobals) + "_" + "identifier")
	FuncGlobals_Type       FieldType = FieldType((FuncGlobals) + "_" + "type")

	// Custom Fields
	Custom_Template FieldType = FieldType((Custom) + "_" + "template")
)

type ExpressionType string

const (
	// Expression is either a value
	Value ExpressionType = "value"

	// an Array of Expressions
	Array ExpressionType = "array"

	// an Table col identifier with alias
	ColId ExpressionType = "col_id"

	// an Expression generated by operation
	Op_Concat       ExpressionType = "op_concat"
	Op_Uppercase    ExpressionType = "op_uppercase"
	Op_Lowercase    ExpressionType = "op_lowercase"
	Op_Snake2Pascal ExpressionType = "op_snake2pascal"
	Op_Snake2Camel  ExpressionType = "op_snake2camel"
	Op_Pascal2Snake ExpressionType = "op_pascal2snake"
	Op_Pascal2Camel ExpressionType = "op_pascal2camel"
	Op_Camel2Snake  ExpressionType = "op_camel2snake"
	Op_Camel2Pascal ExpressionType = "op_camel2pascal"
)

type Expression struct {
	typ   ExpressionType
	arr   []Expression
	value string
}

type Field struct {
	typ FieldType
	val Expression
}

type SubPrimitives struct {
	typ  SubPrimType
	args []Expression
}

type Primitives struct {
	typ       PrimitiveType
	sub_prims []SubPrimitives
	fields    []Field
}

type GenC struct {
	primitives map[string]Primitives
}

func (e *Expression) Print() {
	fmt.Print("Expression: ")
	switch e.typ {

	case Value:
		fmt.Println(e.value)

	case Array:
		fmt.Print("Array: [")
		for _, arr_elem := range e.arr {
			arr_elem.Print()
			fmt.Println()
		}
		fmt.Println("]")

	case ColId:
		fmt.Print("Alias Id: ")
		e.arr[0].Print()
		fmt.Print("Col Id: ")
		e.arr[1].Print()

	case Op_Concat:
		fmt.Printf("Concat: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Uppercase:
		fmt.Printf("Uppercase: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Lowercase:
		fmt.Printf("Lowercase: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Snake2Pascal:
		fmt.Printf("Snake2Pascal: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Snake2Camel:
		fmt.Printf("Snake2Camel: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Pascal2Snake:
		fmt.Printf("Pascal2Snake: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Pascal2Camel:
		fmt.Printf("Pascal2Camel: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Camel2Snake:
		fmt.Printf("Camel2Snake: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	case Op_Camel2Pascal:
		fmt.Printf("Camel2Pascal: \n1 :")
		e.arr[0].Print()
		fmt.Printf("2 : ")
		e.arr[1].Print()

	}
}

func (f *Field) Print() {
	fmt.Println("Type: ", f.typ)
	fmt.Print("Value: ")
	f.val.Print()
}

func (s *SubPrimitives) Print() {
	fmt.Println("Type: ", s.typ)
	for _, exp := range s.args {
		exp.Print()
	}
}

func (p *Primitives) Print() {
	fmt.Println("Type: ", p.typ)
	fmt.Println("SubPrimitves:")
	for _, sp := range p.sub_prims {
		sp.Print()
	}
	fmt.Println("Fields:")
	for _, field := range p.fields {
		field.Print()
	}
}

func ParseGenc(t *Tokenizer) {
	token := t.NextToken()
	if token.Typ != Primitive {
		log.Fatalf("First Token Found should be primitive check the formating of file")
	}

	switch PrimitiveType(token.Str) {
	case Table:
		{

		}

	case Enum:
		{

		}

	case Enum2String:
		{

		}

	case FuncTypes:
		{

		}

	case FuncGlobals:
		{

		}

	case Custom:
		{

		}

	default:
		{
			log.Fatalf("Invalid Primitive type: %s", token.Str)
		}
	}
}

//  (section) ------------------------------------------------------------------------- : parser  //
