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

// Query generates a function that given a sql.Stmt and a struct
// definition, executes the query and returns a slice of instances of that
// struct with the query columns mapped to the struct fields.
//
// If the fields slice is non-empty, only those fields will be filled.
func (g *Generator) Query(name string, s *Struct, fields ...string) error {
	template := g.templates.Get(QueryTmpl)

	if len(fields) == 0 {
		fields = make([]string, len(s.Fields))
		for i, field := range s.Fields {
			fields[i] = field.Name
		}
	} else {
		// Ensure that the provided field names are valid.
		for _, name := range fields {
			field := s.Get(name)
			if field == nil {
				return fmt.Errorf("struct %s has no %s field", s.Type, name)
			}
		}
	}

	err := template.Execute(g.buf, struct {
		Name   string
		Struct *Struct
		Fields []string
	}{
		Name:   name,
		Struct: s,
		Fields: fields,
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
