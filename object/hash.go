package object

import (
	"bytes"
	"fmt"
	"shark/types"
	"strings"
)

type Hash struct {
	Pairs map[HashKey]HashPair
}

type HashKey struct {
	Type  types.ISharkType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hashable interface {
	// TODO: Cache the return values of HashKey() for performance
	HashKey() HashKey
}

func (h *Hash) Type() types.ISharkType {
	var keyType types.ISharkType
	var valueType types.ISharkType

	for _, pair := range h.Pairs {
		if pair.Key.Type() != nil {
			keyType = pair.Key.Type()
		}

		if pair.Value.Type() != nil {
			valueType = pair.Value.Type()
		}
	}

	if keyType == nil {
		keyType = types.TSharkAny{}
	}

	if valueType == nil {
		valueType = types.TSharkAny{}
	}

	return types.TSharkHashMap{Indexes: keyType, Collects: valueType}
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
