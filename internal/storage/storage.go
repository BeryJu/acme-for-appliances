package storage

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

func PathPrefix() string {
	storageDir := viper.GetString("storage")
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	return path.Join(wd, storageDir)
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
