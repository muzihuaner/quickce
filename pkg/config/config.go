package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       int
	MaxConcurr int
	DefaultTimeout int
}

func Load() *Config {
	return &Config{
		Port:       getEnvInt("PORT", 8080),
		MaxConcurr: getEnvInt("MAX_CONCURRENT", 50),
		DefaultTimeout: getEnvInt("DEFAULT_TIMEOUT", 5),
	}
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
