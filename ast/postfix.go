package ast

import (
	"bytes"
	"shark/token"
)

type PostfixExpression struct {
	Left     Expression
	Token    token.Token
	Operator string
}

func (pe *PostfixExpression) expressionNode() {}

func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}
