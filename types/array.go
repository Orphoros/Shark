package types

type TSharkArray struct {
	ISharkType
	Collects ISharkType
}

func (t TSharkArray) SharkTypeString() string {
	if t.Collects == nil {
		return "array<>"
	}
	return "array<" + t.Collects.SharkTypeString() + ">"
}

func (t TSharkArray) Is(sharkType ISharkType) bool {
	switch sharkType := sharkType.(type) {
	case TSharkArray:
		if sharkType.Collects == nil {
			return true
		}
		if sharkType.Collects.Is(t.Collects) {
			return true
		}
		return false
	default:
		return false
	}
}
