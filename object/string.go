package object

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"shark/types"
)

type String struct {
	Value string
}

func (s *String) Inspect() string { return s.Value }

func (s *String) HashKey() HashKey {
	// TODO: Implement separate chaining to handle collisions
	h := fnv.New64a()
	if write, err := h.Write([]byte(s.Value)); write != len(s.Value) || err != nil {
		return HashKey{}
	}
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (s *String) Type() types.ISharkType { return types.TSharkString{} }

func (s *String) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(s.Value)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (s *String) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	return decoder.Decode(&s.Value)
}
