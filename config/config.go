package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"shark/exception"
	"shark/internal"
)

type Config struct {
	OrpVM VmConf `json:"orpVm"`
}

func NewDefaultConfig() Config {
	return Config{
		OrpVM: NewDefaultVmConf(),
	}
}

func NewConfigFromFile(path string) (*Config, error) {
	var conf Config

	file, err := internal.ReadFile(path)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &conf); err != nil {
		return nil, err
	}

	defVmConf := NewDefaultVmConf()
	if conf.OrpVM.GlobalsSize == 0 {
		conf.OrpVM.GlobalsSize = defVmConf.GlobalsSize
	}
	if conf.OrpVM.StackSize == 0 {
		conf.OrpVM.StackSize = defVmConf.StackSize
	}
	if conf.OrpVM.MaxFrames == 0 {
		conf.OrpVM.MaxFrames = defVmConf.MaxFrames
	}

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
