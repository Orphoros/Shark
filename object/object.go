package object

type Type string

type Object interface {
	Type() Type
	Inspect() string
}

const (
	INTEGER_OBJ           Type = "INTEGER"
	BOOLEAN_OBJ           Type = "BOOLEAN"
	NULL_OBJ              Type = "NULL"
	RETURN_VALUE_OBJ      Type = "RETURN_VALUE"
	ERROR_OBJ             Type = "ERROR"
	FUNCTION_OBJ          Type = "FUNCTION"
	STRING_OBJ            Type = "STRING"
	BUILTIN_OBJ           Type = "BUILTIN"
	ARRAY_OBJ             Type = "ARRAY"
	HASH_OBJ              Type = "HASH"
	COMPILED_FUNCTION_OBJ Type = "COMPILED_FUNCTION"
	CLOSURE_OBJ           Type = "CLOSURE"
	TUPLE_OBJ             Type = "TUPLE"
)

func (t Type) String() string {
	if t == "" {
		return "UNDEFINED"
	}

	return string(t)
}

type Null struct{}

type ReturnValue struct {
	Value Object
}

type Builtin struct {
	CanCache bool
	Fn       BuiltinFunction
}

type BuiltinFunction func(args ...Object) Object

func (n *Null) Inspect() string { return "null" }

func (n *Null) Type() Type { return NULL_OBJ }

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

func (rv *ReturnValue) Type() Type { return RETURN_VALUE_OBJ }

func (b *Builtin) Inspect() string { return "builtin function" }

func (b *Builtin) Type() Type { return BUILTIN_OBJ }
