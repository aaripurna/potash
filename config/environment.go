package config

import (
	"os"
)

type AppEnvType string

const (
	AppEnvLocal      AppEnvType = "local"
	AppEnvTest       AppEnvType = "test"
	AppEnvProduction AppEnvType = "production"
)

var AppEnv string
var ViteServerPort string
var ManifestData []byte

var NodeEnv string

func InitEnv() {
	AppEnv = getEnv("APP_ENV", "local")
	ViteServerPort = getEnv("VITE_SERVER_PORT", "5173")
	NodeEnv = getEnv("NODE_ENV", "local")
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)

	if val != "" {
		return val
	}

	return fallback
}
