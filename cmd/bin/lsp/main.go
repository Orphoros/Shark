package main

import (
	"shark/cmd/bin"
	"shark/lsp"

	"github.com/integrii/flaggy"
)

var Version string
var Build string
var Codename string

func main() {
	var port int = 59027

	flaggy.SetName("shark")
	flaggy.SetDescription("The Shark programming language")
	flaggy.SetVersion(bin.FormatVersion(Version, Build, Codename))

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "Language server for the Shark programming language."

	flaggy.Int(&port, "p", "port", "The port to listen on")

	flaggy.Parse()

	lsp.Start(port)
}
