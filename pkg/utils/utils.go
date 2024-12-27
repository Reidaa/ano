package utils

import (
	"log"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// Info writes logs in the color blue with "INFO: " as prefix.
var Info = log.New(os.Stdout, "\u001b[34mINFO: \u001B[0m", log.LstdFlags)

// Warning writes logs in the color yellow with "WARNING: " as prefix.
var Warning = log.New(os.Stdout, "\u001b[33mWARNING: \u001B[0m", log.LstdFlags|log.Lshortfile)

// Error writes logs in the color red with "ERROR: " as prefix.
var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

// Debug writes logs in the color cyan with "DEBUG: " as prefix.
var Debug = log.New(os.Stdout, "\u001b[36mDEBUG: \u001B[0m", log.LstdFlags|log.Lshortfile)

type CliArgumentError struct {
}

func (e *CliArgumentError) Error() string {
	return "error with argument"
}
