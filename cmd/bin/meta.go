package bin

import (
	"os"
	"path/filepath"
	"runtime"
	"shark/config"
	"shark/exception"
	"shark/internal"
)

func FormatVersion(version, build, codename string) string {
	var curVersion string
	if version == "" {
		curVersion = "dev"
	} else {
		curVersion = version
		if build != "" {
			curVersion += " (" + build + ")"
		}
	}
	curVersion += "\nCore: " + runtime.Version()
	curVersion += "\nCodename: " + codename
	return curVersion
}

func LocateConfigFile(givenPath, targetFile string) config.Config {
	const configFileName = "shark.json"

	// if config file is defined, use it
	if givenPath != "" {
		return config.NewConfigFromFile(givenPath)
	}

	// first check if a config file exists next to the target file
	if internal.IsFileExists(filepath.Join(filepath.Dir(targetFile), configFileName)) {
		givenPath = filepath.Join(filepath.Dir(targetFile), configFileName)
		// use the config file next to the target file
		return config.NewConfigFromFile(givenPath)
	}

	// next, check if a config file exists in the current directory
	dir, err := os.Getwd()
	if err != nil {
		exception.PrintExitMsgCtx("Could not get the current working directory", err.Error(), 1)
	}
	if internal.IsFileExists(filepath.Join(dir, configFileName)) {
		// use the config file in the current directory
		return config.NewConfigFromFile(filepath.Join(dir, configFileName))
	}

	// if no config file is found, use the default config
	return config.NewDefaultConfig()
}
