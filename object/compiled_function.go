package object

import (
	"bytes"
	"encoding/gob"
	"shark/code"
	"shark/types"
)

type CompiledFunction struct {
	ObjType       types.ISharkType
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
	NumDefaults   int
}

func (cf *CompiledFunction) Inspect() string { return "CompiledFunction" }

func (cf *CompiledFunction) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	if err := encoder.Encode(cf.NumLocals); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cf.NumParameters); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cf.Instructions); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (cf *CompiledFunction) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&cf.NumLocals); err != nil {
		return err
	}
	if err := decoder.Decode(&cf.NumParameters); err != nil {
		return err
	}
	if err := decoder.Decode(&cf.Instructions); err != nil {
		return err
	}
	return nil
}

func (cf *CompiledFunction) Type() types.ISharkType { return cf.ObjType }
