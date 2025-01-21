package object

import "shark/types"

type Error struct {
	Message string
}

func (e *Error) Inspect() string { return "ERROR: " + e.Message }

func (e *Error) Type() types.ISharkType { return types.TSharkError{} }
