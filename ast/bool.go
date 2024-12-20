package ast

import (
	"shark/object"
	"shark/token"
)

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) Type() object.Type { return object.BOOLEAN_OBJ }

func (b *Boolean) TokenPos() token.Position { return b.Token.Pos }

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) String() string { return b.Token.Literal }
