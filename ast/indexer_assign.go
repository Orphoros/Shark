package ast

import (
	"shark/token"
)

type IndexAssignExpression struct {
	Left  Expression
	Index Expression
	Value Expression
	Token token.Token
}

func (ia *IndexAssignExpression) expressionNode() {}

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
