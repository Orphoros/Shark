package object

import (
	"bytes"
	"strings"
)

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (a *Array) Type() Type { return ARRAY_OBJ }
