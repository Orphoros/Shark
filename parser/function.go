package parser

import (
	"shark/ast"
	"shark/token"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.ARROW) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.curTokenIs(token.RPAREN) {
		return identifiers
	}

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		exp := p.parseExpression(LOWEST)
		ident.DefaultValue = &exp
	}

	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.peekTokenIs(token.ASSIGN) {
			p.nextToken()
			p.nextToken()
			exp := p.parseExpression(LOWEST)
			ident.DefaultValue = &exp
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}
