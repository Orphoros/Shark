package object

import (
	"bytes"
	"strings"
)

type Tuple struct {
	Elements []Object
}

func (t *Tuple) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range t.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(")")

	return out.String()
}

func (a *Tuple) Type() Type { return TUPLE_OBJ }
