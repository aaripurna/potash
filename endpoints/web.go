package endpoints

import (
	"github.com/aaripurna/go-fullstack-template/core"
	"github.com/aaripurna/go-fullstack-template/web"
	"github.com/gofiber/fiber/v2"
)

func RouteWeb(app *fiber.App, pagesWeb *web.PagesWeb) {
	app.Get("/", core.HandleReq(pagesWeb.Index))
}
