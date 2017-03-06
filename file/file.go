package file

import (
	"io/ioutil"
	"os"
)

// CreateDirectory creates a directory
func CreateDirectory(dir string, perms os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, perms)
		return err
	}

	return nil
}

// CreateFile creates a file if it doens't already exist
func CreateFile(file string) (*os.File, bool, error) {
	var err error
	if _, err = os.Stat(file); os.IsNotExist(err) {
		fd, err := os.Create(file)
		if err != nil {
			return fd, false, err
		}

		return fd, true, err
	}

	return nil, false, err
}

// Reads content from a file
func ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}
