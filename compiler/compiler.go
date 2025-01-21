package compiler

import (
	"fmt"
	"shark/ast"
	"shark/bytecode"
	"shark/code"
	"shark/exception"
	"shark/object"
	"shark/token"
	"shark/types"
	"sort"
)

type Compiler struct {
	lastCompiledType types.ISharkType
	symbolTable      *SymbolTable
	upToPos          *token.Position
	scopes           []CompilationScope
	constants        []object.Object
	scopeIndex       int
}

type EmittedInstruction struct {
	opcode   code.Opcode
	position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

func New(upToPos ...token.Position) *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	SymbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		SymbolTable.DefineBuiltin(i, v.Name, v.Builtin.FuncType)
	}

	var pos *token.Position

	if len(upToPos) == 0 {
		pos = nil
	} else {
		pos = &upToPos[0]
	}

	return &Compiler{
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		constants:   []object.Object{},
		symbolTable: SymbolTable,
		upToPos:     pos,
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object, upToPos ...token.Position) *Compiler {
	c := New(upToPos...)
	c.symbolTable = symbolTable
	c.constants = constants

	return c
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) Compile(node ast.Node) (*exception.SharkError, bool) {
	if c.upToPos != nil {
		nodePos := node.TokenPos()

		if nodePos.Line > c.upToPos.Line {
			return nil, true
		}
	}
	switch node := node.(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			if err, stopped := c.Compile(statement); err != nil || stopped {
				return err, stopped
			}
		}
	case *ast.ExpressionStatement:
		if err, stopped := c.Compile(node.Expression); err != nil || stopped {
			return err, stopped
		}
		c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpPop)
	case *ast.PostfixExpression:
		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			return newSharkError(
				exception.SharkErrorIdentifierExpected, node.Token.Literal,
				"Identifier expected for postfix expression",
				exception.NewSharkErrorCause(fmt.Sprintf("cannot use type %T", node.Left), node.Token.Pos),
			), false
		}
		symbol, ok := c.symbolTable.Resolve(ident.Value)
		if !ok {
			return newSharkError(
				exception.SharkErrorIdentifierNotFound, ident.Value,
				"Variable not found for postfix expression",
				exception.NewSharkErrorCause("Variable not found for postfix expression", node.Token.Pos),
			), false
		}
		if !symbol.VariadicType && !symbol.Mutable {
			return newSharkError(exception.SharkErrorImmutableValue, ident.Value,
				"Add the 'mut' keyword before the variable name to make it mutable, or use the 'var' keyword to declare the variable",
				exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Pos),
			), false
		}
		c.loadSymbol(symbol)
		switch node.Operator {
		case "++":
			if symbol.Scope == GlobalScope {
				c.emit(symbol.ObjType, code.OpIncrementGlobal, symbol.Index)
			} else {
				c.emit(symbol.ObjType, code.OpIncrementLocal, symbol.Index)
			}
		case "--":
			if symbol.Scope == GlobalScope {
				c.emit(symbol.ObjType, code.OpDecrementGlobal, symbol.Index)
			} else {
				c.emit(symbol.ObjType, code.OpDecrementLocal, symbol.Index)
			}
		}
	case *ast.InfixExpression:
		if node.Operator == "<" {
			if err, stopped := c.Compile(node.Right); err != nil || stopped {
				return err, stopped
			}
			if err, stopped := c.Compile(node.Left); err != nil || stopped {
				return err, stopped
			}
			if _, ok := c.lastCompiledType.(types.TSharkI64); !ok {
				return newSharkError(exception.SharkErrorTypeMismatch, c.lastCompiledType.SharkTypeString(),
					"Use a number for comparison",
					exception.NewSharkErrorCause(fmt.Sprintf("Cannot use type '%s' for left value comparison", c.lastCompiledType.SharkTypeString()), node.Token.Pos),
				), false
			}
			c.emit(types.TSharkBool{}, code.OpGreaterThan)
			return nil, false
		} else if node.Operator == "<=" {
			if err, stopped := c.Compile(node.Right); err != nil || stopped {
				return err, stopped
			}
			if err, stopped := c.Compile(node.Left); err != nil || stopped {
				return err, stopped
			}
			c.emit(types.TSharkBool{}, code.OpGreaterThanEqual)
			return nil, false
		}
		if node.Operator != "=" &&
			node.Operator != "+=" &&
			node.Operator != "-=" &&
			node.Operator != "*=" &&
			node.Operator != "/=" {
			if err, stopped := c.Compile(node.Left); err != nil || stopped {
				return err, stopped
			}
			if err, stopped := c.Compile(node.Right); err != nil || stopped {
				return err, stopped
			}
		}
		switch node.Operator {
		case "..":
			c.emit(types.TSharkArray{Collects: c.lastCompiledType}, code.OpRange)
		case "+":
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpAdd)
		case "-":
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpSub)
		case "*":
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpMul)
		case "**":
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpPower)
		case "/":
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpDiv)
		case "==":
			c.emit(types.TSharkBool{}, code.OpEqual)
		case "=", "+=", "-=", "*=", "/=":
			identLeft, ok := node.Left.(*ast.Identifier)
			if !ok {
				return newSharkError(
					exception.SharkErrorIdentifierExpected, node.Token.Literal,
					"Make sure to use a variable for reassignment",
					exception.NewSharkErrorCause("left value must be an identifier, but is not", node.Token.Pos),
				), false
			}
			symbolLeft, ok := c.symbolTable.Resolve(identLeft.Value)
			if !ok {
				return newSharkError(
					exception.SharkErrorIdentifierNotFound, identLeft.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found for reassignment", node.Token.Pos),
				), false
			}
			if !symbolLeft.VariadicType && !symbolLeft.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, identLeft.Value,
					"Add the 'mut' keyword before the variable name to make it mutable, or use the 'var' keyword to declare the variable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Pos),
				), false
			}

			identRight, ok := node.Right.(*ast.Identifier)
			if ok {
				symbolRight, ok := c.symbolTable.Resolve(identRight.Value)
				if !ok {
					return newSharkError(
						exception.SharkErrorIdentifierNotFound, identRight.Value,
						"Make sure the variable is defined before using it",
						exception.NewSharkErrorCause("Variable not found for reassignment", node.Token.Pos),
					), false
				}
				if !symbolLeft.ObjType.Is(symbolRight.ObjType) && !symbolLeft.VariadicType {
					return newSharkError(exception.SharkErrorTypeMismatch, symbolRight.ObjType.SharkTypeString(),
						"Declare the variable with 'var' keyword instead",
						exception.NewSharkErrorCause(fmt.Sprintf("Cannot assign type '%s' to type '%s'", symbolRight.ObjType.SharkTypeString(), symbolLeft.ObjType.SharkTypeString()), node.Token.Pos),
					), false
				}
			}

			index := symbolLeft.Index

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

			if symbolLeft.Scope == GlobalScope {
				localitySet = code.OpSetGlobal
				localityGet = code.OpGetGlobal
			} else {
				localitySet = code.OpSetLocal
				localityGet = code.OpGetLocal
			}

			if node.Operator != "=" {
				c.emit(symbolLeft.ObjType, localityGet, index)
				if err, stopped := c.Compile(node.Right); err != nil || stopped {
					return err, stopped
				}
				c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, op)
			} else {
				if err, stopped := c.Compile(node.Right); err != nil || stopped {
					return err, stopped
				}
			}
			if !symbolLeft.ObjType.Is(c.lastCompiledType) && !symbolLeft.VariadicType {
				return newSharkError(exception.SharkErrorTypeMismatch, c.lastCompiledType.SharkTypeString(),
					"Declare the variable with 'var' keyword instead",
					exception.NewSharkErrorCause(fmt.Sprintf("Cannot assign type '%s' to type '%s'", c.lastCompiledType.SharkTypeString(), symbolLeft.ObjType.SharkTypeString()), node.Token.Pos),
				), false
			}
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, localitySet, index)
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, localityGet, index)
		case "!=":
			c.emit(types.TSharkBool{}, code.OpNotEqual)
		case ">":
			c.emit(types.TSharkBool{}, code.OpGreaterThan)
		case ">=":
			c.emit(types.TSharkBool{}, code.OpGreaterThanEqual)
		case "&&":
			c.emit(types.TSharkBool{}, code.OpAnd)
		case "||":
			c.emit(types.TSharkBool{}, code.OpOr)
		default:
			return newSharkError(exception.SharkErrorUnknownOperator, node.Operator,
				"Try using an other operator, such as '&&' or '+'",
				exception.NewSharkErrorCause("Invalid operator for infix expression", node.Token.Pos),
			), false
		}
	case *ast.PrefixExpression:
		switch node.Operator {
		case "!":
			if err, stopped := c.Compile(node.Right); err != nil || stopped {
				return err, stopped
			}
			c.emit(types.TSharkBool{}, code.OpBang)
		case "-":
			if err, stopped := c.Compile(node.Right); err != nil || stopped {
				return err, stopped
			}
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpMinus)
		case "++":
			if node.RightIdent == nil {
				return newSharkError(exception.SharkErrorIdentifierExpected, node.Token.Literal,
					"Only variables can be used for '++' operator",
					exception.NewSharkErrorCause("Operator not followed by variable", node.Token.Pos),
				), false
			}
			symbol, ok := c.symbolTable.Resolve(node.RightIdent.Value)
			if !ok {
				return newSharkError(exception.SharkErrorIdentifierNotFound, node.RightIdent.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found", node.Token.Pos),
				), false
			}
			if !symbol.VariadicType && !symbol.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, node.RightIdent.Value,
					"Add the 'mut' keyword before the variable name to make it mutable, or use the 'var' keyword to declare the variable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Pos),
				), false
			}
			if symbol.Scope == GlobalScope {
				c.emit(symbol.ObjType, code.OpIncrementGlobal, symbol.Index)
			} else {
				c.emit(symbol.ObjType, code.OpIncrementLocal, symbol.Index)
			}
			c.loadSymbol(symbol)
		case "--":
			if node.RightIdent == nil {
				return newSharkError(exception.SharkErrorIdentifierExpected, node.Token.Literal,
					"Only variables can be used for '--' operator",
					exception.NewSharkErrorCause("Operator noy followed by variable", node.Token.Pos),
				), false
			}
			symbol, ok := c.symbolTable.Resolve(node.RightIdent.Value)
			if !ok {
				return newSharkError(exception.SharkErrorIdentifierNotFound, node.RightIdent.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found", node.Token.Pos),
				), false
			}
			if !symbol.VariadicType && !symbol.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, node.RightIdent.Value,
					"Add the 'mut' keyword before the variable name to make it mutable, or use the 'var' keyword to declare the variable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Pos),
				), false
			}
			if symbol.Scope == GlobalScope {
				c.emit(symbol.ObjType, code.OpDecrementGlobal, symbol.Index)
			} else {
				c.emit(symbol.ObjType, code.OpDecrementLocal, symbol.Index)
			}
			c.loadSymbol(symbol)
		case "...":
			if err, stopped := c.Compile(node.Right); err != nil || stopped {
				return err, stopped
			}
			c.emit(types.TSharkArray{Collects: c.lastCompiledType}, code.OpSpread)
		default:
			return newSharkError(exception.SharkErrorUnknownOperator, node.Operator,
				"Try using an other operator",
				exception.NewSharkErrorCause("Invalid operator for prefix expression", node.Token.Pos),
			), false
		}
	case *ast.IntegerLiteral:
		integer := &object.Int64{Value: node.Value}
		c.emit(types.TSharkI64{}, code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(types.TSharkBool{}, code.OpTrue)
		} else {
			c.emit(types.TSharkBool{}, code.OpFalse)
		}
	case *ast.IfExpression:
		//FIXME: No new scope is created for block statements inside "if" and "else"
		if err, stopped := c.Compile(node.Condition); err != nil || stopped {
			return err, stopped
		}
		jumpNotTruthyPos := c.emit(types.TSharkAny{}, code.OpJumpNotTruthy, 9999)

		if err, stopped := c.Compile(node.Consequence); err != nil || stopped {
			return err, stopped
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		if c.lastInstructionIs(code.OpSetGlobal) || c.lastInstructionIs(code.OpSetLocal) || len(node.Consequence.Statements) == 0 {
			c.emit(types.TSharkNull{}, code.OpNull)
		}

		jumpPos := c.emit(types.TSharkAny{}, code.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil || len(node.Alternative.Statements) == 0 {
			c.emit(types.TSharkNull{}, code.OpNull)
		} else {
			if err, stopped := c.Compile(node.Alternative); err != nil || stopped {
				return err, stopped
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
			if c.lastInstructionIs(code.OpSetGlobal) || c.lastInstructionIs(code.OpSetLocal) {
				c.emit(types.TSharkNull{}, code.OpNull)
			}
		}
		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.WhileStatement:
		//FIXME: No new scope is created for block statements inside "while"
		conditionPos := len(c.currentInstructions())
		lastCompiledType := c.lastCompiledType
		if err, stopped := c.Compile(node.Condition); err != nil || stopped {
			return err, stopped
		}
		jumpNotTruthyPos := c.emit(types.TSharkVariadic{Enclosed: lastCompiledType}, code.OpJumpNotTruthy, 9999)

		if err, stopped := c.Compile(node.Body); err != nil || stopped {
			return err, stopped
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		c.emit(nil, code.OpJump, conditionPos)

		afterBodyPos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterBodyPos)
	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			if err, stopped := c.Compile(statement); err != nil || stopped {
				return err, stopped
			}
		}
	case *ast.LetStatement:
		symbol, ok := c.symbolTable.Resolve(node.Name.Value)
		if ok {
			return newSharkError(exception.SharkErrorDuplicateIdentifier, node.Name.Value,
				"Remove 'let' before the variable name",
				exception.NewSharkErrorCause("Cannot use let to reassign value to an existing variable", node.Token.Pos),
			), false
		}

		if err, stopped := c.Compile(node.Value); err != nil || stopped {
			return err, stopped
		}
		symbol = c.symbolTable.Define(node.Name.Value, node.Name.Mutable, node.Name.VariadicType, c.lastCompiledType, &node.Name.Token.Pos)
		if symbol.Scope == GlobalScope {
			c.emit(symbol.ObjType, code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(symbol.ObjType, code.OpSetLocal, symbol.Index)
		}
	case *ast.TupleDeconstruction:
		// check if the right value is an identifier tuple
		rightIdent, ok := node.Value.(*ast.Identifier)
		var tupleType types.TSharkTuple
		if ok {
			// check if ident is a tuple
			symbol, ok := c.symbolTable.Resolve(rightIdent.Value)
			if !ok {
				return newSharkError(exception.SharkErrorIdentifierNotFound, rightIdent.Value,
					"Make sure the variable is defined before using it with the 'let' keyword",
					exception.NewSharkErrorCause("Variable not found for tuple deconstruction", node.Token.Pos),
				), false
			}
			if !symbol.ObjType.Is(types.TSharkTuple{}) {
				return newSharkError(exception.SharkErrorTypeMismatch, symbol.ObjType.SharkTypeString(),
					"Use a tuple for tuple deconstruction",
					exception.NewSharkErrorCause(fmt.Sprintf("Cannot deconstruct type '%s'", symbol.ObjType.SharkTypeString()), node.Token.Pos),
				), false
			}
			tupleType = symbol.ObjType.(types.TSharkTuple)
			c.loadSymbol(symbol)
		} else {
			if err, stopped := c.Compile(node.Value); err != nil || stopped {
				return err, stopped
			}
			if !c.lastCompiledType.Is(types.TSharkTuple{}) {
				return newSharkError(exception.SharkErrorTypeMismatch, c.lastCompiledType.SharkTypeString(),
					"Use a tuple for tuple deconstruction",
					exception.NewSharkErrorCause(fmt.Sprintf("Cannot deconstruct type '%s'", c.lastCompiledType.SharkTypeString()), node.Token.Pos),
				), false
			}
			tupleType = c.lastCompiledType.(types.TSharkTuple)
		}
		c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpTupleDeconstruct, len(node.Names))
		for i, name := range node.Names {
			symbol, ok := c.symbolTable.Resolve(name.Value)
			if ok {
				return newSharkError(exception.SharkErrorDuplicateIdentifier, name.Value,
					"Remove 'let' before the variable name",
					exception.NewSharkErrorCause("Cannot use let to reassign value to an existing variable", node.Token.Pos),
				), false
			}
			symbol = c.symbolTable.Define(name.Value, name.Mutable, name.VariadicType, tupleType.Collects[i], &name.Token.Pos)
			if symbol.Scope == GlobalScope {
				c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpSetGlobal, symbol.Index)
			} else {
				c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpSetLocal, symbol.Index)
			}
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return newSharkError(exception.SharkErrorIdentifierNotFound, node.Value,
				fmt.Sprintf("You must define '%s' before using it with the 'let' keyword", node.Value),
				exception.NewSharkErrorCause(fmt.Sprintf("identifier '%s' is not defined", node.Value), node.Token.Pos),
			), false
		}
		c.loadSymbol(symbol)
		c.lastCompiledType = symbol.ObjType
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(types.TSharkString{}, code.OpConstant, c.addConstant(str))
	case *ast.ArrayLiteral:
		var elementType types.ISharkType
		for _, element := range node.Elements {
			if err, stopped := c.Compile(element); err != nil || stopped {
				return err, stopped
			}
			if elementType == nil {
				elementType = c.lastCompiledType
			} else if !elementType.Is(c.lastCompiledType) && !elementType.Is(types.TSharkAny{}) {
				elementType = types.TSharkAny{}
			} else {
				elementType = c.lastCompiledType
			}

		}
		c.emit(elementType, code.OpArray, len(node.Elements))
	case *ast.TupleLiteral:
		var elementTypes []types.ISharkType
		for _, element := range node.Elements {
			if err, stopped := c.Compile(element); err != nil || stopped {
				return err, stopped
			}
			elementTypes = append(elementTypes, c.lastCompiledType)
		}
		c.emit(types.TSharkTuple{Collects: elementTypes}, code.OpTuple, len(node.Elements))
	case *ast.HashLiteral:
		var keyType, valueType types.ISharkType
		var keys []ast.Expression
		for key := range node.Pairs {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, key := range keys {
			if err, stopped := c.Compile(key); err != nil || stopped {
				return err, stopped
			}
			if keyType == nil {
				keyType = c.lastCompiledType
			} else if !keyType.Is(c.lastCompiledType) && !keyType.Is(types.TSharkAny{}) {
				keyType = types.TSharkAny{}
			} else {
				keyType = c.lastCompiledType
			}
			if err, stopped := c.Compile(node.Pairs[key]); err != nil || stopped {
				return err, stopped
			}
			if valueType == nil {
				valueType = c.lastCompiledType
			} else if !valueType.Is(c.lastCompiledType) && !valueType.Is(types.TSharkAny{}) {
				valueType = types.TSharkAny{}
			} else {
				valueType = c.lastCompiledType
			}
		}

		c.emit(types.TSharkHashMap{Indexes: keyType, Collects: valueType}, code.OpHash, len(node.Pairs)*2)
	case *ast.IndexExpression:
		if err, stopped := c.Compile(node.Left); err != nil || stopped {
			return err, stopped
		}

		indexValueType := c.lastCompiledType

		if err, stopped := c.Compile(node.Index); err != nil || stopped {
			return err, stopped
		}

		c.emit(indexValueType, code.OpIndex)
	case *ast.IndexAssignExpression:
		if err, stopped := c.Compile(node.Value); err != nil || stopped {
			return err, stopped
		}

		if err, stopped := c.Compile(node.Left); err != nil || stopped {
			return err, stopped
		}

		// if value is an identifier, check if it is mutable
		ident, ok := node.Left.(*ast.Identifier)
		if ok {
			symbol, ok := c.symbolTable.Resolve(ident.Value)
			if !ok {
				return newSharkError(exception.SharkErrorIdentifierNotFound, ident.Value,
					"Make sure the variable is defined before using it",
					exception.NewSharkErrorCause("Variable not found for index assignment", node.Token.Pos),
				), false
			}
			if !symbol.VariadicType && !symbol.Mutable {
				return newSharkError(exception.SharkErrorImmutableValue, ident.Value,
					"Add the 'mut' keyword before the variable name to make it mutable, or use the 'var' keyword to declare the variable",
					exception.NewSharkErrorCause("Cannot reassign value to a constant", node.Token.Pos),
				), false
			}
		}

		if err, stopped := c.Compile(node.Index); err != nil || stopped {
			return err, stopped
		}
		c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpIndexAssign)
	case *ast.FunctionLiteral:
		c.enterScope()
		if node.Name != "" {
			c.symbolTable.DefineFunctionName(node.Name, &node.Token.Pos)
		}
		numDefaults := 0
		isOptionalsActive := false
		var paramTypes []types.ISharkType
		for _, param := range node.Parameters {
			if param.DefaultValue != nil {
				isOptionalsActive = true
				numDefaults++
				if err, stopped := c.Compile(*param.DefaultValue); err != nil || stopped {
					return err, stopped
				}
				symbol := c.symbolTable.Define(param.Value, param.Mutable, param.VariadicType, c.lastCompiledType, &param.Token.Pos)
				c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpSetLocalDefault, symbol.Index)
			} else {
				c.symbolTable.Define(param.Value, param.Mutable, param.VariadicType, types.TSharkAny{}, &param.Token.Pos)
			}
			if isOptionalsActive && param.DefaultValue == nil {
				return newSharkError(exception.SharkErrorOptionalParameter, param.Value,
					"Move this parameter before the optional parameters",
					exception.NewSharkErrorCause("Non-optional parameter after optional parameter", param.Token.Pos),
				), false
			}
			if c.lastCompiledType != nil && !c.lastCompiledType.Is(types.TSharkAny{}) {
				paramTypes = append(paramTypes, c.lastCompiledType)
			} else {
				paramTypes = append(paramTypes, types.TSharkAny{})
			}
		}
		if err, stopped := c.Compile(node.Body); err != nil || stopped {
			return err, stopped
		}
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpReturn)
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
			ObjType:       types.TSharkFuncType{ArgsList: paramTypes, ReturnT: types.TSharkAny{}},
		}
		c.emit(types.TSharkFuncType{ArgsList: paramTypes, ReturnT: types.TSharkAny{}}, code.OpClosure, c.addConstant(compiledFn), len(freeSymbols))
	case *ast.ReturnStatement:
		if c.scopeIndex == 0 {
			return newSharkError(exception.SharkErrorTopLeverReturn, nil,
				"Use 'exit(0);' instead",
				exception.NewSharkErrorCause("Unexpected return statement in main scope", node.Token.Pos),
			), false
		}
		if err, stopped := c.Compile(node.ReturnValue); err != nil || stopped {
			return err, stopped
		}

		c.emit(types.TSharkVariadic{Enclosed: c.lastCompiledType}, code.OpReturnValue)
	case *ast.CallExpression:
		if err, stopped := c.Compile(node.Function); err != nil || stopped {
			return err, stopped
		}
		funcType := c.lastCompiledType
		for _, arg := range node.Arguments {
			if err, stopped := c.Compile(arg); err != nil || stopped {
				return err, stopped
			}
		}
		c.emit(funcType, code.OpCall, len(node.Arguments))
	}

	return nil, false
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.position

	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.opcode = code.OpReturnValue
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
	currentSymbolTable := *c.symbolTable
	currentSymbolTable.Outer = nil
	c.symbolTable = c.symbolTable.Outer
	c.symbolTable.Inner = &currentSymbolTable

	return instructions
}

func (c *Compiler) GetSymbolTable() *SymbolTable {
	return c.symbolTable
}

func (c *Compiler) addConstant(obj object.Object) int {
	// TODO: Add detection for duplicate constants for functions and closures
	if !obj.Type().Is(types.TSharkFuncType{}) && !obj.Type().Is(types.TSharkClosure{}) {
		for i, constant := range c.constants {

			if constant.Type().Is(obj.Type()) && constant.Inspect() == obj.Inspect() {
				return i
			}
		}
	}
	c.constants = append(c.constants, obj)

	return len(c.constants) - 1
}

func (c *Compiler) emit(sharkType types.ISharkType, op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)
	if variadic, ok := sharkType.(types.TSharkVariadic); ok {
		sharkType = variadic.Enclosed
	}
	c.lastCompiledType = sharkType

	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{opcode: op, position: pos}
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

	return c.scopes[c.scopeIndex].lastInstruction.opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction
	oldInstructions := c.currentInstructions()
	newInstructions := oldInstructions[:last.position]
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
		c.emit(s.ObjType, code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(s.ObjType, code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(s.ObjType, code.OpGetBuiltin, s.Index)
	case FreeScope:
		c.emit(s.ObjType, code.OpGetFree, s.Index)
	case FunctionScope:
		c.emit(s.ObjType, code.OpCurrentClosure)
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
