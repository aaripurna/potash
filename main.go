/*
Copyright © 2025 Nawa Aripurna <nawa@aaripurna.com>
*/
package main

import (
	"embed"

	"os"

	"github.com/aaripurna/potash/cmd"
	"github.com/aaripurna/potash/config"
	"github.com/joho/godotenv"
)

//go:embed public/*
var publicFS embed.FS

func main() {
	appEnv := os.Getenv("APP_ENV")

	switch appEnv {
	case "test":
		godotenv.Load(".env.test")
	case "production":
		godotenv.Load(".env")
	default:
		godotenv.Load(".env.local")
	}

	config.InitEnv()

	if data, err := publicFS.ReadFile("public/.vite/manifest.json"); err == nil {
		config.ManifestData = data
	}

	cmd.Execute()
}
