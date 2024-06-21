package code

import "testing"

func TestMake(t *testing.T) {
	t.Run("should get the correct bytecode instruction", func(t *testing.T) {
		tests := []struct {
			op       Opcode
			operands []int
			expected []byte
		}{
			{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
			{OpAdd, []int{}, []byte{byte(OpAdd)}},
			{OpPower, []int{}, []byte{byte(OpPower)}},
			{OpAnd, []int{}, []byte{byte(OpAnd)}},
			{OpOr, []int{}, []byte{byte(OpOr)}},
			{OpPop, []int{}, []byte{byte(OpPop)}},
			{OpSub, []int{}, []byte{byte(OpSub)}},
			{OpMul, []int{}, []byte{byte(OpMul)}},
			{OpDiv, []int{}, []byte{byte(OpDiv)}},
			{OpTrue, []int{}, []byte{byte(OpTrue)}},
			{OpFalse, []int{}, []byte{byte(OpFalse)}},
			{OpEqual, []int{}, []byte{byte(OpEqual)}},
			{OpNotEqual, []int{}, []byte{byte(OpNotEqual)}},
			{OpGreaterThan, []int{}, []byte{byte(OpGreaterThan)}},
			{OpGreaterThanEqual, []int{}, []byte{byte(OpGreaterThanEqual)}},
			{OpMinus, []int{}, []byte{byte(OpMinus)}},
			{OpBang, []int{}, []byte{byte(OpBang)}},
			{OpJumpNotTruthy, []int{65534}, []byte{byte(OpJumpNotTruthy), 255, 254}},
			{OpJump, []int{65534}, []byte{byte(OpJump), 255, 254}},
			{OpNull, []int{}, []byte{byte(OpNull)}},
			{OpGetGlobal, []int{65534}, []byte{byte(OpGetGlobal), 255, 254}},
			{OpSetGlobal, []int{65534}, []byte{byte(OpSetGlobal), 255, 254}},
			{OpArray, []int{65534}, []byte{byte(OpArray), 255, 254}},
			{OpHash, []int{65534}, []byte{byte(OpHash), 255, 254}},
			{OpIndex, []int{}, []byte{byte(OpIndex)}},
			{OpCall, []int{255}, []byte{byte(OpCall), 255}},
			{OpReturnValue, []int{}, []byte{byte(OpReturnValue)}},
			{OpReturn, []int{}, []byte{byte(OpReturn)}},
			{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
			{OpSetLocal, []int{255}, []byte{byte(OpSetLocal), 255}},
			{OpGetBuiltin, []int{255}, []byte{byte(OpGetBuiltin), 255}},
			{OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
			{OpIncrementLocal, []int{255}, []byte{byte(OpIncrementLocal), 255}},
			{OpIncrementGlobal, []int{65534}, []byte{byte(OpIncrementGlobal), 255, 254}},
			{OpDecrementGlobal, []int{65534}, []byte{byte(OpDecrementGlobal), 255, 254}},
			{OpDecrementLocal, []int{255}, []byte{byte(OpDecrementLocal), 255}},
			{OpSetLocalDefault, []int{255}, []byte{byte(OpSetLocalDefault), 255}},
		}

		for _, tt := range tests {
			instruction := Make(tt.op, tt.operands...)
			if len(instruction) != len(tt.expected) {
				t.Errorf("instruction has wrong length. want=%d, got=%d", len(tt.expected), len(instruction))
			}

			for i, b := range tt.expected {
				if instruction[i] != tt.expected[i] {
					t.Errorf("wrong byte at position %d. want=%d, got=%d", i, b, instruction[i])
				}
			}
		}

	})
}

func TestInstructionString(t *testing.T) {
	t.Run("should get the correct string representation of the instruction", func(t *testing.T) {
		instructions := []Instructions{
			Make(OpAdd),
			Make(OpConstant, 2),
			Make(OpConstant, 65535),
			Make(OpSub),
			Make(OpMul),
			Make(OpDiv),
			Make(OpTrue),
			Make(OpFalse),
			Make(OpEqual),
			Make(OpNotEqual),
			Make(OpGreaterThan),
			Make(OpMinus),
			Make(OpBang),
			Make(OpJumpNotTruthy, 65534),
			Make(OpJump, 65534),
			Make(OpNull),
			Make(OpGetGlobal, 65534),
			Make(OpSetGlobal, 65534),
			Make(OpArray, 65534),
			Make(OpHash, 65534),
			Make(OpIndex),
			Make(OpCall, 255),
			Make(OpReturnValue),
			Make(OpReturn),
			Make(OpGetLocal, 1),
			Make(OpSetLocal, 2),
			Make(OpGetBuiltin, 255),
			Make(OpClosure, 65534, 255),
			Make(OpAnd),
			Make(OpOr),
			Make(OpPower),
			Make(OpGreaterThanEqual),
			Make(OpIncrementGlobal, 65534),
			Make(OpIncrementLocal, 255),
			Make(OpDecrementGlobal, 65534),
			Make(OpDecrementLocal, 255),
			Make(OpSetLocalDefault, 255),
		}

		expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
0007 OpSub
0008 OpMul
0009 OpDiv
0010 OpTrue
0011 OpFalse
0012 OpEqual
0013 OpNotEqual
0014 OpGreaterThan
0015 OpMinus
0016 OpBang
0017 OpJumpNotTruthy 65534
0020 OpJump 65534
0023 OpNull
0024 OpGetGlobal 65534
0027 OpSetGlobal 65534
0030 OpArray 65534
0033 OpHash 65534
0036 OpIndex
0037 OpCall 255
0039 OpReturnValue
0040 OpReturn
0041 OpGetLocal 1
0043 OpSetLocal 2
0045 OpGetBuiltin 255
0047 OpClosure 65534 255
0051 OpAnd
0052 OpOr
0053 OpPower
0054 OpGreaterThanEqual
0055 OpIncrementGlobal 65534
0058 OpIncrementLocal 255
0060 OpDecrementGlobal 65534
0063 OpDecrementLocal 255
0065 OpSetLocalDefault 255
`

		concatted := Instructions{}

		for _, ins := range instructions {
			concatted = append(concatted, ins...)
		}

		if concatted.String() != expected {
			t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q", expected, concatted.String())
		}
	})
}

func TestReadOperands(t *testing.T) {
	t.Run("should read the correct operands", func(t *testing.T) {
		tests := []struct {
			op        Opcode
			operands  []int
			bytesRead int
		}{
			{OpConstant, []int{65534}, 2},
			{OpGetLocal, []int{255}, 1},
			{OpSetLocal, []int{255}, 1},
			{OpGetBuiltin, []int{255}, 1},
			{OpCall, []int{255}, 1},
			{OpJumpNotTruthy, []int{65534}, 2},
			{OpJump, []int{65534}, 2},
			{OpGetGlobal, []int{65534}, 2},
			{OpSetGlobal, []int{65534}, 2},
			{OpArray, []int{65534}, 2},
			{OpHash, []int{65534}, 2},
			{OpClosure, []int{65534, 255}, 3},
			{OpIncrementGlobal, []int{65534}, 2},
			{OpIncrementLocal, []int{255}, 1},
			{OpDecrementGlobal, []int{65534}, 2},
			{OpDecrementLocal, []int{255}, 1},
			{OpSetLocalDefault, []int{255}, 1},
		}

		for _, tt := range tests {
			instruction := Make(tt.op, tt.operands...)

			def, err := Lookup(byte(tt.op))
			if err != nil {
				t.Fatalf("definition not found: %q", err)
			}

			operandsRead, n := ReadOperands(def, instruction[1:])
			if n != tt.bytesRead {
				t.Fatalf("wrong number of bytes read. want=%d, got=%d", tt.bytesRead, n)
			}

			for i, want := range tt.operands {
				if operandsRead[i] != want {
					t.Errorf("wrong operand at position %d. want=%d, got=%d", i, want, operandsRead[i])
				}
			}
		}
	})
}
