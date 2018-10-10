package packages

import (
	"path"
	"runtime"
)

var fixturesDirectory string

func init() {
	fixturesDirectory = path.Join(path.Dir(currentFile(), "..", "text", "fixtures"))
}

func currentFile() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filename
}
