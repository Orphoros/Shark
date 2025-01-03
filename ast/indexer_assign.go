package ast

import (
	"shark/object"
	"shark/token"
)

type IndexAssignExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
	Value Expression
}

func (ia *IndexAssignExpression) Type() object.Type { return ia.Left.Type() }

func (ia *IndexAssignExpression) TokenPos() token.Position { return ia.Token.Pos }

func (ia *IndexAssignExpression) TokenLiteral() string { return ia.Token.Literal }

func (ia *IndexAssignExpression) String() string {
	var out string
	out += "("
	out += ia.Left.String()
	out += "["
	out += ia.Index.String()
	out += "] = "
	out += ia.Value.String()
	out += ")"
	return out
}
