/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "embed"

	"os"

	"github.com/aaripurna/go-fullstack-template/cmd"
	"github.com/aaripurna/go-fullstack-template/config"
	"github.com/joho/godotenv"
)

//go:embed public/vite/.vite/manifest.json
var manifestData []byte

//go:embed public/vite/.vite/manifest.json
var foo string

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

	config.ManifestData = manifestData
	config.InitEnv()

	cmd.Execute()
}
