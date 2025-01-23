package types

type TSharkClosure struct {
	ISharkType
	FuncType ISharkType
}

func (t TSharkClosure) SharkTypeString() string {
	if t.FuncType == nil {
		return "closure<>"
	}
	return t.FuncType.SharkTypeString()
}

func (t TSharkClosure) Is(sharkType ISharkType) bool {
	if tt, ok := sharkType.(TSharkClosure); ok {
		if tt.FuncType == nil {
			return true
		}
		if t.FuncType == nil {
			return false
		}
		return t.FuncType.Is(tt.FuncType)
	}
	return false
}
