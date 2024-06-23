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

	"github.com/integrii/flaggy"
)

var Version string
var Build string
var Codename string

func main() {

	var file string
	var outName string
	var cnf string

	flaggy.SetName("shark")
	flaggy.SetDescription("The Shark programming language")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpAppend = "A subcommand is required"
	flaggy.DefaultParser.AdditionalHelpPrepend = "SDK for the Shark programming language."

	flaggy.String(&cnf, "c", "config", "The configuration file")

	runCommand := flaggy.NewSubcommand("run")
	runCommand.Description = "Interpret a SharkLang source code file"
	runCommand.AddPositionalValue(&file, "file", 1, true, "The file to interpret")

	compileCommand := flaggy.NewSubcommand("compile")
	compileCommand.Description = "Compile a SharkLang source code file into bytecode"
	compileCommand.AddPositionalValue(&file, "file", 1, true, "The file to compile")
	compileCommand.String(&outName, "o", "out", "The output file name")

	execCommand := flaggy.NewSubcommand("exec")
	execCommand.Description = "Execute a SharkLang bytecode file"
	execCommand.AddPositionalValue(&file, "file", 1, true, "The bytecode file")

	decompileCommand := flaggy.NewSubcommand("decompile")
	decompileCommand.Description = "Decompiles a shark binary"
	decompileCommand.AddPositionalValue(&file, "file", 1, true, "The bytecode file")

	flaggy.AttachSubcommand(runCommand, 1)
	flaggy.AttachSubcommand(compileCommand, 1)
	flaggy.AttachSubcommand(execCommand, 1)
	flaggy.AttachSubcommand(decompileCommand, 1)
	flaggy.Parse()

	if cnf == "" {
		dir, err := os.Getwd()
		if err != nil {
			exception.PrintExitMsgCtx("Could not get the current working directory", err.Error(), 1)
		}
		if internal.IsFileExists(filepath.Join(dir, "Shark.toml")) {
			cnf = filepath.Join(dir, "Shark.toml")
		} else if internal.IsFileExists(filepath.Join(filepath.Dir(file), "Shark.toml")) {
			cnf = filepath.Join(filepath.Dir(file), "Shark.toml")
		}
	}

	var argConfig internal.Config

	if cnf == "" {
		argConfig = internal.GetDefaultConfig()
	} else {
		argConfig = internal.GetConfigFromFile(cnf)
	}

	serializer.RegisterTypes()

	if execCommand.Used {
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
		sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.OrpVM)
		sharkEmitter.Exec(bytecode)
	} else if compileCommand.Used {
		absPath, err := filepath.Abs(file)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", file), err.Error(), 1)
		}
		f := internal.ReadFile(absPath)
		sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.OrpVM)
		file := string(f)
		if bytecode := sharkEmitter.Compile(&file); bytecode != nil {
			fileName := internal.GetFileName(absPath) + ".egg"
			if outName != "" {
				fileName = outName
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
		}
	} else if runCommand.Used {
		absPath, err := filepath.Abs(file)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", file), err.Error(), 1)
		}
		f := internal.ReadFile(absPath)
		sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.OrpVM)
		sharkEmitter.Interpret(string(f))
	} else if decompileCommand.Used {
		gobFile := internal.OpenFile(file)
		defer func(gobFile *os.File) {
			err := gobFile.Close()
			if err != nil {
				exception.PrintExitMsgCtx(fmt.Sprintf("Could not close file '%s'", file), err.Error(), 1)
			}
		}(gobFile)
		var bytecode *compiler.Bytecode
		if err := gob.NewDecoder(gobFile).Decode(&bytecode); err != nil {
			exception.PrintExitMsg("Binary file is not compatible", 1)
		}

		emitter.EmitInstructionsTable(bytecode, os.Stdout)
	} else {
		flaggy.ShowHelpAndExit("Error: No subcommand provided")
	}
}
