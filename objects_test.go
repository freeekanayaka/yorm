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

	str, err := objects.Struct("foo", "Foo")
	require.NoError(t, err)

	assert.Equal(t, "foo", str.Package)
	assert.Equal(t, "Foo", str.Name)
}

func newPackages(t *testing.T) map[string]*ast.Package {
	t.Helper()

	src := `
package foo

type Foo struct {
}
`
	parser := yorm.NewParser()
	err := parser.LoadFile("foo.go", src)
	require.NoError(t, err)

	pkgs, err := parser.Parse()
	require.NoError(t, err)

	return pkgs
}
