package handler

import (
	"fmt"

	"github.com/pkg/errors"
)

func infolLog(err interface{}) {
	fmt.Printf("INFO     %v+\n", errorString(err))
}

func warnLog(err interface{}) {
	fmt.Printf("WARNING  %v+\n", errorString(err))
}

func fatalLog(err interface{}) {
	fmt.Printf("FATAL    %v+\n", errorString(err))
}

func panicLog(err interface{}) {
	fmt.Printf("PANIC    %v+\n", errorString(err))
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
