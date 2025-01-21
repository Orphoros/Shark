package ast

import (
	"bytes"
	"shark/token"
	"strings"
)

type HashLiteral struct {
	Pairs map[Expression]Expression
	Token token.Token
}

func (hl HashLiteral) expressionNode() {}

func (hl HashLiteral) TokenPos() token.Position { return hl.Token.Pos }

func (hl HashLiteral) TokenLiteral() string { return hl.Token.Literal }

func (hl HashLiteral) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
