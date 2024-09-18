package parser

import (
	"shark/ast"
	"shark/token"
)

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	if (p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COMMA)) ||
		p.curTokenIs(token.MUTABLE) ||
		(p.curTokenIs(token.RPAREN) && p.peekTokenIs(token.ARROW)) ||
		(p.curTokenIs(token.IDENT) && p.peekTokenIs(token.ASSIGN)) ||
		(p.curTokenIs(token.IDENT) && p.peekTokenIs(token.RPAREN)) {
		return p.parseFunctionLiteral()
	}

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
