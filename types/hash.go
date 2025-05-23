package types

// Collection Map TShark
type TSharkHashMap struct {
	ISharkType
	Indexes  ISharkType
	Collects ISharkType
}

func (t TSharkHashMap) SharkTypeString() string {
	if t.Indexes == nil || t.Collects == nil {
		return "hashmap<>"
	}
	return "hashmap<" + t.Indexes.SharkTypeString() + "," + t.Collects.SharkTypeString() + ">"
}

func (t TSharkHashMap) Is(sharkType ISharkType) bool {
	switch sharkType := sharkType.(type) {
	case TSharkHashMap:
		if sharkType.Indexes == nil && sharkType.Collects == nil {
			return true
		}
		if t.Indexes == nil || t.Collects == nil {
			return false
		}
		if t.Indexes.Is(sharkType.Indexes) && t.Collects.Is(sharkType.Collects) {
			return true
		}
		return false
	default:
		return false
	}
}
