package parser

import (
	"shark/exception"
	"shark/token"
	"shark/types"
)

var typeMap = map[token.Type]types.ISharkType{
	token.T_I64:      types.TSharkI64{},
	token.T_BOOL:     types.TSharkBool{},
	token.T_ANY:      types.TSharkAny{},
	token.T_STRING:   types.TSharkString{},
	token.T_TUPLE:    types.TSharkTuple{},
	token.T_ARRAY:    types.TSharkArray{},
	token.T_HASHMAP:  types.TSharkHashMap{},
	token.T_FUNCTION: types.TSharkFuncType{},
}

func (p *Parser) parseType() types.ISharkType {
	var sharkType types.ISharkType = nil

	switch p.curToken.Type {
	case token.T_I64:
		sharkType = types.TSharkI64{}
	case token.T_BOOL:
		sharkType = types.TSharkBool{}
	case token.T_ANY:
		sharkType = types.TSharkAny{}
	case token.T_STRING:
		sharkType = types.TSharkString{}
	case token.T_HASHMAP:
		p.nextToken()
		if p.curToken.Type != token.LT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '<' after 'hashmap'",
				exception.NewSharkErrorCause("missing '<' after 'hashmap'", p.curToken.Pos),
			))
			return nil
		}
		p.nextToken()
		if sharkType = typeMap[p.curToken.Type]; sharkType == nil {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeNotFound, p.curToken.Literal,
				"Did you mean 'i64', 'bool', 'string' or 'any'?",
				exception.NewSharkErrorCause("this is not a valid type", p.curToken.Pos),
			))
			return nil
		}

		indexType := p.parseType()

		p.expectPeek(token.COMMA)

		p.nextToken()

		collectionType := p.parseType()

		sharkType = types.TSharkHashMap{
			Indexes:  indexType,
			Collects: collectionType,
		}
		p.nextToken()
		if p.curToken.Type != token.GT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '>' after hashmap type",
				exception.NewSharkErrorCause("missing '>' after hashmap type", p.curToken.Pos),
			))
			return nil
		}
	case token.T_TUPLE:
		p.nextToken()
		if p.curToken.Type != token.LT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '<' after 'tuple'",
				exception.NewSharkErrorCause("missing '<' after 'tuple'", p.curToken.Pos),
			))
			return nil
		}
		p.nextToken()
		if sharkType = typeMap[p.curToken.Type]; sharkType == nil {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeNotFound, p.curToken.Literal,
				"Did you mean 'i64', 'bool', 'string' or 'any'?",
				exception.NewSharkErrorCause("this is not a valid type", p.curToken.Pos),
			))
			return nil
		}

		sharkType = types.TSharkTuple{
			Collection: p.parseTypeList(),
		}
		p.nextToken()
		if p.curToken.Type != token.GT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '>' after tuple type",
				exception.NewSharkErrorCause("missing '>' after tuple type", p.curToken.Pos),
			))
			return nil
		}
	case token.T_ARRAY:
		p.nextToken()
		if p.curToken.Type != token.LT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '<' after 'array'",
				exception.NewSharkErrorCause("missing '<' after 'array'", p.curToken.Pos),
			))
			return nil
		}
		p.nextToken()
		if sharkType = typeMap[p.curToken.Type]; sharkType == nil {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeNotFound, p.curToken.Literal,
				"Did you mean 'i64', 'bool', 'string' or 'any'?",
				exception.NewSharkErrorCause("this is not a valid type", p.curToken.Pos),
			))
			return nil
		}

		sharkType = types.TSharkArray{
			Collection: p.parseType(),
		}
		p.nextToken()
		if p.curToken.Type != token.GT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '>' after array type",
				exception.NewSharkErrorCause("missing '>' after array type", p.curToken.Pos),
			))
			return nil
		}
	case token.T_FUNCTION:
		p.nextToken()
		if p.curToken.Type != token.LT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '<' after 'func'",
				exception.NewSharkErrorCause("missing '<' after 'func'", p.curToken.Pos),
			))
			return nil
		}
		p.expectPeek(token.LPAREN)
		p.nextToken()
		argTypeList := p.parseTypeList()
		p.expectPeek(token.POINTER)
		p.nextToken()
		returnType := p.parseType()
		p.nextToken()
		if p.curToken.Type != token.GT {
			p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeSyntax, p.curToken.Literal,
				"Expected '>' after function type",
				exception.NewSharkErrorCause("missing '>' after function type", p.curToken.Pos),
			))
			return nil
		}

		sharkType = types.TSharkFuncType{
			ReturnT:  returnType,
			ArgsList: argTypeList,
		}
	default:
		p.errors = append(p.errors, newSharkError(exception.SharkErrorTypeNotFound, p.curToken.Literal,
			"Did you mean 'i64' or 'string'?",
			exception.NewSharkErrorCause("this is not a valid type", p.curToken.Pos),
		))
	}

	if p.peekToken.Type == token.QUESTION {
		sharkType = types.TSharkOptional{
			Type: sharkType,
		}
		p.nextToken()
	}

	return sharkType
}

func (p *Parser) parseTypeList() []types.ISharkType {
	sharkTypes := []types.ISharkType{}

	if p.curToken.Type == token.RPAREN {
		return sharkTypes
	}

	sharkTypes = append(sharkTypes, p.parseType())
	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		sharkTypes = append(sharkTypes, p.parseType())
	}

	p.expectPeek(token.RPAREN)

	return sharkTypes
}
