package compiler

import (
	"shark/token"
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
	Name    string
	Mutable bool
	Scope   SymbolScope
	Index   int
	Pos     *token.Position
}

type SymbolTable struct {
	Outer          *SymbolTable
	FreeSymbols    []Symbol
	store          map[string]Symbol
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

func (s *SymbolTable) Define(name string, mutable bool, pos *token.Position) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Mutable: mutable, Pos: pos}

	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

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

		return s.DefineFree(obj, obj.Mutable, obj.Pos), true
	}

	return obj, ok
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Scope: BuiltinScope, Index: index, Mutable: false}

	s.store[name] = symbol

	return symbol
}

func (s *SymbolTable) DefineFree(original Symbol, mutable bool, pos *token.Position) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Scope: FreeScope, Index: len(s.FreeSymbols) - 1, Mutable: mutable, Pos: pos}

	s.store[original.Name] = symbol

	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string, pos *token.Position) Symbol {
	symbol := Symbol{Name: name, Scope: FunctionScope, Index: 0, Pos: pos}

	s.store[name] = symbol

	return symbol
}
