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
