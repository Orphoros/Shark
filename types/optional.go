package types

type TSharkOptional struct {
	ISharkType
	Type ISharkType
}

func (t TSharkOptional) SharkTypeString() string {
	if t.Type == nil {
		return "?"
	}
	return t.Type.SharkTypeString() + "?"
}

func (t TSharkOptional) Is(sharkType ISharkType) bool {
	switch sharkType := sharkType.(type) {
	case TSharkOptional:
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
