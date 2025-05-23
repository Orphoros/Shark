package ast

import (
	"bytes"
	"shark/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	TokenPos() token.Position
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionStatement struct {
	Expression Expression
	Token      token.Token
}

func (es *ExpressionStatement) statementNode()           {}
func (es *ExpressionStatement) expressionNode()          {}
func (es *ExpressionStatement) TokenLiteral() string     { return es.Token.Literal }
func (es *ExpressionStatement) TokenPos() token.Position { return es.Token.Pos }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenPos() token.Position {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenPos()
	} else {
		return token.Position{}
	}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
