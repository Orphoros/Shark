package types

import (
	"bytes"
)

type TSharkFuncType struct {
	ISharkType
	ReturnT  ISharkType
	ArgsList []ISharkType
}

func (t TSharkFuncType) SharkTypeString() string {
	var buf bytes.Buffer
	buf.WriteString("func<(")
	if len(t.ArgsList) == 0 && t.ReturnT == nil {
		buf.WriteString(")>")
		return buf.String()
	}

	if len(t.ArgsList) > 0 {
		for i, arg := range t.ArgsList {
			buf.WriteString(arg.SharkTypeString())
			if i != len(t.ArgsList)-1 {
				buf.WriteString(",")
			}
		}
	}

	if t.ReturnT != nil {
		buf.WriteString(")->")
		buf.WriteString(t.ReturnT.SharkTypeString())
		buf.WriteString(">")
	} else {
		buf.WriteString(")>")
	}

	return buf.String()
}

func (t TSharkFuncType) Is(sharkType ISharkType) bool {
	switch sharkType := sharkType.(type) {
	case TSharkFuncType:
		if sharkType.ArgsList == nil && sharkType.ReturnT == nil {
			return true
		}
		if len(t.ArgsList) != len(sharkType.ArgsList) {
			return false
		}
		if len(t.ArgsList) > 0 {
			for i, arg := range t.ArgsList {
				if !arg.Is(sharkType.ArgsList[i]) {
					return false
				}
			}
		}
		if sharkType.ReturnT == nil && t.ReturnT == nil {
			return true
		}
		return t.ReturnT.Is(sharkType.ReturnT)
	default:
		return false
	}
}
