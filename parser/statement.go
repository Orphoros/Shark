package parser

import (
	"shark/ast"
	"shark/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET, token.VAR:
		if p.peekTokenIs(token.LPAREN) {
			return p.parseTupleDestructuring()
		}
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt

}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	// TODO: Support no value return
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	variadicType := false
	mutable := false

	if p.curTokenIs(token.VAR) {
		variadicType = true
	}

	stmt := &ast.LetStatement{Token: p.curToken}

	// check if mutable for let statement
	if !variadicType && p.peekTokenIs(token.MUTABLE) {
		mutable = true
		p.nextToken()
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: mutable, VariadicType: variadicType}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if stmt.Value == nil {
		return nil
	}

	stmt.Name.ObjType = stmt.Value.Type()

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
		fl.Token.Pos = stmt.Name.Token.Pos
	}

	if p.peekTokenIs(token.SEMICOLON) {

		p.nextToken()

	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseTupleDestructuring() *ast.TupleDeconstruction {
	variadicType := false

	if p.curTokenIs(token.VAR) {
		variadicType = true

	}

	stmt := &ast.TupleDeconstruction{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if variadicType && !p.expectPeek(token.IDENT) {
		return nil
	} else if !variadicType {
		p.nextToken()
	}

	stmt.Names = p.parseIdentifierList(variadicType)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIdentifierList(variadicType bool) []*ast.Identifier {
	var identifiers []*ast.Identifier
	mutable := false

	if p.curTokenIs(token.RPAREN) {
		return identifiers
	}

	if p.curTokenIs(token.MUTABLE) {
		mutable = true
		p.nextToken()
	}

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: mutable, VariadicType: variadicType}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()

		if variadicType && !p.expectPeek(token.IDENT) {
			return nil
		} else if !variadicType {
			p.nextToken()
		}

		if p.curTokenIs(token.MUTABLE) {
			mutable = true
			p.nextToken()
		}
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: mutable, VariadicType: variadicType}
		identifiers = append(identifiers, ident)
	}

	return identifiers
}
