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
