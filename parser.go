package yorm

import (
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type Parser struct {
	fset  *token.FileSet
	files map[string]map[string]*ast.File
}

func (p *Parser) LoadFile(filename string, src interface{}) error {
	file, err := parser.ParseFile(p.fset, filename, src, parser.AllErrors)
	if err != nil {
		return errors.Wrapf(err, "parse %q", filename)
	}

	name := file.Name.Name
	files, ok := p.files[name]
	if !ok {
		files = map[string]*ast.File{}
		p.files[name] = files
	}

	files[name] = file

	return nil
}

func (p *Parser) Parse() (map[string]*ast.Package, error) {
	pkgs := map[string]*ast.Package{}

	for name, files := range p.files {
		pkg, err := ast.NewPackage(p.fset, files, nil, nil)
		if err != nil && !isBenignParseError(err) {
			return nil, errors.Wrapf(err, "package %q", name)
		}
		pkgs[name] = pkg
	}

	return pkgs, nil
}

func NewParser() *Parser {
	return &Parser{
		fset:  token.NewFileSet(),
		files: map[string]map[string]*ast.File{},
	}
}

func isBenignParseError(err error) bool {
	errs, ok := err.(scanner.ErrorList)
	if !ok {
		return false
	}

	for _, err := range errs {
		if !strings.Contains(err.Msg, "undeclared name") {
			return false
		}
	}

	return true
}
