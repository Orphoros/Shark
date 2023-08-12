package parser

import "shark/ast"

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	if ident, ok := expression.Right.(*ast.Identifier); ok {
		expression.RightIdent = ident
	}

	return expression
}
