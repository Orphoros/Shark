package ast

import (
	"bytes"
	"shark/object"
	"shark/token"
	"strings"
)

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) Type() object.Type { return object.RETURN_VALUE_OBJ }

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
