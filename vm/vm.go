package vm

import (
	"fmt"
	"shark/bytecode"
	"shark/code"
	"shark/config"
	"shark/exception"
	"shark/object"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants   []object.Object
	stack       []object.Object
	sp          int // Always points to the next value. Top of stack is stack[sp-1]
	globals     []object.Object
	frames      []*Frame
	framesIndex int
	conf        *config.VmConf
	cache       *expirable.LRU[string, object.Object]
}

func NewDefault(bytecode *bytecode.Bytecode) *VM {
	conf := config.NewDefaultVmConf()

	return New(bytecode, &conf)
}

func New(bytecode *bytecode.Bytecode, conf *config.VmConf) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, conf.MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants:   bytecode.Constants,
		stack:       make([]object.Object, conf.StackSize),
		sp:          0,
		globals:     make([]object.Object, conf.GlobalsSize),
		frames:      frames,
		framesIndex: 1,
		conf:        conf,
		cache:       expirable.NewLRU[string, object.Object](conf.CacheSize, nil, time.Minute*5),
	}
}

func NewWithGlobalsStore(bytecode *bytecode.Bytecode, s []object.Object, conf *config.VmConf) *VM {
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
			value := vm.constants[constIndex]
			if err := vm.push(value); err != nil {
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
				return newSharkError(exception.SharkErrorNonNumberIncrement, vm.globals[globalIndex].Type())
			}
			vm.globals[globalIndex] = &object.Integer{Value: intVal.Value + 1}
		case code.OpIncrementLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			intVal, ok := vm.stack[frame.basePointer+int(localIndex)].(*object.Integer)
			if !ok {
				return newSharkError(exception.SharkErrorNonNumberIncrement, vm.stack[frame.basePointer+int(localIndex)].Type())
			}
			vm.stack[frame.basePointer+int(localIndex)] = &object.Integer{Value: intVal.Value + 1}
		case code.OpDecrementGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			intVal, ok := vm.globals[globalIndex].(*object.Integer)
			if !ok {
				return newSharkError(exception.SharkErrorNonNumberDecrement, vm.globals[globalIndex].Type())
			}
			vm.globals[globalIndex] = &object.Integer{Value: intVal.Value - 1}
		case code.OpDecrementLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			intVal, ok := vm.stack[frame.basePointer+int(localIndex)].(*object.Integer)
			if !ok {
				return newSharkError(exception.SharkErrorNonNumberIncrement, vm.stack[frame.basePointer+int(localIndex)].Type())
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
			// clear the stack between sp and sp-numElements with nil
			for i := vm.sp; i < vm.sp+numElements; i++ {
				vm.stack[i] = nil
			}
			if err := vm.push(array); err != nil {
				return err
			}
		case code.OpTuple:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			tpl := vm.buildTuple(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements
			// clear the stack between sp and sp-numElements with nil
			for i := vm.sp; i < vm.sp+numElements; i++ {
				vm.stack[i] = nil
			}
			if err := vm.push(tpl); err != nil {
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
		case code.OpIndexAssign:
			index := vm.pop()
			left := vm.pop()
			value := vm.pop()
			if err := vm.executeIndexAssign(left, index, value); err != nil {
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
				return newSharkError(exception.SharkErrorTopLeverReturn)
			}
			returnValue := vm.pop()
			frame := vm.popFrame()
			// clear the stack between sp and basePointer with nil
			for i := vm.sp; i < frame.basePointer; i++ {
				vm.stack[i] = nil
			}
			vm.sp = frame.basePointer - 1
			if err := vm.push(returnValue); err != nil {
				return err
			}
			// cache the result if applicable
			if frame.canCache {
				vm.cache.Add(frame.cacheKey, returnValue)
			}
		case code.OpReturn:
			frame := vm.popFrame()
			// clear the stack between sp and basePointer with nil
			for i := vm.sp; i < frame.basePointer; i++ {
				vm.stack[i] = nil
			}
			vm.sp = frame.basePointer - 1
			if err := vm.push(Null); err != nil {
				return err
			}
		case code.OpTupleDeconstruct:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			tpl := vm.pop()
			tuple, ok := tpl.(*object.Tuple)
			if !ok {
				return newSharkError(exception.SharkErrorMismatchedTypes, tpl.Type(), object.TUPLE_OBJ)
			}
			if len(tuple.Elements) != numElements {
				return newSharkError(exception.SharkErrorTupleDeconstructMismatch, len(tuple.Elements), numElements)
			}
			for i := numElements - 1; i >= 0; i-- {
				if err := vm.push(tuple.Elements[i]); err != nil {
					return err
				}
			}
		case code.OpSpread:
			operand := vm.pop()
			strObj, ok := operand.(*object.String)
			if !ok {
				return newSharkError(exception.SharkErrorMismatchedTypes, operand.Type(), object.STRING_OBJ)
			}
			str := strObj.Value
			elements := make([]object.Object, 0, len(str))
			for _, c := range str {
				elements = append(elements, &object.String{Value: string(c)})
			}
			if err := vm.push(&object.Array{Elements: elements}); err != nil {
				return err
			}
		case code.OpRange:
			endObj := vm.pop()
			startObj := vm.pop()

			endInt, endOk := endObj.(*object.Integer)
			startInt, startOk := startObj.(*object.Integer)

			if !startOk || !endOk {
				return newSharkError(exception.SharkErrorMismatchedTypes, "Integer", "Integer")
			}

			startVal := startInt.Value
			endVal := endInt.Value

			var elements []object.Object
			if startVal > endVal {
				elements = make([]object.Object, startVal-endVal+1)
				for i, v := startVal, 0; i >= endVal; i, v = i-1, v+1 {
					elements[v] = &object.Integer{Value: i}
				}
			} else {
				elements = make([]object.Object, endVal-startVal+1)
				for i, v := startVal, 0; i <= endVal; i, v = i+1, v+1 {
					elements[v] = &object.Integer{Value: i}
				}
			}

			if err := vm.push(&object.Array{Elements: elements}); err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case code.OpSetLocalDefault:
			if vm.sp == 0 {
				return newSharkError(exception.SharkErrorNoDefaultValue)
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
		return newSharkError(exception.SharkErrorNonFunction, constant.Type())
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
		return newSharkError(exception.SharkErrorVMStackOverflow)
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

	switch leftValue := left.(type) {
	case *object.Integer:
		switch rightValue := right.(type) {
		case *object.Integer:
			var result int64
			switch op {
			case code.OpAdd:
				result = leftValue.Value + rightValue.Value
			case code.OpSub:
				result = leftValue.Value - rightValue.Value
			case code.OpMul:
				result = leftValue.Value * rightValue.Value
			case code.OpDiv:
				if rightValue.Value == 0 {
					return newSharkError(exception.SharkErrorDivisionByZero)
				}
				result = leftValue.Value / rightValue.Value
			case code.OpPower:
				result = intPow(leftValue.Value, rightValue.Value)
			default:
				return newSharkError(exception.SharkErrorUnknownOperator, op)
			}
			vm.push(&object.Integer{Value: result})
			return nil
		default:
			return newSharkError(exception.SharkErrorMismatchedTypes, "Integer", right.Type())
		}
	case *object.String:
		switch rightValue := right.(type) {
		case *object.String:
			return vm.executeBinaryStringOperation(op, leftValue, rightValue)
		default:
			return newSharkError(exception.SharkErrorMismatchedTypes, "String", right.Type())
		}
	default:
		return newSharkError(exception.SharkErrorUnknownType, left.Type())
	}
}

func intPow(a, b int64) int64 {
	result := int64(1)
	for b > 0 {
		if b&1 == 1 {
			result *= a
		}
		a *= a
		b >>= 1
	}
	return result
}

func (vm *VM) executeComparison(op code.Opcode) *exception.SharkError {
	right := vm.pop()
	left := vm.pop()

	switch leftVal := left.(type) {
	case *object.Integer:
		if rightVal, ok := right.(*object.Integer); ok {
			return vm.executeIntegerComparison(op, leftVal, rightVal)
		}
	case *object.Boolean:
		if rightVal, ok := right.(*object.Boolean); ok {
			return vm.executeBooleanComparison(op, leftVal, rightVal)
		}
	case *object.String:
		if rightVal, ok := right.(*object.String); ok {
			return vm.executeStringComparison(op, leftVal, rightVal)
		}
	}

	return newSharkError(exception.SharkErrorMismatchedTypes, left.Type(), right.Type())
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
		return newSharkError(exception.SharkErrorUnknownBoolOperator, op)
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
		return newSharkError(exception.SharkErrorUnknownOperator, op)
	}
}

func (vm *VM) executeStringComparison(op code.Opcode, left, right object.Object) *exception.SharkError {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	default:
		return newSharkError(exception.SharkErrorUnknownStringOperator, op)
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
		return newSharkError(exception.SharkErrorMismatchedTypes, operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) *exception.SharkError {
	if op != code.OpAdd {
		return newSharkError(exception.SharkErrorUnknownStringOperator, op)
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

func (vm *VM) buildTuple(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Tuple{Elements: elements}
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
	switch left := left.(type) {
	case *object.Array:
		if index, ok := index.(*object.Integer); ok {
			return vm.executeArrayIndex(left, index)
		}
		return newSharkError(exception.SharkErrorNonIndexable, left.Type())
	case *object.String:
		if index, ok := index.(*object.Integer); ok {
			return vm.executeStringIndex(left, index)
		}
		return newSharkError(exception.SharkErrorNonIndexable, left.Type())
	case *object.Tuple:
		if index, ok := index.(*object.Integer); ok {
			return vm.executeTupleIndex(left, index)
		}
		return newSharkError(exception.SharkErrorNonIndexable, left.Type())
	case *object.Hash:
		return vm.executeHashIndex(left, index)
	default:
		return newSharkError(exception.SharkErrorNonIndexable, left.Type())
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

func (vm *VM) executeStringIndex(str, index object.Object) *exception.SharkError {
	strObject := str.(*object.String)
	i := index.(*object.Integer).Value
	m := int64(len(strObject.Value) - 1)

	if i < 0 || i > m {
		return vm.push(Null)
	}

	return vm.push(&object.String{Value: string(strObject.Value[i])})
}

func (vm *VM) executeTupleIndex(tpl, index object.Object) *exception.SharkError {
	tplObject := tpl.(*object.Tuple)
	i := index.(*object.Integer).Value

	if i < 0 || i >= int64(len(tplObject.Elements)) {
		return vm.push(Null)
	}

	return vm.push(tplObject.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) *exception.SharkError {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newSharkError(exception.SharkErrorNonHashable, index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) executeIndexAssign(left, index, value object.Object) *exception.SharkError {
	switch left := left.(type) {
	case *object.Array:
		switch index := index.(type) {
		case *object.Integer:
			return vm.executeArrayIndexAssign(left, index, value)
		default:
			return newSharkError(exception.SharkErrorNonIndexable, index.Type())
		}
	case *object.Hash:
		return vm.executeHashIndexAssign(left, index, value)
	default:
		return newSharkError(exception.SharkErrorNonIndexable, left.Type())
	}
}

func (vm *VM) executeArrayIndexAssign(array, index, value object.Object) *exception.SharkError {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	m := int64(len(arrayObject.Elements) - 1)
	if i < 0 || i > m {
		return newSharkError(exception.SharkErrorIndexOutOfBounds, i)
	}
	arrayObject.Elements[i] = value
	return vm.push(arrayObject)
}

func (vm *VM) executeHashIndexAssign(hash, index, value object.Object) *exception.SharkError {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newSharkError(exception.SharkErrorNonHashable, index.Type())
	}
	hashObject.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: value}
	return vm.push(hashObject)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func newSharkError(code exception.SharkErrorCode, param ...interface{}) *exception.SharkError {
	return exception.NewSharkError(exception.SharkErrorTypeRuntime, code, param...)
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

	args := vm.stack[vm.sp-numArgs : vm.sp]

	key, canCache := vm.createCacheKey(callee, args)

	if canCache {
		if result, ok := vm.cache.Get(key); ok {
			// Pop the callee and arguments off the stack
			vm.sp = vm.sp - numArgs - 1
			// Push the cached result onto the stack
			if err := vm.push(result); err != nil {
				return err
			}
			return nil
		}
	}

	switch callee := callee.(type) {
	case *object.Closure:
		cl := callee

		requiredArgs := cl.Fn.NumParameters
		defaultArgs := cl.Fn.NumDefaults
		minArgs := requiredArgs - defaultArgs

		if numArgs < minArgs || numArgs > requiredArgs {
			return newSharkError(exception.SharkErrorArgumentNumberMismatch, requiredArgs, numArgs)
		}

		if vm.framesIndex >= vm.conf.MaxFrames {
			return newSharkError(exception.SharkErrorVMStackOverflow)
		}

		frame := vm.frames[vm.framesIndex]
		if frame == nil {
			frame = &Frame{}
			vm.frames[vm.framesIndex] = frame
		}

		frame.cl = cl
		frame.ip = -1
		frame.basePointer = vm.sp - numArgs
		frame.cacheKey = key
		frame.canCache = canCache

		vm.framesIndex++
		vm.currentFrame().ip = frame.ip

		vm.sp = frame.basePointer + cl.Fn.NumLocals

		return nil

	case *object.Builtin:
		if err := vm.callBuiltin(callee, numArgs); err != nil {
			return err
		}
		result := vm.stack[vm.sp-1]
		if canCache {
			vm.cache.Add(key, result)
		}
		return nil

	default:
		return newSharkError(exception.SharkErrorNonFunctionCall, callee.Type())
	}
}

func (vm *VM) createCacheKey(callee object.Object, args []object.Object) (string, bool) {
	var keyBuilder strings.Builder

	switch callee := callee.(type) {
	case *object.Closure:
		keyBuilder.WriteString(fmt.Sprintf("%p", callee.Fn))
	case *object.Builtin:
		if !callee.CanCache {
			return "", false
		}
		keyBuilder.WriteString(fmt.Sprintf("%p", callee))
	default:
		keyBuilder.WriteString(fmt.Sprintf("%p", callee))
	}

	keyBuilder.WriteString("(")
	for i, arg := range args {
		if i > 0 {
			keyBuilder.WriteString(",")
		}
		if hashableArg, ok := arg.(object.Hashable); ok {
			keyBuilder.WriteString(fmt.Sprintf("%s:%v", arg.Type(), hashableArg.HashKey()))
		} else {
			return "", false
		}
	}
	keyBuilder.WriteString(")")

	return keyBuilder.String(), true
}
