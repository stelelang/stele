package stele

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
