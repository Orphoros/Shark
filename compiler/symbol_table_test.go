package compiler

import (
	"testing"
)

func TestDefine(t *testing.T) {
	t.Run("should define a new symbol", func(t *testing.T) {
		expected := map[string]Symbol{
			"a": {Name: "a", Scope: GlobalScope, Index: 0},
			"b": {Name: "b", Scope: GlobalScope, Index: 1},
			"c": {Name: "c", Scope: LocalScope, Index: 0},
			"d": {Name: "d", Scope: LocalScope, Index: 1},
			"e": {Name: "e", Scope: LocalScope, Index: 0},
			"f": {Name: "f", Scope: LocalScope, Index: 1},
		}

		global := NewSymbolTable()

		if a := global.Define("a"); a != expected["a"] {
			t.Errorf("expected %s=%+v, got=%+v", "a", expected["a"], a)
		}

		if b := global.Define("b"); b != expected["b"] {
			t.Errorf("expected %s=%+v, got=%+v", "b", expected["b"], b)
		}

		firstLocal := NewEnclosedSymbolTable(global)

		if c := firstLocal.Define("c"); c != expected["c"] {
			t.Errorf("expected %s=%+v, got=%+v", "c", expected["c"], c)
		}

		if d := firstLocal.Define("d"); d != expected["d"] {
			t.Errorf("expected %s=%+v, got=%+v", "d", expected["d"], d)
		}

		secondLocal := NewEnclosedSymbolTable(firstLocal)

		if e := secondLocal.Define("e"); e != expected["e"] {
			t.Errorf("expected %s=%+v, got=%+v", "e", expected["e"], e)
		}

		if f := secondLocal.Define("f"); f != expected["f"] {
			t.Errorf("expected %s=%+v, got=%+v", "f", expected["f"], f)
		}
	})
}

func TestResolveGlobal(t *testing.T) {
	t.Run("should resolve a global symbol", func(t *testing.T) {
		global := NewSymbolTable()

		global.Define("a")
		global.Define("b")

		expected := []Symbol{
			{Name: "a", Scope: GlobalScope, Index: 0},
			{Name: "b", Scope: GlobalScope, Index: 1},
		}

		for _, sym := range expected {
			result, ok := global.Resolve(sym.Name)

			if !ok {
				t.Errorf("expected %s to be resolved", sym.Name)

				continue
			}

			if result != sym {
				t.Errorf("expected %s=%+v, got=%+v", sym.Name, sym, result)
			}
		}
	})
}

func TestResolveLocal(t *testing.T) {
	t.Run("should resolve a local symbol", func(t *testing.T) {
		global := NewSymbolTable()

		global.Define("a")
		global.Define("b")

		local := NewEnclosedSymbolTable(global)
		local.Define("c")
		local.Define("d")

		expected := []Symbol{
			{Name: "a", Scope: GlobalScope, Index: 0},
			{Name: "b", Scope: GlobalScope, Index: 1},
			{Name: "c", Scope: LocalScope, Index: 0},
			{Name: "d", Scope: LocalScope, Index: 1},
		}

		for _, sym := range expected {
			result, ok := local.Resolve(sym.Name)

			if !ok {
				t.Errorf("expected %s to be resolved", sym.Name)

				continue
			}

			if result != sym {
				t.Errorf("expected %s=%+v, got=%+v", sym.Name, sym, result)
			}
		}
	})
}

func TestResolveNestedLocal(t *testing.T) {
	t.Run("should resolve a nested local symbol", func(t *testing.T) {
		global := NewSymbolTable()
		global.Define("a")
		global.Define("b")

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c")
		firstLocal.Define("d")

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e")
		secondLocal.Define("f")

		tests := []struct {
			table           *SymbolTable
			expectedSymbols []Symbol
		}{
			{
				firstLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0},
					{Name: "b", Scope: GlobalScope, Index: 1},
					{Name: "c", Scope: LocalScope, Index: 0},
					{Name: "d", Scope: LocalScope, Index: 1},
				},
			},
			{
				secondLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0},
					{Name: "b", Scope: GlobalScope, Index: 1},
					{Name: "e", Scope: LocalScope, Index: 0},
					{Name: "f", Scope: LocalScope, Index: 1},
				},
			},
		}

		for _, tt := range tests {
			for _, sym := range tt.expectedSymbols {
				result, ok := tt.table.Resolve(sym.Name)

				if !ok {
					t.Errorf("expected %s to be resolved", sym.Name)

					continue
				}

				if result != sym {
					t.Errorf("expected %s=%+v, got=%+v", sym.Name, sym, result)
				}
			}
		}
	})
}

func TestDefineResolveBultins(t *testing.T) {
	t.Run("should define and resolve builtins", func(t *testing.T) {
		global := NewSymbolTable()
		firstLocal := NewEnclosedSymbolTable(global)
		secondLocal := NewEnclosedSymbolTable(firstLocal)

		expected := []Symbol{
			{Name: "a", Scope: BuiltinScope, Index: 0},
			{Name: "c", Scope: BuiltinScope, Index: 1},
			{Name: "e", Scope: BuiltinScope, Index: 2},
			{Name: "f", Scope: BuiltinScope, Index: 3},
		}

		for i, v := range expected {
			global.DefineBuiltin(i, v.Name)
		}

		for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
			for _, sym := range expected {
				result, ok := table.Resolve(sym.Name)

				if !ok {
					t.Errorf("expected %s to be resolved", sym.Name)

					continue
				}

				if result != sym {
					t.Errorf("expected %s=%+v, got=%+v", sym.Name, sym, result)
				}
			}
		}
	})
}

func TestResolveFree(t *testing.T) {
	t.Run("should resolve a free symbol", func(t *testing.T) {
		global := NewSymbolTable()
		global.Define("a")
		global.Define("b")

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c")
		firstLocal.Define("d")

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e")
		secondLocal.Define("f")

		tests := []struct {
			table               *SymbolTable
			expectedSymbols     []Symbol
			expectedFreeSymbols []Symbol
		}{
			{
				firstLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0},
					{Name: "b", Scope: GlobalScope, Index: 1},
					{Name: "c", Scope: LocalScope, Index: 0},
					{Name: "d", Scope: LocalScope, Index: 1},
				},
				[]Symbol{},
			},
			{
				secondLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0},
					{Name: "b", Scope: GlobalScope, Index: 1},
					{Name: "c", Scope: FreeScope, Index: 0},
					{Name: "d", Scope: FreeScope, Index: 1},
					{Name: "e", Scope: LocalScope, Index: 0},
					{Name: "f", Scope: LocalScope, Index: 1},
				},
				[]Symbol{
					{Name: "c", Scope: LocalScope, Index: 0},
					{Name: "d", Scope: LocalScope, Index: 1},
				},
			},
		}

		for _, tt := range tests {
			for _, sym := range tt.expectedSymbols {
				result, ok := tt.table.Resolve(sym.Name)
				if !ok {
					t.Errorf("name %s is not resolvable", sym.Name)

					continue
				}

				if result != sym {
					t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
				}
			}

			if len(tt.table.FreeSymbols) != len(tt.expectedFreeSymbols) {
				t.Errorf("wrong number of free symbols. got=%d, want=%d", len(tt.table.FreeSymbols), len(tt.expectedFreeSymbols))

				continue
			}

			for i, sym := range tt.expectedFreeSymbols {
				result := tt.table.FreeSymbols[i]

				if result != sym {
					t.Errorf("wrong free symbol. got=%+v, want=%+v", result, sym)
				}
			}
		}
	})
}

func TestResolveUnresolvableFree(t *testing.T) {
	t.Run("should resolve an unresolvable free symbol", func(t *testing.T) {
		global := NewSymbolTable()
		global.Define("a")

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c")

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e")
		secondLocal.Define("f")

		expected := []Symbol{
			{Name: "a", Scope: GlobalScope, Index: 0},
			{Name: "c", Scope: FreeScope, Index: 0},
			{Name: "e", Scope: LocalScope, Index: 0},
			{Name: "f", Scope: LocalScope, Index: 1},
		}

		for _, sym := range expected {
			result, ok := secondLocal.Resolve(sym.Name)

			if !ok {
				t.Errorf("name %s is not resolvable", sym.Name)

				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}

		expectedUnresolvable := []string{"b", "d"}

		for _, name := range expectedUnresolvable {
			_, ok := secondLocal.Resolve(name)

			if ok {
				t.Errorf("name %s is resolvable, but was expected not to", name)
			}
		}
	})
}

func TestDefineAndResolveFunctionName(t *testing.T) {
	t.Run("should define and resolve function name", func(t *testing.T) {
		global := NewSymbolTable()
		global.DefineFunctionName("a")
		expected := Symbol{Name: "a", Scope: FunctionScope, Index: 0}

		result, ok := global.Resolve(expected.Name)
		if !ok {
			t.Errorf("function name %s is not resolvable", expected.Name)

			return
		}

		if result != expected {
			t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, result)
		}
	})
}

func TestShadowingFunctionName(t *testing.T) {
	t.Run("should shadow function name", func(t *testing.T) {
		global := NewSymbolTable()
		global.DefineFunctionName("a")
		global.Define("a")

		expected := Symbol{Name: "a", Scope: GlobalScope, Index: 0}

		result, ok := global.Resolve(expected.Name)
		if !ok {
			t.Errorf("function name %s is not resolvable", expected.Name)

			return
		}

		if result != expected {
			t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, result)
		}
	})
}
