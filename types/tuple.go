package types

import "bytes"

type TSharkTuple struct {
	ISharkType
	ISharkCollection
	Collection []ISharkType
}

func (t TSharkTuple) SharkTypeString() string {
	if t.Collection == nil {
		return "tuple<>"
	}
	var buf bytes.Buffer
	buf.WriteString("tuple<")
	for i, collect := range t.Collection {
		buf.WriteString(collect.SharkTypeString())
		if i != len(t.Collection)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(">")
	return buf.String()
}

func (t TSharkTuple) Is(sharkType ISharkType) bool {
	switch sharkType := sharkType.(type) {
	case TSharkTuple:
		if sharkType.Collection == nil {
			return true
		}
		if len(sharkType.Collection) != len(t.Collection) {
			return false
		}
		for i := 0; i < len(t.Collection); i++ {
			if !t.Collection[i].Is(sharkType.Collection[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (t TSharkTuple) Collects() []ISharkType {
	return t.Collection
}
