package ast

import (
	"bytes"
	"shark/token"
)

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Token       token.Token
}

func (ie IfExpression) expressionNode() {}

func (ie IfExpression) TokenPos() token.Position { return ie.Token.Pos }

func (ie IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
