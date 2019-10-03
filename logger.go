package main

import (
	"log"
	"os"
	"strings"
)

var debug bool

func init() {
	debug = os.Getenv("DEBUG") != ""
}

// infoLog example:
// infoLog("timezone %s", timezone)
func infoLog(msg string, vars ...interface{}) {
	log.Printf(strings.Join([]string{"[INFO]", msg}, " "), vars...)
}

// debugLog example:
// debugLog("timezone %s", timezone)
func debugLog(msg string, vars ...interface{}) {
	if debug {
		log.Printf(strings.Join([]string{"[DEBUG]", msg}, " "), vars...)
	}
}

// errorLog example:
// errorLog(errors.Errorf("Invalid timezone %s", timezone))
func errorLog(err error) {
	log.Printf("[ERROR] %+v\n", err)
}
