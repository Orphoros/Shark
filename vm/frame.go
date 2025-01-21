package vm

import (
	"shark/code"
	"shark/object"
)

type Frame struct {
	cl          *object.Closure
	cacheKey    string
	ip          int
	basePointer int
	canCache    bool
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}

	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
