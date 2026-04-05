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
	Primitive           TokenType = "Primitive"
	BraceOpen           TokenType = "BraceOpen"
	BraceClose          TokenType = "BraceClose"
	ParanOpen           TokenType = "ParanOpen"
	ParanClose          TokenType = "ParanClose"
	FieldPoint          TokenType = "FieldPoint"
	BacktickStringValue TokenType = "BacktickStringValue"
	NormalStringValue   TokenType = "NormalStringValue"
	CommaSeperator      TokenType = "CommaSeperator"
	FieldID             TokenType = "FieldID"
	Operation           TokenType = "Operation"
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

// parser : ------------------------------------------------------------------------- (section)  //
type PrimitiveType string

const (
	//  primitive types : ------------------------------------------------------------ (section)  //
	Table       PrimitiveType = "table"
	Enum        PrimitiveType = "enum"
	Enum2String PrimitiveType = "enum_to_string"
	FuncTypes   PrimitiveType = "func_types"
	FuncGlobals PrimitiveType = "func_global"
	Custom      PrimitiveType = "custom"
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
	//  expression types : ----------------------------------------------------------- (section)  //

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
	typ       PrimitiveType
	sub_prims []SubPrimitives
	fields    []Field
}

type GenC struct {
	Primitives map[string]Primitives
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
		branch.AddNode(fmt.Sprintf("Value: %v", e.value))

	case Array:
		b := branch.AddBranch("Array")
		for _, elem := range e.arr {
			elem.addToTree(b)
		}

	case ColId:
		b := branch.AddBranch("ColId")
		e.arr[0].addToTree(b.AddBranch("Alias Id"))
		e.arr[1].addToTree(b.AddBranch("Col Id"))

	default:
		if name, ok := opNames[e.typ]; ok {
			b := branch.AddBranch(name)
			e.arr[0].addToTree(b.AddBranch("1"))
			e.arr[1].addToTree(b.AddBranch("2"))
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
	b := branch.AddBranch(fmt.Sprintf("Primitives: %v", p.typ))

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
	fmt.Println()
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

	p.nextToken()
	switch p.currToken.Typ {

	case BraceOpen:
		exp.typ = Array
		for ; p.currToken.Typ != BraceClose; {
			var array_exp Expression
			p.parseExpression(&array_exp)
			exp.arr = append(exp.arr, array_exp)
		}
		if p.currToken.Typ != BraceClose {
			p.Errorf("Array brace not close here")
		}
		p.nextToken()
		p.currToken.Print()

	case BacktickStringValue:
		exp.typ = Value
		exp.value = p.currToken.Str

	case NormalStringValue:
		exp.value = p.currToken.Str
		if p.peekToken.Typ == FieldPoint {
			p.nextToken()
			exp.typ = ColId
			exp.arr = append(exp.arr, Expression{
				typ:   ExpressionType(NormalStringValue),
				value: p.currToken.Str,
			})
		} else {
			exp.typ = Value
		}

	case TokenType(Op_Concat):

	case TokenType(Op_Uppercase):

	case TokenType(Op_Lowercase):

	case TokenType(Op_Snake2Pascal):

	case TokenType(Op_Snake2Camel):

	case TokenType(Op_Pascal2Snake):

	case TokenType(Op_Pascal2Camel):

	case TokenType(Op_Camel2Snake):

	case TokenType(Op_Camel2Pascal):

	default:
		p.Errorf("This is token cannot start a expression")
	}

}

func (p *Parser) parseTable() (string, Primitives) {

	//  table parsing : -------------------------------------------------------------- (section)  //
	var id string
	var table Primitives
	table.typ = Table

	id = p.parsePrimitiveId()
	fmt.Println("Parsed Id: ", id)

	p.nextToken()
	if p.currToken.Typ == BraceOpen {

		for range 2 {
			p.nextToken()

			var field_type FieldType
			var field_val Expression

			if p.currToken.Typ == FieldID {

				switch FieldType(Table + "_" + PrimitiveType(p.currToken.Str)) {

				case Table_Cols:
					{

						field_type = Table_Cols

						p.nextToken()
						if p.currToken.Typ == Equals {
							p.parseExpression(&field_val)
						} else {
							p.Errorf(
								"The Field is followed by a equals sign followed by the expression")
						}
					}

				case Table_Rows:
					{
						field_type = Table_Rows

						p.nextToken()
						if p.currToken.Typ == Equals {
							p.parseExpression(&field_val)
						} else {
							p.Errorf(
								"The Field is followed by a equals sign followed by the expression")
						}
					}

				default:
					p.Errorf("Expected a table field")
				}

			} else {
				p.Errorf("This should've been a feild instead of whatever")
			}

			table.fields = append(table.fields, Field{
				typ: field_type,
				val: field_val,
			})
		}

	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	return id, table
}

func ParseGenc(t *Tokenizer) *GenC {

	p := NewParser(t)
	prmitives := make(map[string]Primitives)
	genc := &GenC{Primitives: prmitives}

	if p.currToken.Typ != Primitive {
		log.Fatalf("First Token Found should be primitive check the formating of file")
	}

	//  parser core : ---------------------------------------------------------------- (section)  //
	switch PrimitiveType(p.currToken.Str) {

	case Table:
		{
			id, table := p.parseTable()
			genc.Primitives[id] = table
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
			log.Fatalf("Invalid Primitive type: %s", p.currToken.Str)
		}
	}

	return genc
}

//  (section) ------------------------------------------------------------------ : parser proper  //

//  (section) ------------------------------------------------------------------------- : parser  //
