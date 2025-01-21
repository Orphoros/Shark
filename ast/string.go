package ast

import (
	"shark/token"
)

type StringLiteral struct {
	Value string
	Token token.Token
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenPos() token.Position { return sl.Token.Pos }

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

func (sl *StringLiteral) String() string { return sl.Token.Literal }
