package ast

import (
	"bytes"
	"fmt"
	"shark/token"
	"shark/types"
	"strings"
)

type FunctionLiteral struct {
	Body        *BlockStatement
	Name        string
	Parameters  []*Identifier
	Token       token.Token
	DefinedType types.ISharkType
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenPos() token.Position { return fl.Token.Pos }

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		out.WriteString(fmt.Sprintf("<%s>", fl.Name))
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
