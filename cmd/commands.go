package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"shark/bytecode"
	"shark/config"
	"shark/emitter"
	"shark/exception"
	"shark/internal"
)

func ExecuteSharkCodeFile(path string, argConfig *config.Config) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	f := internal.ReadFile(absPath)
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.OrpVM)
	sharkEmitter.Interpret(string(f))
}

func DecompileSharkBinaryFile(path string) {
	gobFile := internal.ReadFile(path)
	bc, err := bytecode.FromBytes(gobFile)
	if err != nil {
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}

	io.WriteString(os.Stdout, bc.ToString())
}

func CompileSharkCodeFile(path, outName, compression string, emitInstructionSet bool, argConfig *config.Config) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	f := internal.ReadFile(absPath)
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.OrpVM)
	file := string(f)
	if bc := sharkEmitter.Compile(&file); bc != nil {
		fileName := internal.GetFileName(absPath) + ".egg"
		if outName != "" {
			fileName = outName
		}
		var bcType bytecode.BytecodeType
		switch compression {
		case "brotli":
			bcType = bytecode.BcTypeCompressedBrotli
		case "none":
			bcType = bytecode.BcTypeNormal
		default:
			bcType = bytecode.BcTypeNormal
		}
		bytes, err := bc.ToBytes(bcType, bytecode.BcVersionOnos1)
		if err != nil {
			exception.PrintExitMsgCtx("Could not convert bytecode to bytes", err.Error(), 1)
		}
		internal.WriteFile(fileName, bytes)
	}
}

func ExecuteSharkBinaryFile(path string, argConfig *config.Config) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	gobFile := internal.ReadFile(absPath)
	bc, err := bytecode.FromBytes(gobFile)
	if err != nil {
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.OrpVM)
	sharkEmitter.Exec(bc)
}
