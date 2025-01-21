package types

type TSharkAny struct {
	ISharkType
}

func (TSharkAny) SharkTypeString() string { return "any" }

func (TSharkAny) Is(sharkType ISharkType) bool { return true }
