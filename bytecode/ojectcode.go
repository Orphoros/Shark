package bytecode

import (
	"encoding/binary"
	"fmt"
)

type ObjCode []byte

func (o *ObjCode) isMagicNumberMatch() bool {
	magicNumber := []byte{0x6e, 0x65, 0x78} // "nex"
	if len(*o) < len(magicNumber) {
		return false
	}
	return string((*o)[:len(magicNumber)]) == string(magicNumber)
}

func (o *ObjCode) CompressionType() (Type, error) {
	if !o.isMagicNumberMatch() {
		return 0, fmt.Errorf("magic number mismatch")
	}

	if len(*o) < 5 {
		return 0, fmt.Errorf("object code is too short to contain a compression type")
	}
	compressionType := Type((*o)[4])
	if compressionType < minType || compressionType > maxType {
		return 0, fmt.Errorf("invalid compression type: %d", compressionType)
	}

	return compressionType, nil
}

func (o *ObjCode) Version() (Version, error) {
	if !o.isMagicNumberMatch() {
		return 0, fmt.Errorf("magic number mismatch")
	}

	if len(*o) < 4 {
		return 0, fmt.Errorf("object code is too short to contain a version")
	}
	version := Version((*o)[3])
	if version != BcVersionOnos1 {
		return 0, fmt.Errorf("invalid version: %d", version)
	}

	return version, nil
}
func (o *ObjCode) InstructionLength() (uint64, error) {
	if !o.isMagicNumberMatch() {
		return 0, fmt.Errorf("magic number mismatch")
	}

	if len(*o) < 13 {
		return 0, fmt.Errorf("object code is too short to contain an instruction length")
	}

	return binary.BigEndian.Uint64((*o)[5:13]), nil
}
