package ast

import (
	"bytes"
	"shark/token"
	"strings"
)

type CallExpression struct {
	Function  Expression
	Arguments []Expression
	Token     token.Token
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenPos() token.Position { return ce.Token.Pos }

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
