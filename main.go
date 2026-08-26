/*
Copyright © 2025 Nawa Aripurna <nawa@aaripurna.com>
*/
package main

import (
	"embed"
	"io/fs"
	"log"
	"os"

	"github.com/aaripurna/potash/cmd"
	"github.com/aaripurna/potash/config"
	"github.com/joho/godotenv"
)

// all: is required so .vite/manifest.json is embedded - a plain public/* pattern
// skips paths beginning with a dot.
//
//go:embed all:public
var publicFS embed.FS

//go:embed views
var viewsFS embed.FS

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

	public, err := fs.Sub(publicFS, "public")

	if err != nil {
		log.Fatalf("unable to open the embedded public directory: %v", err)
	}

	views, err := fs.Sub(viewsFS, "views")

	if err != nil {
		log.Fatalf("unable to open the embedded views directory: %v", err)
	}

	config.PublicFS = public
	config.ViewsFS = views

	if data, err := fs.ReadFile(public, ".vite/manifest.json"); err == nil {
		config.ManifestData = data
	}

	cmd.Execute()
}
