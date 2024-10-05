package bin

import (
	"runtime"
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
