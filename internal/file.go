package internal

import (
	"net/url"
	"os"
	"path/filepath"
	"shark/exception"
	"unicode/utf8"
)

// GetFileName returns the file name without the path and extension.
func GetFileName(fileName string) string {
	return fileName[:utf8.RuneCountInString(fileName)-utf8.RuneCountInString(filepath.Ext(fileName))]
}

func GetFilePathFromURI(uri string) (*string, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}

	return &u.Path, nil
}

// IsFileEndsWith checks if the file name ends with the given extension.
func IsFileEndsWith(fileName string, ext string) bool {
	return filepath.Ext(fileName) == ext
}

func ReadFile(fileName string) ([]byte, error) {
	f, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func WriteFile(fileName string, data []byte) error {
	gobFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(gobFile *os.File) {
		if err := gobFile.Close(); err != nil {
			exception.PrintExitMsgCtx("Could not close file", err.Error(), 1)
		}
	}(gobFile)
	if _, err := gobFile.Write(data); err != nil {
		return err
	}

	return nil
}

func IsFileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
