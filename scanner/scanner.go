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

	buf         strings.Builder
	tok         Token
	tline, tcol int
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	s.tok = Token{}
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
	s.col++
	defer func() {
		if c == '\n' {
			s.line++
			s.col = 0
		}
	}()

	c, _, err := s.r.ReadRune()
	if err != nil {
		s.throw(err)
	}
	return c
}

func (s *Scanner) unread() {
	err := s.r.UnreadRune()
	if err != nil {
		s.throw(err)
	}
}

func (s *Scanner) whitespace(eof bool) state {
	if eof {
		return nil
	}

	c := s.read()
	switch {
	case unicode.IsSpace(c):
		return s.whitespace

	case unicode.IsLetter(c) || (c == '_'):
		s.buf.WriteRune(c)
		s.startToken()
		return s.ident

	case c == '"':
		s.startToken()
		return s.string

	case unicode.IsNumber(c):
		s.startToken()
		s.buf.WriteRune(c)
		return s.int

	default:
		s.throw(fmt.Errorf("unexpected character: %q", c))
		return nil
	}
}

func (s *Scanner) ident(eof bool) state {
	if eof {
		str := s.buf.String()
		s.endToken(tokenType(str), str)
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
		s.endToken(tokenType(str), str)
		return nil

	default:
		s.unread()
		str := s.buf.String()
		s.endToken(tokenType(str), str)
		return nil
	}
}

func (s *Scanner) string(eof bool) state {
	if eof {
		s.throw(errors.New("unterminated string literal"))
		return nil
	}

	c := s.read()
	if c == '"' {
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

func (s *Scanner) startToken() {
	s.tline = s.line + 1
	s.tcol = s.col
}

func (s *Scanner) endToken(t int, v any) {
	s.tok = Token{
		Line: s.tline,
		Col:  s.tcol,
		Type: t,
		Val:  v,
	}
}

type state func(eof bool) state

type stateErr struct{ err error }
