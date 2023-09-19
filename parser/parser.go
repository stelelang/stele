package parser

import (
	"io"

	"deedles.dev/stele"
	"deedles.dev/stele/scanner"
)

func Parse(r io.Reader) (script stele.Script, err error) {
	p := parser{s: scanner.New(r)}
	defer p.catch(&err)
	return p.parseScript(), nil
}

type parser struct {
	s *scanner.Scanner
}

func (p *parser) parseScript() stele.Script {
	panic("not implemented")
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
