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
	f, err := internal.ReadFile(absPath)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.NidumVM)
	sharkEmitter.Interpret(string(f))
}

func DecompileSharkBinaryFile(path string) {
	gobFile, err := internal.ReadFile(path)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	bc, err := bytecode.FromBytes(gobFile)
	if err != nil {
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}

	if _, err := io.WriteString(os.Stdout, bc.ToString()); err != nil {
		return
	}
}

func CompileSharkCodeFile(path, outName, compression string, emitInstructionSet bool, argConfig *config.Config) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	f, err := internal.ReadFile(absPath)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.NidumVM)
	file := string(f)
	if bc := sharkEmitter.Compile(&file); bc != nil {
		fileName := internal.GetFileName(absPath) + ".egg"
		if outName != "" {
			fileName = outName
		}
		var bcType bytecode.Type
		switch compression {
		case "brotli":
			bcType = bytecode.BcTypeCompressedBrotli
		case "none":
			bcType = bytecode.BcTypeNormal
		default:
			bcType = bytecode.BcTypeNormal
		}
		bytes, err := bc.ToObj(bcType, bytecode.BcVersionOnos1)
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
	gobFile, err := internal.ReadFile(absPath)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	bc, err := bytecode.FromBytes(gobFile)
	if err != nil {
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.NidumVM)
	sharkEmitter.Exec(bc)
}

func ShowOjbMeta(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	gobFile, err := internal.ReadFile(absPath)
	if err != nil {
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	objFile := bytecode.ObjCode(gobFile)

	compressionType, err := objFile.CompressionType()

	if err != nil {
		exception.PrintExitMsgCtx("Could not get compression type", err.Error(), 1)
	}

	version, err := objFile.Version()

	if err != nil {
		exception.PrintExitMsgCtx("Could not get version", err.Error(), 1)
	}

	instructSize, err := objFile.InstructionLength()

	if err != nil {
		exception.PrintExitMsgCtx("Could not get instruction length", err.Error(), 1)
	}

	bytecode, err := bytecode.FromBytes(gobFile)

	if err != nil {
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}

	fmt.Printf("COMPRESSION: \t%s\n", compressionType.String())
	fmt.Printf("VERSION: \t%s\n", version.String())
	fmt.Printf("OBJ SIZE: \t%d bytes\n", len(objFile))
	fmt.Printf("INSTRUCT SIZE: \t%d bytes\n", instructSize)
	fmt.Printf("NUM CONSTS: \t%d\n", len(bytecode.Constants))
}
