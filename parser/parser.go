package parser

import (
	"fmt"
	"shark/ast"
	"shark/exception"
	"shark/lexer"
	"shark/token"
	"strconv"
)

const (
	_           int = iota
	LOWEST          // Shark precedence: lowest precedence
	ASSIGN          // Shark precedence: = or += or -= or /= or *=
	COND            // Shark precedence: AND or OR
	EQUALS          // Shark precedence: ==
	LESSGREATER     // Shark precedence: > or < or >= or <=
	SUM             // Shark precedence: + or -
	PRODUCT         // Shark precedence: * or /
	POWER           // Shark precedence: **
	PREFIX          // Shark precedence: -X or !X
	POSTFIX         // Shark precedence: X++ or X--
	CALL            // Shark precedence: myFunction(X)
	INDEX           // Shark precedence: array[index]
)

// Precedence table for Shark operators. The higher the value, the higher the precedence.
var precedence = map[token.TokenType]int{
	token.ASSIGN:      ASSIGN,
	token.PLUS_EQ:     ASSIGN,
	token.MIN_EQ:      ASSIGN,
	token.DIV_EQ:      ASSIGN,
	token.MUL_EQ:      POWER,
	token.EQ:          EQUALS,
	token.NOT_EQ:      EQUALS,
	token.LT:          LESSGREATER,
	token.GT:          LESSGREATER,
	token.LTE:         LESSGREATER,
	token.GTE:         LESSGREATER,
	token.AND:         COND,
	token.OR:          COND,
	token.POW:         POWER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.SLASH:       PRODUCT,
	token.ASTERISK:    PRODUCT,
	token.LPAREN:      CALL,
	token.LBRACKET:    INDEX,
	token.PLUS_PLUS:   POSTFIX,
	token.MINUS_MINUS: POSTFIX,
	token.RANGE:       PRODUCT,
}

type (
	// Function that represents a prefix parse function.
	prefixParseFn func() ast.Expression
	// Function that represents an infix parse function.
	infixParseFn func(ast.Expression) ast.Expression
	// Function that represents a postfix parse function.
	postfixParseFn func(ast.Expression) ast.Expression
)

// Parser is a Shark parser struct.
type Parser struct {
	l               *lexer.Lexer
	curToken        token.Token
	peekToken       token.Token
	errors          []exception.SharkError
	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

// Creates a new Shark parser. It takes a lexer as input.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []exception.SharkError{},
	}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.PLUS_EQ, p.parseInfixExpression)
	p.registerInfix(token.MIN_EQ, p.parseInfixExpression)
	p.registerInfix(token.DIV_EQ, p.parseInfixExpression)
	p.registerInfix(token.MUL_EQ, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.POW, p.parseInfixExpression)
	p.registerInfix(token.RANGE, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpression)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerPrefix(token.PLUS_PLUS, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS_MINUS, p.parsePrefixExpression)

	return p
}

// Parse a Shark program. It returns an AST representation of the program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Registers a shark prefix function for a token. It takes a token type and a prefix parse function as input.
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Registers a shark infix function for a token. It takes a token type and an infix parse function as input.
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Registers a shark postfix function for a token. It takes a token type and a postfix parse function as input.
func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []exception.SharkError {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	if errors := p.l.PopErrors(); errors != nil {
		// add lexer errors to parser errors
		p.errors = append(p.errors, errors...)
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	errorMsg := "Could not parse the code due to an unexpected character."
	causeMsg := fmt.Sprintf("expected '%s', but got '%s' instead", t, p.peekToken.Type)
	var suggestionMsg string
	switch t {
	case token.SEMICOLON:
		suggestionMsg = "Try adding a semicolon at the end of the statement."
	case token.RPAREN:
		suggestionMsg = "Try adding a closing parenthesis."
	case token.RBRACE:
		suggestionMsg = "Try adding a closing brace."
	case token.RBRACKET:
		suggestionMsg = "Try adding a closing bracket."
	case token.EOF:
		suggestionMsg = "Try terminated the statement."
	}

	p.errors = append(p.errors, exception.SharkError{
		ErrCause: []exception.SharkErrorCause{
			{
				CauseMsg: causeMsg,
				Col:      p.peekToken.ColFrom,
				ColTo:    p.peekToken.ColTo,
				Line:     p.peekToken.Line,
				LineTo:   p.peekToken.LineTo,
			},
		},
		ErrMsg:     errorMsg,
		ErrHelpMsg: suggestionMsg,
		ErrCode:    exception.SharkErrorUnexpectedToken,
		ErrType:    exception.SharkErrorTypeParser,
	})
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, exception.SharkError{
			ErrCause: []exception.SharkErrorCause{
				{
					Line:     p.curToken.Line,
					LineTo:   p.curToken.LineTo,
					Col:      p.curToken.ColFrom,
					ColTo:    p.curToken.ColTo,
					CauseMsg: err.Error(),
				},
			},
			ErrMsg:     "Could not parse the code due to an invalid integer literal.",
			ErrHelpMsg: "Try to use a smaller decimal number instead",
			ErrCode:    exception.SharkErrorInteger,
			ErrType:    exception.SharkErrorTypeParser,
		})
		return nil
	}

	lit.Value = value

	return lit
}

// Checks the precedence of the next token. It returns the precedence of the next token.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Checks the precedence of the current token. It returns the precedence of the current token.
func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Parse a boolean expression. It returns an AST representation of the boolean expression.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}
