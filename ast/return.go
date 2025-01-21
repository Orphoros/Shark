package ast

import (
	"bytes"
	"shark/token"
)

type ReturnStatement struct {
	ReturnValue Expression
	Token       token.Token
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenPos() token.Position { return rs.Token.Pos }

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
