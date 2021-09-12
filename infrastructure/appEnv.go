package infrastructure

import (
	"github.com/joho/godotenv"
)

var env = ""

// SetAppEnv sets application envirionment string.
func SetAppEnv(appEnv string) error {
	env = appEnv
	if err := godotenv.Load(envFileName()); err != nil {
		return err
	}
	return nil
}

// AppEnv returns application envirionment of 'staging', 'production'.
func AppEnv() string {
	return env
}

func envFileName() string {
	switch AppEnv() {
	case "production":
		return ".env.production"
	case "staging":
		return ".env.staging"
	default:
		return ".env.development"
	}
}
