package main

import "C"
import (
	"os"
	"shark/config"
	"shark/emitter"
)

//export execute
func execute(source, code string) {
	vmConf := config.NewDefaultVmConf()
	sharkEmitter := emitter.New(&source, os.Stdout, &vmConf)
	sharkEmitter.Interpret(code)
}

func main() {}
