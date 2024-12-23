package parser

import (
	"shark/ast"
	"shark/exception"
	"shark/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		// TODO: add support to assign new value to variable
		if p.curToken.Type == token.EOF {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorEOF, nil,
				"Try removing the invalid syntax.",
				exception.NewSharkErrorCause("expected expression here", p.curToken.Pos),
			))
			return nil
		}
		p.errors = append(p.errors, newSharkError(exception.SharkErrorExpectedExpression, p.curToken.Literal,
			"Try removing the invalid syntax.",
			exception.NewSharkErrorCause("this is not an expression", p.curToken.Pos),
		))
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && (precedence < p.peekPrecedence() || p.isPrecedenceAssign(precedence)) {
		postfix := p.postfixParseFns[p.peekToken.Type]
		if postfix != nil {
			ident, ok := leftExp.(*ast.Identifier)
			if !ok {
				p.errors = append(p.errors, newSharkError(exception.SharkErrorExpectedIdentifier, p.curToken.Literal,
					"Try removing the invalid syntax.",
					exception.NewSharkErrorCause("this is not an identifier", p.curToken.Pos),
				))
				return nil
			}
			p.nextToken()
			leftExp = postfix(ident)
		}
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseTupleLiteral(curToken token.Token, firstItem ast.Expression) ast.Expression {
	tuple := &ast.TupleLiteral{Token: curToken}
	var list = []ast.Expression{firstItem}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return tuple
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	tuple.Elements = list

	return tuple
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) isPrecedenceAssign(precedence int) bool {
	return precedence == p.peekPrecedence() && (p.peekToken.Type == token.ASSIGN ||
		p.peekToken.Type == token.PLUS_EQ ||
		p.peekToken.Type == token.MIN_EQ ||
		p.peekToken.Type == token.MUL_EQ ||
		p.peekToken.Type == token.DIV_EQ)
}
