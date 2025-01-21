package types

type TSharkI64 struct {
	ISharkType
}

func (TSharkI64) SharkTypeString() string { return "i64" }

func (TSharkI64) Is(sharkType ISharkType) bool {
	switch t := sharkType.(type) {
	case TSharkI64:
		return true
	case TSharkVariadic:
		return t.Is(TSharkI64{})
	default:
		return false
	}
}
