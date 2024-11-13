package ast

import (
	"bytes"
	"shark/token"
)

type PrefixExpression struct {
	Token      token.Token
	Operator   string
	Right      Expression
	RightIdent *Identifier
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
