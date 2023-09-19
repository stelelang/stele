package scanner

import (
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		tok   Token
	}{
		{name: "Whitespace", input: "     "},
		{name: "Func", input: "func", tok: Token{1, 1, FUNC, "func"}},
		{name: "SimpleIdent", input: "test ", tok: Token{1, 1, IDENT, "test"}},
		{name: "PrivateIdent", input: " _something_private", tok: Token{1, 2, IDENT, "_something_private"}},
		{name: "ConstIdent", input: "test!=", tok: Token{1, 1, IDENT, "test!"}},
		{name: "String", input: "\"a test\"", tok: Token{1, 1, STRING, "a test"}},
		{name: "Int", input: "123", tok: Token{1, 1, INT, int64(123)}},
		{name: "Float", input: "123.5321", tok: Token{1, 1, FLOAT, 123.5321}},
		{name: "Plus", input: "+!", tok: Token{1, 1, PLUS, "+"}},
		{name: "Left Shift", input: "<<", tok: Token{1, 1, LSHIFT, "<<"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := NewScanner(strings.NewReader(test.input))
			s.Scan()
			if s.Err() != nil {
				t.Fatal(s.Err())
			}
			if s.Tok() != test.tok {
				t.Fatalf("token doesn't match\n\tgot: %+v\n\texpected: %+v", s.Tok(), test.tok)
			}
		})
	}
}
