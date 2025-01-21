package object

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"shark/types"
)

type Int64 struct {
	Value int64
}

func (i *Int64) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Int64) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (i *Int64) Type() types.ISharkType { return types.TSharkI64{} }

func (i *Int64) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(i.Value)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (i *Int64) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	return decoder.Decode(&i.Value)
}
