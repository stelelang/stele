package parser

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"deedles.dev/stele"
	"deedles.dev/stele/parser/ast"
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
	var decls []stele.Declaration
	allowImport := true

	for {
		tok, ok := p.next()
		if !ok {
			script.Scope = script.Scope.AddAll(decls)
			return script
		}

		switch tok.Type {
		case scanner.IMPORT:
			if !allowImport {
				p.throw(errors.New("imports must come before all other top-level declarations"))
			}
			decls = append(decls, p.parseImport())
		case scanner.LET:
			allowImport = false
			decls = append(decls, p.parseLet())
		default:
			p.throw(UnexpectedTokenError{tok})
		}
	}
}

func (p *parser) parseImport() ast.Import {
	path := p.expect(scanner.STRING).Val.(string)

	tok := p.expect(-1)
	switch tok.Type {
	case scanner.AS:
		id := p.expect(scanner.IDENT).Val.(string)
		p.expect(scanner.SEMI)
		return ast.Import{Name: id, Path: path}

	case scanner.SEMI:
		// TODO: Is the basename good enough?
		return ast.Import{Name: filepath.Base(path), Path: path}

	default:
		p.throw(UnexpectedTokenError{tok})
		return ast.Import{}
	}
}

func (p *parser) parseLet() ast.Let {
	id := p.expect(scanner.IDENT).Val.(string)

	tok := p.expect(-1)
	switch tok.Type {
	case scanner.IDENT:
		panic("Not implemented.")
	case scanner.ASSIGN:
		rhs := p.parseExpr()
		return ast.Let{
			Name:   id,
			T:      rhs.Type(),
			Assign: &stele.Assign{ID: id, Val: rhs},
		}
	default:
		p.throw(UnexpectedTokenError{tok})
		return ast.Let{}
	}
}

func (p *parser) parseExpr() stele.Expr {
	panic("Not implemented.")
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
