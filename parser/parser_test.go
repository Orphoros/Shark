package parser

import (
	"fmt"
	"shark/ast"
	"shark/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	t.Run("should parse let statements", func(t *testing.T) {
		input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		let mut z = 27;
		`

		l := lexer.New(&input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 4 {
			t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
		}

		tests := []struct {
			expectedIdentifier string
			mutable            bool
		}{
			{"x", false},
			{"y", false},
			{"foobar", false},
			{"z", true},
		}

		for i, tt := range tests {
			stmt := program.Statements[i]

			if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.mutable) {
				return
			}
		}
	})
}

func TestReturnStatements(t *testing.T) {
	t.Run("should parse return statements", func(t *testing.T) {
		input := `
		return 5;
		return 10;
		return 993322;
		`

		l := lexer.New(&input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 3 {
			t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
		}

		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)

			if !ok {
				t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
				continue
			}

			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
			}
		}
	})
}

func TestIdentifierExpression(t *testing.T) {
	t.Run("should parse identifier expressions", func(t *testing.T) {
		input := "foobar;"

		l := lexer.New(&input)

		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		ident, ok := stmt.Expression.(*ast.Identifier)

		if !ok {
			t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
		}

		if ident.Value != "foobar" {
			t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
		}

		if ident.TokenLiteral() != "foobar" {
			t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
		}
	})
}

func TestIntegerLiteralExpression(t *testing.T) {
	t.Run("should parse integer literal expressions", func(t *testing.T) {
		input := "5;"
		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.IntegerLiteral)

		if !ok {
			t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
		}

		if literal.Value != 5 {
			t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
		}

		if literal.TokenLiteral() != "5" {
			t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
		}
	})
}

func TestParsingPrefixExpressions(t *testing.T) {
	t.Run("should parse prefix expressions", func(t *testing.T) {
		prefixTests := []struct {
			input    string
			operator string
			value    interface{}
		}{
			{"!5;", "!", 5},
			{"-15;", "-", 15},
			{"!true;", "!", true},
			{"!false;", "!", false},
		}

		for _, tt := range prefixTests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.PrefixExpression)

			if !ok {
				t.Fatalf("exp not *ast.PrefixExpression. got=%T", stmt.Expression)
			}

			if exp.Operator != tt.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			}

			if !testLiteralExpression(t, exp.Right, tt.value) {
				return
			}
		}
	})
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	t.Helper()
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)

	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)

		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)

		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())

		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())

		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)

		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)

		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())

		return false
	}

	return true
}

func testPostfixExpression(t *testing.T, exp ast.Expression, left interface{}, right string) bool {
	opExp, ok := exp.(*ast.PostfixExpression)

	if !ok {
		t.Errorf("exp not *ast.PostfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != right {
		t.Errorf("exp.Operator is not '%s'. got=%s", right, opExp.Operator)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp not *ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLetStatement(t *testing.T, s ast.Statement, name string, mutable bool) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	if letStmt.Name.Mutable != mutable {
		t.Errorf("letStmt.Name.Mutable not '%t'. got=%t", mutable, letStmt.Name.Mutable)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()

	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestParsingInfixExpressions(t *testing.T) {
	t.Run("should parse infix expressions", func(t *testing.T) {
		infixTests := []struct {
			input    string
			leftVal  interface{}
			operator string
			rightVal interface{}
		}{
			{"5 + 5;", 5, "+", 5},
			{"5 - 5;", 5, "-", 5},
			{"5 * 5;", 5, "*", 5},
			{"5 ** 5;", 5, "**", 5},
			{"5 / 5;", 5, "/", 5},
			{"5 > 5;", 5, ">", 5},
			{"5 <= 5;", 5, "<=", 5},
			{"5 >= 5;", 5, ">=", 5},
			{"5 < 5;", 5, "<", 5},
			{"5 == 5;", 5, "==", 5},
			{"5 != 5;", 5, "!=", 5},
			{"true == true", true, "==", true},
			{"true != false", true, "!=", false},
			{"false == false", false, "==", false},
			{"true && true", true, "&&", true},
			{"false && true", false, "&&", true},
			{"1..5", 1, "..", 5},
		}

		for _, tt := range infixTests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.InfixExpression)

			if !ok {
				t.Fatalf("exp not *ast.InfixExpression. got=%T", stmt.Expression)
			}

			if !testLiteralExpression(t, exp.Left, tt.leftVal) {
				return
			}

			if exp.Operator != tt.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			}

			if !testLiteralExpression(t, exp.Right, tt.rightVal) {
				return
			}

			if !testInfixExpression(t, stmt.Expression, tt.leftVal, tt.operator, tt.rightVal) {
				return
			}
		}
	})
}

func TestParsingPostfixExpressions(t *testing.T) {
	t.Run("should parse postfix expressions", func(t *testing.T) {
		postfixTests := []struct {
			input    string
			leftVal  interface{}
			operator string
		}{
			{"a++", "a", "++"},
			{"a--", "a", "--"},
		}

		for _, tt := range postfixTests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			if !testPostfixExpression(t, stmt.Expression, tt.leftVal, tt.operator) {
				return
			}
		}
	})
}

func TestParsingVariablePrefixExpressions(t *testing.T) {
	t.Run("should parse variable prefix expressions", func(t *testing.T) {
		postfixTests := []struct {
			input    string
			rightVal interface{}
			operator string
		}{
			{"++a", "a", "++"},
			{"--a", "a", "--"},
			{"...a", "a", "..."},
		}

		for _, tt := range postfixTests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.PrefixExpression)

			if !ok {
				t.Fatalf("exp not *ast.PrefixExpression. got=%T", stmt.Expression)
			}

			if !testLiteralExpression(t, exp.Right, tt.rightVal) {
				return
			}

			if exp.RightIdent == nil {
				t.Fatalf("exp.RightIdent is nil, but was expecting to get '%s'", tt.rightVal)
			}

			if exp.RightIdent.Value != tt.rightVal {
				t.Fatalf("exp.RightIdent.Value is not '%s'. got=%s", tt.rightVal, exp.RightIdent.Value)
			}

			if exp.Operator != tt.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			}
		}
	})
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	t.Run("should parse operator precedence", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"-a * b", "((-a) * b)"},
			{"!-a", "(!(-a))"},
			{"a + b + c", "((a + b) + c)"},
			{"a + b - c", "((a + b) - c)"},
			{"a * b * c", "((a * b) * c)"},
			{"a * b / c", "((a * b) / c)"},
			{"a + b / c", "(a + (b / c))"},
			{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
			{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
			{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
			{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
			{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
			{"true", "true"},
			{"false", "false"},
			{"3 > 5 == false", "((3 > 5) == false)"},
			{"3 < 5 == true", "((3 < 5) == true)"},
			{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
			{"(5 + 5) * 2", "((5 + 5) * 2)"},
			{"2 / (5 + 5)", "(2 / (5 + 5))"},
			{"-(5 + 5)", "(-(5 + 5))"},
			{"!(true == true)", "(!(true == true))"},
			{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
			{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
			{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
			{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
			{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
			{"a++", "(a++)"},
			{"a--", "(a--)"},
			{"-a++", "(-(a++))"},
			{"-a--", "(-(a--))"},
			{"a++ + b", "((a++) + b)"},
			{"--a", "(--a)"},
			{"++a", "(++a)"},
			{"1+2..5", "((1 + 2) .. 5)"},
			{"1..5*2", "(1 .. (5 * 2))"},
			{"let a = ...a", "let a = (...a);"},
			{"a = 1 + 1", "(a = (1 + 1))"},
			{"a = b = c = 3 + 4", "(a = (b = (c = (3 + 4))))"},
			{"let a = b = 1 + 3 * 4", "let a = (b = (1 + (3 * 4)));"},
			{"let x = y = z = 1 + 3 * 4", "let x = (y = (z = (1 + (3 * 4))));"},
			{"a += b -= c", "(a += (b -= c))"},
			{"let a = b /= c *= d", "let a = (b /= (c *= d));"},
		}

		for _, tt := range tests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			actual := program.String()

			if actual != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, actual)
			}
		}
	})
}

func TestIfExpression(t *testing.T) {
	t.Run("should parse if expressions", func(t *testing.T) {
		input := `if (x < y) { x }`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)

		if !ok {
			t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Fatalf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if exp.Alternative != nil {
			t.Errorf("exp.Alternative was not nil. got=%+v", exp.Alternative)
		}
	})

	t.Run("should parse if-else expressions", func(t *testing.T) {
		input := `if (x < y) { x } else { y }`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)

		if !ok {
			t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Fatalf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if len(exp.Alternative.Statements) != 1 {
			t.Fatalf("alternative is not 1 statements. got=%d", len(exp.Alternative.Statements))
		}

		alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
		}

		if !testIdentifier(t, alternative.Expression, "y") {
			return
		}
	})
}

func TestWhileStatement(t *testing.T) {
	t.Run("should parse while statement", func(t *testing.T) {
		input := `while (a < 5) { a++; }`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.WhileStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T", program.Statements[1])
		}

		if !testInfixExpression(t, stmt.Condition, "a", "<", 5) {
			return
		}

		if len(stmt.Body.Statements) != 1 {
			t.Fatalf("body is not 1 statements. got=%d", len(stmt.Body.Statements))
		}

		body, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement in while body. got=%T", stmt.Body.Statements[0])
		}

		if !testPostfixExpression(t, body.Expression, "a", "++") {
			return
		}
	})
}

func TestFunctionLiteralParsing(t *testing.T) {
	t.Run("should parse function literals", func(t *testing.T) {
		input := `(x, y) => { x + y; }`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteral)

		if !ok {
			t.Fatalf("exp not *ast.FunctionLiteral. got=%T", stmt.Expression)
		}

		if len(function.Parameters) != 2 {
			t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
		}

		testLiteralExpression(t, function.Parameters[0], "x")
		testLiteralExpression(t, function.Parameters[1], "y")

		if len(function.Body.Statements) != 1 {
			t.Fatalf("function.Body.Statements has not 1 statements. got=%d", len(function.Body.Statements))
		}

		bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
		}

		testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
	})
}

func TestFunctionLiteralWithName(t *testing.T) {
	t.Run("should parse function literals with name", func(t *testing.T) {
		input := `let myFunction = () => { };`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Body does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", program.Statements[0])
		}
		function, ok := stmt.Value.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Value is not ast.FunctionLiteral. got=%T", stmt.Value)
		}

		if function.Name != "myFunction" {
			t.Errorf("function literal name wrong. want 'myFunction', got=%q", function.Name)
		}
	})
}

func TestFunctionParameterParsing(t *testing.T) {
	t.Run("should parse function parameters", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []string
		}{
			{"() => {};", []string{}},
			{"(x) => {};", []string{"x"}},
			{"(x, y, z) => {};", []string{"x", "y", "z"}},
		}

		for _, tt := range tests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			function := stmt.Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.expected) {
				t.Errorf("length parameters wrong. want %d, got=%d", len(tt.expected), len(function.Parameters))
			}

			for i, ident := range tt.expected {
				testLiteralExpression(t, function.Parameters[i], ident)
			}
		}
	})

	t.Run("should parse function parameters with default values", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []struct {
				identifier string
				value      interface{}
			}
		}{
			{"(x = 1) => {};", []struct {
				identifier string
				value      interface{}
			}{{"x", 1}}},
			{"(x = 1, y = 2) => {};", []struct {
				identifier string
				value      interface{}
			}{{"x", 1}, {"y", 2}}},
			{"(x = 1, y = 2, z = 3) => {};", []struct {
				identifier string
				value      interface{}
			}{{"x", 1}, {"y", 2}, {"z", 3}}},
		}

		for _, tt := range tests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			function := stmt.Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.expected) {
				t.Errorf("length parameters wrong. want %d, got=%d", len(tt.expected), len(function.Parameters))
			}

			for i, param := range tt.expected {
				testLiteralExpression(t, function.Parameters[i], param.identifier)
				testLiteralExpression(t, *function.Parameters[i].DefaultValue, param.value)
			}
		}
	})
}

func TestCallExpressionParsing(t *testing.T) {
	t.Run("should parse call expressions", func(t *testing.T) {
		input := `add(1, 2 * 3, 4 + 5);`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.CallExpression)

		if !ok {
			t.Fatalf("exp not *ast.CallExpression. got=%T", stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 3 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		testLiteralExpression(t, exp.Arguments[0], 1)

		testInfixExpression(t, exp.Arguments[1], 2, "*", 3)

		testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
	})
}

func TestLetStatement(t *testing.T) {
	t.Run("should parse let statements", func(t *testing.T) {
		tests := []struct {
			input              string
			expectedIdentifier string
			expectedValue      interface{}
		}{
			{"let x = 5;", "x", 5},
			{"let y = true;", "y", true},
			{"let foobar = y;", "foobar", "y"},
		}

		for _, tt := range tests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.LetStatement)

			if !ok {
				t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", program.Statements[0])
			}

			if stmt.Name.Value != tt.expectedIdentifier {
				t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, stmt.Name.Value)
			}

			if stmt.Name.TokenLiteral() != tt.expectedIdentifier {
				t.Errorf("stmt.Name.TokenLiteral not '%s'. got=%s", tt.expectedIdentifier, stmt.Name.TokenLiteral())
			}

			testLiteralExpression(t, stmt.Value, tt.expectedValue)
		}

	})
}

func TestReassignStatement(t *testing.T) {
	t.Run("should parse reassign statements", func(t *testing.T) {
		tests := []struct {
			input              string
			expectedIdentifier string
			expectedValue      interface{}
		}{
			{"x = 5;", "x", 5},
			{"y = true;", "y", true},
			{"foobar = y;", "foobar", "y"},
			{"a += 1;", "a", 1},
			{"b -= 2;", "b", 2},
			{"c *= 3;", "c", 3},
			{"d /= 4;", "d", 4},
		}

		for _, tt := range tests {
			l := lexer.New(&tt.input)
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			infix, ok := stmt.Expression.(*ast.InfixExpression)

			if !ok {
				t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
			}

			testLiteralExpression(t, infix.Right, tt.expectedValue)
		}

	})
}

func TestStringLiteralExpression(t *testing.T) {
	t.Run("should parse string literal expression", func(t *testing.T) {
		input := `"hello world";`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.StringLiteral)

		if !ok {
			t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
		}

		if literal.Value != "hello world" {
			t.Errorf("literal.Value not %s. got=%s", "hello world", literal.Value)
		}
	})
}

func TestParsingArrayLiterals(t *testing.T) {
	t.Run("should parse array literals", func(t *testing.T) {
		input := "[1, 2 * 2, 3 + 3]"

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		array := stmt.Expression.(*ast.ArrayLiteral)

		if len(array.Elements) != 3 {
			t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
		}

		testIntegerLiteral(t, array.Elements[0], 1)
		testInfixExpression(t, array.Elements[1], 2, "*", 2)
		testInfixExpression(t, array.Elements[2], 3, "+", 3)
	})
}

func TestParsingIndexExpressions(t *testing.T) {
	t.Run("should parse index expressions", func(t *testing.T) {
		input := "myArray[1 + 1]"

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		indexExp := stmt.Expression.(*ast.IndexExpression)

		if !testIdentifier(t, indexExp.Left, "myArray") {
			return
		}

		if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
			return
		}
	})
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	t.Run("should parse hash literals with string keys", func(t *testing.T) {
		input := `{"one": 1, "two": 2, "three": 3}`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		hash, ok := stmt.Expression.(*ast.HashLiteral)

		if !ok {
			t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
		}

		if len(hash.Pairs) != 3 {
			t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
		}

		expected := map[string]int64{
			"one":   1,
			"two":   2,
			"three": 3,
		}

		for key, value := range hash.Pairs {
			literal, ok := key.(*ast.StringLiteral)

			if !ok {
				t.Errorf("key is not ast.StringLiteral. got=%T", key)
			}

			expectedValue := expected[literal.String()]

			testIntegerLiteral(t, value, expectedValue)
		}
	})

	t.Run("should parse empty hash literals", func(t *testing.T) {
		input := "{}"

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		hash, ok := stmt.Expression.(*ast.HashLiteral)

		if !ok {
			t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
		}

		if len(hash.Pairs) != 0 {
			t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
		}
	})

	t.Run("should parse hash literals with expressions", func(t *testing.T) {
		input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

		l := lexer.New(&input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		hash, ok := stmt.Expression.(*ast.HashLiteral)

		if !ok {
			t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
		}

		if len(hash.Pairs) != 3 {
			t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
		}

		tests := map[string]func(ast.Expression){

			"one": func(e ast.Expression) {
				testInfixExpression(t, e, 0, "+", 1)
			},

			"two": func(e ast.Expression) {
				testInfixExpression(t, e, 10, "-", 8)
			},

			"three": func(e ast.Expression) {
				testInfixExpression(t, e, 15, "/", 5)
			},
		}

		for key, value := range hash.Pairs {
			literal, ok := key.(*ast.StringLiteral)

			if !ok {
				t.Errorf("key is not ast.StringLiteral. got=%T", key)
			}

			testFunc, ok := tests[literal.String()]

			if !ok {
				t.Errorf("No test function for key %q found", literal.String())
			}

			testFunc(value)
		}
	})
}
