package util

import (
	"path/filepath"
	"unicode/utf8"
)

// GetFileName returns the file name without the path and extension.
func GetFileName(fileName string) string {
	return fileName[:utf8.RuneCountInString(fileName)-utf8.RuneCountInString(filepath.Ext(fileName))]
}

// IsFileEndsWith checks if the file name ends with the given extension.
func IsFileEndsWith(fileName string, ext string) bool {
	return filepath.Ext(fileName) == ext
}
