package yorm

import "html/template"

// Templates is a registry of all code templates used by yorm.
type Templates struct {
	templates map[int]*template.Template
}

// NewTemplates returns a new templates registry initialized with default templates.
func NewTemplates() *Templates {
	return &Templates{templates: defaultTemplates}
}

// Get returns the template with the given code.
func (t *Templates) Get(code int) *template.Template {
	return t.templates[code]
}

// Set sets the template with the given code.
func (t *Templates) Set(code int, template *template.Template) {
	t.templates[code] = template
}

// Template codes.
const (
	HeaderTmpl = iota
	QueryTmpl
)

var defaultTemplates = map[int]*template.Template{
	HeaderTmpl: headerTmpl,
	QueryTmpl:  queryTmpl,
}

var headerTmpl = template.Must(template.New("").Parse(`package {{.Package}}

// The code below was automatically generated - DO NOT EDIT!

import (
        {{- range .Imports }}
                "{{ . }}"
        {{- end}}
)
`))

var queryTmpl = template.Must(template.New("").Parse(`
func {{.Name}} (ctx context.Context, stmt *sql.Stmt, args ...interface{}) ([]{{.Struct.Type}}, error) {
        objects := make([]{{.Struct.Type}}, 0)

        rows, err := stmt.QueryContext(ctx, args...)
        if err != nil {
                return nil, errors.Wrap(err, "{{.Name}}: run query")
        }
        defer rows.Close()

        for i := 0; rows.Next(); i++ {
                objects = append(objects, {{.Struct.Type}}{})

                if err := rows.Scan(
                {{- range .Fields }}
                        &objects[i].{{ . }},
                {{- end}}
                ); err != nil {
                        return nil, errors.Wrap(err, "{{.Name}}: scan row")
                }
        }

        if err := rows.Err(); err != nil {
                return nil, errors.Wrap(err, "{{.Name}}: rows error")
        }

        return objects, nil
}
`))
