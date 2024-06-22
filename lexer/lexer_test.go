package lexer

import (
	"shark/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	t.Run("should tokenize individual tokens", func(t *testing.T) {
		input := `=+(){},;`
		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.ASSIGN, "="},
			{token.PLUS, "+"},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}
		l := New(&input)

		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})

	t.Run("should tokenize identifiers", func(t *testing.T) {
		input := `let five = 5;
let ten = 10;


let add = (x, y) => {
	x + y;
};

let result = add(five, ten);

!-/5*;

5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;

10 != 9;

"foobar"

"foo bar"

[1, 2];

{"foo": "bar"}

true && false || true;

2 ** 4;

3 >= 2;
2 <= 1;

i++;
i--;
++i;
--i;
a = 1;
a += 1;
a -= 1;
a *= 1;
a /= 1;
while (true) { let a = 1; }
1..10;
...a;
`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.ARROW, "=>"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.BANG, "!"},
			{token.MINUS, "-"},
			{token.SLASH, "/"},
			{token.INT, "5"},
			{token.ASTERISK, "*"},
			{token.SEMICOLON, ";"},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.GT, ">"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.TRUE, "true"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.ELSE, "else"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.FALSE, "false"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.INT, "10"},
			{token.EQ, "=="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.INT, "10"},
			{token.NOT_EQ, "!="},
			{token.INT, "9"},
			{token.SEMICOLON, ";"},
			{token.STRING, "foobar"},
			{token.STRING, "foo bar"},
			{token.LBRACKET, "["},
			{token.INT, "1"},
			{token.COMMA, ","},
			{token.INT, "2"},
			{token.RBRACKET, "]"},
			{token.SEMICOLON, ";"},
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.STRING, "bar"},
			{token.RBRACE, "}"},
			{token.TRUE, "true"},
			{token.AND, "&&"},
			{token.FALSE, "false"},
			{token.OR, "||"},
			{token.TRUE, "true"},
			{token.SEMICOLON, ";"},
			{token.INT, "2"},
			{token.POW, "**"},
			{token.INT, "4"},
			{token.SEMICOLON, ";"},
			{token.INT, "3"},
			{token.GTE, ">="},
			{token.INT, "2"},
			{token.SEMICOLON, ";"},
			{token.INT, "2"},
			{token.LTE, "<="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.IDENT, "i"},
			{token.PLUS_PLUS, "++"},
			{token.SEMICOLON, ";"},
			{token.IDENT, "i"},
			{token.MINUS_MINUS, "--"},
			{token.SEMICOLON, ";"},
			{token.PLUS_PLUS, "++"},
			{token.IDENT, "i"},
			{token.SEMICOLON, ";"},
			{token.MINUS_MINUS, "--"},
			{token.IDENT, "i"},
			{token.SEMICOLON, ";"},
			{token.IDENT, "a"},
			{token.ASSIGN, "="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.IDENT, "a"},
			{token.PLUS_EQ, "+="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.IDENT, "a"},
			{token.MIN_EQ, "-="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.IDENT, "a"},
			{token.MUL_EQ, "*="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.IDENT, "a"},
			{token.DIV_EQ, "/="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.WHILE, "while"},
			{token.LPAREN, "("},
			{token.TRUE, "true"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.LET, "let"},
			{token.IDENT, "a"},
			{token.ASSIGN, "="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.INT, "1"},
			{token.RANGE, ".."},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.SPREAD, "..."},
			{token.IDENT, "a"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}

		l := New(&input)

		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})
}

func TestComments(t *testing.T) {
	t.Run("should ignore single line comments", func(t *testing.T) {
		input := `
		let a = 1; // this is a comment
		// this is a comment
		2 + 3;
		 `

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.LET, "let"},
			{token.IDENT, "a"},
			{token.ASSIGN, "="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.INT, "2"},
			{token.PLUS, "+"},
			{token.INT, "3"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}
		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})

	t.Run("should ignore multiline comments", func(t *testing.T) {
		input := `
		let a = 1; /* this is a comment */
		1 + /* this is a comment */ 3;
		/* comment
		with
		multiline */
		2 + 3;
		 `

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.LET, "let"},
			{token.IDENT, "a"},
			{token.ASSIGN, "="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.INT, "1"},
			{token.PLUS, "+"},
			{token.INT, "3"},
			{token.SEMICOLON, ";"},
			{token.INT, "2"},
			{token.PLUS, "+"},
			{token.INT, "3"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}
		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})
}

func TestTokenDebugLocators(t *testing.T) {
	t.Run("should count lines", func(t *testing.T) {
		input := `let a = 1;
		// single line comment
		let b = 2;
		/* 
			multiline comment
		 */
		let c = 3;`

		l := New(&input)

		for l.NextToken().Type != token.EOF {
		}

		if l.curLine != 7 {
			t.Fatalf("expected lineCount to be 7, got %d", l.curLine)
		}
	})

	t.Run("tokens should have line count", func(t *testing.T) {
		input := `let a = 1;
		// single line comment
		let b = 2;
		/* 
			multiline comment
		 */
		let c = 3;`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
			line            int
		}{
			{token.LET, "let", 1},
			{token.IDENT, "a", 1},
			{token.ASSIGN, "=", 1},
			{token.INT, "1", 1},
			{token.SEMICOLON, ";", 1},
			{token.LET, "let", 3},
			{token.IDENT, "b", 3},
			{token.ASSIGN, "=", 3},
			{token.INT, "2", 3},
			{token.SEMICOLON, ";", 3},
			{token.LET, "let", 7},
			{token.IDENT, "c", 7},
			{token.ASSIGN, "=", 7},
			{token.INT, "3", 7},
			{token.SEMICOLON, ";", 7},
			{token.EOF, "", 7},
		}
		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
			if tok.Line != tt.line {
				t.Fatalf("tests[%d] - line wrong for token '%s', expected=%d, got=%d", i, tok.Literal, tt.line, tok.Line)
			}
		}
	})

	t.Run("should count columns", func(t *testing.T) {
		input := `let a = 1**2 <= sum( 3 );
		let b = a++;`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
			fromCol         int
			toCol           int
		}{
			{token.LET, "let", 1, 4},
			{token.IDENT, "a", 5, 6},
			{token.ASSIGN, "=", 7, 8},
			{token.INT, "1", 9, 10},
			{token.POW, "**", 10, 12},
			{token.INT, "2", 12, 13},
			{token.LTE, "<=", 14, 16},
			{token.IDENT, "sum", 17, 20},
			{token.LPAREN, "(", 20, 21},
			{token.INT, "3", 22, 23},
			{token.RPAREN, ")", 24, 25},
			{token.SEMICOLON, ";", 25, 26},
			{token.LET, "let", 3, 6},
			{token.IDENT, "b", 7, 8},
			{token.ASSIGN, "=", 9, 10},
			{token.IDENT, "a", 11, 12},
			{token.PLUS_PLUS, "++", 12, 14},
			{token.SEMICOLON, ";", 14, 15},
			{token.EOF, "", 15, 15},
		}

		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
			if tok.ColFrom != tt.fromCol {
				t.Fatalf("tests[%d] - columns wrong for token '%s', expected=(colFrom:%d, colTo:%d), got=(colFrom:%d, colTo:%d)", i, tok.Literal, tt.fromCol, tt.toCol, tok.ColFrom, tok.ColTo)
			}
		}
	})

	t.Run("should count multiline column token", func(t *testing.T) {
		input := `let a = "hello
		world";`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
			fromCol         int
			toCol           int
		}{
			{token.LET, "let", 1, 4},
			{token.IDENT, "a", 5, 6},
			{token.ASSIGN, "=", 7, 8},
			{token.STRING, "hello\n\t\tworld", 9, 9},
		}

		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
			if tok.ColFrom != tt.fromCol {
				t.Fatalf("tests[%d] - colFrom wrong for token '%s', expected=%d, got=%d", i, tok.Literal, tt.fromCol, tok.ColFrom)
			}
			if tok.ColTo != tt.toCol {
				t.Fatalf("tests[%d] - colTo wrong for token '%s', expected=%d, got=%d", i, tok.Literal, tt.toCol, tok.ColTo)
			}
		}
	})

	t.Run("should count multiline line token", func(t *testing.T) {
		input := `let a = "hello
		nice
		world";`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
			fromLine        int
			toLine          int
		}{
			{token.LET, "let", 1, 1},
			{token.IDENT, "a", 1, 1},
			{token.ASSIGN, "=", 1, 1},
			{token.STRING, "hello\n\t\tnice\n\t\tworld", 1, 3},
		}

		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
			if tok.Line != tt.fromLine {
				t.Fatalf("tests[%d] - lines wrong for token '%s', expected=(lineFrom:%d, lineTo:%d), got=(lineFrom:%d, lineTo:%d)", i, tok.Literal, tt.fromLine, tt.toLine, tok.Line, tok.LineTo)
			}
			if tok.LineTo != tt.toLine {
				t.Fatalf("tests[%d] - lines wrong for token '%s', expected=(lineFrom:%d, lineTo:%d), got=(lineFrom:%d, lineTo:%d)", i, tok.Literal, tt.fromLine, tt.toLine, tok.Line, tok.LineTo)
			}
		}
	})
}

func TestSlashString(t *testing.T) {
	t.Run("should parse slash string", func(t *testing.T) {
		input := `let a = "Hello\n\tWorld\\r";`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.LET, "let"},
			{token.IDENT, "a"},
			{token.ASSIGN, "="},
			{token.STRING, "Hello\n\tWorld\\r"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}
		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})
}

func TestSpecialNumbers(t *testing.T) {
	t.Run("should parse non-decimal numbers", func(t *testing.T) {
		input := `let a = 1;
		let b = 0b1010111;
		let c = 0xF4A06;`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.LET, "let"},
			{token.IDENT, "a"},
			{token.ASSIGN, "="},
			{token.INT, "1"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "b"},
			{token.ASSIGN, "="},
			{token.INT, "0b1010111"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "c"},
			{token.ASSIGN, "="},
			{token.INT, "0xF4A06"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}
		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})
}

func TestUnicodeCharacters(t *testing.T) {
	t.Run("should parse unicode characters", func(t *testing.T) {
		input := `let ᚠᚢᚦᚨᚱᚲ = "Γειά σου Κόσμε";`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.LET, "let"},
			{token.IDENT, "ᚠᚢᚦᚨᚱᚲ"},
			{token.ASSIGN, "="},
			{token.STRING, "Γειά σου Κόσμε"},
			{token.SEMICOLON, ";"},
		}
		l := New(&input)
		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})
}
