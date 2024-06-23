package main

import "C"
import (
	"os"
	"shark/emitter"
	"shark/vm"
)

//export execute
func execute(source, code string) {
	vmConf := vm.NewDefaultConf()
	sharkEmitter := emitter.New(&source, os.Stdout, &vmConf)
	sharkEmitter.Interpret(code)
}

func main() {}
