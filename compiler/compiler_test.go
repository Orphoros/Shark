package compiler

import (
	"fmt"
	"shark/ast"
	"shark/code"
	"shark/lexer"
	"shark/object"
	"shark/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	t.Run("should compile integer arithmetic", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "1 + 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1; 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "2 / 1",
				expectedConstants: []interface{}{2, 1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDiv),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "50 / 2 * 2 + 10 - 5",
				expectedConstants: []interface{}{50, 2, 10, 5},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDiv),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpMul),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpAdd),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpSub),
					code.Make(code.OpPop),
				},
			},
			{
				input: "-50",

				expectedConstants: []interface{}{50},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpMinus),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "5 ** 2",
				expectedConstants: []interface{}{5, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPower),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestConstantsPool(t *testing.T) {
	t.Run("should reuse same constants", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "1 + 1",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				},
			},
			{
				input:             `let a = "Hello World"; let b = "Hello World";`,
				expectedConstants: []interface{}{"Hello World"},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
		}

		runCompilerTests(t, tests)

	})
}

func TestBooleanExpressions(t *testing.T) {
	t.Run("should compile boolean expressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "true",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "false",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpFalse),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 > 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 < 2",
				expectedConstants: []interface{}{2, 1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 >= 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThanEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 <= 2",
				expectedConstants: []interface{}{2, 1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThanEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 == 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 != 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpNotEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "true == false",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpFalse),
					code.Make(code.OpEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "true != false",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpFalse),
					code.Make(code.OpNotEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "true && false",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpFalse),
					code.Make(code.OpAnd),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "true || false",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpFalse),
					code.Make(code.OpOr),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "!true",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpBang),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestConditionals(t *testing.T) {
	t.Run("should compile conditionals", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "if (true) { 10 }; 3333",
				expectedConstants: []interface{}{10, 3333},
				expectedInstructions: []code.Instructions{
					// 0000
					code.Make(code.OpTrue),
					// 0001
					code.Make(code.OpJumpNotTruthy, 10),
					// 0004
					code.Make(code.OpConstant, 0),
					// 0007
					code.Make(code.OpJump, 11),
					// 0010
					code.Make(code.OpNull),
					// 0011
					code.Make(code.OpPop),
					// 0012
					code.Make(code.OpConstant, 1),
					// 0015
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})

	t.Run("should compile conditionals with nulls", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "if (true) { }",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					// 0000
					code.Make(code.OpTrue),
					// 0001
					code.Make(code.OpJumpNotTruthy, 8),
					// 0004
					code.Make(code.OpNull),
					// 0005
					code.Make(code.OpJump, 9),
					// 0008
					code.Make(code.OpNull),
					// 0009
					code.Make(code.OpPop),
				},
			},
			{
				input:             "if (false) { } else { }",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					// 0000
					code.Make(code.OpFalse),
					// 0001
					code.Make(code.OpJumpNotTruthy, 8),
					// 0004
					code.Make(code.OpNull),
					// 0005
					code.Make(code.OpJump, 9),
					// 0008
					code.Make(code.OpNull),
					// 0009
					code.Make(code.OpPop),
				},
			},
			{
				input:             "if (true) { let a = 1; }",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					// 0000
					code.Make(code.OpTrue),
					// 0001
					code.Make(code.OpJumpNotTruthy, 14),
					// 0004
					code.Make(code.OpConstant, 0),
					// 0007
					code.Make(code.OpSetGlobal, 0),
					// 0010
					code.Make(code.OpNull),
					// 0011
					code.Make(code.OpJump, 15),
					// 0014
					code.Make(code.OpNull),
					// 0015
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestWhileStatement(t *testing.T) {
	t.Run("should compile while statements", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "if (true) { let a = 1; }",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					// 0000
					code.Make(code.OpTrue),
					// 0001
					code.Make(code.OpJumpNotTruthy, 14),
					// 0004
					code.Make(code.OpConstant, 0),
					// 0007
					code.Make(code.OpSetGlobal, 0),
					// 0010
					code.Make(code.OpNull),
					// 0011
					code.Make(code.OpJump, 15),
					// 0014
					code.Make(code.OpNull),
					// 0015
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut a = 1; while (a < 10) { a++; }",
				expectedConstants: []interface{}{1, 10},
				expectedInstructions: []code.Instructions{
					// 0000
					code.Make(code.OpConstant, 0),
					// 0003
					code.Make(code.OpSetGlobal, 0),
					// 0006
					code.Make(code.OpConstant, 1),
					// 0009
					code.Make(code.OpGetGlobal, 0),
					// 0012
					code.Make(code.OpGreaterThan),
					// 0013
					code.Make(code.OpJumpNotTruthy, 25),
					// 0016
					code.Make(code.OpGetGlobal, 0),
					// 0019
					code.Make(code.OpIncrementGlobal, 0),
					// 0022
					code.Make(code.OpJump, 6),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestGlobalLetStatements(t *testing.T) {
	t.Run("should compile global let statements", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "let one = 1; let two = 2;",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetGlobal, 1),
				},
			},
			{
				input:             "let one = 1; one;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let one = 1; let two = one; two;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
					code.Make(code.OpGetGlobal, 1),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})

	t.Run("should compile global let statements with reassign", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "let mut points = 1; points = 2;",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut a = 1; let mut b = 2; let x = a = b = 3;",
				expectedConstants: []interface{}{1, 2, 3},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetGlobal, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpSetGlobal, 1),
					code.Make(code.OpGetGlobal, 1),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpSetGlobal, 2),
				},
			},
			{
				input:             "let mut a = 1; a += 2;",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut a = 1; a -= 2;",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSub),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut a = 1; a *= 2;",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpMul),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut a = 1; a /= 2;",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDiv),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})

	t.Run("should compile let statement scopes", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
					let num = 55;
					() => { num }
				`,
				expectedConstants: []interface{}{
					55,
					[]code.Instructions{
						code.Make(code.OpGetGlobal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: `
					() => {
						let num = 55;
						num
					}
				`,
				expectedConstants: []interface{}{
					55,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: `
					() => {
						let a = 55;
						let b = 77;
						a + b
					}
				`,
				expectedConstants: []interface{}{
					55,
					77,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpSetLocal, 1),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpGetLocal, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestVariablePostfix(t *testing.T) {
	t.Run("should compile variable postfix expressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "let mut one = 1; one++;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpIncrementGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut one = 1; one--;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpDecrementGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut one = 1; let two = one++;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpIncrementGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
			{
				input: "let a = () => { let mut b = 1; b++; };",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpIncrementLocal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
				},
			},
			{
				input: "let mut a = 1; let b = () => { a++; };",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpGetGlobal, 0),
						code.Make(code.OpIncrementGlobal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestVariablePrefix(t *testing.T) {
	t.Run("should compile variable prefix expressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "let mut one = 1; ++one;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpIncrementGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut one = 1; --one;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpDecrementGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             `let t = "hello"; ...t;`,
				expectedConstants: []interface{}{"hello"},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpSpread),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "let mut one = 1; let two = ++one;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpIncrementGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
			{
				input:             "let mut one = 1; let two = --one;",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpDecrementGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
			{
				input: "let a = () => { let mut b = 1; ++b; };",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpIncrementLocal, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
				},
			},
			{
				input: "let a = () => { let mut b = 1; --b; };",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpDecrementLocal, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
				},
			},
			{
				input: "let mut a = 1; let b = () => { ++a; };",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpIncrementGlobal, 0),
						code.Make(code.OpGetGlobal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
			{
				input: "let mut a = 1; let b = () => { --a; };",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpDecrementGlobal, 0),
						code.Make(code.OpGetGlobal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestStringExpressions(t *testing.T) {
	t.Run("should compile string expressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             `"shark"`,
				expectedConstants: []interface{}{"shark"},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             `"sha" + "rk"`,
				expectedConstants: []interface{}{"sha", "rk"},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
}

func TestArrayLiterals(t *testing.T) {
	t.Run("should compile array literals", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "[]",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpArray, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "[1, 2, 3]",
				expectedConstants: []interface{}{1, 2, 3},

				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpArray, 3),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "[1 + 2, 3 - 4, 5 * 6]",
				expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpSub),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpMul),
					code.Make(code.OpArray, 3),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestTupleLiterals(t *testing.T) {
	t.Run("should compile tuple literals", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "(1, 2)",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpTuple, 2),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "(1 + 2, 3 - 4, 5 * 6)",
				expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpSub),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpMul),
					code.Make(code.OpTuple, 3),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestTupleDeconstructExpressions(t *testing.T) {
	t.Run("should compile tuple deconstruct expressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "let (a, b) = (1, 2);",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpTuple, 2),
					code.Make(code.OpTupleDeconstruct, 2),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
				},
			},
			{
				input:             "let (a, b, c) = (1, 2, 3);",
				expectedConstants: []interface{}{1, 2, 3},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpTuple, 3),
					code.Make(code.OpTupleDeconstruct, 3),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
					code.Make(code.OpSetGlobal, 2),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestRangeOperator(t *testing.T) {
	t.Run("should spread incrementing numbers as range", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "1..10",
				expectedConstants: []interface{}{1, 10},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpRange),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestHashLiterals(t *testing.T) {
	t.Run("should compile hash literals", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "{}",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpHash, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input:             `{1: 2, 3: 4, 5: 6}`,
				expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpHash, 6),
					code.Make(code.OpPop),
				},
			},
			{
				input:             `{1: 2 + 3, 4: 5 * 6}`,
				expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpAdd),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpMul),
					code.Make(code.OpHash, 4),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestIndexExpressions(t *testing.T) {
	t.Run("should compile index expressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "[1, 2, 3][1 + 1]",
				expectedConstants: []interface{}{1, 2, 3},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpArray, 3),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpIndex),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "{1: 2}[2 - 1]",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpHash, 2),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSub),
					code.Make(code.OpIndex),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})

	t.Run("should compile index reassign expressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "[1, 2, 3][1] = 10",
				expectedConstants: []interface{}{10, 1, 2, 3},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpArray, 3),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpIndexAssign),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "{1: 2}[1] = 10",
				expectedConstants: []interface{}{10, 1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpHash, 2),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpIndexAssign),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestFunction(t *testing.T) {
	t.Run("should compile function literals", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: "() => { return 5 + 10}",
				expectedConstants: []interface{}{
					5,
					10,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: "() => { 5 + 10 }",
				expectedConstants: []interface{}{
					5,
					10,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: "() => { 1; 2 }",
				expectedConstants: []interface{}{
					1,
					2,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpPop),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestFunctionWithoutReturnValue(t *testing.T) {
	t.Run("should compile function literals without return value", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: "() => {}",
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpReturn),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestFunctionCalls(t *testing.T) {
	t.Run("should compile function calls", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: "() => { 24 }()",
				expectedConstants: []interface{}{
					24,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: "let a = () => { 24 }; a();",
				expectedConstants: []interface{}{
					24,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: "let oneArg = (b) => { b }; oneArg(24)",
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
					24,
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
				},
			},
			{
				input: "let manyArg = (c, d, e) => { c; d; e; }; manyArg(24, 25, 26)",
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpPop),
						code.Make(code.OpGetLocal, 1),
						code.Make(code.OpPop),
						code.Make(code.OpGetLocal, 2),
						code.Make(code.OpReturnValue),
					},
					24,
					25,
					26,
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpCall, 3),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})

	t.Run("should compile function calls with default arguments", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: "let manyArg = (a = 1) => { a }",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocalDefault, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
				},
			},
			{
				input: "let x = [7, 8, 9, 10]; let manyArg = (a = 1, b = 2) => { let f = 7; a; b; }; manyArg()",
				expectedConstants: []interface{}{
					7,
					8,
					9,
					10,
					1,
					2,
					[]code.Instructions{
						code.Make(code.OpConstant, 4),
						code.Make(code.OpSetLocalDefault, 0),
						code.Make(code.OpConstant, 5),
						code.Make(code.OpSetLocalDefault, 1),
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 2),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpPop),
						code.Make(code.OpGetLocal, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpArray, 4),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpClosure, 6, 0),
					code.Make(code.OpSetGlobal, 1),
					code.Make(code.OpGetGlobal, 1),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: "let manyArg = (a, b = 2) => { a; b; };",
				expectedConstants: []interface{}{
					2,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocalDefault, 1),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpPop),
						code.Make(code.OpGetLocal, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
				},
			},
		}

		runCompilerTests(t, tests)
	})

	t.Run("should compile function calls with default argument overrides", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: "let manyArg = (a = 1) => { a }; manyArg(2)",
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocalDefault, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
					2,
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestCompilerScopes(t *testing.T) {
	t.Run("should compile with scopes", func(t *testing.T) {
		compiler := New()
		if compiler.scopeIndex != 0 {
			t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
		}
		globalSymbolTable := compiler.symbolTable

		compiler.emit(code.OpMul)

		compiler.enterScope()
		if compiler.scopeIndex != 1 {
			t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
		}

		compiler.emit(code.OpSub)

		if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
			t.Errorf("instructions length wrong. got=%d",
				len(compiler.scopes[compiler.scopeIndex].instructions))
		}

		last := compiler.scopes[compiler.scopeIndex].lastInstruction
		if last.opcode != code.OpSub {
			t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
				last.opcode, code.OpSub)
		}

		if compiler.symbolTable.Outer != globalSymbolTable {
			t.Errorf("compiler did not enclose symbolTable")
		}

		compiler.leaveScope()
		if compiler.scopeIndex != 0 {
			t.Errorf("scopeIndex wrong. got=%d, want=%d",
				compiler.scopeIndex, 0)
		}

		if compiler.symbolTable != globalSymbolTable {
			t.Errorf("compiler did not restore global symbol table")
		}
		if compiler.symbolTable.Outer != nil {
			t.Errorf("compiler modified global symbol table incorrectly")
		}

		compiler.emit(code.OpAdd)

		if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
			t.Errorf("instructions length wrong. got=%d",
				len(compiler.scopes[compiler.scopeIndex].instructions))
		}

		last = compiler.scopes[compiler.scopeIndex].lastInstruction
		if last.opcode != code.OpAdd {
			t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
				last.opcode, code.OpAdd)
		}

		previous := compiler.scopes[compiler.scopeIndex].previousInstruction
		if previous.opcode != code.OpMul {
			t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d",
				previous.opcode, code.OpMul)
		}
	})
}

func TestBuiltins(t *testing.T) {
	t.Run("should compile builtins", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
					len([]);
					push([], 1);
				`,
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpGetBuiltin, 2),
					code.Make(code.OpArray, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetBuiltin, 6),
					code.Make(code.OpArray, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpCall, 2),
					code.Make(code.OpPop),
				},
			},
			{
				input: `
					() => { len([]) }`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetBuiltin, 2),
						code.Make(code.OpArray, 0),
						code.Make(code.OpCall, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestClosures(t *testing.T) {
	t.Run("should compile closures", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
				(a) => {
					(b) => {
						a + b
					}
				}
				`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetFree, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
					[]code.Instructions{
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpClosure, 0, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: `
				(a) => {
					(b) => {
						(c) => {
							a + b + c
						}
					}
				};
				`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetFree, 0),
						code.Make(code.OpGetFree, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
					[]code.Instructions{
						code.Make(code.OpGetFree, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpClosure, 0, 2),
						code.Make(code.OpReturnValue),
					},
					[]code.Instructions{
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpClosure, 1, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			},
			{
				input: `
				let global = 55;
	
				() => {
					let a = 66;
	
					() => {
						let b = 77;
	
						() => {
							let c = 88;
	
							global + a + b + c;
						}
					}
				}
				`,
				expectedConstants: []interface{}{
					55,
					66,
					77,
					88,
					[]code.Instructions{
						code.Make(code.OpConstant, 3),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpGetGlobal, 0),
						code.Make(code.OpGetFree, 0),
						code.Make(code.OpAdd),
						code.Make(code.OpGetFree, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
					[]code.Instructions{
						code.Make(code.OpConstant, 2),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpGetFree, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpClosure, 4, 2),
						code.Make(code.OpReturnValue),
					},
					[]code.Instructions{
						code.Make(code.OpConstant, 1),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpClosure, 5, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpClosure, 6, 0),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
}

func TestRecursiveFunctions(t *testing.T) {
	t.Run("should compile recursive functions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
				let countDown = (x) => { countDown(x - 1); };
				countDown(1);
				`,
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpCurrentClosure),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSub),
						code.Make(code.OpCall, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
				},
			},
			{
				input: `
				let wrapper = () => {
					let countDown = (x) => { countDown(x - 1); };
					countDown(1);
				};
				wrapper();
				`,
				expectedConstants: []interface{}{
					1,
					[]code.Instructions{
						code.Make(code.OpCurrentClosure),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSub),
						code.Make(code.OpCall, 1),
						code.Make(code.OpReturnValue),
					},
					[]code.Instructions{
						code.Make(code.OpClosure, 1, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpConstant, 0),
						code.Make(code.OpCall, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)

	})
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()

		if err, _ := compiler.Compile(program); err != nil {
			t.Fatalf("compiler error: %+v", err)
		}

		bytecode := compiler.Bytecode()

		if err := testInstructions(tt.expectedInstructions, bytecode.Instructions); err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		if err := testConstants(t, tt.expectedConstants, bytecode.Constants); err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(&input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q", i, concatted, actual)
		}
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	t.Helper()
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			if err := testIntegerObject(int64(constant), actual[i]); err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case string:
			if err := testStringObject(constant, actual[i]); err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}
		case []code.Instructions:
			fn, ok := actual[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - object is not CompiledFunction. got=%T (%+v)", i, actual[i], actual[i])
			}

			if err := testInstructions(constant, fn.Instructions); err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s", i, err)
			}
		}
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}
