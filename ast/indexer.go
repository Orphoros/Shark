package ast

import (
	"shark/object"
	"shark/token"
)

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) Type() object.Type { return ie.Left.Type() }

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
