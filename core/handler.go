package core

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HandlerFunc func(*AppContext) error

type HtmlResponse struct {
	Layout     string
	Template   string
	Assigns    fiber.Map
	StatusCode int
}

func (h HtmlResponse) Error() string {
	return fmt.Sprintf("Layout = %s, Template = %s, Assign = %v, StatusCode = %d", h.Layout, h.Template, h.Assigns, h.StatusCode)
}

func HandleReq(handlerFn HandlerFunc) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		appCtx := &AppContext{Ctx: ctx}

		result := handlerFn(appCtx)

		if h, ok := result.(HtmlResponse); ok {
			statusCode := h.StatusCode

			if statusCode == 0 {
				statusCode = 200
			}

			ctx.Status(statusCode)

			return ctx.Render(h.Template, h.Assigns, h.Layout)
		}

		return result
	}
}
