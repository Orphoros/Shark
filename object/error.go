package object

type Error struct {
	Message string
}

func (e *Error) Inspect() string { return "ERROR: " + e.Message }

func (e *Error) Type() Type { return ERROR_OBJ }
