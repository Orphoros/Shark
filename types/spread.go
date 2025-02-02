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
	switch sharkType := sharkType.(type) {
	case TSharkSpread:
		if sharkType.Type == nil {
			return true
		}
		if t.Type == nil {
			return false
		}
		if t.Type.Is(sharkType.Type) {
			return true
		}
		return false
	default:
		if t.Type == nil {
			return false
		}

		return t.Type.Is(sharkType)
	}
}
