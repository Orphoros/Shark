package bytecode

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"shark/code"
	"shark/object"
)

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type Type byte

type Version byte

const (
	BcTypeNormal Type = iota
	BcTypeCompressedBrotli
)

var magicNumber []byte = []byte{0x6F, 0x62, 0x63} // "onbc"

const (
	BcVersionOnos1 Version = iota
)

func (b *Bytecode) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	if err := encoder.Encode(b.Instructions); err != nil {
		return nil, err
	}
	if err := encoder.Encode(b.Constants); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (b *Bytecode) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&b.Instructions); err != nil {
		return err
	}
	if err := decoder.Decode(&b.Constants); err != nil {
		return err
	}
	return nil
}

func FromBytes(data []byte) (*Bytecode, error) {
	r := bytes.NewReader(data)

	mn := make([]byte, len(magicNumber))
	if err := binary.Read(r, binary.BigEndian, &mn); err != nil {
		return nil, err
	}

	if !bytes.Equal(mn, magicNumber) {
		return nil, fmt.Errorf("magic number mismatch")
	}

	var version Version
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return nil, err
	}

	var bytecodeType Type
	if err := binary.Read(r, binary.BigEndian, &bytecodeType); err != nil {
		return nil, err
	}

	var bytecodeLength int64
	if err := binary.Read(r, binary.BigEndian, &bytecodeLength); err != nil {
		return nil, err
	}

	bytecodeData := make([]byte, bytecodeLength)
	if _, err := r.Read(bytecodeData); err != nil {
		return nil, err
	}

	switch bytecodeType {
	case BcTypeNormal:
		break
	case BcTypeCompressedBrotli:
		var err error
		bytecodeData, err = decompressBrotli(bytecodeData)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown bytecode type: %d", bytecodeType)
	}

	b := &Bytecode{}
	decoder := gob.NewDecoder(bytes.NewReader(bytecodeData))
	if err := decoder.Decode(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bytecode) ToBytes(bytecodeType Type, version Version) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write magic number
	if err := binary.Write(buf, binary.BigEndian, magicNumber); err != nil {
		return nil, err
	}

	// Write bytecode version
	if err := binary.Write(buf, binary.BigEndian, version); err != nil {
		return nil, err
	}

	// Write bytecode type
	if err := binary.Write(buf, binary.BigEndian, bytecodeType); err != nil {
		return nil, err
	}

	gobEncoded := new(bytes.Buffer)

	encoder := gob.NewEncoder(gobEncoded)
	if err := encoder.Encode(b); err != nil {
		return nil, err
	}

	switch bytecodeType {
	case BcTypeNormal:
		break
	case BcTypeCompressedBrotli:
		var err error
		gobEncodedBytes := gobEncoded.Bytes()
		gobEncodedBytes, err = compressBrotli(gobEncodedBytes)
		if err != nil {
			return nil, err
		}
		gobEncoded = bytes.NewBuffer(gobEncodedBytes)
	default:
		return nil, fmt.Errorf("unknown bytecode type: %d", bytecodeType)
	}

	// Write bytecode length
	if err := binary.Write(buf, binary.BigEndian, int64(gobEncoded.Len())); err != nil {
		return nil, err
	}

	// Write bytecode
	if _, err := buf.Write(gobEncoded.Bytes()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b *Bytecode) ToString() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "@main:\n")

	ins := b.Instructions

	txt := ins.String()

	lines := bytes.Split([]byte(txt), []byte("\n"))

	for i, line := range lines {
		if i == len(lines)-1 {
			break
		}
		fmt.Fprintf(&buf, "| %s\n", line)
	}

	fmt.Fprintf(&buf, "\n@constants:\n")

	for i, constant := range b.Constants {
		switch constant.Type() {
		case object.COMPILED_FUNCTION_OBJ:
			fmt.Fprintf(&buf, "| %04d FUNC {\n", i)
			cf := constant.(*object.CompiledFunction)
			txt := cf.Instructions.String()
			lines := bytes.Split([]byte(txt), []byte("\n"))
			for i, line := range lines {
				if i == len(lines)-1 {
					break
				}
				fmt.Fprintf(&buf, "| \t%s\n", line)
			}
			fmt.Fprintf(&buf, "| }\n")
		case object.STRING_OBJ:
			fmt.Fprintf(&buf, "| %04d %s: \"%s\"\n", i, constant.Type(), constant.Inspect())
		default:
			fmt.Fprintf(&buf, "| %04d %s: %s\n", i, constant.Type(), constant.Inspect())
		}
	}

	return buf.String()
}
