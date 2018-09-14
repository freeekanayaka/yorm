package yorm_test

import (
	"testing"

	"github.com/freeekanayaka/yorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_StructQuery(t *testing.T) {
	templates := yorm.NewTemplates()
	naming := yorm.DefaultNaming()
	generator := yorm.NewGenerator(templates, naming)

	fields := []*yorm.Field{
		{Name: "Email", Type: yorm.String},
		{Name: "Age", Type: yorm.Int},
	}

	s := yorm.AnonymousStruct(fields)
	err := generator.StructQuery("f", s)
	require.NoError(t, err)

	output, err := generator.Output()
	require.NoError(t, err)
	assert.Equal(t, `
func f(ctx context.Context, stmt *sql.Stmt, args ...interface{}) ([]struct {
	Email string
	Age   int
}, error) {
	objects := make([]struct {
		Email string
		Age   int
	}, 0)

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "f: run query")
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		objects = append(objects, struct {
			Email string
			Age   int
		}{})

		if err := rows.Scan(
			&objects[i].Email,
			&objects[i].Age,
		); err != nil {
			return nil, errors.Wrap(err, "f: scan row")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "f: rows error")
	}

	return objects, nil
}
`, string(output))
}

func TestGenerator_StructQuery_Fields(t *testing.T) {
	templates := yorm.NewTemplates()
	naming := yorm.DefaultNaming()
	generator := yorm.NewGenerator(templates, naming)

	fields := []*yorm.Field{
		{Name: "Email", Type: yorm.String},
		{Name: "Age", Type: yorm.Int},
	}

	s := yorm.AnonymousStruct(fields)
	err := generator.StructQuery("f", s, "Email")
	require.NoError(t, err)

	output, err := generator.Output()
	require.NoError(t, err)
	assert.Equal(t, `
func f(ctx context.Context, stmt *sql.Stmt, args ...interface{}) ([]struct {
	Email string
	Age   int
}, error) {
	objects := make([]struct {
		Email string
		Age   int
	}, 0)

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "f: run query")
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		objects = append(objects, struct {
			Email string
			Age   int
		}{})

		if err := rows.Scan(
			&objects[i].Email,
		); err != nil {
			return nil, errors.Wrap(err, "f: scan row")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "f: rows error")
	}

	return objects, nil
}
`, string(output))
}
