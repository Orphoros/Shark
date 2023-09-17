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
			p.errors = append(p.errors, exception.SharkError{
				ErrCause: []exception.SharkErrorCause{
					{
						Line:     p.curToken.Line,
						LineTo:   p.curToken.LineTo,
						ColTo:    p.curToken.ColTo,
						Col:      p.curToken.ColFrom,
						CauseMsg: "expected expression here",
					},
				},
				ErrMsg:     "expected to receive an expression, but got EOF instead",
				ErrHelpMsg: "Try removing the invalid syntax.",
				ErrCode:    exception.SharkErrorEOF,
				ErrType:    exception.SharkErrorTypeParser,
			})
			return nil
		}
		p.errors = append(p.errors, exception.SharkError{
			ErrCause: []exception.SharkErrorCause{
				{
					Line:     p.curToken.Line,
					LineTo:   p.curToken.LineTo,
					ColTo:    p.curToken.ColTo,
					Col:      p.curToken.ColFrom,
					CauseMsg: "this is not an expression",
				},
			},
			ErrHelpMsg: "Try removing the invalid syntax.",
			ErrMsg:     "expected to receive an expression that evaluates to a value, but got '" + p.curToken.Literal + "' instead, which has no value",
			ErrType:    exception.SharkErrorTypeParser,
			ErrCode:    exception.SharkErrorExpectedExpression,
		})
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && (precedence < p.peekPrecedence() || p.isPrecedenceAssign(precedence)) {
		postfix := p.postfixParseFns[p.peekToken.Type]
		if postfix != nil {
			ident, ok := leftExp.(*ast.Identifier)
			if !ok {
				p.errors = append(p.errors, exception.SharkError{
					ErrCause: []exception.SharkErrorCause{
						{
							Line:     p.curToken.Line,
							LineTo:   p.curToken.LineTo,
							ColTo:    p.curToken.ColTo,
							Col:      p.curToken.ColFrom,
							CauseMsg: "this is not an identifier",
						},
					},
					ErrMsg:     "expected to receive an identifier, but got '" + p.curToken.Literal + "' instead",
					ErrHelpMsg: "Try removing the invalid syntax.",
					ErrCode:    exception.SharkErrorExpectedIdentifier,
					ErrType:    exception.SharkErrorTypeParser,
				})
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

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
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
