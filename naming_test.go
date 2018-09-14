package yorm_test

import (
	"testing"

	"github.com/freeekanayaka/yorm"
	"github.com/stretchr/testify/assert"
)

func TestDefaultNaming(t *testing.T) {
	naming := yorm.DefaultNaming{}

	assert.Equal(t, "foo_bar", naming.Column("FooBar"))
	assert.Equal(t, "FooBar", naming.Field("foo_bar"))
}
