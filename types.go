package yorm

import (
	"fmt"
	"strings"
)

// Struct holds information about a Go struct.
type Struct struct {
	Type   string
	Fields []*Field
}

// Get returns the field with the given name, if any.
func (s *Struct) Get(name string) *Field {
	for _, field := range s.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

// AnonymousStruct is a helper to create an anonymous struct with the given fields.
func AnonymousStruct(fields []*Field) *Struct {
	declarations := make([]string, len(fields))

	for i, field := range fields {
		declarations[i] = fmt.Sprintf("%s %s", field.Name, field.Type)
	}

	return &Struct{
		Type:   fmt.Sprintf("struct{%s}", strings.Join(declarations, "\n")),
		Fields: fields,
	}
}

// Field holds information about a field in a Go struct.
type Field struct {
	Name string
	Type interface{}
}

// IsScalar returns true if the Type of the field is a Scalar.
func (f *Field) IsScalar() bool {
	_, ok := f.Type.(Scalar)
	return ok
}

// Scalar types are the ones that can be mapped directly to a SQL column.
type Scalar int

// List of scalar type codes.
const (
	Int Scalar = iota
	String
)

func (s Scalar) String() string {
	switch s {
	case Int:
		return "int"
	case String:
		return "string"
	default:
		panic("unknown scalar type code")
	}
}
