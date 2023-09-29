package stele

import "strings"

type Scope struct {
	parent *Scope
	ids    func() []string
	get    func(string) Declaration
}

var (
	// fakeScope is a signal value to help make the zero value useful.
	fakeScope Scope

	// TODO: Add predeclared identifiers.
	rootScope = Scope{
		parent: &fakeScope,
		get:    func(string) Declaration { return nil },
	}
)

func RootScope() Scope {
	return rootScope
}

func (s Scope) Parent() *Scope {
	switch s.parent {
	case nil:
		return &rootScope
	case &fakeScope:
		return nil
	default:
		return s.parent
	}
}

func (s Scope) Add(d Declaration) Scope {
	return Scope{
		parent: &s,
		get: func(id string) Declaration {
			if id == d.ID() {
				return d
			}
			return nil
		},
	}
}

func (s Scope) AddAll(d []Declaration) Scope {
	m := make(map[string]Declaration, len(d))
	for _, d := range d {
		m[d.ID()] = d
	}

	return Scope{
		parent: &s,
		get: func(id string) Declaration {
			return m[id]
		},
	}
}

func (s Scope) Get(id string) Declaration {
	d := s.get(id)
	if d != nil {
		return d
	}

	p := s.Parent()
	if p == nil {
		return nil
	}
	return p.Get(id)
}

type Declaration interface {
	ID() string
	Type() Type
	Mutable() bool
	Exported() bool
}

type ImportDecl struct {
	Name string
	Path string
}

func (d ImportDecl) ID() string     { return d.Name }
func (d ImportDecl) Type() Type     { panic("Not implemented.") }
func (d ImportDecl) Mutable() bool  { return false }
func (d ImportDecl) Exported() bool { return false }

type LetDecl struct {
	Name string
	T    Type
	RHS  Expr
}

func (d LetDecl) ID() string     { return d.Name }
func (d LetDecl) Type() Type     { return d.T }
func (d LetDecl) Mutable() bool  { return !strings.HasSuffix(d.Name, "!") }
func (d LetDecl) Exported() bool { return !strings.HasPrefix(d.Name, "_") }

type Type struct {
	Name     string
	Features []Feature
}

//go:generate go run golang.org/x/tools/cmd/stringer -type FeatureType

type FeatureType int

const (
	InvalidFeature FeatureType = iota
	LetFeature
	FuncFeature
	MemLayoutFeature
)

type Feature struct {
	Type FeatureType
	Name string

	Args   []Type
	Return Type
}

type Expr interface {
	Type() Type
	//Eval(*State) Value
}
