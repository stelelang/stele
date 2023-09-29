package stele

// Scope represents the declarations available for a given piece of
// code. It is a compile-time structure. The run-time equivalent is
// Frame, which tracks the actual values of declarations.
//
// A zero-value scope is a direct child of the scope retunred by
// [RootScope].
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

// RootScope returns the base scope that all scopes are the child of.
// It is not usually necessary to call this directly, as the
// zero-value of a Scope is considered to be an empty child scope of
// the one returned by this function.
func RootScope() Scope {
	return rootScope
}

// Parent returns the parent of the current Scope, or nil if the
// current scope is the root scope.
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

// Add returns a new child scope continaing d.
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

// AddAll returns a new child scope containing all of the Declarations
// in d.
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

// Get searches up the scope hierarchy, returning the first
// encountered Declaration with the given ID. If no such Declaration
// exists, it returns nil.
func (s Scope) Get(id string) Declaration {
	if s.get != nil {
		d := s.get(id)
		if d != nil {
			return d
		}
	}

	p := s.Parent()
	if p == nil {
		return nil
	}
	return p.Get(id)
}

// A Declaration is something declared in a scope in a Stele program.
// This includes variables declared with let, functions declared with
// func, imports, etc.
type Declaration interface {
	// ID returns the identifier that is being declared.
	ID() string

	// Type returns the type of the value bound to the identifier.
	Type() Type

	// Mutable is true if the value assigned to the identifier can be
	// changed at run-time.
	Mutable() bool

	// Exported is true if the identifier can be used by other scripts
	// importing the current one. This is only applicable to top-level
	// identifiers.
	Exported() bool
}
