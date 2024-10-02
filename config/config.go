package config

import (
	"encoding/json"
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

func NewConfigFromFile(path string) Config {
	var conf Config

	file := internal.ReadFile(path)

	if err := json.Unmarshal(file, &conf); err != nil {
		exception.PrintExitMsgCtx("Could not unmarshal config file", err.Error(), 1)
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

	return conf
}
