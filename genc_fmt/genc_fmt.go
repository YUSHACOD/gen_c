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
					str: "(",
				})
			}

		case t.ch == ')': // ParanClose processor
			{
				tokens = append(tokens, Token{
					typ: ParanClose,
					str: ")",
				})
			}

		case t.ch == '{': // BraceOpen processor
			{
				tokens = append(tokens, Token{
					typ: BraceOpen,
					str: "{",
				})
			}

		case t.ch == '}': // BraceClose processor
			{
				tokens = append(tokens, Token{
					typ: BraceClose,
					str: "}",
				})
			}

		case t.ch == '.': // FieldPoint processor
			{
				tokens = append(tokens, Token{
					typ: FieldPoint,
					str: ".",
				})
			}

		case t.ch == ',': // CommaSeperator processor
			{
				tokens = append(tokens, Token{
					typ: CommaSeperator,
					str: ",",
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
