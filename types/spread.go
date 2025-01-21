package types

type TSharkSpread struct {
	ISharkType
	Type ISharkType
}

func (t TSharkSpread) SharkTypeString() string {
	if t.Type == nil {
		return "..."
	}
	return "..." + t.Type.SharkTypeString()
}

func (t TSharkSpread) Is(sharkType ISharkType) bool {
	if sharkType == nil {
		return true
	}
	return sharkType.Is(t.Type)
}

func (t TSharkSpread) Primitive() {}
