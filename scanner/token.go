package scanner

var (
	keywords = map[string]int{
		"func":   FUNC,
		"import": IMPORT,
		"let":    LET,
		"type":   TYPE,
		"if":     IF,
		"else":   ELSE,
		"switch": SWITCH,
	}

	symbols = map[string]int{
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

// Token types.
const (
	INVALID = iota

	// Keywords
	FUNC
	IMPORT
	LET
	TYPE
	IF
	ELSE
	SWITCH

	// Symbols
	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]
	SEMI     // ;
	PLUS     // +
	MINUS    // -
	MULT     // *
	DIV      // /
	BITNOT   // ^
	BITOR    // |
	BITAND   // &
	NOT      // !
	OR       // ||
	AND      // &&
	EQUAL    // ==
	NOTEQUAL // !=
	ASSIGN   // =
	DOT      // .
	LSHIFT   // <<
	RSHIFT   // >>

	// Other
	IDENT
	STRING
	INT
	FLOAT
)

func tokenType(s string) int {
	if t, ok := keywords[s]; ok {
		return t
	}
	if t, ok := symbols[s]; ok {
		return t
	}
	return IDENT
}

type Token struct {
	Line, Col int
	Type      int
	Val       any
}
