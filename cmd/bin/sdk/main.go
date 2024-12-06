package main

import (
	"shark/cmd"
	"shark/cmd/bin"
	"shark/config"
	"shark/serializer"

	"github.com/integrii/flaggy"
)

var Version string
var Build string
var Codename string

func main() {

	var file string
	var outName string
	var compression string
	var cnf string
	var emitInstructionSet bool

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
	compileCommand.String(&compression, "z", "compression", "The compression algorithm to use (brotli, none)")
	compileCommand.Bool(&emitInstructionSet, "e", "emit", "Emit the instruction set")

	execCommand := flaggy.NewSubcommand("exec")
	execCommand.Description = "Execute a SharkLang bytecode file"
	execCommand.AddPositionalValue(&file, "file", 1, true, "The bytecode file")

	decompileCommand := flaggy.NewSubcommand("decompile")
	decompileCommand.Description = "Decompiles a shark binary"
	decompileCommand.AddPositionalValue(&file, "file", 1, true, "The bytecode file")

	metaViewCommand := flaggy.NewSubcommand("meta")
	metaViewCommand.Description = "View the metadata of a SharkLang bytecode file"
	metaViewCommand.AddPositionalValue(&file, "file", 1, true, "The bytecode file")

	genConf := flaggy.NewSubcommand("genconf")
	genConf.Description = "Generate a default configuration file to the current directory"

	flaggy.AttachSubcommand(runCommand, 1)
	flaggy.AttachSubcommand(compileCommand, 1)
	flaggy.AttachSubcommand(execCommand, 1)
	flaggy.AttachSubcommand(decompileCommand, 1)
	flaggy.AttachSubcommand(metaViewCommand, 1)
	flaggy.AttachSubcommand(genConf, 1)
	flaggy.Parse()

	argConfig, err := config.LocateConfig(&cnf, &file)
	if err != nil {
		flaggy.ShowHelpAndExit("Config Error: " + err.Error())
	}

	serializer.RegisterTypes()

	if execCommand.Used {
		cmd.ExecuteSharkBinaryFile(file, argConfig)
	} else if compileCommand.Used {
		cmd.CompileSharkCodeFile(file, outName, compression, emitInstructionSet, argConfig)
	} else if runCommand.Used {
		cmd.ExecuteSharkCodeFile(file, argConfig)
	} else if decompileCommand.Used {
		cmd.DecompileSharkBinaryFile(file)
	} else if metaViewCommand.Used {
		cmd.ShowOjbMeta(file)
	} else if genConf.Used {
		cmd.GenerateDefaultConfig()
	} else {
		flaggy.ShowHelpAndExit("Error: No subcommand provided")
	}
}
