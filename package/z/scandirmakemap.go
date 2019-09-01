package z

import (
	"io/ioutil"
	"path/filepath"
)

// Scan the directory and make a map using filename without
// extension as key and file contents as value
func ScanDirMakeMap(path string, ext string) map[string]string {
	r := make(map[string]string)
	files, err := ioutil.ReadDir(path)

	if err == nil {
		for _, file := range files {
			var fileName = file.Name()
			var extension = filepath.Ext(fileName)
			var name = fileName[0 : len(fileName)-len(extension)]
			if extension == ext {
				data, err2 := ioutil.ReadFile(path + "/" + fileName)
				if err2 == nil {
					r[name] = string(data)
				}
			}
		}
	}

	return r
}
