package web

import "github.com/aaripurna/go-fullstack-template/core"

type PagesWeb struct{}

func NewPagesWeb() *PagesWeb {
	return &PagesWeb{}
}

func (p *PagesWeb) Index(ctx *core.AppContext) error {
	return core.HtmlResponse{
		Layout:   "layouts/app",
		Template: "pages/index",
	}
}
