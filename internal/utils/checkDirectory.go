package utils

import (
	"fmt"
	"os"
)

func CheckDirectory(directory string) (os.FileInfo, error) {
	directoryInfo, err := os.Stat(directory)
	if err != nil {
		return nil, err
	}
	if directoryInfo.IsDir() {
		return directoryInfo, nil
	}

	return nil, fmt.Errorf("%s is not a directory", directory)

}
