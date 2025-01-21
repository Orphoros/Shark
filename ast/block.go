package ast

import (
	"bytes"
	"shark/token"
)

type BlockStatement struct {
	Statements []Statement
	Token      token.Token
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenPos() token.Position { return bs.Token.Pos }

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
