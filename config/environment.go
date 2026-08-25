package config

import (
	"os"
	"strconv"
	"time"
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

var DatabaseURL string
var DBMaxOpenConns int
var DBMaxIdleConns int
var DBConnMaxLifetime time.Duration

func InitEnv() {
	AppEnv = getEnv("APP_ENV", "local")
	ViteServerPort = getEnv("VITE_SERVER_PORT", "5173")
	NodeEnv = getEnv("NODE_ENV", "local")

	DatabaseURL = getEnv("DATABASE_URL", "postgres://postgres@localhost:5432/potash_dev?sslmode=disable")
	DBMaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25)
	DBMaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 5)
	DBConnMaxLifetime = time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)

	if val != "" {
		return val
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	val, err := strconv.Atoi(os.Getenv(key))

	if err != nil {
		return fallback
	}

	return val
}
