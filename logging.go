package main

import (
	"bytes"
	"log"
)

func LogErrorf(format string, v ...interface{}) {
	logWithPrefix("\033[31mERROR\033[0m\t", format, v...)
}

func LogWarningf(format string, v ...interface{}) {
	logWithPrefix("\033[33mWARN\033[0m\t", format, v...)
}

func LogInfof(format string, v ...interface{}) {
	logWithPrefix("\033[36mINFO\033[0m\t", format, v...)
}

func logWithPrefix(prefix, format string, v ...interface{}) {
	buffer := bytes.NewBufferString(prefix)
	buffer.WriteString(format)
	log.Printf(buffer.String(), v...)
}
