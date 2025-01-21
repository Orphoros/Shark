package serializer

import (
	"encoding/gob"
	"fmt"
	"shark/bytecode"
	"shark/code"
	"shark/compiler"
	"shark/object"
)

const (
	TNull byte = iota
	TInteger
	TBoolean
	TString
	TArray
	TTuple
	THash
	TReturnValue
	TError
	TBuiltin
	TCompiledFunction
	TClosure
	TBytecode
	TInstructions
	TCompilationScope
	TEmittedInstruction
	TObject
)

func RegisterTypes() {
	gob.RegisterName(fmt.Sprintf("%c", TNull), &object.Null{})
	gob.RegisterName(fmt.Sprintf("%c", TInteger), &object.Int64{})
	gob.RegisterName(fmt.Sprintf("%c", TBoolean), &object.Boolean{})
	gob.RegisterName(fmt.Sprintf("%c", TString), &object.String{})
	gob.RegisterName(fmt.Sprintf("%c", TArray), &object.Array{})
	gob.RegisterName(fmt.Sprintf("%c", TTuple), &object.Tuple{})
	gob.RegisterName(fmt.Sprintf("%c", THash), &object.Hash{})
	gob.RegisterName(fmt.Sprintf("%c", TError), &object.Error{})
	gob.RegisterName(fmt.Sprintf("%c", TBuiltin), &object.Builtin{})
	gob.RegisterName(fmt.Sprintf("%c", TCompiledFunction), &object.CompiledFunction{})
	gob.RegisterName(fmt.Sprintf("%c", TClosure), &object.Closure{})
	gob.RegisterName(fmt.Sprintf("%c", TBytecode), &bytecode.Bytecode{})
	gob.RegisterName(fmt.Sprintf("%c", TInstructions), code.Instructions{})
	gob.RegisterName(fmt.Sprintf("%c", TCompilationScope), &compiler.CompilationScope{})
	gob.RegisterName(fmt.Sprintf("%c", TEmittedInstruction), &compiler.EmittedInstruction{})
	gob.RegisterName(fmt.Sprintf("%c", TObject), []object.Object{})
}
