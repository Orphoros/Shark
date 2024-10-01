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
	var outName string
	var compression string
	emitInstructionSet := false

	flaggy.SetName("sharkc")
	flaggy.SetDescription("The Shark programming language compiler")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	// flaggy.DefaultParser.AdditionalHelpAppend = "A subcommand is required"
	// flaggy.DefaultParser.AdditionalHelpPrepend = "Shark can interpret and execute SharkLang code."

	flaggy.String(&outName, "o", "out", "The output file name")
	flaggy.String(&compression, "z", "compression", "The compression algorithm to use (brotli, none)")
	flaggy.Bool(&emitInstructionSet, "e", "emit", "Emit the instruction set")

	flaggy.AddPositionalValue(&file, "file", 1, true, "The Shark file to compile")
	flaggy.Parse()

	serializer.RegisterTypes()

	compile(file, outName, emitInstructionSet, compression)
}

func compile(file, outName string, emitInstructionSet bool, compression string) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", file), err.Error(), 1)
	}
	f := internal.ReadFile(absPath)
	//TODO: accept cnf from file
	vmConf := vm.NewDefaultConf()
	sharkEmitter := emitter.New(&absPath, os.Stdout, &vmConf)
	fileContents := string(f)
	if bc := sharkEmitter.Compile(&fileContents); bc != nil {
		fileName := internal.GetFileName(absPath) + ".egg"
		if outName != "" {
			fileName = internal.GetFileName(outName) + ".egg"
		}
		var bcType bytecode.BytecodeType
		switch compression {
		case "brotli":
			bcType = bytecode.BcTypeCompressedBrotli
		case "none":
			bcType = bytecode.BcTypeNormal
		default:
			bcType = bytecode.BcTypeNormal
		}
		bytes, err := bc.ToBytes(bcType, bytecode.BcVersionOnos1)
		if err != nil {
			exception.PrintExitMsgCtx("Could not convert bytecode to bytes", err.Error(), 1)
		}
		internal.WriteFile(fileName, bytes)

		if emitInstructionSet {
			fmt.Println(bc.Instructions.String())
		}
	}
}
