package ast

import (
	"shark/token"
	"strings"
)

type ArrayLiteral struct {
	Elements []Expression
	Token    token.Token
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenPos() token.Position { return al.Token.Pos }

func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

func (al *ArrayLiteral) String() string {
	var out string

	var elements []string

	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}

	out += "["
	out += strings.Join(elements, ", ")
	out += "]"

	return out
}
