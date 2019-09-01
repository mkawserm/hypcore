package z

import (
	"os"
	"path"
)

// Get the running executable name
func ExeName() string {
	_, fileName := path.Split(os.Args[0])
	return fileName
}

// Alias of ExecName
func ExecutableName() string {
	return ExeName()
}
