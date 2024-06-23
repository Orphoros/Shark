package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
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
	var outName string
	emitInstructionSet := false

	flaggy.SetName("sharkc")
	flaggy.SetDescription("The Shark programming language compiler")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	// flaggy.DefaultParser.AdditionalHelpAppend = "A subcommand is required"
	// flaggy.DefaultParser.AdditionalHelpPrepend = "Shark can interpret and execute SharkLang code."

	flaggy.String(&outName, "o", "out", "The output file name")
	flaggy.Bool(&emitInstructionSet, "e", "emit", "Emit the instruction set")

	flaggy.AddPositionalValue(&file, "file", 1, true, "The Shark file to compile")
	flaggy.Parse()

	serializer.RegisterTypes()

	compile(file, outName, emitInstructionSet)
}

func compile(file, outName string, emitInstructionSet bool) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", file), err.Error(), 1)
	}
	f := internal.ReadFile(absPath)
	vmConf := vm.NewDefaultConf()
	sharkEmitter := emitter.New(&absPath, os.Stdout, &vmConf)
	fileContents := string(f)
	if bytecode := sharkEmitter.Compile(&fileContents); bytecode != nil {
		fileName := internal.GetFileName(absPath) + ".egg"
		if outName != "" {
			fileName = internal.GetFileName(outName) + ".egg"
		}
		gobFile, err := os.Create(fileName)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not create file '%s'", file), err.Error(), 1)
		}
		defer func(gobFile *os.File) {
			err := gobFile.Close()
			if err != nil {
				exception.PrintExitMsgCtx(fmt.Sprintf("Could not close file '%s'", file), err.Error(), 1)
			}
		}(gobFile)
		if err = gob.NewEncoder(gobFile).Encode(bytecode); err != nil {
			exception.PrintExitMsgCtx("Compiler bytecode could not be serialized", err.Error(), 1)
		}

		if emitInstructionSet {
			// emit instruction set to file
			instructionSetFileName := internal.GetFileName(absPath) + ".scc"
			if outName != "" {
				instructionSetFileName = internal.GetFileName(outName) + ".scc"
			}

			instructionSetFile, err := os.Create(instructionSetFileName)
			if err != nil {
				exception.PrintExitMsgCtx(fmt.Sprintf("Could not create file '%s'", file), err.Error(), 1)
			}
			defer func(instructionSetFile *os.File) {
				err := instructionSetFile.Close()
				if err != nil {
					exception.PrintExitMsgCtx(fmt.Sprintf("Could not close file '%s'", file), err.Error(), 1)
				}
			}(instructionSetFile)

			emitter.EmitInstructionsTable(bytecode, instructionSetFile)
		}
	}
}
