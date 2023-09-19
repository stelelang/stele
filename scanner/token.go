package scanner

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
		"=":  ASSIGN,
		".":  DOT,
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
	ASSIGN      // =
	DOT         // .
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

type Token struct {
	Line, Col int
	Type      Type
	Val       any
}
