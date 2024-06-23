package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"shark/exception"
	"shark/vm"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
)

// GetFileName returns the file name without the path and extension.
func GetFileName(fileName string) string {
	return fileName[:utf8.RuneCountInString(fileName)-utf8.RuneCountInString(filepath.Ext(fileName))]
}

// IsFileEndsWith checks if the file name ends with the given extension.
func IsFileEndsWith(fileName string, ext string) bool {
	return filepath.Ext(fileName) == ext
}

func ReadFile(file string) []byte {
	f, err := os.ReadFile(file)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read contents of file '%s'", file), err.Error(), 1)
	}
	return f
}

func OpenFile(file string) *os.File {
	f, err := os.Open(file)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not open file '%s'", file), err.Error(), 1)
	}
	return f
}

func IsFileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func GetConfigFromFile(cnf string) Config {
	var conf Config

	file := ReadFile(cnf)

	if _, err := toml.Decode(string(file), &conf); err != nil {
		exception.PrintExitMsgCtx("Could not unmarshal config file", err.Error(), 1)
	}

	defVmConf := vm.NewDefaultConf()
	if conf.OrpVM.GlobalsSize == 0 {
		conf.OrpVM.GlobalsSize = defVmConf.GlobalsSize
	}
	if conf.OrpVM.StackSize == 0 {
		conf.OrpVM.StackSize = defVmConf.StackSize
	}
	if conf.OrpVM.MaxFrames == 0 {
		conf.OrpVM.MaxFrames = defVmConf.MaxFrames
	}

	return conf
}

func GetDefaultConfig() Config {
	return Config{
		OrpVM: vm.NewDefaultConf(),
	}
}
