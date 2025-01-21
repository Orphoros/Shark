package object

import (
	"fmt"
	"shark/types"
)

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Inspect() string { return fmt.Sprintf("Closure[%p]", c) }

func (c *Closure) Type() types.ISharkType {
	return types.TSharkClosure{FuncType: c.Fn.Type()}
}
