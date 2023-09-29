package ast

import (
	"strings"

	"deedles.dev/stele"
)

type Import struct {
	Name string
	Path string
}

func (d Import) ID() string       { return d.Name }
func (d Import) Type() stele.Type { panic("Not implemented.") }
func (d Import) Mutable() bool    { return false }
func (d Import) Exported() bool   { return false }

type Let struct {
	Name   string
	T      stele.Type
	Assign *stele.Assign
}

func (d Let) ID() string       { return d.Name }
func (d Let) Type() stele.Type { return d.T }
func (d Let) Mutable() bool    { return !strings.HasSuffix(d.Name, "!") }
func (d Let) Exported() bool   { return !strings.HasPrefix(d.Name, "_") }
