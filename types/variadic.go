package types

type TSharkVariadic struct {
	ISharkType
	Enclosed ISharkType
}

func (TSharkVariadic) SharkTypeString() string { return "T" }

func (t TSharkVariadic) Is(sharkType ISharkType) bool {
	if t.Enclosed == nil {
		return true
	}
	return t.Enclosed.Is(sharkType)

}
