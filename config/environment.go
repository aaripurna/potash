package config

import (
	"embed"
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
var PublicFS embed.FS
var ManifestData []byte

func InitEnv() {
	AppEnv = getEnv("APP_ENV", "local")
	ViteServerPort = getEnv("VITE_SERVER_PORT", "5173")
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)

	if val != "" {
		return val
	}

	return fallback
}
