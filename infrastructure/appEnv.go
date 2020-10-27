package infrastructure

var env = ""

// SetAppEnv sets application envirionment string.
func SetAppEnv(appEnv string) {
	env = appEnv
}

// AppEnv returns application envirionment of 'staging', 'production'.
func AppEnv() string {
	return env
}
