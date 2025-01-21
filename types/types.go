package types

type ISharkType interface {
	SharkTypeString() string
	Is(sharkType ISharkType) bool
}
