package ast

import (
	"shark/object"
	"shark/token"
)

type Identifier struct {
	Token        token.Token
	Value        string
	DefaultValue *Expression
	Mutable      bool
	VariadicType bool
	ObjType      object.Type
}

func (i *Identifier) Type() object.Type { return i.ObjType }

func (i *Identifier) TokenPos() token.Position { return i.Token.Pos }

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }
