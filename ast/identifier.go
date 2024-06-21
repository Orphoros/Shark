package ast

import "shark/token"

type Identifier struct {
	Token        token.Token
	Value        string
	DefaultValue *Expression
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }
