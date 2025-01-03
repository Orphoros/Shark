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
	var cnf string
	var logLevel string

	flaggy.SetName("nidum")
	flaggy.SetDescription("The Nidum Virtual Machine (NVM)")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	flaggy.AddPositionalValue(&file, "file", 1, true, "The Shark object file (.egg) to execute")
	flaggy.String(&cnf, "c", "config", "The configuration file")
	flaggy.Parse()

	cmd.RegisterLogger(logLevel)

	argConfig, err := config.LocateConfig(&cnf, &file)
	if err != nil {
		flaggy.ShowHelpAndExit("Error: " + err.Error())
	}

	serializer.RegisterTypes()

	cmd.ExecuteSharkCodeFile(file, argConfig)
}
