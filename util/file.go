package util

import (
	"path/filepath"

	chars "github.com/rivo/uniseg"
)

func GetFileName(fileName string) string {
	return fileName[:chars.GraphemeClusterCount(fileName)-chars.GraphemeClusterCount(filepath.Ext(fileName))]
}
