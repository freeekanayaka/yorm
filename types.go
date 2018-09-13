package yorm

type Struct struct {
	Fields []*Field
}

type Field struct {
	Name string
	Type interface{}
}

type Scalar int

const (
	Int Scalar = iota
	String
)
