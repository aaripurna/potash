package cmd

import (
	"github.com/aaripurna/potash/config"
	"github.com/aaripurna/potash/database"
	"github.com/aaripurna/potash/web"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v3"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

var Container *dig.Container

func init() {
	Container = dig.New()

	Container.Provide(func() *html.Engine {
		return html.New("./views", ".html")
	})

	Container.Provide(func(engine *html.Engine) *fiber.App {
		return fiber.New(fiber.Config{Views: engine})
	})

	// DATABASE

	// Passing the dsn explicitly keeps a second database purely additive: another
	// Provide with its own dsn and a dig.Name, no changes here.
	Container.Provide(func() (*gorm.DB, error) {
		return database.Open(config.DatabaseURL)
	})

	// WEB

	Container.Provide(web.NewPagesWeb)
}
