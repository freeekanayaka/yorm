package yorm

import (
	"fmt"
	"go/ast"
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
		return nil, fmt.Errorf("package %q not found", pkg)
	}

	obj := p.Scope.Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("struct %q not found in package %q", name, pkg)
	}

	typ, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("identifier %q in package %q is not a type", name, pkg)
	}

	_, ok = typ.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("type %q in package %q is not a struct", name, pkg)
	}

	return &Struct{Package: p.Name, Name: obj.Name}, nil
}

type Struct struct {
	Package string
	Name    string
}
