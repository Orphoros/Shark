package object

import "shark/types"

type Null struct{}

func (n *Null) Inspect() string { return "null" }

func (n *Null) Type() types.ISharkType { return types.TSharkNull{} }
