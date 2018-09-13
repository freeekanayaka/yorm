package yorm_test

import (
	"testing"

	"github.com/freeekanayaka/yorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_LoadFile(t *testing.T) {
	src := `
package foo

type Foo struct {
}
`
	parser := yorm.NewParser()
	err := parser.LoadFile("foo.go", src)
	require.NoError(t, err)
}

func TestParser_Parse(t *testing.T) {
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

	assert.Len(t, pkgs, 1)
	assert.Contains(t, pkgs, "foo")
}
