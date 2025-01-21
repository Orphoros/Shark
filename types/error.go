package types

type TSharkError struct {
	ISharkType
}

func (TSharkError) SharkTypeString() string { return "error" }

func (TSharkError) Is(sharkType ISharkType) bool {
	switch sharkType.(type) {
	case TSharkError:
		return true
	default:
		return false
	}
}
