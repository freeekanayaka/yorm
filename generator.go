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

// Header generates the header of a source file.
func (g *Generator) Header(pkg string, imports []string) error {
	template := g.templates.Get(HeaderTmpl)
	err := template.Execute(g.buf, struct {
		Package string
		Imports []string
	}{
		Package: pkg,
		Imports: imports,
	})
	if err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}

// Query generates a function that given a sql.Stmt and a struct
// definition, executes the query and returns a slice of instances of that
// struct with the query columns mapped to the struct fields.
//
// If the fields slice is non-empty, only those fields will be filled.
func (g *Generator) Query(name string, s *Struct, fields ...string) error {
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

	template := g.templates.Get(QueryTmpl)
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
