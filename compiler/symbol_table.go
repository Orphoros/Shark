package compiler

import (
	"shark/token"
	"shark/types"

	"github.com/phuslu/log"
)

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	ObjType      types.ISharkType
	Pos          *token.Position
	Name         string
	Scope        SymbolScope
	Index        int
	Mutable      bool
	VariadicType bool
}

type SymbolTable struct {
	Outer          *SymbolTable
	Inner          *SymbolTable
	store          map[string]Symbol
	FreeSymbols    []Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	free := []Symbol{}

	return &SymbolTable{store: s, FreeSymbols: free}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	st := NewSymbolTable()
	st.Outer = outer

	return st
}

func (s *SymbolTable) Define(name string, mutable, variadicType bool, objType types.ISharkType, pos *token.Position) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Mutable: mutable, Pos: pos, VariadicType: variadicType, ObjType: objType}

	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	log.Trace().
		Str("name", name).
		Bool("mutable", mutable).
		Bool("variadicType", variadicType).
		Str("objType", objType.SharkTypeString()).
		Str("scope", string(symbol.Scope)).
		Int("index", s.numDefinitions).Msg("Define new symbol")

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)

		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		return s.DefineFree(obj, obj.Mutable, obj.VariadicType, obj.ObjType, obj.Pos), true
	}

	return obj, ok
}

func (s *SymbolTable) FindIdent(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	if ok {
		return obj, ok
	}

	if s.Inner != nil {
		return s.Inner.FindIdent(name)
	}

	if s.Outer != nil {
		return s.Outer.FindIdent(name)
	}

	return obj, ok
}

func (s *SymbolTable) DefineBuiltin(index int, name string, objType types.ISharkType) Symbol {
	symbol := Symbol{Name: name, Scope: BuiltinScope, Index: index, Mutable: false, VariadicType: false, ObjType: objType}

	s.store[name] = symbol

	return symbol
}

func (s *SymbolTable) DefineFree(original Symbol, mutable, variadicType bool, objType types.ISharkType, pos *token.Position) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Scope: FreeScope, Index: len(s.FreeSymbols) - 1, Mutable: mutable, Pos: pos, VariadicType: variadicType, ObjType: objType}

	s.store[original.Name] = symbol

	log.Trace().
		Str("name", original.Name).
		Bool("mutable", mutable).
		Bool("variadicType", variadicType).
		Str("objType", objType.SharkTypeString()).
		Str("scope", string(symbol.Scope)).
		Int("index", symbol.Index).Msg("Define free symbol")

	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string, objType types.ISharkType, pos *token.Position) Symbol {
	symbol := Symbol{Name: name, Scope: FunctionScope, Index: 0, Pos: pos, ObjType: objType}

	s.store[name] = symbol

	log.Trace().
		Str("name", name).
		Str("scope", string(symbol.Scope)).
		Int("index", symbol.Index).Msg("Define function name")

	return symbol
}
