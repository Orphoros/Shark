package parser

import (
	"shark/ast"
	"shark/token"
)

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	tkn := p.curToken
	leftExpr := left
	p.nextToken()
	index := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		value := p.parseExpression(LOWEST)
		return &ast.IndexAssignExpression{Token: tkn, Left: leftExpr, Index: index, Value: value}
	}

	return &ast.IndexExpression{Token: tkn, Left: leftExpr, Index: index}
}
