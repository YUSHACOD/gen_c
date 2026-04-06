package genc_fmt

import (
	"fmt"
	"log"
	// "os"
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
	Enum2String PrimitiveType = "enum_to_string_table"
	FuncTypes   PrimitiveType = "func_types"
	FuncGlobals PrimitiveType = "func_globals"
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
		if p.peekToken.Typ == FieldPoint {

			p.nextToken() // setting curr token to field point
			p.nextToken() // skipping the field point token

			exp.typ = ColId
			exp.arr = append(exp.arr, Expression{
				typ:   Value,
				value: p.currToken.Str,
			})
		} else {
			exp.typ = Value
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
			exp.typ = Op_Uppercase
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

		case Op_Lowercase:
			exp.typ = Op_Lowercase
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

		case Op_Snake2Pascal:
			exp.typ = Op_Snake2Pascal
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

		case Op_Snake2Camel:
			exp.typ = Op_Snake2Camel
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

		case Op_Pascal2Snake:
			exp.typ = Op_Pascal2Snake
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

		case Op_Pascal2Camel:
			exp.typ = Op_Pascal2Camel
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

		case Op_Camel2Snake:
			exp.typ = Op_Camel2Snake
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

		case Op_Camel2Pascal:
			exp.typ = Op_Camel2Pascal
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

		default:
			p.Errorf("Not a valid Expression")
		}

	default:
		p.Errorf("This token cannot start a expression")
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

		p.nextToken()
		if p.currToken.Typ != BraceClose {
			p.Errorf("Table primitives scope should end here.")
		}

	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	return id, table
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

func (p *Parser) parseEnum() (string, Primitives) {

	var id string
	var enum Primitives

	enum.typ = Enum

	id = p.parsePrimitiveId()
	fmt.Println("Parsed Id: ", id)

	p.nextToken()
	if p.currToken.Typ == BraceOpen {
		p.nextToken()
		if p.currToken.Typ == Primitive {
			if p.currToken.Str == string(Requires) {

				enum.sub_prims = append(enum.sub_prims, p.parseRequires())

			} else {
				p.Errorf("Wrong sub prim type, expected requires")
			}
		} else {
			p.Errorf("There should be a requires sub prim here")
		}

	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	p.nextToken()
	if p.currToken.Typ == FieldID {
		if FieldType(Enum)+"_"+FieldType(p.currToken.Str) == Enum_ValueName {
			var value_name_field Field
			value_name_field.typ = Enum_ValueName

			p.nextToken()
			if p.currToken.Typ == Equals {
				p.parseExpression(&value_name_field.val)

			} else {
				p.Errorf(
					"The Field is followed by a equals sign followed by the expression")
			}

			enum.fields = append(enum.fields, value_name_field)
		} else {
			p.Errorf("Expected to be a value name field for enum")
		}
	} else {
		p.Errorf("This is expeceted to be a enum field")
	}

	p.nextToken()
	if p.currToken.Typ != BraceClose {
		p.Errorf("enum definition scope not closed properly")
	}

	return id, enum
}

func (p *Parser) parseEnumToString() (string, Primitives) {
	var id string
	var enum_to_string Primitives

	enum_to_string.typ = Enum

	id = p.parsePrimitiveId()
	fmt.Println("Parsed Id: ", id)

	p.nextToken()
	if p.currToken.Typ == BraceOpen {
		p.nextToken()
		if p.currToken.Typ == FieldID {
			if FieldType(Enum2String)+"_"+FieldType(p.currToken.Str) == Enum2String_Enum {
				var value_name_field Field
				value_name_field.typ = Enum_ValueName

				p.nextToken()
				if p.currToken.Typ == Equals {
					p.parseExpression(&value_name_field.val)

				} else {
					p.Errorf(
						"The Field is followed by a equals sign followed by the expression")
				}

				enum_to_string.fields = append(enum_to_string.fields, value_name_field)
			} else {
				p.Errorf("Expected to be a value name field for enum")
			}
		} else {
			p.Errorf("This is expeceted to be a enum field")
		}
	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	p.nextToken()
	if p.currToken.Typ != BraceClose {
		p.Errorf("enum definition scope not closed properly")
	}

	return id, enum_to_string
}

func (p *Parser) parseFuncTypes() (string, Primitives) {

	var id string
	var func_types Primitives

	func_types.typ = Enum

	id = p.parsePrimitiveId()
	fmt.Println("Parsed Id: ", id)

	p.nextToken()
	if p.currToken.Typ == BraceOpen {
		p.nextToken()
		if p.currToken.Typ == Primitive {
			if p.currToken.Str == string(Requires) {

				func_types.sub_prims = append(func_types.sub_prims, p.parseRequires())

			} else {
				p.Errorf("Wrong sub prim type, expected requires")
			}
		} else {
			p.Errorf("There should be a requires sub prim here")
		}

	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	for range 3 {

		p.nextToken()

		var field_type FieldType
		var field_val Expression

		if p.currToken.Typ == FieldID {

			switch FieldType(FuncTypes) + "_" + FieldType(p.currToken.Str) {
			case FuncTypes_Identifier:
				{
					field_type = FuncTypes_Identifier

					p.nextToken()
					if p.currToken.Typ == Equals {
						p.parseExpression(&field_val)

					} else {
						p.Errorf(
							"The Field is followed by a equals sign followed by the expression")
					}
				}

			case FuncTypes_Args:
				{
					field_type = FuncTypes_Args

					p.nextToken()
					if p.currToken.Typ == Equals {
						p.parseExpression(&field_val)

					} else {
						p.Errorf(
							"The Field is followed by a equals sign followed by the expression")
					}
				}

			case FuncTypes_Ret:
				{
					field_type = FuncTypes_Ret

					p.nextToken()
					if p.currToken.Typ == Equals {
						p.parseExpression(&field_val)

					} else {
						p.Errorf(
							"The Field is followed by a equals sign followed by the expression")
					}
				}

			default:
				p.Errorf("Expected to be a value name field for func type")
			}

		} else {
			p.Errorf("This is expeceted to be a func type field")
		}

		func_types.fields = append(func_types.fields, Field{
			typ: field_type,
			val: field_val,
		})

	}

	p.nextToken()
	if p.currToken.Typ != BraceClose {
		p.Errorf("Func Types definition scope not closed properly")
	}

	return id, func_types
}

func (p *Parser) parseFuncGlobals() (string, Primitives) {

	var id string
	var func_globals Primitives

	func_globals.typ = Enum

	id = p.parsePrimitiveId()
	fmt.Println("Parsed Id: ", id)

	p.nextToken()
	if p.currToken.Typ == BraceOpen {
		p.nextToken()
		if p.currToken.Typ == Primitive {
			if p.currToken.Str == string(Requires) {

				func_globals.sub_prims = append(func_globals.sub_prims, p.parseRequires())

			} else {
				p.Errorf("Wrong sub prim type, expected requires")
			}
		} else {
			p.Errorf("There should be a requires sub prim here")
		}

	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	for range 2 {

		p.nextToken()

		var field_type FieldType
		var field_val Expression

		if p.currToken.Typ == FieldID {

			switch FieldType(FuncGlobals) + "_" + FieldType(p.currToken.Str) {

			case FuncGlobals_Identifier:
				{
					field_type = FuncGlobals_Identifier

					p.nextToken()
					if p.currToken.Typ == Equals {
						p.parseExpression(&field_val)

					} else {
						p.Errorf(
							"The Field is followed by a equals sign followed by the expression")
					}
				}

			case FuncGlobals_Type:
				{
					field_type = FuncGlobals_Type

					p.nextToken()
					if p.currToken.Typ == Equals {
						p.parseExpression(&field_val)

					} else {
						p.Errorf(
							"The Field is followed by a equals sign followed by the expression")
					}
				}

			default:
				p.Errorf("Expected to be a value name field for func globals")
			}

		} else {
			p.Errorf("This is expeceted to be a func globals field")
		}

		func_globals.fields = append(func_globals.fields, Field{
			typ: field_type,
			val: field_val,
		})

	}

	p.nextToken()
	if p.currToken.Typ != BraceClose {
		p.Errorf("Func Types definition scope not closed properly")
	}

	return id, func_globals
}

func (p *Parser) parseCustom() (string, Primitives) {

	var id string
	var custom Primitives

	custom.typ = Enum

	id = p.parsePrimitiveId()
	fmt.Println("Parsed Id: ", id)

	p.nextToken()
	if p.currToken.Typ == BraceOpen {
		p.nextToken()
		if p.currToken.Typ == Primitive {
			if p.currToken.Str == string(Requires) {

				custom.sub_prims = append(custom.sub_prims, p.parseRequires())

			} else {
				p.Errorf("Wrong sub prim type, expected requires")
			}
		} else {
			p.Errorf("There should be a requires sub prim here")
		}

	} else {
		p.Errorf("No brace open for primtive definition scope")
	}

	for range 1 {

		p.nextToken()

		var field_type FieldType
		var field_val Expression

		if p.currToken.Typ == FieldID {

			switch FieldType(Custom) + "_" + FieldType(p.currToken.Str) {

			case Custom_Template:
				{
					field_type = Custom_Template

					p.nextToken()
					if p.currToken.Typ == Equals {
						p.parseExpression(&field_val)

					} else {
						p.Errorf(
							"The Field is followed by a equals sign followed by the expression")
					}
				}

			default:
				p.Errorf("Expected to be a field for custom")
			}

		} else {
			p.Errorf("This is expeceted to be a custom prim field")
		}

		custom.fields = append(custom.fields, Field{
			typ: field_type,
			val: field_val,
		})

	}

	p.nextToken()
	if p.currToken.Typ != BraceClose {
		p.Errorf("Func Types definition scope not closed properly")
	}

	return id, custom
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
		switch PrimitiveType(p.currToken.Str) {

		case Table:
			{
				id, table := p.parseTable()
				genc.Primitives[id] = table
			}

		case Enum:
			{
				id, enum := p.parseEnum()
				genc.Primitives[id] = enum

			}

		case Enum2String:
			{
				id, enum_to_string := p.parseEnumToString()
				genc.Primitives[id] = enum_to_string
			}

		case FuncTypes:
			{
				id, func_types := p.parseFuncTypes()
				genc.Primitives[id] = func_types
			}

		case FuncGlobals:
			{
				id, func_globals := p.parseFuncGlobals()
				genc.Primitives[id] = func_globals
			}

		case Custom:
			{
				id, custom := p.parseCustom()
				genc.Primitives[id] = custom

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
	}

	return genc
}

//  (section) ------------------------------------------------------------------ : parser proper  //

//  (section) ------------------------------------------------------------------------- : parser  //
