package vm

import (
	"shark/code"
	"shark/object"
)

type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
	cacheKey    string
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
