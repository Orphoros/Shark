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
	var logLevel string
	emitInstructionSet := false

	flaggy.SetName("sharkc")
	flaggy.SetDescription("The Shark programming language compiler")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	flaggy.String(&outName, "o", "out", "The output file name")
	flaggy.String(&compression, "z", "compression", "The compression algorithm to use (brotli, none)")
	flaggy.Bool(&emitInstructionSet, "e", "emit", "Emit the instruction set")
	flaggy.String(&cnf, "c", "config", "The configuration file")
	flaggy.String(&logLevel, "l", "loglevel", "The log level (trace, debug, info, warn, error, fatal, panic)")

	flaggy.AddPositionalValue(&file, "file", 1, true, "The Shark file to compile")
	flaggy.Parse()

	cmd.RegisterLogger(logLevel)

	argConfig, err := config.LocateConfig(&cnf, &file)
	if err != nil {
		flaggy.ShowHelpAndExit("Error: " + err.Error())
	}

	serializer.RegisterTypes()

	cmd.CompileSharkCodeFile(file, outName, compression, emitInstructionSet, argConfig)
}
