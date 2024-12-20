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

	mutable := false
	variadic := false

	if p.curTokenIs(token.MUTABLE) {
		mutable = true
		p.nextToken()
	} else if p.curTokenIs(token.VAR) {
		variadic = true
		p.nextToken()
	}

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: mutable, VariadicType: variadic}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		exp := p.parseExpression(LOWEST)
		ident.DefaultValue = &exp
		ident.ObjType = exp.Type()
	}

	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		mutable = false
		variadic = false

		if p.curTokenIs(token.MUTABLE) {
			mutable = true
			p.nextToken()
		} else if p.curTokenIs(token.VAR) {
			variadic = true
			p.nextToken()
		}

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: mutable, VariadicType: variadic}

		if p.peekTokenIs(token.ASSIGN) {
			p.nextToken()
			p.nextToken()
			exp := p.parseExpression(LOWEST)
			ident.DefaultValue = &exp
			ident.ObjType = exp.Type()
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}
