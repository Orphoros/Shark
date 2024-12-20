package ast

import (
	"bytes"
	"shark/object"
	"shark/token"
	"strings"
)

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl HashLiteral) Type() object.Type { return object.HASH_OBJ }

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
