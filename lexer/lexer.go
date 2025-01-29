package lexer

import (
	"shark/exception"
	"shark/token"
	"strings"
	"unicode"
)

// Lexer is a struct that is used to tokenize the input string.
type Lexer struct {
	characters   []rune
	errors       []exception.SharkError
	position     int
	readPosition int
	curLine      int
	prevLine     int
	curCol       int
	prevCol      int
	ch           rune
}

// Creates a new Lexer struct and initializes it with the input string.
// The currentLine field is set to 1. The parameter input is the string
// that will be tokenized.
func New(input *string) *Lexer {
	l := &Lexer{characters: []rune(*input), curLine: 1, curCol: 1, prevLine: 1, prevCol: 1}
	l.readChar()
	return l
}

// Reads the next character and advances the position in the input string.
// The read values are stored in the Lexer struct.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.characters) {
		l.ch = rune(0)
	} else {
		l.ch = l.characters[l.readPosition]
	}
	l.advancePosition()
}

// Advances to the next token in the input string and returns the new token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	if l.ch == '/' && l.peekChar() == '/' {
		l.skipSingleLineComment()
		return l.NextToken()
	}
	if l.ch == '/' && l.peekChar() == '*' {
		l.skipMultiLineComment()
	}
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.EQ, string(ch)+string(l.ch))
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.ARROW, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ASSIGN, string(l.ch))
		}
	case ';':
		tok = l.newToken(token.SEMICOLON, string(l.ch))
	case '(':
		tok = l.newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.ch))
	case '"':
		tok = l.newToken(token.STRING, l.readString())
	case ',':
		tok = l.newToken(token.COMMA, string(l.ch))
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.PLUS_PLUS, string(ch)+string(l.ch))
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.PLUS_EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.PLUS, string(l.ch))
		}
	case '{':
		tok = l.newToken(token.LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(token.RBRACE, string(l.ch))
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.ch))
	case '?':
		tok = l.newToken(token.QUESTION, string(l.ch))
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.MINUS_MINUS, string(ch)+string(l.ch))
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.MIN_EQ, string(ch)+string(l.ch))
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.POINTER, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.MINUS, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.BANG, string(l.ch))
		}
	case '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.POW, string(ch)+string(l.ch))
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.MUL_EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ASTERISK, string(l.ch))
		}
	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			// check for spread operator, else it's a range operator
			if l.peekChar() == '.' {
				l.readChar()
				tok = l.newToken(token.SPREAD, string(ch)+string(ch)+string(l.ch))
			} else {
				tok = l.newToken(token.RANGE, string(ch)+string(l.ch))
			}
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.DIV_EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.SLASH, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.LTE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.LT, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.GTE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.GT, string(l.ch))
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.AND, string(ch)+string(l.ch))
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.OR, string(ch)+string(l.ch))
		}
	case ':':
		tok = l.newToken(token.COLON, string(l.ch))
	case 0:
		tok = l.newToken(token.EOF, "")
	default:
		if unicode.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Pos = token.Position{Line: l.prevLine, ColFrom: l.prevCol, LineTo: l.curLine, ColTo: l.curCol}
			l.registerPosition()
			l.readChar()
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			tok.Pos = token.Position{Line: l.prevLine, ColFrom: l.prevCol, LineTo: l.curLine, ColTo: l.curCol}
			l.registerPosition()
			l.readChar()
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}
	l.readChar()
	return tok
}

// Reads a number from the input and returns it as a string.
func (l *Lexer) readNumber() string {
	str := ""
	accept := "0123456789_"
	if l.ch == '0' && l.peekChar() == 'x' {
		accept = "0x123456789abcdefABCDEF_"
	}
	if l.ch == '0' && l.peekChar() == 'b' {
		accept = "0b01_"
	}
	if l.ch == '0' && l.peekChar() == 'o' {
		accept = "o01234567_"
	}
	for strings.Contains(accept, string(l.ch)) {
		if l.ch == '_' {
			if !strings.Contains(accept, string(l.peekChar())) {
				l.errors = append(l.errors, newSharkError(exception.SharkErrorInvalidNumber, nil,
					"Remove the underscore from the number",
					exception.NewSharkErrorCause("Underscore is not permitted here", token.Position{Line: l.curLine, ColFrom: l.curCol - 1, LineTo: l.curLine, ColTo: l.curCol}),
				))
				break
			}
		} else {
			str += string(l.ch)
			if !strings.Contains(accept, string(l.peekChar())) {
				break
			}
		}
		l.readChar()
	}
	return str
}

// Reads an identifier and advances the lexer's position until it encounters a non-letter character.
// Returns the identifier as a string.
func (l *Lexer) readIdentifier() string {
	id := ""
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		id += string(l.ch)
		if !unicode.IsLetter(l.peekChar()) && !unicode.IsDigit(l.peekChar()) && l.peekChar() != '_' {
			break
		}
		l.readChar()
	}

	return id
}

// Reads a string between two double quotes and returns it.
func (l *Lexer) readString() string {
	out := ""
	multilineMode := false
	for {
		l.readChar()
		if isNewLine(l.ch) {
			l.registerNewlinePosition()
			multilineMode = true
		}
		if l.ch == 0 {
			l.errors = append(l.errors, newSharkError(exception.SharkErrorUnterminatedString, nil,
				"Add a closing double quote to the end of the string",
				exception.NewSharkErrorCause("There is no closing double quote before the end of the file", token.Position{Line: l.prevLine, ColFrom: l.prevCol, LineTo: l.curLine, ColTo: l.curCol}),
			))
			break
		}
		if l.ch == '"' {
			break
		}
		if l.ch == '\\' {
			l.readChar()
			if l.ch == 'n' {
				l.ch = '\n'
			}
			if l.ch == 'r' {
				l.ch = '\r'
			}
			if l.ch == 't' {
				l.ch = '\t'
			}
			if l.ch == '"' {
				l.ch = '"'
			}
			if l.ch == '\\' {
				l.ch = '\\'
			}
		}
		out = out + string(l.ch)
	}
	if multilineMode {
		l.errors = append(l.errors, newSharkError(exception.SharkErrorUnterminatedString, nil,
			"Use a \\n instead to create multiline strings in double quoted strings",
			exception.NewSharkErrorCause("String is started here", token.Position{Line: l.prevLine, ColFrom: l.prevCol, LineTo: l.prevLine, ColTo: l.prevCol + 1}),
			exception.NewSharkErrorCause("String ends here", token.Position{Line: l.curLine, ColFrom: l.curCol - 1, LineTo: l.curLine, ColTo: l.curCol}),
		))
	}
	return out
}

// Advances the lexer's position until it encounters a non-digit character.
// This function will increment the current line number if it encounters a isNewLine character.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if isNewLine(l.ch) {
			l.registerNewlinePosition()
		}
		l.registerPosition()
		l.readChar()
	}
}

// Creates a new token with the given type, character, and line number.
func (l *Lexer) newToken(tokenType token.Type, literal string) token.Token {
	colTo := l.curCol
	lineTo := l.curLine
	colFrom := l.prevCol
	lineFrom := l.prevLine
	l.registerPosition()
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Pos: token.Position{
			Line:    lineFrom,
			ColFrom: colFrom,
			LineTo:  lineTo,
			ColTo:   colTo,
		},
	}
}

// Peeks at the next character in the input without incrementing the current position.
// If the current position is at the end of the input, this function will return 0.
// Returns the next character in the input.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.characters) {
		return rune(0)
	}
	return l.characters[l.readPosition]
}

// Skips a single line comment. This is a comment that starts with // and ends with a isNewLine.
// This function will increment the current line number.
func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.skipWhitespace()
}

// Skips a multiline comment. This is a comment that starts with /* and ends with */.
// This function will increment the current line number.
func (l *Lexer) skipMultiLineComment() {
	foundEndOfComment := false
	for !foundEndOfComment {
		if l.ch == 0 {
			foundEndOfComment = true
		}
		if isNewLine(l.ch) {
			l.registerPosition()
			l.registerNewlinePosition()
		}
		if l.ch == '*' && l.peekChar() == '/' {
			foundEndOfComment = true
			l.readChar()
		}
		l.readChar()
	}
	l.skipWhitespace()
}

// Registers a newline position. This function will increment the current line number
// and reset the current column number.
func (l *Lexer) registerNewlinePosition() {
	l.curLine++
	l.curCol = 1
}

// Advances the lexer's position by one character.
func (l *Lexer) advancePosition() {
	l.position = l.readPosition
	l.readPosition++
	l.curCol++
}

// Registers the current position as the previous position.
func (l *Lexer) registerPosition() {
	l.prevLine = l.curLine
	l.prevCol = l.curCol
}

// Returns the errors that occurred during lexing and pops them from the lexer.
func (l *Lexer) PopErrors() []exception.SharkError {
	errors := l.errors
	l.errors = nil
	return errors
}

// isNewLine returns true if the given rune is a newline character.
func isNewLine(r rune) bool {
	return r == '\n' || r == '\r'
}

// Checks if a given byte (character) is a digit. A digit is defined as 0-9.
// Returns true if the byte is a digit, false otherwise.
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func newSharkError(code exception.SharkErrorCode, param interface{}, helpMsg string, cause ...exception.SharkErrorCause) exception.SharkError {
	var err exception.SharkError

	if param == nil {
		err = *exception.NewSharkError(exception.SharkErrorTypeLexer, code)
	} else {
		err = *exception.NewSharkError(exception.SharkErrorTypeLexer, code, param)
	}

	if helpMsg != "" {
		err.SetHelpMsg(helpMsg)
	}

	for _, c := range cause {
		err.AddCause(c)
	}

	return err
}
