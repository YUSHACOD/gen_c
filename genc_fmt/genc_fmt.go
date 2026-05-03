package genc_fmt

import (
	"fmt"
	"log"
	"strings"

	"github.com/xlab/treeprint"
)

// tokenizer : ---------------------------------------------------------------------- (section)  //
type TokenType string

const (
	//  token types : ---------------------------------------------------------------- (section)  //
	TT_Primitive           TokenType = "Primitive"
	TT_BraceOpen           TokenType = "BraceOpen"
	TT_BraceClose          TokenType = "BraceClose"
	TT_ParanOpen           TokenType = "ParanOpen"
	TT_ParanClose          TokenType = "ParanClose"
	TT_FieldPoint          TokenType = "FieldPoint"
	TT_BacktickStringValue TokenType = "BacktickStringValue"
	TT_NormalStringValue   TokenType = "NormalStringValue"
	TT_CommaSeperator      TokenType = "CommaSeperator"
	TT_Operation           TokenType = "Operation"
	TT_FieldID             TokenType = "FieldID"
	TT_Equals              TokenType = "Equals"
	TT_Eof                 TokenType = "Eof"
)

type Token struct {
	Typ  TokenType
	Str  string
	line int32
	col  int32
}

type Tokenizer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	lines        []string
	line         int32
	col          int32
}

func (t Token) Print() {
	fmt.Println("Type: ", t.Typ)
	fmt.Println("Value: ", t.Str)
	fmt.Printf("%d:%d\n", t.line, t.col)
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.lines = strings.Split(t.input, "\n")
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
		if t.ch == '\n' {
			t.line += 1
			t.col = 0
		} else {
			t.col += 1
		}
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

	case t.ch == '>' && t.input[t.readPosition] == '>':
		{
			// t.col != 0 because t.col == 0 indicates
			// line increment
			for t.readChar(); t.col != 0; t.readChar() {
			}
		}

	case t.ch == '@': // Primitives processor
		{
			primitive := Token{
				Typ:  TT_Primitive,
				line: t.line,
				col:  t.col,
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
				Typ:  TT_FieldID,
				line: t.line,
				col:  t.col,
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
				Typ:  TT_Operation,
				line: t.line,
				col:  t.col,
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
				Typ:  TT_BacktickStringValue,
				line: t.line,
				col:  t.col,
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
				Typ:  TT_ParanOpen,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == ')': // ParanClose processor
		{
			token = Token{
				Typ:  TT_ParanClose,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '{': // BraceOpen processor
		{
			token = Token{
				Typ:  TT_BraceOpen,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '}': // BraceClose processor
		{
			token = Token{
				Typ:  TT_BraceClose,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '.': // FieldPoint processor
		{
			token = Token{
				Typ:  TT_FieldPoint,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == ',': // CommaSeperator processor
		{
			token = Token{
				Typ:  TT_CommaSeperator,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '=': // Equals processor
		{
			token = Token{
				Typ:  TT_Equals,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case isValidIdentifierChar(t.ch): // NormalStringValue Processor
		{
			normal_value := Token{
				Typ:  TT_NormalStringValue,
				line: t.line,
				col:  t.col,
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
				Typ:  TT_Eof,
				line: t.line,
				col:  t.col,
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

// parser : -------------------------------------------------------------------------- (section)  //
type PrimitiveType string

const (
	//  primitive types : ------------------------------------------------------------ (section)  //
	PT_Table       PrimitiveType = "table"
	PT_Enum        PrimitiveType = "enum"
	PT_Enum2String PrimitiveType = "enum_to_string_table"
	PT_Struct      PrimitiveType = "struct"
	PT_FuncTypes   PrimitiveType = "func_types"
	PT_FuncGlobals PrimitiveType = "func_globals"
	PT_Custom      PrimitiveType = "custom"
	PT_GenCFile    PrimitiveType = "genc"
	PT_GenHFile    PrimitiveType = "genh"
	PT_GenCPPFile  PrimitiveType = "gencpp"
	PT_GenHPPFile  PrimitiveType = "genhpp"
)

type SubPrimType string

const (
	//  sub_primitive types : -------------------------------------------------------- (section)  //
	ST_Requires SubPrimType = "requires"
)

type FieldType string

const (
	//  field types : ---------------------------------------------------------------- (section)  //

	// Table Fields
	FT_Table_Cols FieldType = FieldType((PT_Table) + "_" + "cols")
	FT_Table_Rows FieldType = FieldType((PT_Table) + "_" + "rows")

	// Enum Fields
	FT_Enum_ValueName FieldType = FieldType((PT_Enum) + "_" + "value_name")

	// Enum2String Fields
	FT_Enum2String_Enum FieldType = FieldType((PT_Enum2String) + "_" + "enum")

	// Struct Fields
	FT_Struct_FieldTypes FieldType = FieldType(PT_Struct + "_" + "field_types")
	FT_Struct_FieldIds   FieldType = FieldType(PT_Struct + "_" + "field_ids")

	// FunTypes Fields
	FT_FuncTypes_Identifier FieldType = FieldType((PT_FuncTypes) + "_" + "identifier")
	FT_FuncTypes_Args       FieldType = FieldType((PT_FuncTypes) + "_" + "args")
	FT_FuncTypes_Ret        FieldType = FieldType((PT_FuncTypes) + "_" + "ret")

	// FuncGlobals Fields
	FT_FuncGlobals_Identifier FieldType = FieldType((PT_FuncGlobals) + "_" + "identifier")
	FT_FuncGlobals_Typ       FieldType = FieldType((PT_FuncGlobals) + "_" + "type")

	// Custom Fields
	FT_Custom_Template FieldType = FieldType((PT_Custom) + "_" + "template")

	// Gen Primitive Fields
	FT_GenCFile_Primitives   FieldType = FieldType((PT_GenCFile) + "_" + "primitives")
	FT_GenHFile_Primitives   FieldType = FieldType((PT_GenHFile) + "_" + "primitives")
	FT_GenCPPFile_Primitives FieldType = FieldType((PT_GenCPPFile) + "_" + "primitives")
	FT_GenHPPFile_Primitives FieldType = FieldType((PT_GenHPPFile) + "_" + "primitives")
)

type ExpressionType string

const (
	//  expression types : ----------------------------------------------------------- (section)  //

	// Expression is either a value
	ET_Value ExpressionType = "value"

	// an ET_Array of Expressions
	ET_Array ExpressionType = "array"

	// an Table col identifier with alias
	ET_ColId ExpressionType = "col_id"

	// Prim Id alias
	ET_PrimIdAlias ExpressionType = "prim_id_alias"

	// an Expression generated by operation
	ET_OP_Concat       ExpressionType = "op_concat"
	ET_OP_Uppercase    ExpressionType = "op_uppercase"
	ET_OP_Lowercase    ExpressionType = "op_lowercase"
	ET_OP_Snake2Pascal ExpressionType = "op_snake2pascal"
	ET_OP_Snake2Camel  ExpressionType = "op_snake2camel"
	ET_OP_Pascal2Snake ExpressionType = "op_pascal2snake"
	ET_OP_Pascal2Camel ExpressionType = "op_pascal2camel"
	ET_OP_Camel2Snake  ExpressionType = "op_camel2snake"
	ET_OP_Camel2Pascal ExpressionType = "op_camel2pascal"
)

// ast element structs : ------------------------------------------------------------ (section)  //
type Expression struct {
	typ   ExpressionType
	arr   []Expression
	value string
}

type Field struct {
	typ FieldType
	val Expression
}

type SubPrimitive struct {
	typ  SubPrimType
	args []Expression
}

type Primitive struct {
	Typ       PrimitiveType
	sub_prims []SubPrimitive
	fields    []Field
}

type GenC struct {
	Primitives map[string]Primitive
	Ids        []string
}

//  (section) ------------------------------------------------------------ : ast element structs  //

// ast element print helpers : ------------------------------------------------------ (section)  //

var opNames = map[ExpressionType]string{
	ET_OP_Concat:       "Concat",
	ET_OP_Uppercase:    "Uppercase",
	ET_OP_Lowercase:    "Lowercase",
	ET_OP_Snake2Pascal: "Snake2Pascal",
	ET_OP_Snake2Camel:  "Snake2Camel",
	ET_OP_Pascal2Snake: "Pascal2Snake",
	ET_OP_Pascal2Camel: "Pascal2Camel",
	ET_OP_Camel2Snake:  "Camel2Snake",
	ET_OP_Camel2Pascal: "Camel2Pascal",
}

func (e *Expression) addToTree(branch treeprint.Tree) {
	switch e.typ {

	case ET_Value:
		branch.AddNode(fmt.Sprintf("value: %v", e.value))

	case ET_Array:
		b := branch.AddBranch("Array")
		for _, elem := range e.arr {
			elem.addToTree(b)
		}

	case ET_ColId:
		branch.AddNode(fmt.Sprintf("alias: %v", e.value))
		branch.AddNode(fmt.Sprintf("col: %s", e.arr[0].value))

	case ET_PrimIdAlias:
		b := branch.AddBranch("PrimIdAlias")
		e.arr[0].addToTree(b)
		e.arr[1].addToTree(b)

	default:
		if name, ok := opNames[e.typ]; ok {
			b := branch.AddBranch(name)
			if name == opNames[ET_OP_Concat] {
				for i, e := range e.arr {
					num := fmt.Sprintf("%d", i)
					e.addToTree(b.AddBranch(num))
				}
			} else {
				e.arr[0].addToTree(b.AddBranch("Exp"))
			}
		}
	}
}

func (e *Expression) Print() {
	tree := treeprint.New()
	e.addToTree(tree.AddBranch("Expression"))
	fmt.Println(tree)
}

func (f *Field) addToTree(branch treeprint.Tree) {
	b := branch.AddBranch(fmt.Sprintf("Field: %v", f.typ))
	f.val.addToTree(b.AddBranch("Value"))
}

func (f *Field) Print() {
	tree := treeprint.New()
	f.addToTree(tree)
	fmt.Println(tree)
}

func (s *SubPrimitive) addToTree(branch treeprint.Tree) {
	b := branch.AddBranch(fmt.Sprintf("SubPrimitives: %v", s.typ))
	for _, exp := range s.args {
		exp.addToTree(b)
	}
}

func (s *SubPrimitive) Print() {
	tree := treeprint.New()
	s.addToTree(tree)
	fmt.Println(tree)
}

func (p *Primitive) addToTree(branch treeprint.Tree) {
	b := branch.AddBranch(fmt.Sprintf("Primitives: %v", p.Typ))

	subs := b.AddBranch("SubPrimitives")
	for _, sp := range p.sub_prims {
		sp.addToTree(subs)
	}

	fields := b.AddBranch("Fields")
	for _, field := range p.fields {
		field.addToTree(fields)
	}
}

func (p *Primitive) Print() {
	tree := treeprint.New()
	p.addToTree(tree)
	fmt.Println(tree)
}

//  (section) ------------------------------------------------------ : ast element print helpers  //

//  parser proper : ------------------------------------------------------------------ (section)  //

type Parser struct {
	t         *Tokenizer
	currToken Token
	peekToken Token
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.t.NextToken()
}

func NewParser(t *Tokenizer) *Parser {
	p := &Parser{t: t}
	p.nextToken()
	p.nextToken()
	return p
}

// error helpers : ------------------------------------------------------------------ (section)  //
const (
	Pad int32 = 3
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func printParseError(t *Tokenizer, tok Token, error_string string) {

	start := max(0, tok.line-Pad)
	end := min(int32(len(t.lines)-1), tok.line+Pad)

	length := len(fmt.Sprintf("%d", len(t.lines)))

	for i := start; i <= end; i += 1 {
		fmt.Printf("%0*d: %s\n", length, i, t.lines[i])
		if i == tok.line {
			fmt.Print(strings.Repeat(" ", int(tok.col+1)+length), "^")
			fmt.Print(Red, error_string, Reset)
			fmt.Println()
		}
	}
}

func (p *Parser) Errorf(fmt_string string, a ...any) {
	error_string := fmt.Sprintf(fmt_string, a...) + fmt.Sprintf("%v", p.currToken)
	printParseError(p.t, p.currToken, error_string)
	// os.Exit(1)
	panic(1)
}

//  (section) ------------------------------------------------------------------ : error helpers  //

func (p *Parser) parsePrimitiveId() string {

	id := ""
	p.nextToken()
	if p.currToken.Typ == TT_ParanOpen {

		p.nextToken()
		if p.currToken.Typ == TT_NormalStringValue {

			id = p.currToken.Str

			p.nextToken()
			if p.currToken.Typ != TT_ParanClose {
				p.Errorf("Field parantheses is not closed")
			}
		} else {
			p.Errorf("No id specified for the table ?")
		}
	} else {
		p.Errorf("No id specifier parantheses opened")
	}

	return id
}

func (p *Parser) parseExpression(exp *Expression) {

	parseSingleArgOperation := func(exp *Expression, typ ExpressionType) {
		exp.typ = typ
		p.nextToken()
		if p.currToken.Typ == TT_ParanOpen {

			var arg Expression
			p.parseExpression(&arg)
			exp.arr = append(exp.arr, arg)

			p.nextToken()
			if p.currToken.Typ != TT_ParanClose {
				p.Errorf("this operation takes only one expression close the scope with ')'")
			}
		} else {
			p.Errorf("Expected a start of operations args scope with '('")
		}
	}

	//  expression parsing : --------------------------------------------------------- (section)  //
	p.nextToken()
	switch p.currToken.Typ {

	case TT_BraceOpen:
		exp.typ = ET_Array
		for p.peekToken.Typ != TT_BraceClose {
			var array_exp Expression
			p.parseExpression(&array_exp)
			exp.arr = append(exp.arr, array_exp)
		}
		p.nextToken()
		if p.currToken.Typ != TT_BraceClose {
			p.Errorf("Array brace not close here")
		}

	case TT_BacktickStringValue:
		exp.typ = ET_Value
		exp.value = p.currToken.Str

	case TT_NormalStringValue:
		exp.value = p.currToken.Str
		exp.typ = ET_Value
		if p.peekToken.Typ == TT_FieldPoint {

			p.nextToken() // setting curr token to field point
			p.nextToken() // skipping the field point token

			exp.typ = ET_ColId
			exp.arr = append(exp.arr, Expression{
				typ:   ET_Value,
				value: p.currToken.Str,
			})
		}

	case TT_Operation:
		switch ExpressionType("op_" + p.currToken.Str) {

		case ET_OP_Concat:
			exp.typ = ET_OP_Concat
			p.nextToken()
			if p.currToken.Typ == TT_ParanOpen {
				arg := Expression{}
				for {
					p.parseExpression(&arg)
					exp.arr = append(exp.arr, arg)
					if p.peekToken.Typ == TT_CommaSeperator {
						p.nextToken()
					}
					if p.peekToken.Typ == TT_ParanClose {
						break
					}
				}

				p.nextToken()
				if p.currToken.Typ != TT_ParanClose {
					p.Errorf("This is not reachable think about this")
				}
			} else {
				p.Errorf("Expected a start of operations args scope with '('")
			}

		case ET_OP_Uppercase:
			parseSingleArgOperation(exp, ET_OP_Uppercase)

		case ET_OP_Lowercase:
			parseSingleArgOperation(exp, ET_OP_Lowercase)

		case ET_OP_Snake2Pascal:
			parseSingleArgOperation(exp, ET_OP_Snake2Pascal)

		case ET_OP_Snake2Camel:
			parseSingleArgOperation(exp, ET_OP_Snake2Camel)

		case ET_OP_Pascal2Snake:
			parseSingleArgOperation(exp, ET_OP_Pascal2Snake)

		case ET_OP_Pascal2Camel:
			parseSingleArgOperation(exp, ET_OP_Pascal2Camel)

		case ET_OP_Camel2Snake:
			parseSingleArgOperation(exp, ET_OP_Camel2Snake)

		case ET_OP_Camel2Pascal:
			parseSingleArgOperation(exp, ET_OP_Camel2Pascal)

		default:
			p.Errorf("Not a valid Expression")
		}

	default:
		p.Errorf("This token cannot start a expression")
	}

}

func (p *Parser) parseRequires() SubPrimitive {

	//  requires parsing : ----------------------------------------------------------- (section)  //
	var requires SubPrimitive
	requires.typ = ST_Requires

	p.nextToken()
	if p.currToken.Typ == TT_ParanOpen {

		var arg Expression
		arg.typ = ET_Array

		for p.nextToken(); p.currToken.Typ != TT_ParanClose; p.nextToken() {

			var id_alias Expression
			id_alias.typ = ET_PrimIdAlias

			switch p.currToken.Typ {

			case TT_NormalStringValue:

				id_alias.arr = append(id_alias.arr, Expression{
					typ:   ET_Value,
					value: p.currToken.Str,
				})

				p.nextToken()
				if p.currToken.Typ == TT_NormalStringValue {
					id_alias.arr = append(id_alias.arr, Expression{
						typ:   ET_Value,
						value: p.currToken.Str,
					})
				} else {
					p.Errorf("An alias for a primtive field is required")
				}

			case TT_CommaSeperator:
				continue

			default:
				p.Errorf("This token is invalid for an alias element in Requires sub prim")

			}

			arg.arr = append(arg.arr, id_alias)
		}

		requires.args = append(requires.args, arg)
	} else {
		p.Errorf("Required a Parantheses scope open for Primtive Id aliases(Requires sub prim)")
	}

	return requires
}

func (p *Parser) parsePrimitive(
	typ PrimitiveType,
	fields map[FieldType]struct{},
) (string, Primitive) {

	//  primitive parsing : ---------------------------------------------------------- (section)  //
	var id string
	var prim Primitive

	prim.Typ = typ

	id = p.parsePrimitiveId()

	p.nextToken()
	if p.currToken.Typ == TT_BraceOpen {
		if p.peekToken.Typ == TT_Primitive {
			p.nextToken()
			if p.currToken.Typ == TT_Primitive {
				if p.currToken.Str == string(ST_Requires) {

					prim.sub_prims = append(prim.sub_prims, p.parseRequires())

				} else {
					p.Errorf("Wrong sub prim type, expected requires")
				}
			} else {
				p.Errorf("There should be a requires sub prim here")
			}
		}
	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	for range fields {

		p.nextToken()

		var field_type FieldType
		var field_val Expression

		if p.currToken.Typ == TT_FieldID {

			field_type_id := FieldType(string(typ) + "_" + p.currToken.Str)
			if _, ok := fields[field_type_id]; ok {

				field_type = field_type_id

				p.nextToken()
				if p.currToken.Typ == TT_Equals {
					p.parseExpression(&field_val)

				} else {
					p.Errorf(
						"The Field is followed by a equals sifollowed by the expression")
				}

			} else {
				p.Errorf("Expected to be a value name field for %s", typ)
			}
		} else {
			p.Errorf("This is expeceted to be a %s field", typ)
		}

		prim.fields = append(prim.fields, Field{
			typ: field_type,
			val: field_val,
		})
	}

	p.nextToken()
	if p.currToken.Typ != TT_BraceClose {
		p.Errorf("%s definition scope not closed properly", typ)
	}

	return id, prim
}

func (p *Parser) parseTable() (string, Primitive) {

	//  table parsing : -------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_Table,
		map[FieldType]struct{}{
			FT_Table_Cols: {},
			FT_Table_Rows: {},
		},
	)
}

func (p *Parser) parseEnum() (string, Primitive) {
	//  enum parsing : --------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_Enum,

		map[FieldType]struct{}{
			FT_Enum_ValueName: {},
		},
	)
}

func (p *Parser) parseEnumToString() (string, Primitive) {
	//  enum2string parsing : -------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_Enum2String,
		map[FieldType]struct{}{
			FT_Enum2String_Enum: {},
		},
	)
}

func (p *Parser) parseStruct() (string, Primitive) {
	//  struct parsing : ------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_Struct,
		map[FieldType]struct{}{
			FT_Struct_FieldIds:   {},
			FT_Struct_FieldTypes: {},
		},
	)
}

func (p *Parser) parseFuncTypes() (string, Primitive) {
	//  func_types parsing : --------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_FuncTypes,
		map[FieldType]struct{}{
			FT_FuncTypes_Identifier: {},
			FT_FuncTypes_Args:       {},
			FT_FuncTypes_Ret:        {},
		},
	)
}

func (p *Parser) parseFuncGlobals() (string, Primitive) {
	//  func_globals parsing : ------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_FuncGlobals,
		map[FieldType]struct{}{
			FT_FuncGlobals_Identifier: {},
			FT_FuncGlobals_Typ:       {},
		},
	)
}

func (p *Parser) parseCustom() (string, Primitive) {
	//  custom parsing : ------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_Custom,
		map[FieldType]struct{}{
			FT_Custom_Template: {},
		},
	)
}

func (p *Parser) parseGenCPrim() (string, Primitive) {
	//  file_type_parsing : ---------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		PT_GenCFile,

		map[FieldType]struct{}{
			FT_GenCFile_Primitives: {},
		},
	)
}

func (p *Parser) parseGenHPrim() (string, Primitive) {
	return p.parsePrimitive(
		PT_GenHFile,
		map[FieldType]struct{}{
			FT_GenHFile_Primitives: {},
		},
	)
}

func (p *Parser) parseGenCppPrim() (string, Primitive) {
	return p.parsePrimitive(
		PT_GenCPPFile,
		map[FieldType]struct{}{
			FT_GenCPPFile_Primitives: {},
		},
	)
}

func (p *Parser) parseGenHppPrim() (string, Primitive) {
	return p.parsePrimitive(
		PT_GenHPPFile,
		map[FieldType]struct{}{
			FT_GenHPPFile_Primitives: {},
		},
	)
}

func ParseGenc(t *Tokenizer) *GenC {

	p := NewParser(t)
	prmitives := make(map[string]Primitive)
	genc := &GenC{Primitives: prmitives}

	if p.currToken.Typ != TT_Primitive {
		log.Fatalf("First Token Found should be primitive check the formating of file")
	}

	//  parser core : ---------------------------------------------------------------- (section)  //
	for ; p.currToken.Typ != TT_Eof; p.nextToken() {
		var id_string string
		switch PrimitiveType(p.currToken.Str) {

		case PT_Table:
			{
				id, table := p.parseTable()
				genc.Primitives[id] = table
				id_string = id
			}

		case PT_Enum:
			{
				id, enum := p.parseEnum()
				genc.Primitives[id] = enum
				id_string = id
			}

		case PT_Enum2String:
			{
				id, enum_to_string := p.parseEnumToString()
				genc.Primitives[id] = enum_to_string
				id_string = id
			}

		case PT_Struct:
			{
				id, struct_prim := p.parseStruct()
				genc.Primitives[id] = struct_prim
				id_string = id
			}

		case PT_FuncTypes:
			{
				id, func_types := p.parseFuncTypes()
				genc.Primitives[id] = func_types
				id_string = id
			}

		case PT_FuncGlobals:
			{
				id, func_globals := p.parseFuncGlobals()
				genc.Primitives[id] = func_globals
				id_string = id
			}

		case PT_Custom:
			{
				id, custom := p.parseCustom()
				genc.Primitives[id] = custom
				id_string = id
			}

		case PT_GenCFile:
			{
				id, gen_c_file := p.parseGenCPrim()
				genc.Primitives[id] = gen_c_file
				id_string = id
			}

		case PT_GenHFile:
			{
				id, gen_h_file := p.parseGenHPrim()
				genc.Primitives[id] = gen_h_file
				id_string = id
			}

		case PT_GenCPPFile:
			{
				id, gen_cpp_file := p.parseGenCppPrim()
				genc.Primitives[id] = gen_cpp_file
				id_string = id
			}

		case PT_GenHPPFile:
			{
				id, gen_hpp_file := p.parseGenHppPrim()
				genc.Primitives[id] = gen_hpp_file
				id_string = id
			}

		default:
			{
				// for k, v := range genc.Primitives {
				// 	fmt.Println("Primitive Id: ", k)
				// 	fmt.Println("Primitive Val:")
				// 	v.Print()
				// }
				p.Errorf("Invalid Primitive type: %s", p.currToken.Str)
			}
		}

		genc.Ids = append(genc.Ids, id_string)
	}

	return genc
}

//  (section) ------------------------------------------------------------------ : parser proper  //

//  (section) ------------------------------------------------------------------------- : parser  //

