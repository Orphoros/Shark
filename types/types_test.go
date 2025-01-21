package types

import "testing"

func TestPrimitiveTypes(t *testing.T) {
	t.Run("TSharkBool", func(t *testing.T) {
		var sharkBool TSharkBool
		if sharkBool.SharkTypeString() != "bool" {
			t.Errorf("Expected TSharkBool.SharkTypeString() to return 'bool', got %s", sharkBool.SharkTypeString())
		}
		if !sharkBool.Is(TSharkBool{}) {
			t.Error("Expected TSharkBool.Is(TSharkBool{}) to return true")
		}
		if sharkBool.Is(TSharkAny{}) {
			t.Error("Expected TSharkBool.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkAny", func(t *testing.T) {
		var sharkAny TSharkAny
		if sharkAny.SharkTypeString() != "any" {
			t.Errorf("Expected TSharkAny.SharkTypeString() to return 'any', got %s", sharkAny.SharkTypeString())
		}
		if !sharkAny.Is(TSharkAny{}) {
			t.Error("Expected TSharkAny.Is(TSharkAny{}) to return true")
		}
		if !sharkAny.Is(TSharkBool{}) {
			t.Error("Expected TSharkAny.Is(TSharkBool{}) to return true")
		}
		if !sharkAny.Is(TSharkArray{}) {
			t.Error("Expected TSharkAny.Is(TSharkArray{}) to return true")
		}
	})

	t.Run("TSharkString", func(t *testing.T) {
		var sharkString TSharkString
		if sharkString.SharkTypeString() != "string" {
			t.Errorf("Expected TSharkString.SharkTypeString() to return 'string', got %s", sharkString.SharkTypeString())
		}
		if !sharkString.Is(TSharkString{}) {
			t.Error("Expected TSharkString.Is(TSharkString{}) to return true")
		}
		if sharkString.Is(TSharkAny{}) {
			t.Error("Expected TSharkString.Is(TSharkAny{}) to return false")
		}

		if sharkString.Is(TSharkArray{}) {
			t.Error("Expected TSharkString.Is(TSharkArray{}) to return false")
		}
	})

	t.Run("TSharkI64", func(t *testing.T) {
		var sharkI64 TSharkI64
		if sharkI64.SharkTypeString() != "i64" {
			t.Errorf("Expected TSharkI64.SharkTypeString() to return 'i64', got %s", sharkI64.SharkTypeString())
		}
		if !sharkI64.Is(TSharkI64{}) {
			t.Error("Expected TSharkI64.Is(TSharkI64{}) to return true")
		}
		if sharkI64.Is(TSharkAny{}) {
			t.Error("Expected TSharkI64.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkError", func(t *testing.T) {
		var sharkError TSharkError
		if sharkError.SharkTypeString() != "error" {
			t.Errorf("Expected TSharkError.SharkTypeString() to return 'error', got %s", sharkError.SharkTypeString())
		}
		if !sharkError.Is(TSharkError{}) {
			t.Error("Expected TSharkError.Is(TSharkError{}) to return true")
		}
		if sharkError.Is(TSharkAny{}) {
			t.Error("Expected TSharkError.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkNull", func(t *testing.T) {
		var sharkNull TSharkNull
		if sharkNull.SharkTypeString() != "null" {
			t.Errorf("Expected TSharkNull.SharkTypeString() to return 'null', got %s", sharkNull.SharkTypeString())
		}
		if !sharkNull.Is(TSharkNull{}) {
			t.Error("Expected TSharkNull.Is(TSharkNull{}) to return true")
		}
		if sharkNull.Is(TSharkAny{}) {
			t.Error("Expected TSharkNull.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkSpread", func(t *testing.T) {
		var sharkSpread TSharkSpread
		if sharkSpread.SharkTypeString() != "..." {
			t.Errorf("Expected TSharkSpread.SharkTypeString() to return '...', got %s", sharkSpread.SharkTypeString())
		}
		if !sharkSpread.Is(TSharkSpread{}) {
			t.Error("Expected TSharkSpread.Is(TSharkSpread{}) to return true")
		}
		if !sharkSpread.Is(TSharkAny{}) {
			t.Error("Expected TSharkSpread.Is(TSharkAny{}) to return true")
		}
	})

	t.Run("TSharkSpread Typed", func(t *testing.T) {
		sharkSpread := TSharkSpread{Type: TSharkI64{}}

		if sharkSpread.SharkTypeString() != "...i64" {
			t.Errorf("Expected TSharkSpread.SharkTypeString() to return '...i64', got %s", sharkSpread.SharkTypeString())
		}
		if !sharkSpread.Is(TSharkSpread{Type: TSharkI64{}}) {
			t.Error("Expected TSharkSpread.Is(TSharkSpread{Type: TSharkI64{}}) to return true")
		}
		if sharkSpread.Is(TSharkSpread{Type: TSharkBool{}}) {
			t.Error("Expected TSharkSpread.Is(TSharkSpread{Type: TSharkBool{}}) to return false")
		}
		if !sharkSpread.Is(TSharkAny{}) {
			t.Error("Expected TSharkSpread.Is(TSharkAny{}) to return true")
		}
		if sharkSpread.Is(TSharkSpread{}) {
			t.Error("Expected TSharkSpread.Is(TSharkSpread{}) to return false")
		}
		if sharkSpread.Is(TSharkSpread{Type: TSharkString{}}) {
			t.Error("Expected TSharkSpread.Is(TSharkSpread{Type: TSharkString{}}) to return false")
		}
	})
}

func TestCollectionTypes(t *testing.T) {
	t.Run("TSharkArray", func(t *testing.T) {
		var sharkArray TSharkArray
		if sharkArray.SharkTypeString() != "array<>" {
			t.Errorf("Expected TSharkArray.SharkTypeString() to return 'array<>', got %s", sharkArray.SharkTypeString())
		}
		if !sharkArray.Is(TSharkArray{}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{}) to return true")
		}
		if sharkArray.Is(TSharkAny{}) {
			t.Error("Expected TSharkArray.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkArray Typed", func(t *testing.T) {
		sharkArray := TSharkArray{Collects: TSharkI64{}}

		if sharkArray.SharkTypeString() != "array<i64>" {
			t.Errorf("Expected TSharkArray.SharkTypeString() to return 'array<i64>', got %s", sharkArray.SharkTypeString())
		}
		if !sharkArray.Is(TSharkArray{Collects: TSharkI64{}}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{Collects: TSharkI64{}}) to return true")
		}
		if sharkArray.Is(TSharkArray{Collects: TSharkBool{}}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{Collects: TSharkBool{}}) to return false")
		}
		if sharkArray.Is(TSharkAny{}) {
			t.Error("Expected TSharkArray.Is(TSharkAny{}) to return false")
		}
		if !sharkArray.Is(TSharkArray{}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{}) to return true")
		}
		if !sharkArray.Is(TSharkArray{Collects: TSharkAny{}}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{Collects: TSharkAny{}}) to return true")
		}
		if sharkArray.Is(TSharkArray{Collects: TSharkString{}}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{Collects: TSharkString{}}) to return false")
		}
	})

	t.Run("TSharkMap", func(t *testing.T) {
		var sharkMap TSharkHashMap
		if sharkMap.SharkTypeString() != "hashmap<>" {
			t.Errorf("Expected TSharkMap.SharkTypeString() to return 'hashmap<>', got %s", sharkMap.SharkTypeString())
		}
		if !sharkMap.Is(TSharkHashMap{}) {
			t.Error("Expected TSharkMap.Is(TSharkMap{}) to return true")
		}
		if sharkMap.Is(TSharkAny{}) {
			t.Error("Expected TSharkMap.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkMap Typed", func(t *testing.T) {
		sharkMap := TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}

		if sharkMap.SharkTypeString() != "hashmap<i64,string>" {
			t.Errorf("Expected TSharkMap.SharkTypeString() to return 'hashmap<i64,string>', got %s", sharkMap.SharkTypeString())
		}
		if !sharkMap.Is(TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}) {
			t.Error("Expected TSharkMap.Is(TSharkMap{Indexes: TSharkI64{}, Collects: TSharkString{}}) to return true")
		}
		if sharkMap.Is(TSharkHashMap{Indexes: TSharkBool{}, Collects: TSharkString{}}) {
			t.Error("Expected TSharkMap.Is(TSharkMap{Indexes: TSharkBool{}, Collects: TSharkString{}}) to return false")
		}
		if !sharkMap.Is(TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkAny{}}) {
			t.Error("Expected TSharkMap.Is(TSharkMap{Indexes: TSharkI64{}, Collects: TSharkAny{}}) to return true")
		}
		if !sharkMap.Is(TSharkHashMap{Indexes: TSharkAny{}, Collects: TSharkString{}}) {
			t.Error("Expected TSharkMap.Is(TSharkMap{Indexes: TSharkAny{}, Collects: TSharkString{}}) to return true")
		}
		if !sharkMap.Is(TSharkHashMap{}) {
			t.Error("Expected TSharkMap.Is(TSharkMap{}) to return true")
		}
	})

	t.Run("TSharkTuple", func(t *testing.T) {
		var sharkTuple TSharkTuple
		if sharkTuple.SharkTypeString() != "tuple<>" {
			t.Errorf("Expected TSharkTuple.SharkTypeString() to return 'tuple<>', got %s", sharkTuple.SharkTypeString())
		}
		if !sharkTuple.Is(TSharkTuple{}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{}) to return true")
		}
		if sharkTuple.Is(TSharkAny{}) {
			t.Error("Expected TSharkTuple.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkTuple Typed", func(t *testing.T) {
		sharkTuple := TSharkTuple{Collects: []ISharkType{TSharkI64{}}}

		if sharkTuple.SharkTypeString() != "tuple<i64>" {
			t.Errorf("Expected TSharkTuple.SharkTypeString() to return 'tuple<i64>', got %s", sharkTuple.SharkTypeString())
		}
		if !sharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkI64{}}}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkI64{}}}) to return true")
		}
		if sharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkBool{}}}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkBool{}}}) to return false")
		}
		if sharkTuple.Is(TSharkAny{}) {
			t.Error("Expected TSharkTuple.Is(TSharkAny{}) to return false")
		}
		if !sharkTuple.Is(TSharkTuple{}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{}) to return true")
		}
		if sharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkString{}}}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkString{}}}) to return false")
		}
	})

	t.Run("TSharkTuple Typed with multiple types", func(t *testing.T) {
		sharkTuple := TSharkTuple{Collects: []ISharkType{TSharkI64{}, TSharkString{}}}

		if sharkTuple.SharkTypeString() != "tuple<i64,string>" {
			t.Errorf("Expected TSharkTuple.SharkTypeString() to return 'tuple<i64,string>', got %s", sharkTuple.SharkTypeString())
		}
		if !sharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkI64{}, TSharkString{}}}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkI64{}, TSharkString{}}}) to return true")
		}
		if sharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkBool{}}}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkBool{}}}) to return false")
		}
		if sharkTuple.Is(TSharkAny{}) {
			t.Error("Expected TSharkTuple.Is(TSharkAny{}) to return false")
		}
		if !sharkTuple.Is(TSharkTuple{}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{}) to return true")
		}
		if sharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkString{}, TSharkI64{}}}) {
			t.Error("Expected TSharkTuple.Is(TSharkTuple{Collects: []ISharkType{TSharkString{}, TSharkI64{}}}) to return false")
		}
	})
}

func TestFunctionTypes(t *testing.T) {
	t.Run("TSharkFunction", func(t *testing.T) {
		var sharkFunction TSharkFuncType
		if sharkFunction.SharkTypeString() != "func<()>" {
			t.Errorf("Expected TSharkFunction.SharkTypeString() to return 'func<()>', got %s", sharkFunction.SharkTypeString())
		}
		if !sharkFunction.Is(TSharkFuncType{}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{}) to return true")
		}
		if sharkFunction.Is(TSharkAny{}) {
			t.Error("Expected TSharkFunction.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkFunction with return value", func(t *testing.T) {
		sharkFunction := TSharkFuncType{ReturnT: TSharkI64{}}

		if sharkFunction.SharkTypeString() != "func<()->i64>" {
			t.Errorf("Expected TSharkFunction.SharkTypeString() to return 'func<()->i64>', got %s", sharkFunction.SharkTypeString())
		}
		if !sharkFunction.Is(TSharkFuncType{ReturnT: TSharkI64{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Returns: TSharkI64{}}) to return true")
		}
		if sharkFunction.Is(TSharkFuncType{ReturnT: TSharkBool{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Returns: TSharkBool{}}) to return false")
		}
		if sharkFunction.Is(TSharkAny{}) {
			t.Error("Expected TSharkFunction.Is(TSharkAny{}) to return false")
		}
		if !sharkFunction.Is(TSharkFuncType{}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{}) to return true")
		}
		if sharkFunction.Is(TSharkFuncType{ReturnT: TSharkAny{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Returns: TSharkAny{}}) to return false")
		}
		if sharkFunction.Is(TSharkFuncType{ReturnT: TSharkString{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Returns: TSharkString{}}) to return false")
		}
	})

	t.Run("TSharkFunction with one argument", func(t *testing.T) {
		sharkFunction := TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}

		if sharkFunction.SharkTypeString() != "func<(i64)>" {
			t.Errorf("Expected TSharkFunction.SharkTypeString() to return 'func<(i64)>', got %s", sharkFunction.SharkTypeString())
		}
		if !sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkI64{}}}) to return true")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkBool{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkBool{}}}) to return false")
		}
		if sharkFunction.Is(TSharkAny{}) {
			t.Error("Expected TSharkFunction.Is(TSharkAny{}) to return false")
		}
		if !sharkFunction.Is(TSharkFuncType{}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{}) to return true")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkAny{}}}) to return false")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkString{}, TSharkI64{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkString{}, TSharkI64{}}}) to return false")
		}
		if sharkFunction.Is(TSharkFuncType{ReturnT: TSharkI64{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Returns: TSharkI64{}}) to return false")
		}
	})

	t.Run("TSharkFunction with multiple arguments", func(t *testing.T) {
		sharkFunction := TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkString{}}}

		if sharkFunction.SharkTypeString() != "func<(i64,string)>" {
			t.Errorf("Expected TSharkFunction.SharkTypeString() to return 'func<(i64,string)>', got %s", sharkFunction.SharkTypeString())
		}
		if !sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkString{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkI64{}, TSharkString{}}}) to return true")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkBool{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkBool{}}}) to return false")
		}
		if sharkFunction.Is(TSharkAny{}) {
			t.Error("Expected TSharkFunction.Is(TSharkAny{}) to return false")
		}
		if !sharkFunction.Is(TSharkFuncType{}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{}) to return true")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkAny{}}}) to return false")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkString{}, TSharkI64{}}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkString{}, TSharkI64{}}}) to return false")
		}
		if sharkFunction.Is(TSharkFuncType{ReturnT: TSharkI64{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Returns: TSharkI64{}}) to return false")
		}
	})

	t.Run("TSharkFunction with multiple arguments and return value", func(t *testing.T) {
		sharkFunction := TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkString{}}, ReturnT: TSharkBool{}}

		if sharkFunction.SharkTypeString() != "func<(i64,string)->bool>" {
			t.Errorf("Expected TSharkFunction.SharkTypeString() to return 'func<(i64,string)->bool>', got %s", sharkFunction.SharkTypeString())
		}
		if !sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkString{}}, ReturnT: TSharkBool{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkI64{}, TSharkString{}}, Returns: TSharkBool{}}) to return true")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkString{}}, ReturnT: TSharkI64{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkI64{}, TSharkString{}}, Returns: TSharkI64{}}) to return false")
		}
		if sharkFunction.Is(TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkString{}}, ReturnT: TSharkAny{}}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{Args: []ISharkType{TSharkI64{}, TSharkString{}}, Returns: TSharkAny{}}) to return false")
		}
		if !sharkFunction.Is(TSharkFuncType{}) {
			t.Error("Expected TSharkFunction.Is(TSharkFunction{}) to return true")
		}
	})
}

func TestVariadicTypes(t *testing.T) {
	t.Run("TSharkVariadic", func(t *testing.T) {
		var sharkVariadic TSharkVariadic
		if sharkVariadic.SharkTypeString() != "T" {
			t.Errorf("Expected TSharkVariadic.SharkTypeString() to return 'T', got %s", sharkVariadic.SharkTypeString())
		}
		if !sharkVariadic.Is(TSharkVariadic{}) {
			t.Error("Expected TSharkVariadic.Is(TSharkVariadic{}) to return true")
		}
		if !sharkVariadic.Is(TSharkAny{}) {
			t.Error("Expected TSharkVariadic.Is(TSharkAny{}) to return true")
		}
		if !sharkVariadic.Is(TSharkI64{}) {
			t.Error("Expected TSharkVariadic.Is(TSharkI64{}) to return true")
		}
	})

	t.Run("TSharkVariadic Typed", func(t *testing.T) {
		sharkVariadic := TSharkVariadic{Enclosed: TSharkI64{}}

		if sharkVariadic.SharkTypeString() != "T" {
			t.Errorf("Expected TSharkVariadic.SharkTypeString() to return 'T', got %s", sharkVariadic.SharkTypeString())
		}
		if !sharkVariadic.Is(TSharkVariadic{Enclosed: TSharkI64{}}) {
			t.Error("Expected TSharkVariadic.Is(TSharkVariadic{Enclosed: TSharkI64{}}) to return true")
		}
		if sharkVariadic.Is(TSharkVariadic{Enclosed: TSharkBool{}}) {
			t.Error("Expected TSharkVariadic.Is(TSharkVariadic{Enclosed: TSharkBool{}}) to return false")
		}
		if sharkVariadic.Is(TSharkAny{}) {
			t.Error("Expected TSharkVariadic.Is(TSharkAny{}) to return false")
		}
		if !sharkVariadic.Is(TSharkI64{}) {
			t.Error("Expected TSharkVariadic.Is(TSharkI64{}) to return true")
		}
	})

	t.Run("Array with TSharkVariadic unenclosed", func(t *testing.T) {
		sharkArray := TSharkArray{Collects: TSharkVariadic{}}

		if sharkArray.SharkTypeString() != "array<T>" {
			t.Errorf("Expected TSharkArray.SharkTypeString() to return 'array<T>', got %s", sharkArray.SharkTypeString())
		}
		if !sharkArray.Is(TSharkArray{Collects: TSharkVariadic{Enclosed: TSharkI64{}}}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{Collects: TSharkVariadic{Enclosed: TSharkI64{}}}) to return true")
		}
		if !sharkArray.Is(TSharkArray{Collects: TSharkVariadic{Enclosed: TSharkBool{}}}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{Collects: TSharkVariadic{Enclosed: TSharkBool{}}}) to return false")
		}
		if sharkArray.Is(TSharkAny{}) {
			t.Error("Expected TSharkArray.Is(TSharkAny{}) to return false")
		}
		if !sharkArray.Is(TSharkArray{Collects: TSharkI64{}}) {
			t.Error("Expected TSharkArray.Is(TSharkArray{Collects: TSharkI64{}}) to return true")
		}
	})
}

func TestClosureTypes(t *testing.T) {
	t.Run("TSharkClosure", func(t *testing.T) {
		var sharkClosure TSharkClosure
		if sharkClosure.SharkTypeString() != "closure<>" {
			t.Errorf("Expected TSharkClosure.SharkTypeString() to return 'closure<>', got %s", sharkClosure.SharkTypeString())
		}
		if !sharkClosure.Is(TSharkClosure{}) {
			t.Error("Expected TSharkClosure.Is(TSharkClosure{}) to return true")
		}
		if sharkClosure.Is(TSharkAny{}) {
			t.Error("Expected TSharkClosure.Is(TSharkAny{}) to return false")
		}
	})

	t.Run("TSharkClosure Typed", func(t *testing.T) {
		sharkClosure := TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkI64{}}}

		if sharkClosure.SharkTypeString() != "func<()->i64>" {
			t.Errorf("Expected TSharkClosure.SharkTypeString() to return 'func<()->i64>', got %s", sharkClosure.SharkTypeString())
		}
		if !sharkClosure.Is(TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkI64{}}}) {
			t.Error("Expected TSharkClosure.Is(TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkI64{}}}) to return true")
		}
		if sharkClosure.Is(TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkBool{}}}) {
			t.Error("Expected TSharkClosure.Is(TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkBool{}}}) to return false")
		}
		if sharkClosure.Is(TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkAny{}}}) {
			t.Error("Expected TSharkClosure.Is(TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkAny{}}}) to return false")
		}
		if !sharkClosure.Is(TSharkClosure{}) {
			t.Error("Expected TSharkClosure.Is(TSharkClosure{}) to return true")
		}
	})
}
