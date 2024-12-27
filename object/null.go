package object

type Null struct{}

func (n *Null) Inspect() string { return "null" }

func (n *Null) Type() Type { return NULL_OBJ }
