package main

import "C"
import (
	"os"
	"shark/emitter"
)

//export execute
func execute(source, code string) {
	sharkEmitter := emitter.New(&source, os.Stdout)
	sharkEmitter.Interpret(code)
}

func main() {}
