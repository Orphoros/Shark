package object

import (
	"fmt"
	"os"
	"shark/types"
)

type Builtin struct {
	CanCache bool
	Fn       BuiltinFunction
	FuncType types.ISharkType
}

type BuiltinFunction func(args ...Object) Object

func (b *Builtin) Inspect() string { return "builtin function" }

func (b *Builtin) Type() types.ISharkType { return b.FuncType }

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{"exit",
		&Builtin{
			Fn:       Exit,
			CanCache: true,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkI64{}}, ReturnT: types.TSharkNull{}},
		},
	},
	{"puts",
		&Builtin{
			Fn:       Puts,
			CanCache: false,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkSpread{Type: types.TSharkAny{}}}, ReturnT: types.TSharkNull{}},
		},
	},
	{"len",
		&Builtin{
			Fn:       Len,
			CanCache: true,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkCollection{Collection: []types.ISharkType{types.TSharkSpread{Type: types.TSharkAny{}}}}}, ReturnT: types.TSharkI64{}},
		},
	},
	{"first",
		&Builtin{
			Fn:       First,
			CanCache: true,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkCollection{Collection: []types.ISharkType{types.TSharkSpread{Type: types.TSharkVariadic{}}}}}, ReturnT: types.TSharkVariadic{}},
		},
	},
	{"last",
		&Builtin{
			Fn:       Last,
			CanCache: true,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkCollection{Collection: []types.ISharkType{types.TSharkSpread{Type: types.TSharkVariadic{}}}}}, ReturnT: types.TSharkVariadic{}},
		},
	},
	{"rest",
		&Builtin{
			Fn:       Rest,
			CanCache: true,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkCollection{Collection: []types.ISharkType{types.TSharkSpread{Type: types.TSharkVariadic{}}}}}, ReturnT: types.TSharkArray{Collection: types.TSharkVariadic{}}},
		},
	},
	{"push",
		&Builtin{
			Fn:       Push,
			CanCache: true,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkArray{Collection: types.TSharkVariadic{}}, types.TSharkVariadic{}}, ReturnT: types.TSharkArray{Collection: types.TSharkVariadic{}}},
		},
	},
	{"type",
		&Builtin{
			Fn:       ObjType,
			CanCache: true,
			FuncType: types.TSharkFuncType{ArgsList: []types.ISharkType{types.TSharkAny{}}, ReturnT: types.TSharkString{}},
		},
	},
}

func ObjType(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	// FIXME: when ident is defined as :bool?, type() returns only bool
	return &String{Value: args[0].Type().SharkTypeString()}
}

func Len(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		return &Int64{Value: int64(len(arg.Value))}
	case *Array:
		return &Int64{Value: int64(len(arg.Elements))}
	// FIXME: Make new type collection to support hash maps for length
	// case *Hash:
	// 	return &Int64{Value: int64(len(arg.Pairs))}
	case *Tuple:
		return &Int64{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type().SharkTypeString())
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
		return newError("argument to `first` not supported, got %s", args[0].Type().SharkTypeString())
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
		return newError("argument to `last` not supported, got %s", args[0].Type().SharkTypeString())
	}
}

func Rest(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
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

	acceptedType := types.TSharkI64{}

	if !args[0].Type().Is(acceptedType) {
		return newError("argument to exit() must be %s, got %s", acceptedType.SharkTypeString(), args[0].Type().SharkTypeString())
	}

	integer := args[0].(*Int64)
	os.Exit(int(integer.Value))
	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
