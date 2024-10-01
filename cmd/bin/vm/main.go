package main

import (
	"fmt"
	"os"
	"path/filepath"
	"shark/bytecode"
	"shark/cmd/bin"
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
	gobFile := internal.ReadFile(absPath)
	bc, err := bytecode.FromBytes(gobFile)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read bytecode file '%s'", file), err.Error(), 1)
	}
	//TODO: accept cnf from file
	vmConf := vm.NewDefaultConf()
	sharkEmitter := emitter.New(&absPath, os.Stdout, &vmConf)
	sharkEmitter.Exec(bc)
}
