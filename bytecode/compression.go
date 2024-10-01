package bytecode

import (
	"bytes"

	"github.com/andybalholm/brotli"
)

func compressBrotli(data []byte) ([]byte, error) {
	var encoder *brotli.Writer
	var tmpBuffer bytes.Buffer

	encoder = brotli.NewWriterLevel(
		&tmpBuffer,
		brotli.BestCompression,
	)

	_, err := encoder.Write(data)
	if err != nil {
		return nil, err
	}

	if err := encoder.Close(); err != nil {
		return nil, err
	}

	return tmpBuffer.Bytes(), nil
}

func decompressBrotli(data []byte) ([]byte, error) {
	var decoder *brotli.Reader
	var tmpBuffer bytes.Buffer

	decoder = brotli.NewReader(bytes.NewReader(data))

	_, err := tmpBuffer.ReadFrom(decoder)
	if err != nil {
		return nil, err
	}

	return tmpBuffer.Bytes(), nil
}
