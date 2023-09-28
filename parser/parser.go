package parser

import (
	"fmt"
	"io"
	"path/filepath"

	"deedles.dev/stele"
	"deedles.dev/stele/scanner"
)

func Parse(r io.Reader) (script stele.Script, err error) {
	p := parser{s: scanner.New(r)}
	defer p.catch(&err)
	return p.parseScript(), nil
}

type parser struct {
	s   *scanner.Scanner
	buf scanner.Token
}

func (p *parser) next() (scanner.Token, bool) {
	if p.buf.Type != scanner.INVALID {
		tok := p.buf
		p.buf = scanner.Token{}
		return tok, true
	}

	ok := p.s.Scan()
	if err := p.s.Err(); err != nil {
		p.throw(fmt.Errorf("scan for next token: %w", err))
	}
	return p.s.Tok(), ok
}

func (p *parser) expect(t scanner.Type) scanner.Token {
	tok, ok := p.next()
	if !ok {
		p.throw(fmt.Errorf("expected %v but found end of input", t))
	}
	if (t >= 0) && (tok.Type != t) {
		p.throw(fmt.Errorf("expected %v but found %v", t, tok.Type))
	}

	return tok
}

func (p *parser) parseScript() stele.Script {
	var script stele.Script
	for {
		tok, ok := p.next()
		if !ok {
			return script
		}

		switch tok.Type {
		case scanner.IMPORT:
			script.Declarations.Add(p.parseImport())
		default:
			p.throw(UnexpectedTokenError{tok})
		}
	}
}

func (p *parser) parseImport() (id, path string) {
	path = p.expect(scanner.STRING).Val.(string)

	tok := p.expect(-1)
	switch tok.Type {
	case scanner.AS:
		id = p.expect(scanner.IDENT).Val.(string)
		p.expect(scanner.SEMI)
		return id, path

	case scanner.SEMI:
		return filepath.Base(path), path

	default:
		p.throw(UnexpectedTokenError{tok})
		return
	}
}

func (p *parser) throw(err error) {
	panic(parseErr{err})
}

func (p *parser) catch(err *error) {
	switch r := recover().(type) {
	case parseErr:
		*err = r.err
	case nil:
		return
	default:
		panic(r)
	}
}

type parseErr struct{ err error }

type UnexpectedTokenError struct {
	tok scanner.Token
}

func (err UnexpectedTokenError) Error() string {
	return fmt.Sprintf("unexpected token: %v", err.tok.Val)
}
