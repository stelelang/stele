package ast

import "deedles.dev/stele"

type Int struct {
	Val int64
}

func (i Int) Type() stele.Type {
	// TODO: Return a type for int literals.
	return stele.Type{}
}

func (i Int) Eval(state *stele.State) stele.Value {
	panic("Not implemented.")
}
