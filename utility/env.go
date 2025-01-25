package utility

import (
	"log"
	"os"
)

func ParseEnv(key string, fallbackEnabled bool, fallbackValue string) string {
	var env = os.Getenv(key)
	if env == "" && !fallbackEnabled {
		log.Panicf("%s environment variable is required.", key)
	}
	if env == "" && fallbackEnabled {
		env = fallbackValue
	}
	return env
}
