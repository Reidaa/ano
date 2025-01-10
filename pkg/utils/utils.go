package utils

import (
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

type CliArgumentError struct {
}

func (e *CliArgumentError) Error() string {
	return "error with argument"
}
