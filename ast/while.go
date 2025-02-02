package ast

import (
	"bytes"
	"shark/token"
)

type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
	Token     token.Token
}

func (w *WhileStatement) statementNode() {}

func (w *WhileStatement) TokenPos() token.Position { return w.Token.Pos }

func (w *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(w.Condition.String())
	out.WriteString(" -> ")
	out.WriteString(w.Body.String())

	return out.String()
}

func (w *WhileStatement) TokenLiteral() string { return w.Token.Literal }
