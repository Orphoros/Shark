package types

import (
	"testing"
)

func TestPrimitiveTypes(t *testing.T) {
	t.Run("should validate type bool", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkBool{}, TSharkBool{}, true},
			{TSharkBool{}, TSharkAny{}, false},
			{TSharkBool{}, TSharkI64{}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkBool{}, "bool"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}

	})

	t.Run("should validate type any", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkAny{}, TSharkAny{}, true},
			{TSharkAny{}, TSharkBool{}, true},
			{TSharkAny{}, TSharkArray{}, true},
			{TSharkAny{}, TSharkOptional{Type: TSharkBool{}}, false},
			{TSharkAny{}, TSharkSpread{Type: TSharkBool{}}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkAny{}, "any"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type string", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkString{}, TSharkString{}, true},
			{TSharkString{}, TSharkAny{}, false},
			{TSharkString{}, TSharkI64{}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkString{}, "string"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type i64", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkI64{}, TSharkI64{}, true},
			{TSharkI64{}, TSharkAny{}, false},
			{TSharkI64{}, TSharkString{}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkI64{}, "i64"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type error", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkError{}, TSharkError{}, true},
			{TSharkError{}, TSharkAny{}, false},
			{TSharkError{}, TSharkI64{}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkError{}, "error"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type null", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkNull{}, TSharkNull{}, true},
			{TSharkNull{}, TSharkAny{}, false},
			{TSharkNull{}, TSharkI64{}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkNull{}, "null"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type spread", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkSpread{}, TSharkSpread{}, true},
			{TSharkSpread{}, TSharkAny{}, false},
			{TSharkSpread{}, TSharkI64{}, false},
			{TSharkSpread{Type: TSharkI64{}}, TSharkSpread{Type: TSharkI64{}}, true},
			{TSharkSpread{}, TSharkSpread{Type: TSharkI64{}}, false},
			{TSharkSpread{Type: TSharkI64{}}, TSharkSpread{Type: TSharkAny{}}, false},
			{TSharkSpread{Type: TSharkI64{}}, TSharkSpread{}, true},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkSpread{}, "..."},
			{TSharkSpread{Type: TSharkI64{}}, "...i64"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type optional", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkOptional{}, TSharkOptional{}, true},
			{TSharkOptional{}, TSharkAny{}, false},
			{TSharkOptional{}, TSharkI64{}, false},
			{TSharkOptional{Type: TSharkI64{}}, TSharkOptional{Type: TSharkI64{}}, true},
			{TSharkOptional{Type: TSharkI64{}}, TSharkI64{}, true},
			{TSharkOptional{}, TSharkOptional{Type: TSharkI64{}}, false},
			{TSharkOptional{Type: TSharkI64{}}, TSharkOptional{Type: TSharkAny{}}, false},
			{TSharkOptional{Type: TSharkI64{}}, TSharkOptional{}, true},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkOptional{}, "?"},
			{TSharkOptional{Type: TSharkI64{}}, "i64?"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})
}

func TestCollectionTypes(t *testing.T) {
	t.Run("should validate type array", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkArray{}, TSharkArray{}, true},
			{TSharkArray{}, TSharkAny{}, false},
			{TSharkArray{}, TSharkI64{}, false},
			{TSharkArray{Collection: TSharkI64{}}, TSharkArray{Collection: TSharkI64{}}, true},
			{TSharkArray{Collection: TSharkI64{}}, TSharkArray{Collection: TSharkAny{}}, false},
			{TSharkArray{}, TSharkArray{Collection: TSharkI64{}}, false},
			{TSharkArray{Collection: TSharkI64{}}, TSharkArray{}, true},
			{TSharkArray{Collection: TSharkAny{}}, TSharkArray{}, true},
			{TSharkArray{Collection: TSharkAny{}}, TSharkArray{Collection: TSharkBool{}}, true},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkArray{}, "array<>"},
			{TSharkArray{Collection: TSharkI64{}}, "array<i64>"},
			{TSharkArray{Collection: TSharkAny{}}, "array<any>"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type hashmap", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{

			{TSharkHashMap{}, TSharkHashMap{}, true},
			{TSharkHashMap{}, TSharkAny{}, false},
			{TSharkHashMap{}, TSharkI64{}, false},
			{TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, true},
			{TSharkHashMap{}, TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, false},
			{TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, TSharkHashMap{Indexes: TSharkAny{}, Collects: TSharkString{}}, false},
			{TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkAny{}}, false},
			{TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, TSharkHashMap{Indexes: TSharkI64{}}, false},
			{TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, TSharkHashMap{Collects: TSharkString{}}, false},
			{TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, TSharkHashMap{}, true},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkHashMap{}, "hashmap<>"},
			{TSharkHashMap{Indexes: TSharkI64{}, Collects: TSharkString{}}, "hashmap<i64,string>"},
			{TSharkHashMap{Indexes: TSharkAny{}, Collects: TSharkAny{}}, "hashmap<any,any>"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate type tuple", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkTuple{}, TSharkTuple{}, true},
			{TSharkTuple{}, TSharkAny{}, false},
			{TSharkTuple{}, TSharkI64{}, false},
			{TSharkTuple{Collection: []ISharkType{TSharkSpread{Type: TSharkAny{}}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, true},
			{TSharkTuple{Collection: []ISharkType{TSharkSpread{Type: TSharkAny{}}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}, TSharkString{}}}, false},
			{TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, true},
			{TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, TSharkTuple{Collection: []ISharkType{TSharkAny{}}}, false},
			{TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}, TSharkString{}}}, false},
			{TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, true},
			{TSharkTuple{Collection: []ISharkType{TSharkAny{}}}, TSharkTuple{Collection: []ISharkType{TSharkAny{}}}, true},
			{TSharkTuple{Collection: []ISharkType{TSharkAny{}}}, TSharkTuple{Collection: []ISharkType{TSharkBool{}}}, true},
			{TSharkTuple{Collection: []ISharkType{TSharkAny{}}}, TSharkTuple{Collection: []ISharkType{TSharkAny{}, TSharkBool{}}}, false},
			{TSharkTuple{Collection: []ISharkType{TSharkString{}, TSharkI64{}}}, TSharkTuple{Collection: []ISharkType{TSharkString{}, TSharkI64{}}}, true},
			{TSharkTuple{Collection: []ISharkType{TSharkString{}, TSharkI64{}}}, TSharkTuple{Collection: []ISharkType{TSharkString{}, TSharkAny{}}}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkTuple{}, "tuple<>"},
			{TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, "tuple<i64>"},
			{TSharkTuple{Collection: []ISharkType{TSharkAny{}}}, "tuple<any>"},
			{TSharkTuple{Collection: []ISharkType{TSharkAny{}, TSharkBool{}}}, "tuple<any,bool>"},
			{TSharkTuple{Collection: []ISharkType{TSharkAny{}, TSharkBool{}, TSharkI64{}}}, "tuple<any,bool,i64>"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})

	t.Run("should validate collection interface type", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkCollection{}, TSharkCollection{}, true},
			{TSharkCollection{}, TSharkAny{}, false},
			{TSharkCollection{}, TSharkI64{}, false},
			{TSharkCollection{Collection: []ISharkType{TSharkSpread{Type: TSharkAny{}}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, true},
			{TSharkCollection{Collection: []ISharkType{TSharkSpread{Type: TSharkAny{}}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}, TSharkString{}}}, true},
			{TSharkCollection{Collection: []ISharkType{TSharkI64{}}}, TSharkCollection{Collection: []ISharkType{TSharkI64{}}}, true},
			{TSharkCollection{}, TSharkArray{}, true},
			{TSharkCollection{}, TSharkTuple{}, true},
			{TSharkCollection{}, TSharkString{}, true},
			{TSharkCollection{Collection: []ISharkType{TSharkI64{}}}, TSharkArray{}, false},
			{TSharkCollection{Collection: []ISharkType{TSharkI64{}}}, TSharkArray{Collection: TSharkI64{}}, true},
			{TSharkCollection{Collection: []ISharkType{TSharkI64{}}}, TSharkTuple{}, false},
			{TSharkCollection{Collection: []ISharkType{TSharkI64{}}}, TSharkTuple{Collection: []ISharkType{TSharkI64{}}}, true},
			{TSharkArray{}, TSharkCollection{}, false},
			{TSharkTuple{}, TSharkCollection{}, false},
			{TSharkString{}, TSharkCollection{}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkCollection{}, "collection<>"},
			{TSharkCollection{Collection: []ISharkType{TSharkI64{}}}, "collection<i64>"},
			{TSharkCollection{Collection: []ISharkType{TSharkAny{}}}, "collection<any>"},
			{TSharkCollection{Collection: []ISharkType{TSharkAny{}, TSharkBool{}}}, "collection<any,bool>"},
			{TSharkCollection{Collection: []ISharkType{TSharkAny{}, TSharkBool{}, TSharkI64{}}}, "collection<any,bool,i64>"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})
}

func TestFunctionTypes(t *testing.T) {
	t.Run("should validate type func", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkFuncType{}, TSharkFuncType{}, true},
			{TSharkFuncType{}, TSharkAny{}, false},
			{TSharkFuncType{}, TSharkI64{}, false},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}, TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}, true},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}, TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}, true},
			{TSharkFuncType{ReturnT: TSharkString{}}, TSharkFuncType{ReturnT: TSharkString{}}, true},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}, TSharkFuncType{}, true},
			{TSharkFuncType{}, TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}, false},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}, TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkBool{}}, ReturnT: TSharkString{}}, false},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}, TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkBool{}}, false},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkBool{}}, ReturnT: TSharkString{}}, TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}, TSharkError{}}, ReturnT: TSharkString{}}, false},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkFuncType{}, "func<()>"},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}, "func<(i64)>"},
			{TSharkFuncType{ReturnT: TSharkString{}}, "func<()->string>"},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}, "func<(i64)->string>"},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}}, ReturnT: TSharkAny{}}, "func<(any)->any>"},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}, TSharkBool{}}, ReturnT: TSharkAny{}}, "func<(any,bool)->any>"},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}, TSharkBool{}, TSharkI64{}}, ReturnT: TSharkAny{}}, "func<(any,bool,i64)->any>"},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}, TSharkBool{}, TSharkSpread{Type: TSharkI64{}}}, ReturnT: TSharkAny{}}, "func<(any,bool,...i64)->any>"},
			{TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}, TSharkBool{}, TSharkOptional{Type: TSharkI64{}}}, ReturnT: TSharkAny{}}, "func<(any,bool,i64?)->any>"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})
}

func TestVariadicTypes(t *testing.T) {
	t.Run("should validate type variadic", func(t *testing.T) {
		tests_matching := []struct {
			givenType ISharkType
			otherType ISharkType
			expected  bool
		}{
			{TSharkVariadic{}, TSharkVariadic{}, true},
			{TSharkVariadic{}, TSharkAny{}, true},
			{TSharkVariadic{}, TSharkI64{}, true},
			{TSharkVariadic{Enclosed: TSharkI64{}}, TSharkVariadic{Enclosed: TSharkI64{}}, true},
			{TSharkVariadic{Enclosed: TSharkI64{}}, TSharkVariadic{Enclosed: TSharkAny{}}, true},
			{TSharkVariadic{Enclosed: TSharkI64{}}, TSharkVariadic{}, true},
		}

		for _, test := range tests_matching {
			validateTypeMatching(t, test.givenType, test.otherType, test.expected)
		}

		tests_rep := []struct {
			givenType ISharkType
			stringRep string
		}{
			{TSharkVariadic{}, "T"},
			{TSharkVariadic{Enclosed: TSharkI64{}}, "T"},
		}

		for _, test := range tests_rep {
			validateTypeStringRepresentation(t, test.givenType, test.stringRep)
		}
	})
}

func TestClosureTypes(t *testing.T) {
	t.Run("should validate type closure", func(t *testing.T) {
		t.Run("should validate type closure", func(t *testing.T) {
			tests_matching := []struct {
				givenType ISharkType
				otherType ISharkType
				expected  bool
			}{
				{TSharkClosure{}, TSharkClosure{}, true},
				{TSharkClosure{}, TSharkAny{}, false},
				{TSharkClosure{}, TSharkI64{}, false},
				{TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}}, TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}}, true},
				{TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkString{}}}, TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkString{}}}, true},
				{TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}}, TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}}, true},
			}

			for _, test := range tests_matching {
				validateTypeMatching(t, test.givenType, test.otherType, test.expected)
			}

			tests_rep := []struct {
				givenType ISharkType
				stringRep string
			}{
				{TSharkClosure{}, "closure<>"},
				{TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}}}, "func<(i64)>"},
				{TSharkClosure{FuncType: TSharkFuncType{ReturnT: TSharkString{}}}, "func<()->string>"},
				{TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkI64{}}, ReturnT: TSharkString{}}}, "func<(i64)->string>"},
				{TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}}, ReturnT: TSharkAny{}}}, "func<(any)->any>"},
				{TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}, TSharkBool{}}, ReturnT: TSharkAny{}}}, "func<(any,bool)->any>"},
				{TSharkClosure{FuncType: TSharkFuncType{ArgsList: []ISharkType{TSharkAny{}, TSharkBool{}, TSharkI64{}}, ReturnT: TSharkAny{}}}, "func<(any,bool,i64)->any>"},
			}

			for _, test := range tests_rep {
				validateTypeStringRepresentation(t, test.givenType, test.stringRep)
			}
		})
	})
}

func validateTypeMatching(t *testing.T, givenType ISharkType, otherType ISharkType, expected bool) {
	t.Helper()

	if givenType.Is(otherType) != expected {
		t.Errorf("Type matching from %s to %s should be %v", givenType.SharkTypeString(), otherType.SharkTypeString(), expected)
	}
}

func validateTypeStringRepresentation(t *testing.T, givenType ISharkType, expected string) {
	t.Helper()

	if givenType.SharkTypeString() != expected {
		t.Errorf("Expected SharkTypeString() to return %v, got %v", expected, givenType.SharkTypeString())
	}
}
