package parser

import (
	"shark/ast"
	"shark/token"
)

func (p *Parser) parseGroupedExpression() ast.Expression {
	curToken := p.curToken
	p.nextToken()

	if (p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COMMA)) ||
		p.curTokenIs(token.MUTABLE) ||
		(p.curTokenIs(token.RPAREN) && p.peekTokenIs(token.ARROW)) ||
		(p.curTokenIs(token.IDENT) && p.peekTokenIs(token.ASSIGN)) ||
		(p.curTokenIs(token.VAR) && p.peekTokenIs(token.IDENT)) ||
		(p.curTokenIs(token.IDENT) && p.peekTokenIs(token.RPAREN)) {
		return p.parseFunctionLiteral()
	}

	if p.curTokenIs(token.VAR) {
		return nil
	}

	exp := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.COMMA) {
		return p.parseTupleLiteral(curToken, exp)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
