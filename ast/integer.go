package ast

import (
	"shark/object"
	"shark/token"
)

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) Type() object.Type { return object.INTEGER_OBJ }

func (il *IntegerLiteral) TokenPos() token.Position { return il.Token.Pos }

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string { return il.Token.Literal }
