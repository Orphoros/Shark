package ast

import (
	"bytes"
	"shark/token"
)

type LetStatement struct {
	Value Expression
	Name  *Identifier
	Token token.Token
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenPos() token.Position { return ls.Token.Pos }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()

}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
