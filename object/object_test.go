package object

import (
	"testing"
)

func TestStringHashKey(t *testing.T) {
	t.Run("should return the same hash key for the same string", func(t *testing.T) {
		hello1 := &String{Value: "Hello World"}
		hello2 := &String{Value: "Hello World"}
		diff1 := &String{Value: "My name is johnny"}
		diff2 := &String{Value: "My name is johnny"}

		if hello1.HashKey() != hello2.HashKey() {
			t.Errorf("strings with same content have different hash keys")
		}

		if diff1.HashKey() != diff2.HashKey() {
			t.Errorf("strings with same content have different hash keys")
		}

		if hello1.HashKey() == diff1.HashKey() {
			t.Errorf("strings with different content have same hash keys")
		}
	})
}

func TestObjects(t *testing.T) {
	t.Run("should return the correct array object", func(t *testing.T) {
		arrayObj := &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}} // [1, 2, 3]

		if arrayObj.Type() != ARRAY_OBJ {
			t.Errorf("wrong type. expected=%s, got=%s", ARRAY_OBJ, arrayObj.Type())
		}

		if len(arrayObj.Elements) != 3 {
			t.Errorf("wrong number of elements. expected=%d, got=%d", 3, len(arrayObj.Elements))
		}

		if arrayObj.Elements[0].(*Integer).Value != 1 {
			t.Errorf("wrong value for first element. expected=%d, got=%d", 1, arrayObj.Elements[0].(*Integer).Value)
		}

		if arrayObj.Elements[1].(*Integer).Value != 2 {
			t.Errorf("wrong value for second element. expected=%d, got=%d", 2, arrayObj.Elements[1].(*Integer).Value)
		}

		if arrayObj.Elements[2].(*Integer).Value != 3 {
			t.Errorf("wrong value for third element. expected=%d, got=%d", 3, arrayObj.Elements[2].(*Integer).Value)
		}

		if arrayObj.Inspect() != "[1, 2, 3]" {
			t.Errorf("wrong inspect. expected=%s, got=%s", "[1, 2, 3]", arrayObj.Inspect())
		}
	})

	t.Run("should return the correct boolean object", func(t *testing.T) {
		boolObj := &Boolean{Value: true}

		if boolObj.Type() != BOOLEAN_OBJ {
			t.Errorf("wrong type. expected=%s, got=%s", BOOLEAN_OBJ, boolObj.Type())
		}

		if boolObj.Value != true {
			t.Errorf("wrong value. expected=%t, got=%t", true, boolObj.Value)
		}

		if boolObj.Inspect() != "true" {
			t.Errorf("wrong inspect. expected=%s, got=%s", "true", boolObj.Inspect())
		}
	})

	t.Run("should return the correct error object", func(t *testing.T) {
		errorObj := &Error{Message: "something went wrong"}

		if errorObj.Type() != ERROR_OBJ {
			t.Errorf("wrong type. expected=%s, got=%s", ERROR_OBJ, errorObj.Type())
		}

		if errorObj.Message != "something went wrong" {
			t.Errorf("wrong message. expected=%s, got=%s", "something went wrong", errorObj.Message)
		}

		if errorObj.Inspect() != "ERROR: something went wrong" {
			t.Errorf("wrong inspect. expected=%s, got=%s", "ERROR: something went wrong", errorObj.Inspect())
		}
	})

	t.Run("should return the correct hash object", func(t *testing.T) {
		hashObj := &Hash{Pairs: map[HashKey]HashPair{
			(&String{Value: "one"}).HashKey(): {Key: &String{Value: "one"}, Value: &Integer{Value: 1}},
			(&String{Value: "two"}).HashKey(): {Key: &String{Value: "two"}, Value: &Integer{Value: 2}},
		}}

		if hashObj.Type() != HASH_OBJ {
			t.Errorf("wrong type. expected=%s, got=%s", HASH_OBJ, hashObj.Type())
		}

		if len(hashObj.Pairs) != 2 {
			t.Errorf("wrong number of pairs. expected=%d, got=%d", 2, len(hashObj.Pairs))
		}

		if hashObj.Inspect() != `{one: 1, two: 2}` {
			t.Errorf("wrong inspect. expected=%s, got=%s", `{one: 1, two: 2}`, hashObj.Inspect())
		}
	})

	t.Run("should return the correct int object", func(t *testing.T) {
		intObj := &Integer{Value: 1}

		if intObj.Type() != INTEGER_OBJ {
			t.Errorf("wrong type. expected=%s, got=%s", INTEGER_OBJ, intObj.Type())
		}

		if intObj.Value != 1 {
			t.Errorf("wrong value. expected=%d, got=%d", 1, intObj.Value)
		}

		if intObj.Inspect() != "1" {
			t.Errorf("wrong inspect. expected=%s, got=%s", "1", intObj.Inspect())
		}
	})

	t.Run("should return the correct string object", func(t *testing.T) {
		strObj := &String{Value: "Hello World"}

		if strObj.Type() != STRING_OBJ {
			t.Errorf("wrong type. expected=%s, got=%s", STRING_OBJ, strObj.Type())
		}

		if strObj.Value != "Hello World" {
			t.Errorf("wrong value. expected=%s, got=%s", "Hello World", strObj.Value)
		}

		if strObj.Inspect() != "Hello World" {
			t.Errorf("wrong inspect. expected=%s, got=%s", "Hello World", strObj.Inspect())
		}
	})
}
