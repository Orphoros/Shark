package ast

import (
	"bytes"
	"shark/token"
)

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) statementNode() {}

func (w *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(w.Condition.String())
	out.WriteString(" -> ")
	out.WriteString(w.Body.String())

	return out.String()
}

func (w *WhileStatement) TokenLiteral() string { return w.Token.Literal }
