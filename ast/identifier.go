package ast

import (
	"shark/token"
)

type Identifier struct {
	DefaultValue *Expression
	Value        string
	Token        token.Token
	Mutable      bool
	VariadicType bool
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenPos() token.Position { return i.Token.Pos }

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }
