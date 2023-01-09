package main

import (
	"fmt"
	"os"
)

//go:generate go run golang.org/x/tools/cmd/goyacc stele.y

type lexer struct {
	tokens []any
	i      int
}

func (l *lexer) Lex(lval *yySymType) int {
	if l.i >= len(l.tokens) {
		return yyEofCode
	}

	tok := l.tokens[l.i]
	l.i++

	switch tok := tok.(type) {
	case int:
		lval.val = tok
		return yyINT

	case byte:
		return int(tok)

	default:
		return yyErrCode
	}
}

func (l *lexer) Error(err string) {
	fmt.Fprintf(os.Stderr, "error (%v): %v\n", l.i-1, err)
}

func main() {
	yyErrorVerbose = true
	yyParse(&lexer{tokens: []any{1, byte('+'), 2}})
}
