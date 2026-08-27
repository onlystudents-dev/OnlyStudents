package helpers

import (
	"os"
	"strconv"
)

func GetEnvFallback(name string, default_val string) string {
	val, exists := os.LookupEnv(name)

	if exists {
		return val
	} else {
		return default_val
	}
}

func GetUintEnvFallback(name string, default_val uint64) uint64 {
	value := GetEnvFallback(name, strconv.FormatUint(default_val, 10))

	uint_value, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return default_val
	}

	return uint_value
}
