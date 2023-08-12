package object

import (
	"fmt"
	"os"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{"exit", &Builtin{Fn: Exit}},
	{"puts", &Builtin{Fn: Puts}},
	{"len", &Builtin{Fn: Len}},
	{"first", &Builtin{Fn: First}},
	{"last", &Builtin{Fn: Last}},
	{"rest", &Builtin{Fn: Rest}},
	{"push", &Builtin{Fn: Push}},
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
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func First(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return nil
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

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return nil
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

func GetBuiltinByName(name string) *Builtin {
	for _, b := range Builtins {
		if b.Name == name {
			return b.Builtin
		}
	}

	return nil
}
