package vm

import (
	"fmt"
	"shark/ast"
	"shark/compiler"
	"shark/exception"
	"shark/lexer"
	"shark/object"
	"shark/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	t.Run("should evaluate integer arithmetic", func(t *testing.T) {
		tests := []vmTestCase{
			{"1", 1},
			{"2", 2},
			{"1 + 2", 3},
			{"1 - 2", -1},
			{"1 * 2", 2},
			{"4 ** 2", 16},
			{"4 / 2", 2},
			{"50 / 2 * 2 + 10 - 5", 55},
			{"5 + 5 + 5 + 5 - 10", 10},
			{"2 * 2 * 2 * 2 * 2", 32},
			{"5 * 2 + 10", 20},
			{"5 + 2 * 10", 25},
			{"5 * (2 + 10)", 60},
			{"-5", -5},
			{"-10", -10},
			{"-50 + 100 + -50", 0},
			{"5 * -10", -50},
			{"5 * (-10)", -50},
		}

		runVmTests(t, tests)
	})
}

func TestBooleanExpressions(t *testing.T) {
	t.Run("should evaluate boolean expressions", func(t *testing.T) {
		tests := []vmTestCase{
			{"true", true},
			{"false", false},
			{"1 < 2", true},
			{"1 > 2", false},
			{"1 < 1", false},
			{"1 > 1", false},
			{"3 <= 3", true},
			{"3 >= 3", true},
			{"3 <= 2", false},
			{"3 >= 2", true},
			{"1 == 1", true},
			{"1 != 1", false},
			{"1 == 2", false},
			{"1 != 2", true},
			{"true == true", true},
			{"false == false", true},
			{"true == false", false},
			{"true != false", true},
			{"false != true", true},
			{"(1 < 2) == true", true},
			{"(1 < 2) == false", false},
			{"(1 > 2) == true", false},
			{"(1 > 2) == false", true},
			{"!(if (false) { 5 })", true},
			{"true && true", true},
			{"true && false", false},
			{"false && true", false},
			{"false && false", false},
			{"true || true", true},
			{"true || false", true},
			{"false || true", true},
			{"false || false", false},
		}

		runVmTests(t, tests)
	})
}

func TestConditionals(t *testing.T) {
	t.Run("should evaluate conditionals", func(t *testing.T) {
		tests := []vmTestCase{
			{"if (true) { 10 }", 10},
			{"if (true) { 10 } else { 20 }", 10},
			{"if (false) { 10 } else { 20 }", 20},
			{"if (1) { 10 }", 10},
			{"if (1 < 2) { 10 }", 10},
			{"if (1 > 2) { 10 } else { 20 }", 20},
			{"if (1 < 2) { 10 } else { 20 }", 10},
			{"if (1 > 2) { 10 }", Null},
			{"if (false) { 10 }", Null},
			{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
		}

		runVmTests(t, tests)
	})
}

func TestWhileStatement(t *testing.T) {
	t.Run("should evaluate while statements", func(t *testing.T) {
		tests := []vmTestCase{
			{"let mut a = 0; while (a < 10) { a = a + 1; }; a", 10},
			{"let mut a = 5; while (a > 0) { a = a - 1; }; a", 0},
			{"let mut a = 10; while (a >= 5) { a--; }; a", 4},
			{"let mut a = 10; while (a <= 20) { a += 1; }; a", 21},
			{"let mut sum = 2; let mut i = 1; while (i <= 5) { sum *= 2; i++; }; sum", 64},
		}

		runVmTests(t, tests)
	})
}

func TestGlobalLetStatements(t *testing.T) {
	t.Run("should evaluate global let statements", func(t *testing.T) {
		tests := []vmTestCase{
			{"let one = 1; one", 1},
			{"let one = 1; let two = 2; one + two", 3},
			{"let one = 1; let two = one + one; one + two", 3},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate global let reassign expressions", func(t *testing.T) {
		tests := []vmTestCase{
			{"let mut one = 1; one = 3", 3},
			{"let mut x = 1; let y = 2; x = y", 2},
			{"let mut one = 1; let mut two = 2; let three = one = two = 3; three", 3},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate global let reassign expressions with value change", func(t *testing.T) {
		tests := []vmTestCase{
			{"let mut one = 1; one += 3", 4},
			{"let mut x = 5; let y = 2; x -= y", 3},
			{"let mut z = 2; z *= 3", 6},
			{"let mut y = 10; y /= 2", 5},
		}

		runVmTests(t, tests)
	})
}

func TestStringExpressions(t *testing.T) {
	t.Run("should evaluate string expressions", func(t *testing.T) {
		tests := []vmTestCase{
			{`"shark"`, "shark"},
			{`"sha" + "rk"`, "shark"},
			{`"sha" + "rk" + "y"`, "sharky"},
		}

		runVmTests(t, tests)
	})
}

func TestArrayLiterals(t *testing.T) {
	t.Run("should evaluate array literals", func(t *testing.T) {
		tests := []vmTestCase{
			{"[]", []int{}},
			{"[1, 2, 3]", []int{1, 2, 3}},
			{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
		}

		runVmTests(t, tests)
	})
}

func TestRangeOperator(t *testing.T) {
	t.Run("should evaluate range operator between numbers", func(t *testing.T) {
		tests := []vmTestCase{
			{"1..5", []int{1, 2, 3, 4, 5}},
			{"1..1", []int{1}},
			{"1..0", []int{1, 0}},
			{"0..1", []int{0, 1}},
			{"0..0", []int{0}},
			{"2..-2", []int{2, 1, 0, -1, -2}},
		}
		runVmTests(t, tests)
	})
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()

		if err := comp.Compile(program); err != nil {
			t.Fatalf("compiler error: %+v", err)
		}

		vm := NewDefault(comp.Bytecode())

		if err := vm.Run(); err != nil {
			t.Fatalf("vm error: %+v", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func TestHashLiterals(t *testing.T) {
	t.Run("should evaluate hash literals", func(t *testing.T) {
		tests := []vmTestCase{
			{"{}", map[object.HashKey]int64{}},
			{"{1: 2, 2: 3}", map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			}},
			{"{1 + 1: 2 * 2, 3 + 3: 4 * 4}", map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			}},
		}

		runVmTests(t, tests)
	})
}

func TestIndexExpressions(t *testing.T) {
	t.Run("should evaluate index expressions", func(t *testing.T) {
		tests := []vmTestCase{
			{"[1, 2, 3][1]", 2},
			{"[1, 2, 3][0 + 2]", 3},
			{"[[1, 1, 1]][0][0]", 1},
			{"[][0]", Null},
			{"[1, 2, 3][99]", Null},
			{"[1][-1]", Null},
			{"{1: 1, 2: 2}[1]", 1},
			{"{1: 1, 2: 2}[2]", 2},
			{"{1: 1}[0]", Null},
			{"{}[0]", Null},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate index reassign expressions", func(t *testing.T) {
		tests := []vmTestCase{
			{"let mut x = [1, 2, 3]; x[1] = 10; x[1]", 10},
			{"let mut x = {1: 1, 2: 2}; x[1] = 10; x[1]", 10},
		}

		runVmTests(t, tests)
	})
}

func TestCallingFunctionsWithDefaultArguments(t *testing.T) {
	t.Run("should evaluate calling functions with default arguments", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let test = (a = 1) => { a; };
			test();
			`,
				expected: 1,
			},
			{
				input: `
			let test = (a = 1) => { a; };
			test(2);
			`,
				expected: 2,
			},
			{
				input: `
			let test = (a = 1, b = 2) => { a + b; };
			test();
			`,
				expected: 3,
			},
			{
				input: `
			let test = (a = 1, b = 2) => { a + b; };
			test(2);
			`,
				expected: 4,
			},
			{
				input: `
			let test = (a = 1, b = 2) => { a + b; };
			test(2, 3);
			`,
				expected: 5,
			},
		}

		runVmTests(t, tests)
	})
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	t.Run("should evaluate calling functions without arguments", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let fivePlusTen = () => { 5 + 10; };
			fivePlusTen();
			`,
				expected: 15,
			},
			{
				input: `
			let one = () => { 1; };
			let two = () => { 2; };
			one() + two();
			`,
				expected: 3,
			},
			{
				input: `
			let a = () => { 1; };
			let b = () => { a() + 1; };
			let c = () => { b() + 1; };
			c();
			`,
				expected: 3,
			},
			{
				input: `
			let earlyExit = () => { return 99; 100; };
			earlyExit();
			`,
				expected: 99,
			},
			{
				input: `
			let earlyExit = () => { return 99; return 100; };
			earlyExit();
			`,
				expected: 99,
			},
		}

		runVmTests(t, tests)
	})
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	t.Run("should evaluate calling functions with arguments and bindings", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let identity = (a) => { a; };
			identity(4);
			`,
				expected: 4,
			},
			{
				input: `
			let sum = (a, b) => { a + b; };
			sum(1, 2);
			`,
				expected: 3,
			},
			{
				input: `
			let sum = (a, b) => {
				let c = a + b;
				c;
			};
			sum(1, 2);
			`,
				expected: 3,
			},
			{
				input: `
			let sum = (a, b) => {
				let c = a + b;
				c;
				return c;
			};
			sum(1, 2);
			`,
				expected: 3,
			},
			{
				input: `
				let sum = (a, b) => {
					let c = a + b;
					return c;
					c;
				};
				let outer = () => {
					sum(1, 2) + sum(3, 4);
				};
				outer();
				`,
				expected: 10,
			},
			{
				input: `
			let globalNum = 10;
			let sum = (a, b) => {
				let c = a + b;
				c + globalNum;
			};
			let outer = () => {
				sum(1, 2) + sum(3, 4) + globalNum;
			};
			outer() + globalNum;
			`,
				expected: 50,
			},
		}

		runVmTests(t, tests)
	})
}

func TestCallingFunctionWithWrongArguments(t *testing.T) {
	t.Run("should evaluate calling functions with wrong arguments", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			() => { 1; }(1);
			`,
				expected: exception.SharkErrorArgumentNumberMismatch,
			},
			{
				input: `
			(a) => { a; }();
			`,
				expected: exception.SharkErrorArgumentNumberMismatch,
			},
			{
				input: `
			(a, b) => { a + b; }(1);
			`,
				expected: exception.SharkErrorArgumentNumberMismatch,
			},
		}

		for _, tt := range tests {
			program := parse(tt.input)
			comp := compiler.New()
			if err := comp.Compile(program); err != nil {
				t.Fatalf("compiler error: %+v", err)
			}
			vm := NewDefault(comp.Bytecode())
			if err := vm.Run(); err == nil {
				t.Fatalf("expected error, got none")
			} else {
				if err.ErrCode != tt.expected {
					t.Fatalf("expected error %q, got %q", tt.expected, err)
				}
			}
		}
	})
}

func TestFunctionWithoutReturnValue(t *testing.T) {
	t.Run("should evaluate function without return value", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let noReturn = () => { };
			noReturn();
			`,
				expected: Null,
			},
			{
				input: `
			let noReturn = () => { };
			let noReturnTwo = () => { noReturn(); };
			noReturn();
			noReturnTwo();
			`,
				expected: Null,
			},
		}

		runVmTests(t, tests)
	})
}

func TestVariablePostfix(t *testing.T) {
	t.Run("should evaluate global variable increment postfix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let mut a = 5;
			let b = a++;
			a;
			`,
				expected: 6,
			},
			{
				input: `
			let mut a = 5;
			let b = a++;
			b;
			`,
				expected: 5,
			},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate local variable increment postfix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let test = () => {
				let mut b = 5;
				b++;
			};
			test();
			`,
				expected: 5,
			},
			{
				input: `
			let test = () => {
				let mut b = 5;
				b++;
				b;
			};
			test();
			`,
				expected: 6,
			},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate global variable increment postfix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let mut a = 5;
			let b = a--;
			a;
			`,
				expected: 4,
			},
			{
				input: `
			let mut a = 5;
			let b = a--;
			b;
			`,
				expected: 5,
			},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate local variable decrement postfix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let test = () => {
				let mut b = 5;
				b--;
			};
			test();
			`,
				expected: 5,
			},
			{
				input: `
			let test = () => {
				let mut b = 5;
				b--;
				b;
			};
			test();
			`,
				expected: 4,
			},
		}

		runVmTests(t, tests)
	})

}

func TestVariablePrefix(t *testing.T) {
	t.Run("should spread string to character array", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let a = "shark";
			[...a];
			`,
				expected: []string{"s", "h", "a", "r", "k"},
			},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate global variable decrement prefix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let mut a = 5;
			let b = --a;
			a;
			`,
				expected: 4,
			},
			{
				input: `
			let mut a = 5;
			let b = --a;
			b;
			`,
				expected: 4,
			},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate global variable increment prefix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let mut a = 5;
			let b = ++a;
			a;
			`,
				expected: 6,
			},
			{
				input: `
			let mut a = 5;
			let b = ++a;
			b;
			`,
				expected: 6,
			},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate local variable decrement prefix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let test = () => {
				let mut b = 5;
				--b;
			};
			test();
			`,
				expected: 4,
			},
			{
				input: `
			let test = () => {
				let mut b = 5;
				--b;
				b;
			};
			test();
			`,
				expected: 4,
			},
		}

		runVmTests(t, tests)
	})

	t.Run("should evaluate local variable increment prefix", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let test = () => {
				let mut b = 5;
				++b;
			};
			test();
			`,
				expected: 6,
			},
			{
				input: `
			let test = () => {
				let mut b = 5;
				++b;
				b;
			};
			test();
			`,
				expected: 6,
			},
		}

		runVmTests(t, tests)
	})

}

func TestCallingFunctionsWithBindings(t *testing.T) {
	t.Run("should evaluate calling functions with bindings", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let one = () => { let num = 1; num };
			one();
			`,
				expected: 1,
			},
			{
				input: `
			let oneAndTwo = () => { let one = 1; let two = 2; one + two; };
			oneAndTwo();
			`,
				expected: 3,
			},
			{
				input: `
			let oneAndTwo = () => { let one = 1; let two = 2; one + two; };
			let threeAndFour = () => { let three = 3; let four = 4; three + four; };
			oneAndTwo() + threeAndFour();
			`,
				expected: 10,
			},
			{
				input: `
			let firstFoobar = () => { let foobar = 50; foobar; };
			let secondFoobar = () => { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
				expected: 150,
			},
			{
				input: `
			let globalSeed = 50;
			let minusOne = () => {
				let num = 1;
				globalSeed - num;
			}
			let minusTwo = () => {
				let num = 2;
				globalSeed - num;
			}
			minusOne() + minusTwo();
			`,
				expected: 97,
			},
		}

		runVmTests(t, tests)
	})
}

func TestFirstClassFunctions(t *testing.T) {
	t.Run("should evaluate first class functions", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let returnsOne = () => { 1; };
			let returnsOneReturner = () => { returnsOne; };
			returnsOneReturner()();
			`,
				expected: 1,
			},
			{
				input: `
			let returnsOneReturner = () => {
				let returnsOne = () => { 1; };
				returnsOne;
			};
			returnsOneReturner()();
			`,
				expected: 1,
			},
		}

		runVmTests(t, tests)
	})
}

func TestBuiltinFunctions(t *testing.T) {
	t.Run("should evaluate builtin functions", func(t *testing.T) {
		tests := []vmTestCase{
			{`len("")`, 0},
			{`len("four")`, 4},
			{`len("hello world")`, 11},
			{`len(1)`, &object.Error{Message: "argument to `len` not supported, got INTEGER"}},
			{`len("one", "two")`, &object.Error{Message: "wrong number of arguments. got=2, want=1"}},
			{`len([1, 2, 3])`, 3},
			{`len([])`, 0},
			{`puts("hello", "world!")`, Null},
			{`first([1, 2, 3])`, 1},
			{`first([])`, Null},
			{`first(1)`, &object.Error{Message: "argument to `first` must be ARRAY, got INTEGER"}},
			{`last([1, 2, 3])`, 3},
			{`last([])`, Null},
			{`last(1)`, &object.Error{Message: "argument to `last` must be ARRAY, got INTEGER"}},
			{`rest([1, 2, 3])`, []int{2, 3}},
			{`rest([])`, Null},
			{`push([], 1)`, []int{1}},
			{`push(1, 1)`, &object.Error{Message: "argument to `push` must be ARRAY, got INTEGER"}},
		}

		runVmTests(t, tests)
	})
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let newClosure = (a) => {
			() => { a; };
		};
		let closure = newClosure(99);
		closure();
		`,
			expected: 99,
		},
		{
			input: `
		let newAdder = (a, b) => {
			(c) => { a + b + c };
		};
		let adder = newAdder(1, 2);
		adder(8);
		`,
			expected: 11,
		},
		{
			input: `
		let newAdder = (a, b) => {
			let c = a + b;
			(d) => { c + d };
		};
		let adder = newAdder(1, 2);

		adder(8);
		`,
			expected: 11,
		},
		{
			input: `
			let newAdderOuter = (a, b) => {
				let c = a + b;
				(d) => {
					let e = d + c;
					(f) => { e + f };
				};
			};
			let newAdderInner = newAdderOuter(1, 2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let a = 1;
			let newAdderOuter = (b) => {
				(c) => {
					(d) => { a + b + c + d };
				};
			};
			let newAdderInner = newAdderOuter(2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let newClosure = (a, b) => {
				let one = () => { a; };
				let two = () => { b; };
				() => { one() + two(); };
			};
			let closure = newClosure(9, 90);
			closure();
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	t.Run("should evaluate recursive functions", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let countDown = (x) => {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			countDown(1);
			`,
				expected: 0,
			},
			{
				input: `
			let countDown = (x) => {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			let wrapper = () => {
				countDown(1);
			};
			wrapper();
			`,
				expected: 0,
			},
			{
				input: `
			let wrapper = () => {
				let countDown = (x) => {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				};
				countDown(1);
			};
			wrapper();
			`,
				expected: 0,
			},
		}

		runVmTests(t, tests)
	})
}

func TestRecursiveFibonacci(t *testing.T) {
	t.Run("should evaluate recursive fibonacci", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let fibonacci = (x) => {
				if (x == 0) {
					return 0;
				} else {
					if (x == 1) {
						return 1;
					} else {
						fibonacci(x - 1) + fibonacci(x - 2);
					}
				}
			};
			fibonacci(15);
			`,
				expected: 610,
			},
		}

		runVmTests(t, tests)
	})
}

func TestNullIfStatements(t *testing.T) {
	t.Run("should evaluate null when consequence or alternative is not an expression", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			if (false) { 10 };
			`,
				expected: Null,
			},
			{
				input: `
			if (1 > 2) { 10 };
			`,
				expected: Null,
			},
			{
				input: `
			if (1 < 2) {  };
			`,
				expected: Null,
			},
			{
				input: `
			if (1 > 2) {  } else {  };
			`,
				expected: Null,
			},
			{
				input: `
			if (1 > 2) { 10 } else {  };
			`,
				expected: Null,
			},
			{
				input: `
			if (true) { let a = 1; };
			`,
				expected: Null,
			},
		}

		runVmTests(t, tests)
	})
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		if err := testIntegerObject(int64(expected), actual); err != nil {
			t.Fatalf("testIntegerObject failed: %s", err)
		}
	case bool:
		if err := testBooleanObject(expected, actual); err != nil {
			t.Fatalf("testBooleanObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Fatalf("object is not NULL. got=%T (%+v)", actual, actual)
		}
	case string:
		if err := testStringObject(expected, actual); err != nil {
			t.Fatalf("testStringObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", actual, actual)

			return
		}
		if len(array.Elements) != len(expected) {
			t.Fatalf("wrong number of elements. want=%d, got=%d", len(expected), len(array.Elements))

			return
		}
		for i, expectedElem := range expected {
			if err := testIntegerObject(int64(expectedElem), array.Elements[i]); err != nil {
				t.Fatalf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Fatalf("object is not Hash. got=%T (%+v)", actual, actual)

			return
		}
		if len(hash.Pairs) != len(expected) {
			t.Fatalf("wrong number of elements. want=%d, got=%d", len(expected), len(hash.Pairs))

			return
		}
		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Fatalf("no pair for given key in Pairs")

				return
			}
			if err := testIntegerObject(expectedValue, pair.Value); err != nil {
				t.Fatalf("testIntegerObject failed: %s", err)
			}
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Fatalf("object is not Error. got=%T (%+v)", actual, actual)

			return
		}
		if errObj.Message != expected.Message {
			t.Fatalf("wrong error message. expected=%q, got=%q", expected.Message, errObj.Message)
		}
	}
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

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(&input)
	p := parser.New(l)
	return p.ParseProgram()
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
