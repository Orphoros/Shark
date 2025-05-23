package ast

import (
	"shark/token"
)

type IndexExpression struct {
	Left  Expression
	Index Expression
	Token token.Token
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenPos() token.Position { return ie.Token.Pos }

func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IndexExpression) String() string {
	var out string
	out += "("
	out += ie.Left.String()
	out += "["
	out += ie.Index.String()
	out += "])"
	return out
}
