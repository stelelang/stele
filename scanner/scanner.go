package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Scanner struct {
	r   *bufio.Reader
	err error

	line, col int

	wasUnread   bool
	buf         strings.Builder
	tok         Token
	tline, tcol int
}

func New(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	s.tok = Token{}
	s.buf.Reset()
	state := s.whitespace

	defer func() {
		switch err := recover().(type) {
		case stateErr:
			s.err = fmt.Errorf("(%v:%v) %w", s.line, s.col, err.err)
			if errors.Is(s.err, io.EOF) {
				state(true)
			}
		case nil:
			return
		default:
			panic(err)
		}
	}()

	for state != nil {
		state = state(false)
	}

	return s.err != nil
}

func (s *Scanner) throw(err error) {
	panic(stateErr{err})
}

func (s *Scanner) Tok() Token {
	return s.tok
}

func (s *Scanner) Err() error {
	if errors.Is(s.err, io.EOF) {
		return nil
	}
	return s.err
}

func (s *Scanner) read() (c rune) {
	if !s.wasUnread {
		s.col++
	}

	defer func() {
		if !s.wasUnread && (c == '\n') {
			s.line++
			s.col = 0
		}
		s.wasUnread = false
	}()

	c, _, err := s.r.ReadRune()
	if err != nil {
		s.throw(err)
	}
	return c
}

func (s *Scanner) unread() {
	s.wasUnread = true
	err := s.r.UnreadRune()
	if err != nil {
		s.throw(err)
	}
}

func (s *Scanner) readEscapeSeq() rune {
	switch c := s.read(); c {
	case 't':
		return '\t'
	case 'n':
		return '\n'
	case 'r':
		return '\r'

	case 'x':
		str := string([]rune{s.read(), s.read()})
		v, err := strconv.ParseInt(str, 16, 0)
		if err != nil {
			s.throw(err)
		}
		return rune(v)

	default:
		return c
	}
}

func (s *Scanner) whitespace(eof bool) state {
	if eof {
		return nil
	}

	c := s.read()
	switch {
	case c == '\n':
		if preventSemi(s.Tok()) {
			return s.whitespace
		}

		s.startToken()
		s.endToken(SEMI, ";")
		return nil

	case unicode.IsSpace(c):
		return s.whitespace

	case unicode.IsLetter(c) || (c == '_'):
		s.buf.Reset()
		s.buf.WriteRune(c)
		s.startToken()
		return s.ident

	case c == '"':
		s.startToken()
		return s.string

	case unicode.IsNumber(c):
		s.buf.Reset()
		s.buf.WriteRune(c)
		s.startToken()
		return s.int

	case c == '\'':
		s.startToken()
		return s.char

	default:
		s.buf.Reset()
		s.buf.WriteRune(c)
		s.startToken()
		return s.symbol
	}
}

func (s *Scanner) ident(eof bool) state {
	if eof {
		str := s.buf.String()
		s.endToken(keywordOrIdent(str), str)
		return nil
	}

	c := s.read()
	switch {
	case unicode.IsLetter(c) || (c == '_') || unicode.IsNumber(c):
		s.buf.WriteRune(c)
		return s.ident

	case c == '!':
		s.buf.WriteRune(c)
		str := s.buf.String()
		s.endToken(keywordOrIdent(str), str)
		return nil

	default:
		s.unread()
		str := s.buf.String()
		s.endToken(keywordOrIdent(str), str)
		return nil
	}
}

func (s *Scanner) string(eof bool) state {
	if eof {
		s.throw(errors.New("unterminated string literal"))
		return nil
	}

	c := s.read()
	switch c {
	case '\\':
		s.buf.WriteRune(s.readEscapeSeq())
		return s.string
	case '"':
		s.endToken(STRING, s.buf.String())
		return nil
	}

	s.buf.WriteRune(c)
	return s.string
}

func (s *Scanner) int(eof bool) state {
	if eof {
		v, err := strconv.ParseInt(s.buf.String(), 0, 64)
		if err != nil {
			s.throw(err)
		}
		s.endToken(INT, v)
		return nil
	}

	c := s.read()
	switch {
	case unicode.IsNumber(c):
		s.buf.WriteRune(c)
		return s.int

	case c == '.':
		s.buf.WriteRune(c)
		return s.float

	default:
		s.unread()
		v, err := strconv.ParseInt(s.buf.String(), 0, 64)
		if err != nil {
			s.throw(err)
		}
		s.endToken(INT, v)
		return nil
	}
}

func (s *Scanner) float(eof bool) state {
	if eof {
		v, err := strconv.ParseFloat(s.buf.String(), 64)
		if err != nil {
			s.throw(err)
		}
		s.endToken(FLOAT, v)
		return nil
	}

	c := s.read()
	switch {
	case unicode.IsNumber(c):
		s.buf.WriteRune(c)
		return s.float

	default:
		s.unread()
		v, err := strconv.ParseFloat(s.buf.String(), 64)
		if err != nil {
			s.throw(err)
		}
		s.endToken(FLOAT, v)
		return nil
	}
}

func (s *Scanner) char(eof bool) state {
	if eof {
		s.throw(errors.New("unterimnated char literal"))
		return nil
	}

	c := s.read()
	if c == '\\' {
		c = s.readEscapeSeq()
	}

	if s.read() != '\'' {
		s.throw(errors.New("char literal is too long"))
	}

	s.endToken(INT, c)
	return nil
}

func (s *Scanner) symbol(eof bool) state {
	if eof {
		str := s.buf.String()
		t, ok := symbols[str]
		if !ok {
			s.throw(fmt.Errorf("unexpected characters: %q", str))
		}
		s.endToken(t, str)
		return nil
	}

	str := s.buf.String()
	if str == "#" {
		return s.singleLineComment
	}

	s.buf.WriteRune(s.read())
	str = s.buf.String()

	if t, ok := symbols[str]; ok {
		s.endToken(t, str)
		return nil
	}
	if t, ok := symbols[str[:1]]; ok {
		s.unread()
		s.endToken(t, str[:1])
		return nil
	}

	s.throw(fmt.Errorf("unexpected characters: %q", str))
	return nil
}

func (s *Scanner) singleLineComment(eof bool) state {
	if eof || (s.read() == '\n') {
		if s.read() != '#' {
			s.unread()
			return s.whitespace
		}

		return s.singleLineComment
	}

	// TODO: Yield comments?

	return s.singleLineComment
}

func (s *Scanner) startToken() {
	s.tline = s.line + 1
	s.tcol = s.col
}

func (s *Scanner) endToken(t Type, v any) {
	s.tok = Token{
		Line: s.tline,
		Col:  s.tcol,
		Type: t,
		Val:  v,
	}
}

type state func(eof bool) state

type stateErr struct{ err error }

func (s stateErr) Error() string { return s.err.Error() }
func (s stateErr) Unwrap() error { return s.err }
