package vm

import (
	"fmt"
	"math"
	"shark/code"
	"shark/compiler"
	"shark/exception"
	"shark/object"
)

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VmConf struct {
	StackSize   int
	GlobalsSize int
	MaxFrames   int
}

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]

	globals []object.Object

	frames      []*Frame
	framesIndex int

	conf *VmConf
}

func NewDefault(bytecode *compiler.Bytecode) *VM {
	conf := NewDefaultConf()

	return New(bytecode, &conf)
}

func NewDefaultConf() VmConf {
	return VmConf{
		StackSize:   2048,
		GlobalsSize: 65536,
		MaxFrames:   1024,
	}
}

func New(bytecode *compiler.Bytecode, conf *VmConf) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, conf.MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, conf.StackSize),
		sp:    0,

		globals: make([]object.Object, conf.GlobalsSize),

		frames:      frames,
		framesIndex: 1,

		conf: conf,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object, conf *VmConf) *VM {
	vm := New(bytecode, conf)
	vm.globals = s
	return vm
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() *exception.SharkError {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			if err := vm.push(vm.constants[constIndex]); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpPower:
			if err := vm.executeBinaryOperation(op); err != nil {
				return err
			}
		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanEqual, code.OpAnd, code.OpOr:
			if err := vm.executeComparison(op); err != nil {
				return err
			}
		case code.OpBang:
			if err := vm.executeBangOperator(); err != nil {
				return err
			}
		case code.OpMinus:
			if err := vm.executeMinusOperator(); err != nil {
				return err
			}
		case code.OpIncrementGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			intVal, ok := vm.globals[globalIndex].(*object.Integer)
			if !ok {
				return &exception.SharkError{
					ErrMsg:  "cannot increment a global non-integer value",
					ErrCode: exception.SharkErrorNonNumberIncrement,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			vm.globals[globalIndex] = &object.Integer{Value: intVal.Value + 1}
		case code.OpIncrementLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			intVal, ok := vm.stack[frame.basePointer+int(localIndex)].(*object.Integer)
			if !ok {
				return &exception.SharkError{
					ErrMsg:  "cannot increment a local non-integer value",
					ErrCode: exception.SharkErrorNonNumberIncrement,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			vm.stack[frame.basePointer+int(localIndex)] = &object.Integer{Value: intVal.Value + 1}
		case code.OpDecrementGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			intVal, ok := vm.globals[globalIndex].(*object.Integer)
			if !ok {
				return &exception.SharkError{
					ErrMsg:  "cannot decrement a global non-integer value",
					ErrCode: exception.SharkErrorNonNumberIncrement,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			vm.globals[globalIndex] = &object.Integer{Value: intVal.Value - 1}
		case code.OpDecrementLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			intVal, ok := vm.stack[frame.basePointer+int(localIndex)].(*object.Integer)
			if !ok {
				return &exception.SharkError{
					ErrMsg:  "cannot decrement a local non-integer value",
					ErrCode: exception.SharkErrorNonNumberIncrement,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			vm.stack[frame.basePointer+int(localIndex)] = &object.Integer{Value: intVal.Value - 1}
		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case code.OpNull:
			if err := vm.push(Null); err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			if err := vm.push(vm.globals[globalIndex]); err != nil {
				return err
			}
		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements
			if err := vm.push(array); err != nil {
				return err
			}
		case code.OpHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp -= numElements
			if err = vm.push(hash); err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()
			if err := vm.executeIndexExpression(left, index); err != nil {
				return err
			}
		case code.OpCall:
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			if err := vm.executeCall(int(numArgs)); err != nil {
				return err
			}
		case code.OpReturnValue:
			if vm.currentFrame().basePointer == 0 {
				return &exception.SharkError{
					ErrMsg:  "return statement outside of function is not allowed",
					ErrCode: exception.SharkErrorTopLeverReturn,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			returnValue := vm.pop()
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1
			if err := vm.push(returnValue); err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1
			if err := vm.push(Null); err != nil {
				return err
			}
		case code.OpSpread:
			operand := vm.pop()
			if operand.Type() != object.STRING_OBJ {
				return &exception.SharkError{
					ErrMsg:  "spread operator only supports string operands",
					ErrCode: exception.SharkErrorMismatchedTypes,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			str := operand.(*object.String).Value
			elements := make([]object.Object, len(str))
			for i, c := range str {
				elements[i] = &object.String{Value: string(c)}
			}
			if err := vm.push(&object.Array{Elements: elements}); err != nil {
				return err
			}
		case code.OpRange:
			end := vm.pop()
			start := vm.pop()
			if start.Type() != object.INTEGER_OBJ || end.Type() != object.INTEGER_OBJ {
				return &exception.SharkError{
					ErrMsg:  "range operator only supports integer operands",
					ErrCode: exception.SharkErrorMismatchedTypes,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			startVal := start.(*object.Integer).Value
			endVal := end.(*object.Integer).Value
			if startVal > endVal {
				elements := make([]object.Object, startVal-endVal+1)
				for i := startVal; i >= endVal; i-- {
					elements[startVal-i] = &object.Integer{Value: i}
				}
				if err := vm.push(&object.Array{Elements: elements}); err != nil {
					return err
				}
			} else {
				elements := make([]object.Object, endVal-startVal+1)
				for i := startVal; i <= endVal; i++ {
					elements[i-startVal] = &object.Integer{Value: i}
				}
				if err := vm.push(&object.Array{Elements: elements}); err != nil {
					return err
				}
			}
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case code.OpSetLocalDefault:
			if vm.sp == 0 {
				return &exception.SharkError{
					ErrMsg:  "no value to set",
					ErrCode: exception.SharkErrorNoDefaultValue,
					ErrType: exception.SharkErrorTypeRuntime,
				}
			}
			// only set the local if the local's value is null
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			if vm.stack[frame.basePointer+int(localIndex)] == nil {
				vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
			}
		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			if err := vm.push(vm.stack[frame.basePointer+int(localIndex)]); err != nil {
				return err
			}
		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			definition := object.Builtins[builtinIndex]
			if err := vm.push(definition.Builtin); err != nil {
				return err
			}
		case code.OpClosure:
			constIndex := code.ReadUint16(ins[ip+1:])
			numFree := code.ReadUint8(ins[ip+3:])
			vm.currentFrame().ip += 3
			if err := vm.pushClosure(int(constIndex), int(numFree)); err != nil {
				return err
			}
		case code.OpGetFree:
			freeIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			currentClosure := vm.currentFrame().cl
			if err := vm.push(currentClosure.Free[freeIndex]); err != nil {
				return err
			}
		case code.OpCurrentClosure:
			currentClosure := vm.currentFrame().cl
			if err := vm.push(currentClosure); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) pushClosure(constIndex, numFree int) *exception.SharkError {
	constant := vm.constants[constIndex]
	function, ok := constant.(*object.CompiledFunction)
	if !ok {
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("not a function: %+v", constant),
			ErrCode: exception.SharkErrorNonFunction,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}

	free := make([]object.Object, numFree)

	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	vm.sp = vm.sp - numFree

	closure := &object.Closure{Fn: function, Free: free}

	return vm.push(closure)
}

func (vm *VM) push(o object.Object) *exception.SharkError {
	if vm.sp >= vm.conf.StackSize {
		return &exception.SharkError{
			ErrMsg:  "stack overflow",
			ErrCode: exception.SharkErrorVMStackOverflow,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) executeBinaryOperation(op code.Opcode) *exception.SharkError {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return &exception.SharkError{
			ErrMsg:     fmt.Sprintf("cannot execute binary operation between %s and %s", leftType, rightType),
			ErrHelpMsg: "Types must be the same",
			ErrCode:    exception.SharkErrorNonNumberIncrement,
			ErrType:    exception.SharkErrorTypeRuntime,
		}
	}
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) *exception.SharkError {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	case code.OpPower:
		result = int64(math.Pow(float64(leftValue), float64(rightValue)))
	default:
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("unknown integer operator: %d", op),
			ErrCode: exception.SharkErrorUnknownOperator,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) *exception.SharkError {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return vm.executeBooleanComparison(op, left, right)
	}

	return &exception.SharkError{
		ErrMsg:  fmt.Sprintf("unknown operator: %d (%s %s)", op, left.Type(), right.Type()),
		ErrCode: exception.SharkErrorUnknownOperator,
		ErrType: exception.SharkErrorTypeRuntime,
	}

}

func (vm *VM) executeBooleanComparison(op code.Opcode, left, right object.Object) *exception.SharkError {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	case code.OpAnd:
		return vm.push(nativeBoolToBooleanObject(leftValue && rightValue))
	case code.OpOr:
		return vm.push(nativeBoolToBooleanObject(leftValue || rightValue))
	default:
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("unknown boolean operator integers: %d", op),
			ErrCode: exception.SharkErrorUnknownBoolOperator,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) *exception.SharkError {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	case code.OpGreaterThanEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue >= rightValue))
	default:
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("unknown boolean operator between integers: %d", op),
			ErrCode: exception.SharkErrorUnknownBoolOperator,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}
}

func (vm *VM) executeBangOperator() *exception.SharkError {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False, Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() *exception.SharkError {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("Type %s is unsupported for negation", operand.Type()),
			ErrCode: exception.SharkErrorMismatchedTypes,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) *exception.SharkError {
	if op != code.OpAdd {
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("unknown string operator: %d", op),
			ErrCode: exception.SharkErrorUnknownStringOperator,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, *exception.SharkError) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, &exception.SharkError{
				ErrMsg:  fmt.Sprintf("cannot hash %s as hash key", key.Type()),
				ErrCode: exception.SharkErrorNonHashable,
				ErrType: exception.SharkErrorTypeRuntime,
			}
		}

		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) executeIndexExpression(left, index object.Object) *exception.SharkError {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("cannot index %s", index.Type()),
			ErrCode: exception.SharkErrorNonIndexable,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}
}

func (vm *VM) executeArrayIndex(array, index object.Object) *exception.SharkError {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	m := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > m {
		return vm.push(Null)
	}

	return vm.push(arrayObject.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) *exception.SharkError {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("cannot hash %s for indexing", index.Type()),
			ErrCode: exception.SharkErrorNonHashable,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) *exception.SharkError {
	if vm.framesIndex >= vm.conf.MaxFrames {
		return &exception.SharkError{
			ErrMsg:  "cannot push more frames to the stack",
			ErrCode: exception.SharkErrorVMStackOverflow,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++

	return nil
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) callClosure(cl *object.Closure, numArgs int) *exception.SharkError {
	if numArgs > cl.Fn.NumParameters || numArgs < cl.Fn.NumParameters-cl.Fn.NumDefaults {
		return &exception.SharkError{
			ErrMsg:  fmt.Sprintf("expected %d arguments for function call, but got %d", cl.Fn.NumParameters, numArgs),
			ErrCode: exception.SharkErrorArgumentNumberMismatch,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}

	frame := NewFrame(cl, vm.sp-numArgs)

	if err := vm.pushFrame(frame); err != nil {
		return err
	}

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {

	case *object.Boolean:
		return obj.Value

	case *object.Null:
		return false

	default:
		return true
	}
}

func (vm *VM) callBuiltin(builtin *object.Builtin, numArgs int) *exception.SharkError {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := builtin.Fn(args...)

	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		if err := vm.push(result); err != nil {
			return err
		}
	} else {
		if err := vm.push(Null); err != nil {
			return err
		}
	}

	return nil
}

func (vm *VM) executeCall(numArgs int) *exception.SharkError {
	callee := vm.stack[vm.sp-1-numArgs]

	switch callee := callee.(type) {
	case *object.Closure:
		return vm.callClosure(callee, numArgs)
	case *object.Builtin:
		return vm.callBuiltin(callee, numArgs)
	default:
		return &exception.SharkError{
			ErrMsg:  "cannot call a non-builtin or closure function",
			ErrCode: exception.SharkErrorNonFunctionCall,
			ErrType: exception.SharkErrorTypeRuntime,
		}
	}
}
