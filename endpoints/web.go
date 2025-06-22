package endpoints

import (
	"github.com/aaripurna/potash/core"
	"github.com/aaripurna/potash/web"
	"github.com/gofiber/fiber/v2"
)

func RouteWeb(app *fiber.App, pagesWeb *web.PagesWeb) {
	app.Get("/", core.HandleReq(pagesWeb.Index))
}
