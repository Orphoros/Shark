package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"shark/cmd/bin"
	"shark/compiler"
	"shark/emitter"
	"shark/exception"
	"shark/internal"
	"shark/serializer"
	"shark/vm"

	"github.com/integrii/flaggy"
)

var Version string
var Build string
var Codename string

func main() {
	var file string

	flaggy.SetName("orpvm")
	flaggy.SetDescription("The Orphoros Virtual Machine")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	flaggy.AddPositionalValue(&file, "file", 1, true, "The Shark bytecode file (.egg) to execute")
	flaggy.Parse()

	serializer.RegisterTypes()

	execute(file)
}

func execute(file string) {
	serializer.RegisterTypes()

	absPath, err := filepath.Abs(file)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", file), err.Error(), 1)
	}
	gobFile := internal.OpenFile(absPath)
	defer func(gobFile *os.File) {
		err := gobFile.Close()
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not close file '%s'", file), err.Error(), 1)
		}
	}(gobFile)
	decoder := gob.NewDecoder(gobFile)
	var bytecode *compiler.Bytecode
	err = decoder.Decode(&bytecode)
	if err != nil {
		exception.PrintExitMsg("Binary file is not compatible", 1)
	}
	vmConf := vm.NewDefaultConf()
	sharkEmitter := emitter.New(&absPath, os.Stdout, &vmConf)
	sharkEmitter.Exec(bytecode)
}
