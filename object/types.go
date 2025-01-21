package object

import "shark/types"

type Object interface {
	Inspect() string
	Type() types.ISharkType
}
