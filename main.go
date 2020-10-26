package main

import (
	"os"

	"github.com/super-dog-human/teraconnectgo/interface/handler"
)

func main() {
	if appEnv := os.Args[1]; appEnv != "" {
		handler.Main(appEnv)
	}
}
