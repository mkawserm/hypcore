package app

import "os"

func EnsureDir(dirName string) error {
	err := os.Mkdir(dirName, 0777)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
