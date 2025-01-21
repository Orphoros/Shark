package ast

import (
	"bytes"
	"shark/token"
)

type InfixExpression struct {
	Left     Expression
	Right    Expression
	Operator string
	Token    token.Token
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenPos() token.Position { return ie.Token.Pos }

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
