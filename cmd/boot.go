package cmd

import (
	"log"
	"net/http"

	"github.com/aaripurna/potash/config"
	"github.com/aaripurna/potash/database"
	"github.com/aaripurna/potash/web"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v3"
	"go.uber.org/dig"
)

var Container *dig.Container

func init() {
	Container = dig.New()

	must(Container.Provide(func() *html.Engine {
		return html.NewFileSystem(http.FS(config.ViewsFS), ".html")
	}))

	must(Container.Provide(func(engine *html.Engine) *fiber.App {
		return fiber.New(fiber.Config{Views: engine})
	}))

	// DATABASE
	must(Container.Provide(func() (*database.DB, error) {
		return database.NewDB(config.DatabaseURL)
	}))

	// WEB

	must(Container.Provide(web.NewPagesWeb))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
