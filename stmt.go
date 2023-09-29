package stele

// A Stmt is an executable piece of code.
type Stmt interface {
	// Eval evaluates the Stmt in the context of the given State and
	// possibly returns a value. If the Stmt does not return values, the
	// returned Value will be the zero-value Value.
	Eval(*State) Value
}

// An Expr is an executable piece of code that returns a value. All
// Exprs are Stmts, but not all Stmts are Exprs.
type Expr interface {
	// Type is the type of the value returned by the Expr.
	Type() Type

	// Eval evaluates the Expr in the context of the given State and
	// returns a value. Unlike a Stmt, an Expr _must_ return a valid
	// Value.
	Eval(*State) Value
}

// A Block represents a series of statements.
type Block struct {
	Stmts []Stmt
}

func (b Block) Eval(state *State) Value {
	for _, stmt := range b.Stmts {
		stmt.Eval(state)
	}
	return Value{}
}

// An Assign is an assignment statement. It evaluates an expression
// and assigns it to a variable.
type Assign struct {
	Recv string
	ID   string
	Val  Expr
}

func (a Assign) Eval(state *State) Value {
	panic("Not implemented.")
}
