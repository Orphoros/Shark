package main

import (
	"shark/cmd"
	"shark/cmd/bin"
	"shark/serializer"

	"github.com/integrii/flaggy"
)

var Version string
var Build string
var Codename string

func main() {
	var file string
	var cnf string

	flaggy.SetName("orpvm")
	flaggy.SetDescription("The Orphoros Virtual Machine")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	flaggy.AddPositionalValue(&file, "file", 1, true, "The Shark bytecode file (.egg) to execute")
	flaggy.String(&cnf, "c", "config", "The configuration file")
	flaggy.Parse()

	argConfig := bin.LocateConfigFile(cnf, file)

	serializer.RegisterTypes()

	cmd.ExecuteSharkCodeFile(file, &argConfig)
}
