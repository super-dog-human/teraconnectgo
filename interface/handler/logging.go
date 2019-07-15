package handler

import "fmt"

func infolLog(message string) {
	fmt.Printf("INFO     %v+\n", message)
}

func warnLog(message string) {
	fmt.Printf("WARNING  %v+\n", message)
}

func fatalLog(err error) {
	fmt.Printf("FATAL    %v+\n", err)
}

func panicLog(err error) {
	fmt.Printf("PANIC    %v+\n", err)
}
