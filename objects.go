package yorm

import (
	"fmt"
	"go/ast"

	"github.com/pkg/errors"
)

type Objects struct {
	pkgs map[string]*ast.Package
}

func NewObjects(pkgs map[string]*ast.Package) *Objects {
	return &Objects{
		pkgs: pkgs,
	}
}

func (c *Objects) Struct(pkg string, name string) (*Struct, error) {
	p, ok := c.pkgs[pkg]
	if !ok {
		return nil, fmt.Errorf("package not found")
	}

	obj := p.Scope.Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("identifier not found")
	}

	return astToStruct(obj)
}

func astToStruct(obj *ast.Object) (*Struct, error) {
	typ, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("identifier is not a type")
	}

	str, ok := typ.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("identifier is not a struct type")
	}

	s := &Struct{
		Type:   obj.Name,
		Fields: []*Field{},
	}

	for i, field := range str.Fields.List {
		if len(field.Names) == 0 {
			// Check if this is a parent struct.
			ident, ok := field.Type.(*ast.Ident)
			if !ok {
				continue
			}

			p, err := astToStruct(ident.Obj)
			if err != nil {
				return nil, err
			}
			s.Fields = append(s.Fields, p.Fields...)

			continue
		}

		f, err := astToField(field)
		if err != nil {
			return nil, errors.Wrapf(err, "field %d", i)
		}
		s.Fields = append(s.Fields, f)
	}

	return s, nil
}

func astToField(field *ast.Field) (*Field, error) {
	if len(field.Names) != 1 {
		return nil, fmt.Errorf("invalid struct field")
	}

	name := field.Names[0]

	var t interface{}
	var err error

	switch typ := field.Type.(type) {
	case *ast.Ident:
		t, err = astToScalar(typ)
	}

	if err != nil {
		return nil, err
	}

	f := &Field{Name: name.Name, Type: t}

	return f, nil
}

func astToScalar(ident *ast.Ident) (Scalar, error) {
	var scalar Scalar

	switch ident.String() {
	case "int":
		scalar = Int
	case "string":
		scalar = String
	default:
		return -1, fmt.Errorf("invalid scalar type")
	}

	return scalar, nil
}
