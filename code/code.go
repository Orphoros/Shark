package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpGetLocal
	OpSetLocal
	OpGetBuiltin
	OpClosure
	OpGetFree
	OpCurrentClosure
	OpAnd
	OpOr
	OpPower
	OpGreaterThanEqual
	OpIncrementGlobal
	OpIncrementLocal
	OpDecrementGlobal
	OpDecrementLocal
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:         {"OpConstant", []int{2}},
	OpAdd:              {"OpAdd", []int{}},
	OpPop:              {"OpPop", []int{}},
	OpSub:              {"OpSub", []int{}},
	OpMul:              {"OpMul", []int{}},
	OpDiv:              {"OpDiv", []int{}},
	OpTrue:             {"OpTrue", []int{}},
	OpFalse:            {"OpFalse", []int{}},
	OpEqual:            {"OpEqual", []int{}},
	OpNotEqual:         {"OpNotEqual", []int{}},
	OpGreaterThan:      {"OpGreaterThan", []int{}},
	OpGreaterThanEqual: {"OpGreaterThanEqual", []int{}},
	OpMinus:            {"OpMinus", []int{}},
	OpBang:             {"OpBang", []int{}},
	OpJumpNotTruthy:    {"OpJumpNotTruthy", []int{2}},
	OpJump:             {"OpJump", []int{2}},
	OpNull:             {"OpNull", []int{}},
	OpGetGlobal:        {"OpGetGlobal", []int{2}},
	OpSetGlobal:        {"OpSetGlobal", []int{2}},
	OpArray:            {"OpArray", []int{2}},
	OpHash:             {"OpHash", []int{2}},
	OpIndex:            {"OpIndex", []int{}},
	OpCall:             {"OpCall", []int{1}},
	OpReturnValue:      {"OpReturnValue", []int{}},
	OpReturn:           {"OpReturn", []int{}},
	OpGetLocal:         {"OpGetLocal", []int{1}},
	OpSetLocal:         {"OpSetLocal", []int{1}},
	OpGetBuiltin:       {"OpGetBuiltin", []int{1}},
	OpClosure:          {"OpClosure", []int{2, 1}},
	OpGetFree:          {"OpGetFree", []int{1}},
	OpCurrentClosure:   {"OpCurrentClosure", []int{}},
	OpAnd:              {"OpAnd", []int{}},
	OpOr:               {"OpOr", []int{}},
	OpPower:            {"OpPower", []int{}},
	OpIncrementGlobal:  {"OpIncrementGlobal", []int{2}},
	OpIncrementLocal:   {"OpIncrementLocal", []int{1}},
	OpDecrementGlobal:  {"OpDecrementGlobal", []int{2}},
	OpDecrementLocal:   {"OpDecrementLocal", []int{1}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	instructionLength := 1

	for _, width := range def.OperandWidths {
		instructionLength += width
	}

	instruction := make([]byte, instructionLength)

	instruction[0] = byte(op)

	offset := 1

	for i, operand := range operands {
		width := def.OperandWidths[i]

		switch width {
		case 1:
			instruction[offset] = byte(operand)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}

		offset += width
	}

	return instruction
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0

	for i < len(ins) {
		def, err := Lookup(ins[i])

		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.formatInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) formatInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s", def.Name)
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))

	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }
