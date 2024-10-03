package compiler

import (
	"fmt"
	"shark/ast"
	"shark/bytecode"
	"shark/code"
	"shark/exception"
	"shark/object"
	"sort"
)

type Compiler struct {
	scopes      []CompilationScope
	scopeIndex  int
	constants   []object.Object
	symbolTable *SymbolTable
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	SymbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		SymbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		constants:   []object.Object{},
		symbolTable: SymbolTable,
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	c := New()
	c.symbolTable = symbolTable
	c.constants = constants

	return c
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) Compile(node ast.Node) *exception.SharkError {
	switch node := node.(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			if err := c.Compile(statement); err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.PostfixExpression:
		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			return newSharkError(
				exception.SharkErrorIdentifierExpected, node.Token.Literal,
				"Identifier expected for postfix expression",
				exception.NewSharkErrorCause(fmt.Sprintf("cannot use type %T", node.Left), node.Token.Line, node.Token.Line, node.Token.ColFrom, node.Token.ColTo),
			)
		}
		symbol, ok := c.symbolTable.Resolve(ident.Value)
		if !ok {
			return newSharkError(
				exception.SharkErrorIdentifierNotFound, ident.Value,
				"Variable not found for postfix expression",
				exception.NewSharkErrorCause("Variable not found for postfix expression", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
			)
		}
		if !symbol.Mutable {
			return newSharkError(exception.SharkErrorImmutableValue, ident.Value,
				"Add the 'mut' keyword before the variable name to make it mutable",
				exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
			)
		}
		c.loadSymbol(symbol)
		switch node.Operator {
		case "++":
			if symbol.Scope == GlobalScope {
				c.emit(code.OpIncrementGlobal, symbol.Index)
			} else {
				c.emit(code.OpIncrementLocal, symbol.Index)
			}
		case "--":
			if symbol.Scope == GlobalScope {
				c.emit(code.OpDecrementGlobal, symbol.Index)
			} else {
				c.emit(code.OpDecrementLocal, symbol.Index)
			}
		}
	case *ast.InfixExpression:
		if node.Operator == "<" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		} else if node.Operator == "<=" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			c.emit(code.OpGreaterThanEqual)
			return nil
		}
		if node.Operator != "=" &&
			node.Operator != "+=" &&
			node.Operator != "-=" &&
			node.Operator != "*=" &&
			node.Operator != "/=" {
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			if err := c.Compile(node.Right); err != nil {
				return err
			}
		}
		switch node.Operator {
		case "..":
			c.emit(code.OpRange)
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "**":
			c.emit(code.OpPower)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEqual)
		case "=", "+=", "-=", "*=", "/=":
			ident, ok := node.Left.(*ast.Identifier)
			if !ok {
				return newSharkError(
					exception.SharkErrorIdentifierExpected, node.Token.Literal,
					"Make sure to use a variable for reassignment",
					exception.NewSharkErrorCause("left value must be an identifier, but is not", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			symbol, ok := c.symbolTable.Resolve(ident.Value)
			if !ok {
				return newSharkError(
					exception.SharkErrorIdentifierNotFound, ident.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found for reassignment", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			if !symbol.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, ident.Value,
					"Add the 'mut' keyword before the variable name to make it mutable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			index := symbol.Index

			var op code.Opcode

			switch node.Operator {
			case "+=":
				op = code.OpAdd
			case "-=":
				op = code.OpSub
			case "*=":
				op = code.OpMul
			case "/=":
				op = code.OpDiv
			}

			var localitySet, localityGet code.Opcode

			if symbol.Scope == GlobalScope {
				localitySet = code.OpSetGlobal
				localityGet = code.OpGetGlobal
			} else {
				localitySet = code.OpSetLocal
				localityGet = code.OpGetLocal
			}

			if node.Operator != "=" {
				c.emit(localityGet, index)
				if err := c.Compile(node.Right); err != nil {
					return err
				}
				c.emit(op)
			} else {
				if err := c.Compile(node.Right); err != nil {
					return err
				}
			}
			c.emit(localitySet, index)
			c.emit(localityGet, index)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">":
			c.emit(code.OpGreaterThan)
		case ">=":
			c.emit(code.OpGreaterThanEqual)
		case "&&":
			c.emit(code.OpAnd)
		case "||":
			c.emit(code.OpOr)
		default:
			return newSharkError(exception.SharkErrorUnknownOperator, node.Operator,
				"Try using an other operator, such as '&&' or '+'",
				exception.NewSharkErrorCause("Invalid operator for infix expression", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
			)
		}
	case *ast.PrefixExpression:
		switch node.Operator {
		case "!":
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			c.emit(code.OpBang)
		case "-":
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			c.emit(code.OpMinus)
		case "++":
			if node.RightIdent == nil {
				return newSharkError(exception.SharkErrorIdentifierExpected, node.Token.Literal,
					"Only variables can be used for '++' operator",
					exception.NewSharkErrorCause("Operator not followed by variable", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			symbol, ok := c.symbolTable.Resolve(node.RightIdent.Value)
			if !ok {
				return newSharkError(exception.SharkErrorIdentifierNotFound, node.RightIdent.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			if !symbol.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, node.RightIdent.Value,
					"Add the 'mut' keyword before the variable name to make it mutable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			if symbol.Scope == GlobalScope {
				c.emit(code.OpIncrementGlobal, symbol.Index)
			} else {
				c.emit(code.OpIncrementLocal, symbol.Index)
			}
			c.loadSymbol(symbol)
		case "--":
			if node.RightIdent == nil {
				return newSharkError(exception.SharkErrorIdentifierExpected, node.Token.Literal,
					"Only variables can be used for '--' operator",
					exception.NewSharkErrorCause("Operator noy followed by variable", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			symbol, ok := c.symbolTable.Resolve(node.RightIdent.Value)
			if !ok {
				return newSharkError(exception.SharkErrorIdentifierNotFound, node.RightIdent.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			if !symbol.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, node.RightIdent.Value,
					"Add the 'mut' keyword before the variable name to make it mutable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			if symbol.Scope == GlobalScope {
				c.emit(code.OpDecrementGlobal, symbol.Index)
			} else {
				c.emit(code.OpDecrementLocal, symbol.Index)
			}
			c.loadSymbol(symbol)
		case "...":
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			c.emit(code.OpSpread)
		default:
			return newSharkError(exception.SharkErrorUnknownOperator, node.Operator,
				"Try using an other operator",
				exception.NewSharkErrorCause("Invalid operator for prefix expression", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
			)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.IfExpression:
		//FIXME: No new scope is created for block statements inside "if" and "else"
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		if err := c.Compile(node.Consequence); err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		if c.lastInstructionIs(code.OpSetGlobal) || c.lastInstructionIs(code.OpSetLocal) || len(node.Consequence.Statements) == 0 {
			c.emit(code.OpNull)
		}

		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil || len(node.Alternative.Statements) == 0 {
			c.emit(code.OpNull)
		} else {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
			if c.lastInstructionIs(code.OpSetGlobal) || c.lastInstructionIs(code.OpSetLocal) {
				c.emit(code.OpNull)
			}
		}
		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.WhileStatement:
		//FIXME: No new scope is created for block statements inside "while"
		conditionPos := len(c.currentInstructions())
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		c.emit(code.OpJump, conditionPos)

		afterBodyPos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterBodyPos)
	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			if err := c.Compile(statement); err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		symbol, ok := c.symbolTable.Resolve(node.Name.Value)
		if ok {
			return newSharkError(exception.SharkErrorDuplicateIdentifier, node.Name.Value,
				"Remove 'let' before the variable name",
				exception.NewSharkErrorCause("Cannot use let to reassign value to an existing variable", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
			)
		}
		symbol = c.symbolTable.Define(node.Name.Value, node.Name.Mutable)
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return newSharkError(exception.SharkErrorIdentifierNotFound, node.Value,
				fmt.Sprintf("You must define '%s' before using it with the 'let' keyword", node.Value),
				exception.NewSharkErrorCause(fmt.Sprintf("identifier '%s' is not defined", node.Value), node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
			)
		}
		c.loadSymbol(symbol)
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))
	case *ast.ArrayLiteral:
		for _, element := range node.Elements {
			if err := c.Compile(element); err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.HashLiteral:
		var keys []ast.Expression
		for key := range node.Pairs {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, key := range keys {
			if err := c.Compile(key); err != nil {
				return err
			}
			if err := c.Compile(node.Pairs[key]); err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.IndexExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}

		if err := c.Compile(node.Index); err != nil {
			return err
		}

		c.emit(code.OpIndex)
	case *ast.IndexAssignExpression:
		if err := c.Compile(node.Value); err != nil {
			return err
		}

		if err := c.Compile(node.Left); err != nil {
			return err
		}

		// if value is an identifier, check if it is mutable
		ident, ok := node.Left.(*ast.Identifier)
		if ok {
			symbol, ok := c.symbolTable.Resolve(ident.Value)
			if !ok {
				return newSharkError(exception.SharkErrorIdentifierNotFound, ident.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found for index assignment", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
			if !symbol.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, ident.Value,
					"Add the 'mut' keyword before the variable name to make it mutable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Line, node.Token.LineTo, node.Token.ColFrom, node.Token.ColTo),
				)
			}
		}

		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(code.OpIndexAssign)
	case *ast.FunctionLiteral:
		c.enterScope()
		if node.Name != "" {
			c.symbolTable.DefineFunctionName(node.Name)
		}
		numDefaults := 0
		for _, param := range node.Parameters {
			symbol := c.symbolTable.Define(param.Value, param.Mutable)
			if param.DefaultValue != nil {
				numDefaults++
				if err := c.Compile(*param.DefaultValue); err != nil {
					return err
				}
				c.emit(code.OpSetLocalDefault, symbol.Index)
			}
		}
		if err := c.Compile(node.Body); err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}
		freeSymbols := c.symbolTable.FreeSymbols
		NumLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     NumLocals,
			NumParameters: len(node.Parameters),
			NumDefaults:   numDefaults,
		}

		c.emit(code.OpClosure, c.addConstant(compiledFn), len(freeSymbols))
	case *ast.ReturnStatement:
		if c.scopeIndex == 0 {
			return newSharkError(exception.SharkErrorTopLeverReturn, nil,
				"Use 'exit(0);' instead",
				exception.NewSharkErrorCause("Unexpected return statement in main scope", node.Token.Line, node.Token.Line, node.Token.ColFrom, node.Token.ColTo),
			)
		}
		if err := c.Compile(node.ReturnValue); err != nil {
			return err
		}

		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil {
			return err
		}
		for _, arg := range node.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Arguments))
	}

	return nil
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position

	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) Bytecode() *bytecode.Bytecode {
	return &bytecode.Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) addConstant(obj object.Object) int {
	// TODO: Add detection for duplicate constants for functions and closures
	if obj.Type() != object.FUNCTION_OBJ && obj.Type() != object.CLOSURE_OBJ && obj.Type() != object.COMPILED_FUNCTION_OBJ {
		for i, constant := range c.constants {

			if constant.Type() == obj.Type() && constant.Inspect() == obj.Inspect() {
				return i
			}
		}
	}
	c.constants = append(c.constants, obj)

	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)
	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction
	oldInstructions := c.currentInstructions()
	newInstructions := oldInstructions[:last.Position]
	c.scopes[c.scopeIndex].instructions = newInstructions
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	case FreeScope:
		c.emit(code.OpGetFree, s.Index)
	case FunctionScope:
		c.emit(code.OpCurrentClosure)
	}
}

func newSharkError(code exception.SharkErrorCode, param interface{}, helpMsg string, cause ...exception.SharkErrorCause) *exception.SharkError {
	var err exception.SharkError
	if param == nil {
		err = *exception.NewSharkError(exception.SharkErrorTypeCompiler, code)
	} else {
		err = *exception.NewSharkError(exception.SharkErrorTypeCompiler, code, param)
	}

	if helpMsg != "" {
		err.SetHelpMsg(helpMsg)
	}

	for _, c := range cause {
		err.AddCause(c)
	}

	return &err
}
