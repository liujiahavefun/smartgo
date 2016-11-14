package log

import (
	"os"
)

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func Write(file *os.File, s string) error {
	_, err := file.WriteString(s)
	return err
}
