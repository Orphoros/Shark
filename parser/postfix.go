package parser

import "shark/ast"

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.PostfixExpression{
		Left:     left,
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	return expression
}
