package object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ           = "INTEGER"
	BOOLEAN_OBJ           = "BOOLEAN"
	NULL_OBJ              = "NULL"
	RETURN_VALUE_OBJ      = "RETURN_VALUE"
	ERROR_OBJ             = "ERROR"
	FUNCTION_OBJ          = "FUNCTION"
	STRING_OBJ            = "STRING"
	BUILTIN_OBJ           = "BUILTIN"
	ARRAY_OBJ             = "ARRAY"
	HASH_OBJ              = "HASH"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
	CLOSURE_OBJ           = "CLOSURE"
)

type Null struct{}

type ReturnValue struct {
	Value Object
}

type Builtin struct {
	Fn BuiltinFunction
}

type BuiltinFunction func(args ...Object) Object

func (n *Null) Inspect() string { return "null" }

func (n *Null) Type() ObjectType { return NULL_OBJ }

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

func (b *Builtin) Inspect() string { return "builtin function" }

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
