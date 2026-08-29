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

func GetUint8EnvFallback(name string, defaultVal uint8, max uint8) uint8 {
	v := min(GetUintEnvFallback(name, uint64(defaultVal)), uint64(max))
	return uint8(v) // #nosec G115 -- value clamped to <= max, which fits uint8
}

func GetUint32EnvFallback(name string, defaultVal uint32, max uint32) uint32 {
	v := min(GetUintEnvFallback(name, uint64(defaultVal)), uint64(max))
	return uint32(v) // #nosec G115 -- value clamped to <= max, which fits uint32
}

func getInt64Env(name string, defaultVal int64) int64 {
	val, exists := os.LookupEnv(name)
	if !exists {
		return defaultVal
	}
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	if parsed < 0 {
		return defaultVal
	}
	return parsed
}

func GetInt64EnvFallback(name string, defaultVal int64, max int64) int64 {
	return min(getInt64Env(name, defaultVal), max)
}

func GetInt32EnvFallback(name string, defaultVal int32, max int32) int32 {
	v := min(getInt64Env(name, int64(defaultVal)), int64(max))
	return int32(v) // #nosec G115 -- clamped to <= max, which fits int32
}

func GetIntEnvFallback(name string, defaultVal int, max int) int {
	v := min(getInt64Env(name, int64(defaultVal)), int64(max))
	return int(v) // #nosec G115 -- clamped to <= max, which fits int
}
