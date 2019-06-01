package teraconnectgo

import (
	"github.com/SuperDogHuman/teraconnectgo/interface/handler"
	"google.golang.org/appengine"
)

// Main serve Teraconnect API
func Main(appEnv string) {
	handler.Main(appEnv)
	appengine.Main()
}