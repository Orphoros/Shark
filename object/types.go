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
