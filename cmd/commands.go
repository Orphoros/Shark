package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"shark/bytecode"
	"shark/config"
	"shark/emitter"
	"shark/exception"
	"shark/internal"

	"github.com/phuslu/log"
)

func ExecuteSharkCodeFile(path string, argConfig *config.Config) {
	log.Debug().Msg("Executing Shark code file")
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get abs path of file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	log.Trace().Str("path", absPath).Msg("Getting abs path")
	f, err := internal.ReadFile(absPath)
	if err != nil {
		log.Error().Err(err).Msgf("Could not locate file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.NidumVM)
	sharkEmitter.Interpret(string(f))
}

func DecompileSharkBinaryFile(path string) {
	log.Debug().Msg("Decompiling Shark binary file")
	gobFile, err := internal.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msgf("Could not read file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	log.Trace().Str("path", path).Msg("Reading file")
	bc, err := bytecode.FromBytes(gobFile)
	if err != nil {
		log.Error().Err(err).Msg("Could not decompile bytecode")
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}

	if _, err := io.WriteString(os.Stdout, bc.ToString()); err != nil {
		log.Error().Err(err).Msg("Could not write bytecode to stdout")
		return
	}
}

func CompileSharkCodeFile(path, outName, compression string, emitInstructionSet bool, argConfig *config.Config) {
	log.Debug().Msg("Compiling Shark code file")
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get abs path of file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	log.Trace().Str("path", absPath).Msg("Getting abs path")
	f, err := internal.ReadFile(absPath)
	if err != nil {
		log.Error().Err(err).Msgf("Could not read file '%s'", path)
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
		log.Debug().Str("file", fileName).Str("compression", bcType.String()).Msg("Writing bytecode to file")
		bytes, err := bc.ToObj(bcType, bytecode.BcVersionOnos1)
		if err != nil {
			log.Error().Err(err).Msg("Could not convert bytecode to bytes")
			exception.PrintExitMsgCtx("Could not convert bytecode to bytes", err.Error(), 1)
		}
		if err := internal.WriteFile(fileName, bytes); err != nil {
			log.Error().Err(err).Msg("Could not write bytecode to file")
			exception.PrintExitMsgCtx("Could not write bytecode to file", err.Error(), 1)
		}
	}
}

func ExecuteSharkBinaryFile(path string, argConfig *config.Config) {
	log.Debug().Msg("Executing Shark binary file")
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get abs path of file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	log.Trace().Str("path", absPath).Msg("Getting abs path")
	gobFile, err := internal.ReadFile(absPath)
	if err != nil {
		log.Error().Err(err).Msgf("Could not read file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	bc, err := bytecode.FromBytes(gobFile)
	if err != nil {
		log.Error().Err(err).Msg("Could not decompile bytecode")
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}
	sharkEmitter := emitter.New(&absPath, os.Stdout, &argConfig.NidumVM)
	sharkEmitter.Exec(bc)
}

func ShowOjbMeta(path string) {
	log.Debug().Msg("Showing object metadata")
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get abs path of file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not locate file '%s'", path), err.Error(), 1)
	}
	log.Trace().Str("path", absPath).Msg("Getting abs path")
	gobFile, err := internal.ReadFile(absPath)
	if err != nil {
		log.Error().Err(err).Msgf("Could not read file '%s'", path)
		exception.PrintExitMsgCtx(fmt.Sprintf("Could not read file '%s'", path), err.Error(), 1)
	}
	objFile := bytecode.ObjCode(gobFile)

	compressionType, err := objFile.CompressionType()

	if err != nil {
		log.Error().Err(err).Msg("Could not get compression type")
		exception.PrintExitMsgCtx("Could not get compression type", err.Error(), 1)
	}

	version, err := objFile.Version()

	if err != nil {
		log.Error().Err(err).Msg("Could not get version")
		exception.PrintExitMsgCtx("Could not get version", err.Error(), 1)
	}

	instructSize, err := objFile.InstructionLength()

	if err != nil {
		log.Error().Err(err).Msg("Could not get instruction length")
		exception.PrintExitMsgCtx("Could not get instruction length", err.Error(), 1)
	}

	bytecode, err := bytecode.FromBytes(gobFile)

	if err != nil {
		log.Error().Err(err).Msg("Could not decompile bytecode")
		exception.PrintExitMsgCtx("Could not decompile bytecode", err.Error(), 1)
	}

	fmt.Printf("COMPRESSION: \t%s\n", compressionType.String())
	fmt.Printf("VERSION: \t%s\n", version.String())
	fmt.Printf("OBJ SIZE: \t%d bytes\n", len(objFile))
	fmt.Printf("INSTRUCT SIZE: \t%d bytes\n", instructSize)
	fmt.Printf("NUM CONSTS: \t%d\n", len(bytecode.Constants))
}

func GenerateDefaultConfig() {
	log.Debug().Msg("Generating default config")
	def := config.NewDefaultConfig()

	json, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal default config")
		exception.PrintExitMsgCtx("Could not marshal default config", err.Error(), 1)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Could not get the current working directory")
		exception.PrintExitMsgCtx("Could not get the current working directory", err.Error(), 1)
	}

	file := filepath.Join(dir, "shark.json")

	log.Debug().Str("file", file).Msg("Writing config to file")

	if err := internal.WriteFile(file, json); err != nil {
		log.Error().Err(err).Msg("Could not write to file")
		exception.PrintExitMsgCtx("Could not write to file", err.Error(), 1)
	}
}
