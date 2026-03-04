package config

import "../../../../go-mariadb-benchmark/internal/config/os"

func getEnvWithFallback[T any](key string, parse func(string) (T, error), fallback T) T {
	str, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	if parse == nil {
		return any(str).(T)
	}

	value, err := parse(str)
	if err != nil {
		return fallback
	}
	return value
}
