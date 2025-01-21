package types

type TSharkBool struct {
	ISharkType
}

func (TSharkBool) SharkTypeString() string { return "bool" }

func (TSharkBool) Is(sharkType ISharkType) bool {
	switch t := sharkType.(type) {
	case TSharkBool:
		return true
	case TSharkVariadic:
		return t.Is(TSharkBool{})
	default:
		return false
	}
}
