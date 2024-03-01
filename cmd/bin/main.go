package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"shark/compiler"
	"shark/emitter"
	"shark/exception"
	"shark/serializer"
	"shark/util"
	"time"

	"github.com/integrii/flaggy"
)

var Version string
var Build string

func main() {

	var file string
	var outName string
	measureDuration := false

	startTime := time.Now()

	flaggy.SetName("Shark")
	flaggy.SetDescription("The Shark programming language")
	var curVersion = formatVersion()
	flaggy.SetVersion(curVersion)
	flaggy.Bool(&measureDuration, "d", "duration", "Measure the duration of the execution")

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpAppend = "A subcommand is required"
	flaggy.DefaultParser.AdditionalHelpPrepend = "Shark can interpret and execute SharkLang code."

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

	serializer.RegisterTypes()

	if execCommand.Used {
		absPath, err := filepath.Abs(file)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", file), err.Error(), 1)
		}
		gobFile, err := os.Open(absPath)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not open file '%s'", file), err.Error(), 1)
		}
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
		sharkEmitter := emitter.New(&absPath, os.Stdout)
		sharkEmitter.Exec(bytecode)
	} else if compileCommand.Used {
		absPath, err := filepath.Abs(file)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", file), err.Error(), 1)
		}
		f, err := os.ReadFile(absPath)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not open file '%s'", file), err.Error(), 1)
		}
		sharkEmitter := emitter.New(&absPath, os.Stdout)
		file := string(f)
		if bytecode := sharkEmitter.Compile(&file); bytecode != nil {
			fileName := util.GetFileName(absPath) + ".egg"
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
		f, err := os.ReadFile(absPath)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not open file '%s'", file), err.Error(), 1)
		}
		sharkEmitter := emitter.New(&absPath, os.Stdout)
		sharkEmitter.Interpret(string(f))
	} else if decompileCommand.Used {
		gobFile, err := os.Open(file)
		if err != nil {
			exception.PrintExitMsgCtx(fmt.Sprintf("Could not open file '%s'", file), err.Error(), 1)
		}
		defer func(gobFile *os.File) {
			err := gobFile.Close()
			if err != nil {
				exception.PrintExitMsgCtx(fmt.Sprintf("Could not close file '%s'", file), err.Error(), 1)
			}
		}(gobFile)
		var bytecode *compiler.Bytecode
		if err = gob.NewDecoder(gobFile).Decode(&bytecode); err != nil {
			exception.PrintExitMsg("Binary file is not compatible", 1)
		}

		emitter.EmitInstructionsTable(bytecode, os.Stdout)
	} else {
		flaggy.ShowHelpAndExit("Error: No subcommand provided")
	}

	if measureDuration {
		duration := time.Since(startTime)
		fmt.Printf("\t~ Execution took %s\n", duration)
	}
}

func formatVersion() string {
	var curVersion string
	if Version == "" {
		curVersion = "dev"
	} else {
		curVersion = Version
		if Build != "" {
			curVersion += " (" + Build + ")"
		}
	}
	curVersion += "\nCore: " + runtime.Version()
	return curVersion
}
