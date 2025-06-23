package web

import (
	"github.com/aaripurna/potash/core"
	"github.com/aaripurna/potash/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PagesWeb struct{}

func NewPagesWeb() *PagesWeb {
	return &PagesWeb{}
}

func (p *PagesWeb) Index(ctx *core.AppContext) error {
	return core.HtmlResponse{
		Layout:   "layouts/app",
		Template: "pages/index",
		Assigns: fiber.Map{
			"Alert": dto.AlertDialog{
				ButtonText:  "Confirm Changes",
				ID:          uuid.NewString(),
				Title:       "Are you absolutely sure?",
				Description: "This action cannot be undone. This will permanently delete your account and remove your data from our servers.",
			},
		},
	}
}
