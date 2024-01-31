package logger

import (
	"log"
	"os"
)

const (
	red    = "\033[31m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	reset  = "\033[0m"
)

func InfoLogger() *log.Logger {
	return log.New(os.Stdout, yellow+"INFO: "+reset, log.Lshortfile)
}

func DebugLogger() *log.Logger {
	return log.New(os.Stdout, cyan+"DEBUG: "+reset, log.Lshortfile)
}

func FatalLogger() *log.Logger {
	return log.New(os.Stdout, red+"FATAL: "+reset, log.Lshortfile)
}
