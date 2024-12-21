package compiler

import (
	"shark/object"
	"testing"
)

func TestDefine(t *testing.T) {
	t.Run("should define a new symbol", func(t *testing.T) {
		expected := map[string]Symbol{
			"a": {Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			"b": {Name: "b", Scope: GlobalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
			"c": {Name: "c", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			"d": {Name: "d", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
			"e": {Name: "e", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			"f": {Name: "f", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
		}

		global := NewSymbolTable()

		if a := global.Define("a", true, true, object.INTEGER_OBJ, nil); a != expected["a"] {
			t.Errorf("expected %s=%+v, got=%+v", "a", expected["a"], a)
		}

		if b := global.Define("b", false, false, object.INTEGER_OBJ, nil); b != expected["b"] {
			t.Errorf("expected %s=%+v, got=%+v", "b", expected["b"], b)
		}

		firstLocal := NewEnclosedSymbolTable(global)

		if c := firstLocal.Define("c", true, true, object.INTEGER_OBJ, nil); c != expected["c"] {
			t.Errorf("expected %s=%+v, got=%+v", "c", expected["c"], c)
		}

		if d := firstLocal.Define("d", false, false, object.INTEGER_OBJ, nil); d != expected["d"] {
			t.Errorf("expected %s=%+v, got=%+v", "d", expected["d"], d)
		}

		secondLocal := NewEnclosedSymbolTable(firstLocal)

		if e := secondLocal.Define("e", true, true, object.INTEGER_OBJ, nil); e != expected["e"] {
			t.Errorf("expected %s=%+v, got=%+v", "e", expected["e"], e)
		}

		if f := secondLocal.Define("f", false, false, object.INTEGER_OBJ, nil); f != expected["f"] {
			t.Errorf("expected %s=%+v, got=%+v", "f", expected["f"], f)
		}
	})
}

func TestResolveGlobal(t *testing.T) {
	t.Run("should resolve a global symbol", func(t *testing.T) {
		global := NewSymbolTable()

		global.Define("a", true, true, object.INTEGER_OBJ, nil)
		global.Define("b", true, true, object.INTEGER_OBJ, nil)

		expected := []Symbol{
			{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "b", Scope: GlobalScope, Index: 1, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
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

		global.Define("a", true, true, object.INTEGER_OBJ, nil)
		global.Define("b", false, false, object.INTEGER_OBJ, nil)

		local := NewEnclosedSymbolTable(global)
		local.Define("c", true, true, object.INTEGER_OBJ, nil)
		local.Define("d", false, false, object.INTEGER_OBJ, nil)

		expected := []Symbol{
			{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "b", Scope: GlobalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
			{Name: "c", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "d", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
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
		global.Define("a", true, true, object.INTEGER_OBJ, nil)
		global.Define("b", false, false, object.INTEGER_OBJ, nil)

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c", true, true, object.INTEGER_OBJ, nil)
		firstLocal.Define("d", false, false, object.INTEGER_OBJ, nil)

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e", true, true, object.INTEGER_OBJ, nil)
		secondLocal.Define("f", false, false, object.INTEGER_OBJ, nil)

		tests := []struct {
			table           *SymbolTable
			expectedSymbols []Symbol
		}{
			{
				firstLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "b", Scope: GlobalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
					{Name: "c", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "d", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
				},
			},
			{
				secondLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "b", Scope: GlobalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
					{Name: "e", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "f", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
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
			{Name: "a", Scope: BuiltinScope, Index: 0, Mutable: false, VariadicType: false},
			{Name: "c", Scope: BuiltinScope, Index: 1, Mutable: false, VariadicType: false},
			{Name: "e", Scope: BuiltinScope, Index: 2, Mutable: false, VariadicType: false},
			{Name: "f", Scope: BuiltinScope, Index: 3, Mutable: false, VariadicType: false},
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
		global.Define("a", true, true, object.INTEGER_OBJ, nil)
		global.Define("b", false, false, object.INTEGER_OBJ, nil)

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c", true, true, object.INTEGER_OBJ, nil)
		firstLocal.Define("d", false, false, object.INTEGER_OBJ, nil)

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e", true, true, object.INTEGER_OBJ, nil)
		secondLocal.Define("f", false, false, object.INTEGER_OBJ, nil)

		tests := []struct {
			table               *SymbolTable
			expectedSymbols     []Symbol
			expectedFreeSymbols []Symbol
		}{
			{
				firstLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "b", Scope: GlobalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
					{Name: "c", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "d", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
				},
				[]Symbol{},
			},
			{
				secondLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "b", Scope: GlobalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
					{Name: "c", Scope: FreeScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "d", Scope: FreeScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
					{Name: "e", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "f", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
				},
				[]Symbol{
					{Name: "c", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
					{Name: "d", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
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
		global.Define("a", true, true, object.INTEGER_OBJ, nil)

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c", true, true, object.INTEGER_OBJ, nil)

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e", true, true, object.INTEGER_OBJ, nil)
		secondLocal.Define("f", false, false, object.INTEGER_OBJ, nil)

		expected := []Symbol{
			{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "c", Scope: FreeScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "e", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "f", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
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
		global.DefineFunctionName("a", nil)
		expected := Symbol{Name: "a", Scope: FunctionScope, Index: 0, Mutable: false, VariadicType: false}

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
		global.DefineFunctionName("a", nil)
		global.Define("a", true, true, object.INTEGER_OBJ, nil)

		expected := Symbol{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ}

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

func TestFindSymbolByName(t *testing.T) {
	t.Run("should find symbol by name", func(t *testing.T) {
		global := NewSymbolTable()
		global.Define("a", true, true, object.INTEGER_OBJ, nil)
		global.Define("b", false, false, object.INTEGER_OBJ, nil)

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c", true, true, object.INTEGER_OBJ, nil)
		firstLocal.Define("d", false, false, object.INTEGER_OBJ, nil)

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e", true, true, object.INTEGER_OBJ, nil)
		secondLocal.Define("f", false, false, object.INTEGER_OBJ, nil)

		expected := []Symbol{
			{Name: "a", Scope: GlobalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "b", Scope: GlobalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
			{Name: "c", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "d", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
			{Name: "e", Scope: LocalScope, Index: 0, Mutable: true, VariadicType: true, ObjType: object.INTEGER_OBJ},
			{Name: "f", Scope: LocalScope, Index: 1, Mutable: false, VariadicType: false, ObjType: object.INTEGER_OBJ},
		}

		for _, sym := range expected {
			result, ok := secondLocal.FindIdent(sym.Name)

			if !ok {
				t.Errorf("expected %s to be found", sym.Name)

				continue
			}

			if result != sym {
				t.Errorf("expected %s=%+v, got=%+v", sym.Name, sym, result)
			}
		}
	})
}
