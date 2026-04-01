package genc_fmt

import (
	"fmt"
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
)

type Token struct {
	typ TokenType
	str string
}

type Tokenizer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func (t Token) Print() {
	fmt.Println("Type: ", t.typ)
	fmt.Println("Value: ", t.str)
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

func Tokenize(input string) []Token {
	t := NewTokenizer(input)

	tokens := make([]Token, 0)

	for ; t.ch != 0; t.readChar() {

	process_new_char: // goto is required for clear tokenizing process

		switch {

		case t.ch == '@': // Primitives processor
			{
				primitive := Token{
					typ: Primitive,
				}

				start := t.readPosition
				for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
				}
				end := t.position

				primitive.str = t.input[start:end]
				tokens = append(tokens, primitive)

				goto process_new_char
			}

		case t.ch == '$': // Field Identifier pprocessor
			{
				field := Token{
					typ: FieldToken,
				}

				start := t.readPosition
				for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
				}
				end := t.position

				field.str = t.input[start:end]
				tokens = append(tokens, field)

				goto process_new_char
			}

		case t.ch == '#': // Operation Identifier processor
			{
				operation := Token{
					typ: Operation,
				}

				start := t.readPosition
				for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
				}
				end := t.position

				operation.str = t.input[start:end]
				tokens = append(tokens, operation)

				goto process_new_char
			}

		case t.ch == '`': // BacktickStringValue processor
			{
				bt_string := Token{
					typ: BacktickStringValue,
				}

				start := t.readPosition
				for t.readChar(); t.ch != '`'; t.readChar() {
				}
				end := t.position

				bt_string.str = t.input[start:end]
				tokens = append(tokens, bt_string)
			}

		case t.ch == '(': // ParanOpen processor
			{
				tokens = append(tokens, Token{
					typ: ParanOpen,
				})
			}

		case t.ch == ')': // ParanClose processor
			{
				tokens = append(tokens, Token{
					typ: ParanClose,
				})
			}

		case t.ch == '{': // BraceOpen processor
			{
				tokens = append(tokens, Token{
					typ: BraceOpen,
				})
			}

		case t.ch == '}': // BraceClose processor
			{
				tokens = append(tokens, Token{
					typ: BraceClose,
				})
			}

		case t.ch == '.': // FieldPoint processor
			{
				tokens = append(tokens, Token{
					typ: FieldPoint,
				})
			}

		case t.ch == ',': // CommaSeperator processor
			{
				tokens = append(tokens, Token{
					typ: CommaSeperator,
				})
			}

		case t.ch == '=': // Equals processor
			{
				tokens = append(tokens, Token{
					typ: Equals,
				})
			}

		case isValidIdentifierChar(t.ch): // NormalStringValue Processor
			{
				normal_value := Token{
					typ: NormalStringValue,
				}

				start := t.position
				for t.readChar(); isValidIdentifierChar(t.ch); t.readChar() {
				}
				end := t.position

				normal_value.str = t.input[start:end]
				tokens = append(tokens, normal_value)

				goto process_new_char
			}

		}

	}

	return tokens
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
	typ         FieldType
	expressions []Expression
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

//  (section) ------------------------------------------------------------------------- : parser  //
