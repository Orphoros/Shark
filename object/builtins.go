package object

import (
	"fmt"
	"os"
)

type Builtin struct {
	CanCache bool
	Fn       BuiltinFunction
}

type BuiltinFunction func(args ...Object) Object

func (b *Builtin) Inspect() string { return "builtin function" }

func (b *Builtin) Type() Type { return BUILTIN_OBJ }

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{"exit", &Builtin{Fn: Exit, CanCache: true}},
	{"puts", &Builtin{Fn: Puts, CanCache: false}},
	{"len", &Builtin{Fn: Len, CanCache: true}},
	{"first", &Builtin{Fn: First, CanCache: true}},
	{"last", &Builtin{Fn: Last, CanCache: true}},
	{"rest", &Builtin{Fn: Rest, CanCache: true}},
	{"push", &Builtin{Fn: Push, CanCache: true}},
	{"type", &Builtin{Fn: ObjType, CanCache: true}},
}

func ObjType(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	return &String{Value: string(args[0].Type())}
}

func Len(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	case *Array:
		return &Integer{Value: int64(len(arg.Elements))}
	case *Hash:
		return &Integer{Value: int64(len(arg.Pairs))}
	case *Tuple:
		return &Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func First(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		if len(arg.Value) > 0 {
			return &String{Value: string(arg.Value[0])}
		}
		return nil
	case *Array:
		if len(arg.Elements) > 0 {
			return arg.Elements[0]
		}
		return nil
	case *Tuple:
		if len(arg.Elements) > 0 {
			return arg.Elements[0]
		}
		return nil
	default:
		return newError("argument to `first` not supported, got %s", args[0].Type())
	}
}

func Puts(args ...Object) Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}

	fmt.Println()

	return nil
}

func Last(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		length := len(arg.Value)
		if length > 0 {
			return &String{Value: string(arg.Value[length-1])}
		}
		return nil
	case *Array:
		length := len(arg.Elements)
		if length > 0 {
			return arg.Elements[length-1]
		}
		return nil
	case *Tuple:
		length := len(arg.Elements)
		if length > 0 {
			return arg.Elements[length-1]
		}
		return nil
	default:
		return newError("argument to `last` not supported, got %s", args[0].Type())
	}
}

func Rest(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &Array{Elements: newElements}
	}

	return nil
}

func Push(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)

	newElements := make([]Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &Array{Elements: newElements}
}

func Exit(args ...Object) Object {
	if len(args) != 1 {
		os.Exit(0)
	}

	if args[0].Type() != INTEGER_OBJ {
		return newError("argument to `exit` must be INTEGER, got %s", args[0].Type())
	}

	integer := args[0].(*Integer)
	os.Exit(int(integer.Value))
	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
