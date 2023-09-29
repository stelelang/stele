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
