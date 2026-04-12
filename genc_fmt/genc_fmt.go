package genc_fmt

import (
	"fmt"
	"log"
	"strings"

	"github.com/xlab/treeprint"

	gn "github.com/YUSHACOD/gen_c/gnrtr"
)

// tokenizer : ---------------------------------------------------------------------- (section)  //
type TokenType string

const (
	//  token types : ---------------------------------------------------------------- (section)  //
	Primitive           TokenType = "Primitive"
	BraceOpen           TokenType = "BraceOpen"
	BraceClose          TokenType = "BraceClose"
	ParanOpen           TokenType = "ParanOpen"
	ParanClose          TokenType = "ParanClose"
	FieldPoint          TokenType = "FieldPoint"
	BacktickStringValue TokenType = "BacktickStringValue"
	NormalStringValue   TokenType = "NormalStringValue"
	CommaSeperator      TokenType = "CommaSeperator"
	Operation           TokenType = "Operation"
	FieldID             TokenType = "FieldID"
	Equals              TokenType = "Equals"
	Eof                 TokenType = "Eof"
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
				Typ:  Primitive,
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
				Typ:  FieldID,
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
				Typ:  Operation,
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
				Typ:  BacktickStringValue,
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
				Typ:  ParanOpen,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == ')': // ParanClose processor
		{
			token = Token{
				Typ:  ParanClose,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '{': // BraceOpen processor
		{
			token = Token{
				Typ:  BraceOpen,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '}': // BraceClose processor
		{
			token = Token{
				Typ:  BraceClose,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '.': // FieldPoint processor
		{
			token = Token{
				Typ:  FieldPoint,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == ',': // CommaSeperator processor
		{
			token = Token{
				Typ:  CommaSeperator,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case t.ch == '=': // Equals processor
		{
			token = Token{
				Typ:  Equals,
				line: t.line,
				col:  t.col,
			}
			t.readChar()
		}

	case isValidIdentifierChar(t.ch): // NormalStringValue Processor
		{
			normal_value := Token{
				Typ:  NormalStringValue,
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
				Typ:  Eof,
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
	Table       PrimitiveType = "table"
	Enum        PrimitiveType = "enum"
	Enum2String PrimitiveType = "enum_to_string_table"
	Struct      PrimitiveType = "struct"
	FuncTypes   PrimitiveType = "func_types"
	FuncGlobals PrimitiveType = "func_globals"
	Custom      PrimitiveType = "custom"
	GenCFile    PrimitiveType = "genc"
	GenHFile    PrimitiveType = "genh"
	GenCPPFile  PrimitiveType = "gencpp"
	GenHPPFile  PrimitiveType = "genhpp"
)

type SubPrimType string

const (
	//  sub_primitive types : -------------------------------------------------------- (section)  //
	Requires SubPrimType = "requires"
)

type FieldType string

const (
	//  field types : ---------------------------------------------------------------- (section)  //

	// Table Fields
	Table_Cols FieldType = FieldType((Table) + "_" + "cols")
	Table_Rows FieldType = FieldType((Table) + "_" + "rows")

	// Enum Fields
	Enum_ValueName FieldType = FieldType((Enum) + "_" + "value_name")

	// Enum2String Fields
	Enum2String_Enum FieldType = FieldType((Enum2String) + "_" + "enum")

	// Struct Fields
	Struct_FieldTypes FieldType = FieldType(Struct + "_" + "field_types")
	Struct_FieldIds   FieldType = FieldType(Struct + "_" + "field_ids")

	// FunTypes Fields
	FuncTypes_Identifier FieldType = FieldType((FuncTypes) + "_" + "identifier")
	FuncTypes_Args       FieldType = FieldType((FuncTypes) + "_" + "args")
	FuncTypes_Ret        FieldType = FieldType((FuncTypes) + "_" + "ret")

	// FuncGlobals Fields
	FuncGlobals_Identifier FieldType = FieldType((FuncGlobals) + "_" + "identifier")
	FuncGlobals_Type       FieldType = FieldType((FuncGlobals) + "_" + "type")

	// Custom Fields
	Custom_Template FieldType = FieldType((Custom) + "_" + "template")

	// Gen Primitive Fields
	GenCFile_Primitives   FieldType = FieldType((GenCFile) + "_" + "primitives")
	GenHFile_Primitives   FieldType = FieldType((GenHFile) + "_" + "primitives")
	GenCPPFile_Primitives FieldType = FieldType((GenCPPFile) + "_" + "primitives")
	GenHPPFile_Primitives FieldType = FieldType((GenHPPFile) + "_" + "primitives")
)

type ExpressionType string

const (
	//  expression types : ----------------------------------------------------------- (section)  //

	// Expression is either a value
	Value ExpressionType = "value"

	// an Array of Expressions
	Array ExpressionType = "array"

	// an Table col identifier with alias
	ColId ExpressionType = "col_id"

	// Prim Id alias
	PrimIdAlias ExpressionType = "prim_id_alias"

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

type SubPrimitives struct {
	typ  SubPrimType
	args []Expression
}

type Primitives struct {
	Typ       PrimitiveType
	sub_prims []SubPrimitives
	fields    []Field
}

type GenC struct {
	Primitives map[string]Primitives
	Ids        []string
}

//  (section) ------------------------------------------------------------ : ast element structs  //

// ast element print helpers : ------------------------------------------------------ (section)  //

var opNames = map[ExpressionType]string{
	Op_Concat:       "Concat",
	Op_Uppercase:    "Uppercase",
	Op_Lowercase:    "Lowercase",
	Op_Snake2Pascal: "Snake2Pascal",
	Op_Snake2Camel:  "Snake2Camel",
	Op_Pascal2Snake: "Pascal2Snake",
	Op_Pascal2Camel: "Pascal2Camel",
	Op_Camel2Snake:  "Camel2Snake",
	Op_Camel2Pascal: "Camel2Pascal",
}

func (e *Expression) addToTree(branch treeprint.Tree) {
	switch e.typ {

	case Value:
		branch.AddNode(fmt.Sprintf("value: %v", e.value))

	case Array:
		b := branch.AddBranch("Array")
		for _, elem := range e.arr {
			elem.addToTree(b)
		}

	case ColId:
		branch.AddNode(fmt.Sprintf("alias: %v", e.value))
		branch.AddNode(fmt.Sprintf("col: %s", e.arr[0].value))

	case PrimIdAlias:
		b := branch.AddBranch("PrimIdAlias")
		e.arr[0].addToTree(b)
		e.arr[1].addToTree(b)

	default:
		if name, ok := opNames[e.typ]; ok {
			b := branch.AddBranch(name)
			if name == opNames[Op_Concat] {
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

func (s *SubPrimitives) addToTree(branch treeprint.Tree) {
	b := branch.AddBranch(fmt.Sprintf("SubPrimitives: %v", s.typ))
	for _, exp := range s.args {
		exp.addToTree(b)
	}
}

func (s *SubPrimitives) Print() {
	tree := treeprint.New()
	s.addToTree(tree)
	fmt.Println(tree)
}

func (p *Primitives) addToTree(branch treeprint.Tree) {
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

func (p *Primitives) Print() {
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
	if p.currToken.Typ == ParanOpen {

		p.nextToken()
		if p.currToken.Typ == NormalStringValue {

			id = p.currToken.Str

			p.nextToken()
			if p.currToken.Typ != ParanClose {
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
		if p.currToken.Typ == ParanOpen {

			var arg Expression
			p.parseExpression(&arg)
			exp.arr = append(exp.arr, arg)

			p.nextToken()
			if p.currToken.Typ != ParanClose {
				p.Errorf("this operation takes only one expression close the scope with ')'")
			}
		} else {
			p.Errorf("Expected a start of operations args scope with '('")
		}
	}

	//  expression parsing : --------------------------------------------------------- (section)  //
	p.nextToken()
	switch p.currToken.Typ {

	case BraceOpen:
		exp.typ = Array
		for p.peekToken.Typ != BraceClose {
			var array_exp Expression
			p.parseExpression(&array_exp)
			exp.arr = append(exp.arr, array_exp)
		}
		p.nextToken()
		if p.currToken.Typ != BraceClose {
			p.Errorf("Array brace not close here")
		}

	case BacktickStringValue:
		exp.typ = Value
		exp.value = p.currToken.Str

	case NormalStringValue:
		exp.value = p.currToken.Str
		exp.typ = Value
		if p.peekToken.Typ == FieldPoint {

			p.nextToken() // setting curr token to field point
			p.nextToken() // skipping the field point token

			exp.typ = ColId
			exp.arr = append(exp.arr, Expression{
				typ:   Value,
				value: p.currToken.Str,
			})
		}

	case Operation:
		switch ExpressionType("op_" + p.currToken.Str) {

		case Op_Concat:
			exp.typ = Op_Concat
			p.nextToken()
			if p.currToken.Typ == ParanOpen {
				arg := Expression{}
				for {
					p.parseExpression(&arg)
					exp.arr = append(exp.arr, arg)
					if p.peekToken.Typ == CommaSeperator {
						p.nextToken()
					}
					if p.peekToken.Typ == ParanClose {
						break
					}
				}

				p.nextToken()
				if p.currToken.Typ != ParanClose {
					p.Errorf("This is not reachable think about this")
				}
			} else {
				p.Errorf("Expected a start of operations args scope with '('")
			}

		case Op_Uppercase:
			parseSingleArgOperation(exp, Op_Uppercase)

		case Op_Lowercase:
			parseSingleArgOperation(exp, Op_Lowercase)

		case Op_Snake2Pascal:
			parseSingleArgOperation(exp, Op_Snake2Pascal)

		case Op_Snake2Camel:
			parseSingleArgOperation(exp, Op_Snake2Camel)

		case Op_Pascal2Snake:
			parseSingleArgOperation(exp, Op_Pascal2Snake)

		case Op_Pascal2Camel:
			parseSingleArgOperation(exp, Op_Pascal2Camel)

		case Op_Camel2Snake:
			parseSingleArgOperation(exp, Op_Camel2Snake)

		case Op_Camel2Pascal:
			parseSingleArgOperation(exp, Op_Camel2Pascal)

		default:
			p.Errorf("Not a valid Expression")
		}

	default:
		p.Errorf("This token cannot start a expression")
	}

}

func (p *Parser) parseRequires() SubPrimitives {

	//  requires parsing : ----------------------------------------------------------- (section)  //
	var requires SubPrimitives

	p.nextToken()
	if p.currToken.Typ == ParanOpen {

		var arg Expression
		arg.typ = Array

		for p.nextToken(); p.currToken.Typ != ParanClose; p.nextToken() {

			var id_alias Expression
			id_alias.typ = PrimIdAlias

			switch p.currToken.Typ {

			case NormalStringValue:

				id_alias.arr = append(id_alias.arr, Expression{
					typ:   Value,
					value: p.currToken.Str,
				})

				p.nextToken()
				if p.currToken.Typ == NormalStringValue {
					id_alias.arr = append(id_alias.arr, Expression{
						typ:   Value,
						value: p.currToken.Str,
					})
				} else {
					p.Errorf("An alias for a primtive field is required")
				}

			case CommaSeperator:
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
) (string, Primitives) {

	//  primitive parsing : ---------------------------------------------------------- (section)  //
	var id string
	var prim Primitives

	prim.Typ = typ

	id = p.parsePrimitiveId()

	p.nextToken()
	if p.currToken.Typ == BraceOpen {
		if p.peekToken.Typ == Primitive {
			p.nextToken()
			if p.currToken.Typ == Primitive {
				if p.currToken.Str == string(Requires) {

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

		if p.currToken.Typ == FieldID {

			field_type_id := FieldType(string(typ) + "_" + p.currToken.Str)
			if _, ok := fields[field_type_id]; ok {

				field_type = field_type_id

				p.nextToken()
				if p.currToken.Typ == Equals {
					p.parseExpression(&field_val)

				} else {
					p.Errorf(
						"The Field is followed by a equals sign followed by the expression")
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
	if p.currToken.Typ != BraceClose {
		p.Errorf("%s definition scope not closed properly", typ)
	}

	return id, prim
}

func (p *Parser) parseTable() (string, Primitives) {

	//  table parsing : -------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		Table,
		map[FieldType]struct{}{
			Table_Cols: {},
			Table_Rows: {},
		},
	)
}

func (p *Parser) parseEnum() (string, Primitives) {
	//  enum parsing : --------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		Enum,

		map[FieldType]struct{}{
			Enum_ValueName: {},
		},
	)
}

func (p *Parser) parseEnumToString() (string, Primitives) {
	//  enum2string parsing : -------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		Enum2String,
		map[FieldType]struct{}{
			Enum2String_Enum: {},
		},
	)
}

func (p *Parser) parseStruct() (string, Primitives) {
	//  struct parsing : ------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		Struct,
		map[FieldType]struct{}{
			Struct_FieldIds:   {},
			Struct_FieldTypes: {},
		},
	)
}

func (p *Parser) parseFuncTypes() (string, Primitives) {
	//  func_types parsing : --------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		FuncTypes,
		map[FieldType]struct{}{
			FuncTypes_Identifier: {},
			FuncTypes_Args:       {},
			FuncTypes_Ret:        {},
		},
	)
}

func (p *Parser) parseFuncGlobals() (string, Primitives) {
	//  func_globals parsing : ------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		FuncGlobals,
		map[FieldType]struct{}{
			FuncGlobals_Identifier: {},
			FuncGlobals_Type:       {},
		},
	)
}

func (p *Parser) parseCustom() (string, Primitives) {
	//  custom parsing : ------------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		Custom,
		map[FieldType]struct{}{
			Custom_Template: {},
		},
	)
}

func (p *Parser) parseGenCPrim() (string, Primitives) {
	//  file_type_parsing : ---------------------------------------------------------- (section)  //
	return p.parsePrimitive(
		GenCFile,

		map[FieldType]struct{}{
			GenCFile_Primitives: {},
		},
	)
}

func (p *Parser) parseGenHPrim() (string, Primitives) {
	return p.parsePrimitive(
		GenCFile,
		map[FieldType]struct{}{
			GenCFile_Primitives: {},
		},
	)
}

func (p *Parser) parseGenCppPrim() (string, Primitives) {
	return p.parsePrimitive(
		GenCFile,
		map[FieldType]struct{}{
			GenCFile_Primitives: {},
		},
	)
}

func (p *Parser) parseGenHppPrim() (string, Primitives) {
	return p.parsePrimitive(
		GenCFile,
		map[FieldType]struct{}{
			GenCFile_Primitives: {},
		},
	)
}

func ParseGenc(t *Tokenizer) *GenC {

	p := NewParser(t)
	prmitives := make(map[string]Primitives)
	genc := &GenC{Primitives: prmitives}

	if p.currToken.Typ != Primitive {
		log.Fatalf("First Token Found should be primitive check the formating of file")
	}

	//  parser core : ---------------------------------------------------------------- (section)  //
	for ; p.currToken.Typ != Eof; p.nextToken() {
		var id_string string
		switch PrimitiveType(p.currToken.Str) {

		case Table:
			{
				id, table := p.parseTable()
				genc.Primitives[id] = table
				id_string = id
			}

		case Enum:
			{
				id, enum := p.parseEnum()
				genc.Primitives[id] = enum
				id_string = id
			}

		case Enum2String:
			{
				id, enum_to_string := p.parseEnumToString()
				genc.Primitives[id] = enum_to_string
				id_string = id
			}

		case Struct:
			{
				id, struct_prim := p.parseStruct()
				genc.Primitives[id] = struct_prim
				id_string = id
			}

		case FuncTypes:
			{
				id, func_types := p.parseFuncTypes()
				genc.Primitives[id] = func_types
				id_string = id
			}

		case FuncGlobals:
			{
				id, func_globals := p.parseFuncGlobals()
				genc.Primitives[id] = func_globals
				id_string = id
			}

		case Custom:
			{
				id, custom := p.parseCustom()
				genc.Primitives[id] = custom
				id_string = id
			}

		case GenCFile:
			{
				id, gen_c_file := p.parseGenCPrim()
				genc.Primitives[id] = gen_c_file
				id_string = id
			}

		case GenHFile:
			{
				id, gen_h_file := p.parseGenHPrim()
				genc.Primitives[id] = gen_h_file
				id_string = id
			}

		case GenCPPFile:
			{
				id, gen_cpp_file := p.parseGenCppPrim()
				genc.Primitives[id] = gen_cpp_file
				id_string = id
			}

		case GenHPPFile:
			{
				id, gen_hpp_file := p.parseGenHppPrim()
				genc.Primitives[id] = gen_hpp_file
				id_string = id
			}

		default:
			{
				for k, v := range genc.Primitives {
					fmt.Println("Primitive Id: ", k)
					fmt.Println("Primitive Val:")
					v.Print()
				}
				p.Errorf("Invalid Primitive type: %s", p.currToken.Str)
			}
		}

		genc.Ids = append(genc.Ids, id_string)
	}

	return genc
}

//  (section) ------------------------------------------------------------------ : parser proper  //

//  (section) ------------------------------------------------------------------------- : parser  //

//  expression evaluation : ---------------------------------------------------------- (section)  //

func (e *Expression) evaluate() string {
	var res string

	switch e.typ {

	case Value:
		return e.value

	case Array:
	case ColId:
	case PrimIdAlias:
	case Op_Concat:
	case Op_Uppercase:
	case Op_Lowercase:
	case Op_Snake2Pascal:
	case Op_Snake2Camel:
	case Op_Pascal2Snake:
	case Op_Pascal2Camel:
	case Op_Camel2Snake:
	case Op_Camel2Pascal:

	}

	return res
}

func (e *Expression) evaluateArray() []string {

	if e.typ != Array {
		log.Panicf("This is not a Array Expression %s", e.typ)
	}

	res := make([]string, 0)

	for _, exp := range e.arr {
		res = append(res, exp.evaluate())
	}

	return res
}

//  (section) ---------------------------------------------------------- : expression evaluation  //

// gen writables : ------------------------------------------------------------------ (section)  //
func genTable(p Primitives) gn.Table {
	table := gn.Table{
		Rows: make(map[string][]string),
	}

	for i := range 2 {
		field := p.fields[i]
		switch field.typ {
		case Table_Cols:
			{
				for _, exp := range field.val.arr {
					table.Cols = append(table.Cols, exp.evaluate())
				}
			}

		case Table_Rows:
			{
				for _, exp := range field.val.arr {
					row := exp.evaluateArray()
					for i, row_elem := range row {
						table.Rows[table.Cols[i]] = append(
							table.Rows[table.Cols[i]],
							row_elem,
						)
					}
				}
			}

		default:
			log.Panicf("This is invalid field type for table %s", field.typ)
		}
	}

	return table
}

func GenerateWritables(genc *GenC) gn.GencWritables {
	wrtb := gn.GencWritables{
		Tables:      make(map[string]gn.Table),
		Enums:       make(map[string]gn.Enum),
		FuncTypes:   make(map[string]gn.FuncType),
		FuncGlobals: make(map[string]gn.FuncGlobal),
		Customs:     make(map[string]gn.Custom),
	}

	for _, id := range genc.Ids {
		prim := genc.Primitives[id]

		switch prim.Typ {

		case Table:
			{
				wrtb.Tables[id] = genTable(prim)
			}

		case Enum:
			{

				//  todo  : ---------------------------------------------------------- (section)  //
				// do this, need gen<something> type of function to be a method of writables
				// to gains access of tables and access its data for field evaluation
				// you will have to think hard about it
				// wrtb.Table[id] = genEnum
			}

		case Enum2String:
			{

			}

		case Struct:
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

		case GenCFile:
			{
			}

		case GenHFile:
			{
			}

		case GenCPPFile:
			{
			}

		case GenHPPFile:
			{
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
