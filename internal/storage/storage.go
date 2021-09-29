package storage

import (
	"os"
	"path"
)

func PathPrefix(base string) string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	return path.Join(wd, base)
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
