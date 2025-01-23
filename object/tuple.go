package object

import (
	"bytes"
	"shark/types"
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

func (t *Tuple) Type() types.ISharkType {
	if len(t.Elements) == 0 {
		return types.TSharkTuple{}
	}

	if t.Elements[0] == nil {
		return types.TSharkTuple{Collection: []types.ISharkType{types.TSharkNull{}}}
	}

	var collects []types.ISharkType
	for _, e := range t.Elements {
		collects = append(collects, e.Type())
	}

	return types.TSharkTuple{Collection: collects}
}
