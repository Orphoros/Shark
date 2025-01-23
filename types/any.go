package types

type TSharkAny struct {
	ISharkType
}

func (TSharkAny) SharkTypeString() string { return "any" }

func (TSharkAny) Is(sharkType ISharkType) bool {
	switch sharkType.(type) {
	case TSharkSpread, TSharkOptional:
		return false
	default:
		return true
	}
}
