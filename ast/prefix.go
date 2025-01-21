package ast

import (
	"bytes"
	"shark/token"
)

type PrefixExpression struct {
	Right      Expression
	RightIdent *Identifier
	Operator   string
	Token      token.Token
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenPos() token.Position { return pe.Token.Pos }

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
