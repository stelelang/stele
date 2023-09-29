package stele

type Stmt interface {
	Eval(*State) Value
}

type Expr interface {
	Type() Type
	Eval(*State) Value
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

type Assign struct {
	ID  string
	Val Expr
}

func (a Assign) Eval(state *State) Value {
	panic("Not implemented.")
}
