package yorm_test

import (
	"go/ast"
	"testing"

	"github.com/freeekanayaka/yorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjects_Struct(t *testing.T) {
	pkgs := newPackages(t)

	objects := yorm.NewObjects(pkgs)

	s, err := objects.Struct("model", "User")
	require.NoError(t, err)

	assert.Equal(t, "User", s.Type)

	assert.Len(t, s.Fields, 2)

	f0 := s.Fields[0]
	f1 := s.Fields[1]

	assert.Equal(t, "Email", f0.Name)
	assert.Equal(t, yorm.String, f0.Type)

	assert.Equal(t, "Age", f1.Name)
	assert.Equal(t, yorm.Int, f1.Type)
}

func newPackages(t *testing.T) map[string]*ast.Package {
	t.Helper()

	src := `
package model

type User struct {
        Email string
        Age   int
}
`
	parser := yorm.NewParser()
	err := parser.LoadFile("foo.go", src)
	require.NoError(t, err)

	pkgs, err := parser.Parse()
	require.NoError(t, err)

	return pkgs
}
