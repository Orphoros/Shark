package object

import (
	"bytes"
	"shark/types"
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

func (a *Array) Type() types.ISharkType {
	if len(a.Elements) == 0 {
		return types.TSharkArray{Collects: types.TSharkAny{}}
	}

	if a.Elements[0] == nil {
		return types.TSharkArray{Collects: types.TSharkNull{}}
	}

	if len(a.Elements) > 1 {
		for i := 1; i < len(a.Elements); i++ {
			if !a.Elements[i].Type().Is(a.Elements[i-1].Type()) {
				return types.TSharkArray{Collects: types.TSharkAny{}}
			}
		}
		return types.TSharkArray{Collects: a.Elements[0].Type()}
	}

	return types.TSharkArray{Collects: a.Elements[0].Type()}
}
