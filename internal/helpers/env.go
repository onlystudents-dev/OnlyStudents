package helpers

import "os"

func GetEnvFallback(name string, default_val string) string {
	val, exists := os.LookupEnv(name)

	if exists {
		return val
	} else {
		return default_val
	}
}
