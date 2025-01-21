package ast

import (
	"bytes"
	"shark/token"
	"strings"
)

type TupleLiteral struct {
	Elements []Expression
	Token    token.Token
}

func (tl *TupleLiteral) expressionNode() {}

func (tl *TupleLiteral) TokenPos() token.Position { return tl.Token.Pos }

func (tl *TupleLiteral) TokenLiteral() string { return tl.Token.Literal }

func (tl *TupleLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range tl.Elements {
		elements = append(elements, e.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(")")

	return out.String()
}

type TupleDeconstruction struct {
	Value Expression
	Names []*Identifier
	Token token.Token
}

func (td *TupleDeconstruction) statementNode() {}

func (td *TupleDeconstruction) TokenPos() token.Position { return td.Token.Pos }

func (td *TupleDeconstruction) TokenLiteral() string { return td.Token.Literal }

func (td *TupleDeconstruction) String() string {
	var out bytes.Buffer

	var names []string
	for _, n := range td.Names {
		names = append(names, n.String())
	}

	out.WriteString("let (")
	out.WriteString(strings.Join(names, ", "))
	out.WriteString(") = ")
	out.WriteString(td.Value.String())
	out.WriteString(";")

	return out.String()
}
