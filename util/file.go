package util

import (
	"path/filepath"
	"unicode/utf8"
)

func GetFileName(fileName string) string {
	return fileName[:utf8.RuneCountInString(fileName)-utf8.RuneCountInString(filepath.Ext(fileName))]
}
