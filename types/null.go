package types

type TSharkNull struct {
	ISharkType
}

func (TSharkNull) SharkTypeString() string { return "null" }

func (TSharkNull) Is(sharkType ISharkType) bool {
	switch t := sharkType.(type) {
	case TSharkNull:
		return true
	case TSharkVariadic:
		return t.Is(TSharkNull{})
	default:
		return false
	}
}
