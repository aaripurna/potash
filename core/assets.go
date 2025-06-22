package core

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/aaripurna/potash/config"
	"github.com/gofiber/template/html/v2"
)

type manifestItem struct {
	File string   `json:"file"`
	Name string   `json:"name"`
	Src  string   `json:"src"`
	Css  []string `json:"css"`
}

func AssetHtml(engine *html.Engine) {
	engine.AddFunc(
		"development_dependencies", func() template.HTML {
			if config.AppEnv != string(config.AppEnvProduction) {
				return template.HTML(fmt.Sprintf(`
					<script type="module" src="http://localhost:%s/@vite/client"></script>
				`, config.ViteServerPort))
			} else {
				return template.HTML("")
			}
		},
	)

	engine.AddFunc(
		"vite_asset", func(name string) template.HTML {
			var result template.HTML

			if config.AppEnv != string(config.AppEnvProduction) {
				result = template.HTML(fmt.Sprintf(`
					<script type="module" src="http://localhost:%s/%s"></script>
				`, config.ViteServerPort, strings.TrimSpace(name)))
			} else {
				manifestItem := manifestEntry(name)

				css := ""

				if len(manifestItem.Css) > 0 {
					for _, cssFile := range manifestItem.Css {
						css = fmt.Sprintf(`%s<link rel="stylesheet" href="/%s">`, css, cssFile)
					}

					result = template.HTML(fmt.Sprintf(`
					%s
					<script src="/%s"></script>
					`, css, manifestItem.File))

				}
			}

			return result
		},
	)

	engine.AddFunc(
		"asset_path", func(name string) string {
			return assetsFinder(name)
		},
	)
}

func parseManifestData() map[string]manifestItem {
	var result map[string]manifestItem

	err := json.Unmarshal(config.ManifestData, &result)

	if err != nil {
		log.Fatal("Unable to read the manifest.json\nPlease ensure you run `bunx vite build`")
		panic(err)
	}

	return result
}

func manifestEntry(name string) manifestItem {
	manifestData := parseManifestData()
	if item, ok := manifestData[strings.TrimSpace(name)]; ok {
		return item
	} else {
		panic(fmt.Sprintf("Unable to find %s in your assets list", name))
	}
}

func assetsFinder(name string) string {
	if config.AppEnv != string(config.AppEnvProduction) {
		return fmt.Sprintf("http://localhost:%s/%s", config.ViteServerPort, strings.TrimSpace(name))
	} else {
		manifestData := parseManifestData()

		if item, ok := manifestData[strings.TrimSpace(name)]; ok {
			return fmt.Sprintf("/%s", item.File)
		} else {
			panic(fmt.Sprintf("Unable to find %s in your assets list", name))
		}
	}
}
