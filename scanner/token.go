package scanner

import (
	"fmt"
	"slices"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type Type

var (
	keywords = map[string]Type{
		"func":   FUNC,
		"import": IMPORT,
		"let":    LET,
		"type":   TYPE,
		"if":     IF,
		"else":   ELSE,
		"switch": SWITCH,
		"as":     AS,
	}

	symbols = map[string]Type{
		"(":  LPAREN,
		")":  RPAREN,
		"{":  LBRACE,
		"}":  RBRACE,
		"[":  LBRACKET,
		"]":  RBRACKET,
		";":  SEMI,
		"+":  PLUS,
		"-":  MINUS,
		"*":  MULT,
		"/":  DIV,
		"+=": PLUSASSIGN,
		"-=": MINUSASSIGN,
		"*=": MULTASSIGN,
		"/=": DIVASSIGN,
		"^":  BITNOT,
		"|":  BITOR,
		"&":  BITAND,
		"!":  NOT,
		"||": OR,
		"&&": AND,
		"==": EQUAL,
		"!=": NOTEQUAL,
		"<":  LT,
		">":  GT,
		"<=": LE,
		">=": GE,
		"=":  ASSIGN,
		".":  DOT,
		"|>": PIPE,
		",":  COMMA,
		"<<": LSHIFT,
		">>": RSHIFT,
	}
)

type Type int

// Token types.
const (
	INVALID Type = iota

	// Keywords
	FUNC
	IMPORT
	LET
	TYPE
	IF
	ELSE
	SWITCH
	AS

	// Symbols
	LPAREN      // (
	RPAREN      // )
	LBRACE      // {
	RBRACE      // }
	LBRACKET    // [
	RBRACKET    // ]
	SEMI        // ;
	PLUS        // +
	MINUS       // -
	MULT        // *
	DIV         // /
	PLUSASSIGN  // +=
	MINUSASSIGN // -=
	MULTASSIGN  // *=
	DIVASSIGN   // /=
	BITNOT      // ^
	BITOR       // |
	BITAND      // &
	NOT         // !
	OR          // ||
	AND         // &&
	EQUAL       // ==
	NOTEQUAL    // !=
	LT          // <
	GT          // >
	LE          // <=
	GE          // >=
	ASSIGN      // =
	DOT         // .
	PIPE        // |>
	COMMA       // ,
	LSHIFT      // <<
	RSHIFT      // >>

	// Other
	IDENT
	STRING
	INT
	FLOAT
)

func keywordOrIdent(s string) Type {
	if t, ok := keywords[s]; ok {
		return t
	}
	return IDENT
}

func preventSemi(tok Token) bool {
	return slices.Contains([]Type{
		SEMI,
		DOT,
		PIPE,
		COMMA,
		LPAREN,
		LBRACE,
		LBRACKET,
	}, tok.Type)
}

type Token struct {
	Line, Col int
	Type      Type
	Val       any
}

func (t Token) String() string {
	return fmt.Sprintf("%v (%v:%v)", t.Type, t.Line, t.Col)
}
