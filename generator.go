package yorm

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/pkg/errors"
)

// Generator generates source code for mapping Go objects to SQL tables.
type Generator struct {
	templates *Templates
	naming    Naming
	buf       *bytes.Buffer // Buffer for accumulating generated source code.
}

// NewGenerator create a code generator.
func NewGenerator(templates *Templates, naming Naming) *Generator {
	return &Generator{
		templates: templates,
		naming:    naming,
		buf:       bytes.NewBuffer(nil),
	}
}

// StructQuery generates a function that given a sql.Stmt and a slice
// of field names, executes the query and returns a slice of anonymous
// structures with those fields.
func (g *Generator) StructQuery(name string, s *Struct) error {
	template := g.templates.Get(StructQueryTmpl)

	err := template.Execute(g.buf, struct {
		Name   string
		Struct *Struct
	}{
		Name:   name,
		Struct: s,
	})
	if err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}

// A appends code to the current line.
func (g *Generator) A(format string, a ...interface{}) {
	fmt.Fprintf(g.buf, format, a...)
}

// N accumulates a single new line.
func (g *Generator) N() {
	fmt.Fprintf(g.buf, "\n")
}

// L accumulates a single line of source code.
func (g *Generator) L(format string, a ...interface{}) {
	g.A(format, a...)
	g.N()
}

// Output returns the generated source code.
func (g *Generator) Output() ([]byte, error) {
	data := g.buf.Bytes()
	code, err := format.Source(data)
	if err != nil {
		msg := "Can't format generated source code:\n\n%s"
		return nil, errors.Wrapf(err, msg, data)
	}
	return code, nil
}
