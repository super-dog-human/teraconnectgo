package handler

import (
	"log"

	"github.com/pkg/errors"
)

func infolLog(err interface{}) {
	log.Printf("INFO     %v+\n", errorString(err))
}

func warnLog(err interface{}) {
	log.Printf("WARNING  %v+\n", errorString(err))
}

func fatalLog(err interface{}) {
	log.Printf("FATAL    %v+\n", errorString(err))
}

func panicLog(err interface{}) {
	log.Printf("PANIC    %v+\n", errorString(err))
}

func errorString(err interface{}) string {
	switch err.(type) {
	case error:
		return errors.WithStack(err.(error)).Error()
	case string:
		return err.(string)
	default:
		return ""
	}
}
