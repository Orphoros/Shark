package parser

import (
	"shark/ast"
	"shark/token"
	"shark/types"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	lit.Parameters = p.parseFunctionParameters()

	var returnType types.ISharkType
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		returnType = p.parseType()
	}

	var argTypes []types.ISharkType

	for _, param := range lit.Parameters {
		argTypes = append(argTypes, param.DefinedType)
	}

	lit.DefinedType = &types.TSharkFuncType{ArgsList: argTypes, ReturnT: returnType}

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

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: mutable, IsVariadic: variadic}

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		ident.DefinedType = p.parseType()
	}

	// if p.peekTokenIs(token.ASSIGN) && !ident.DefinedType.Is(types.TSharkOptional{}) {
	// 	p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
	// 		"Default values are only allowed for optional types",
	// 		exception.NewSharkErrorCause("default values are only allowed for optional types", p.curToken.Pos),
	// 	))
	// 	return nil
	// }

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

		mutable = false
		variadic = false

		if p.curTokenIs(token.MUTABLE) {
			mutable = true
			p.nextToken()
		} else if p.curTokenIs(token.VAR) {
			variadic = true
			p.nextToken()
		}

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: mutable, IsVariadic: variadic}

		if p.peekTokenIs(token.COLON) {
			p.nextToken()
			p.nextToken()
			ident.DefinedType = p.parseType()
		}

		// if p.peekTokenIs(token.ASSIGN) && !ident.DefinedType.Is(types.TSharkOptional{}) {
		// 	p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
		// 		"Default values are only allowed for optional types",
		// 		exception.NewSharkErrorCause("default values are only allowed for optional types", p.curToken.Pos),
		// 	))
		// 	return nil
		// }

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
