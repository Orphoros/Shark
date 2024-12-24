package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"shark/exception"
	"shark/internal"

	"github.com/phuslu/log"
)

type Config struct {
	NidumVM VmConf `json:"nvm"`
}

func NewDefaultConfig() Config {
	return Config{
		NidumVM: NewDefaultVmConf(),
	}
}

func NewConfigFromFile(path string) (*Config, error) {
	var conf Config

	file, err := internal.ReadFile(path)

	if err != nil {
		log.Error().Err(err).Msg("Could not read config file")
		return nil, err
	}

	if err := json.Unmarshal(file, &conf); err != nil {
		log.Error().Err(err).Msg("Could not unmarshal config file")
		return nil, err
	}

	defVmConf := NewDefaultVmConf()
	if conf.NidumVM.GlobalsSize == 0 {
		conf.NidumVM.GlobalsSize = defVmConf.GlobalsSize
	}
	if conf.NidumVM.StackSize == 0 {
		conf.NidumVM.StackSize = defVmConf.StackSize
	}
	if conf.NidumVM.MaxFrames == 0 {
		conf.NidumVM.MaxFrames = defVmConf.MaxFrames
	}

	log.Debug().Str("path", path).Any("config", conf).Msg("Conf loaded")

	return &conf, nil
}

func LocateConfig(givenPath, targetFile *string) (*Config, error) {
	const configFileName = "shark.json"

	if givenPath != nil && internal.IsFileExists(*givenPath) {
		return NewConfigFromFile(*givenPath)
	}

	// first check if a config file exists next to the target file
	if targetFile == nil {
		def := NewDefaultConfig()
		return &def, nil
	}

	if internal.IsFileExists(filepath.Join(filepath.Dir(*targetFile), configFileName)) {
		// use the config file next to the target file
		return NewConfigFromFile(filepath.Join(filepath.Dir(*targetFile), configFileName))
	}

	// next, check if a config file exists in the current directory
	dir, err := os.Getwd()
	if err != nil {
		exception.PrintExitMsgCtx("Could not get the current working directory", err.Error(), 1)
	}
	if internal.IsFileExists(filepath.Join(dir, configFileName)) {
		// use the config file in the current directory
		return NewConfigFromFile(filepath.Join(dir, configFileName))
	}

	// if no config file is found, use the default config
	def := NewDefaultConfig()
	return &def, nil
}
