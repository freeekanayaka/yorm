package yorm

import (
	"bytes"
	"strings"
	"unicode"
)

// Naming conventions to use when converting between Go struct/field names and
// SQL tables/columns.
type Naming interface {
	Column(field string) string // Convert a Go field name to a SQL column
	Field(column string) string // Convert a SQL column to a Go field name
}

// DefaultNaming implements yorm's default naming convention, which is to use
// camel case names for Go struct/fields and snake case names for SQL names.
func DefaultNaming() Naming {
	return &defaultNaming{}
}

type defaultNaming struct {
}

// Column converts a Go field name (camel case) to a SQL column (snake case)
func (n *defaultNaming) Column(field string) string {
	return Snake(field)
}

// Field converts a SQL column name (snake case) to a Go field name (camel case)
func (n *defaultNaming) Field(column string) string {
	return Camel(column)
}

// Camel converts s from sname case to camel case ("foo_bar" -> "FooBar")
func Camel(s string) string {
	words := strings.Split(s, "_")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// Snake converts to snake case ("FooBar" -> "foo_bar")
func Snake(name string) string {
	var ret bytes.Buffer

	multipleUpper := false
	var lastUpper rune
	var beforeUpper rune

	for _, c := range name {
		// Non-lowercase character after uppercase is considered to be uppercase too.
		isUpper := (unicode.IsUpper(c) || (lastUpper != 0 && !unicode.IsLower(c)))

		if lastUpper != 0 {
			// Output a delimiter if last character was either the
			// first uppercase character in a row, or the last one
			// in a row (e.g. 'S' in "HTTPServer").  Do not output
			// a delimiter at the beginning of the name.
			firstInRow := !multipleUpper
			lastInRow := !isUpper

			if ret.Len() > 0 && (firstInRow || lastInRow) && beforeUpper != '_' {
				ret.WriteByte('_')
			}
			ret.WriteRune(unicode.ToLower(lastUpper))
		}

		// Buffer uppercase char, do not output it yet as a delimiter
		// may be required if the next character is lowercase.
		if isUpper {
			multipleUpper = (lastUpper != 0)
			lastUpper = c
			continue
		}

		ret.WriteRune(c)
		lastUpper = 0
		beforeUpper = c
		multipleUpper = false
	}

	if lastUpper != 0 {
		ret.WriteRune(unicode.ToLower(lastUpper))
	}
	return string(ret.Bytes())
}
