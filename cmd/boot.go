package cmd

import (
	"github.com/aaripurna/potash/web"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v3"
	"go.uber.org/dig"
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

	// WEB

	Container.Provide(web.NewPagesWeb)
}
