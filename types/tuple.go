package types

import "bytes"

type TSharkTuple struct {
	ISharkType
	Collects []ISharkType
}

func (t TSharkTuple) CollectionOf() ISharkType {
	if t.Collects == nil {
		return TSharkNull{}
	}
	for i := 1; i < len(t.Collects); i++ {
		if !t.Collects[i].Is(t.Collects[i-1]) {
			return TSharkAny{}
		}
	}
	return t.Collects[0]
}

func (t TSharkTuple) SharkTypeString() string {
	if t.Collects == nil {
		return "tuple<>"
	}
	var buf bytes.Buffer
	buf.WriteString("tuple<")
	for i, collect := range t.Collects {
		buf.WriteString(collect.SharkTypeString())
		if i != len(t.Collects)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(">")
	return buf.String()
}

func (t TSharkTuple) Is(sharkType ISharkType) bool {
	switch sharkType := sharkType.(type) {
	case TSharkTuple:
		if sharkType.Collects == nil {
			return true
		}
		if len(sharkType.Collects) != len(t.Collects) {
			return false
		}
		for i := 0; i < len(t.Collects); i++ {
			if !t.Collects[i].Is(sharkType.Collects[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
