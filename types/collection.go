package types

import (
	"bytes"
)

type TSharkCollection struct {
	ISharkType
	ISharkCollection
	Collection []ISharkType
}

func (t TSharkCollection) SharkTypeString() string {
	if t.Collection == nil {
		return "collection<>"
	}
	var buf bytes.Buffer
	buf.WriteString("collection<")
	for i, collect := range t.Collection {
		buf.WriteString(collect.SharkTypeString())
		if i != len(t.Collection)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(">")
	return buf.String()
}

func (t TSharkCollection) Is(sharkType ISharkType) bool {
	if collection, ok := sharkType.(ISharkCollection); ok {
		if t.Collection == nil {
			return true
		}
		if len(t.Collection) == 1 && t.Collection[0].Is(TSharkSpread{Type: TSharkAny{}}) {
			return true
		}
		if len(t.Collects()) == len(collection.Collects()) {
			for i, collect := range t.Collects() {
				if !collect.Is(collection.Collects()[i]) {
					return false
				}
			}
			return true
		}
	}

	return false
}

func (t TSharkCollection) Collects() []ISharkType {
	return t.Collection
}
