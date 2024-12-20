package ast

import (
	"shark/object"
	"shark/token"
)

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) Type() object.Type { return object.STRING_OBJ }

func (sl *StringLiteral) TokenPos() token.Position { return sl.Token.Pos }

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

func (sl *StringLiteral) String() string { return sl.Token.Literal }
