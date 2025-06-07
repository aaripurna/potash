/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"

	"os"

	"github.com/aaripurna/go-fullstack-template/cmd"
	"github.com/aaripurna/go-fullstack-template/config"
	"github.com/joho/godotenv"
)

//go:embed public/*
var publicFS embed.FS

func main() {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "" {
		appEnv = string(config.AppEnvLocal)
	}

	if appEnv == "test" {
		godotenv.Load(".env.test")
	} else {
		godotenv.Load(".env.local", ".env")
	}

	config.InitEnv()

	if data, err := publicFS.ReadFile("public/.vite/manifest.json"); err == nil {
		config.ManifestData = data
	}

	cmd.Execute()
}
