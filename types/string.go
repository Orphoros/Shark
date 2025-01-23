package types

type TSharkString struct {
	ISharkType
	ISharkCollection
}

func (TSharkString) SharkTypeString() string { return "string" }

func (TSharkString) Is(sharkType ISharkType) bool {
	switch t := sharkType.(type) {
	case TSharkString:
		return true
	case TSharkVariadic:
		return t.Is(TSharkString{})
	default:
		return false
	}
}

func (t TSharkString) Collects() []ISharkType {
	return []ISharkType{t}
}
