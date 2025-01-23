package types

type TSharkArray struct {
	ISharkType
	ISharkCollection
	Collection ISharkType
}

func (t TSharkArray) SharkTypeString() string {
	if t.Collection == nil {
		return "array<>"
	}
	return "array<" + t.Collection.SharkTypeString() + ">"
}

func (t TSharkArray) Is(sharkType ISharkType) bool {
	switch sharkType := sharkType.(type) {
	case TSharkArray:
		if sharkType.Collection == nil {
			return true
		}

		if t.Collection == nil {
			return false
		}

		if t.Collection.Is(sharkType.Collection) {
			return true
		}
		return false
	default:
		return false
	}
}

func (t TSharkArray) Collects() []ISharkType {
	return []ISharkType{t.Collection}
}
