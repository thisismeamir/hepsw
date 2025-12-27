package utils

import (
	"errors"
	"fmt"
	"os"
)

func CreateDirectory(directoryPath string) error {
	info, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist, must create it.
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if !info.IsDir() {
			return errors.New(fmt.Sprintf("%s is not a directory", directoryPath))
		} else {
			return nil
		}
	}
	return nil
}
