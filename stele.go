package stele

type State struct{}

type Type struct {
	Name     string
	Features []Feature
}

func (t Type) Valid() bool {
	return t.Name != ""
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

type Value struct {
	Type Type
	Val  any
}

func (v Value) Valid() bool {
	return v.Val != nil
}

type Block struct {
	T     Type
	Stmts []Stmt
}

func (b Block) Type() Type {
	return b.T
}

func (b Block) Eval(state *State) Value {
	for _, stmt := range b.Stmts {
		v := stmt.Eval(state)
		if v.Valid() {
			return v
		}
	}
	return Value{}
}

type Stmt interface {
	Eval(*State) Value
}

type Expr interface {
	Type() Type
	Eval(*State) Value
}

type Assign struct {
	ID  string
	Val Expr
}
