package core

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strings"
	"sync"

	"github.com/aaripurna/potash/config"
	"github.com/gofiber/template/html/v3"
)

type manifestItem struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	Css     []string `json:"css"`
	Imports []string `json:"imports"`
}

var (
	manifestOnce  sync.Once
	manifestCache map[string]manifestItem
)

func AssetHtml(engine *html.Engine) {
	engine.AddFunc(
		"development_dependencies", func() template.HTML {
			if config.NodeEnv != string(config.AppEnvProduction) {
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
			if config.NodeEnv != string(config.AppEnvProduction) {
				return template.HTML(fmt.Sprintf(`
					<script type="module" src="http://localhost:%s/%s"></script>
				`, config.ViteServerPort, strings.TrimSpace(name)))
			}

			entry := manifestEntry(name)
			cssFiles, preloads := entryDependencies(entry)

			var out strings.Builder

			for _, cssFile := range cssFiles {
				fmt.Fprintf(&out, `<link rel="stylesheet" href="/%s">`, cssFile)
			}

			// type="module" is required: any entry that gets code split emits
			// `import "./chunk.js"`, which is a syntax error in a classic script.
			fmt.Fprintf(&out, `<script type="module" crossorigin src="/%s"></script>`, entry.File)

			for _, preload := range preloads {
				fmt.Fprintf(&out, `<link rel="modulepreload" crossorigin href="/%s">`, preload)
			}

			return template.HTML(out.String())
		},
	)

	engine.AddFunc(
		"asset_path", func(name string) string {
			return assetsFinder(name)
		},
	)
}

// entryDependencies walks an entry and its static imports, collecting the
// stylesheets to link and the shared chunks to preload. Without the preloads
// the browser only discovers a shared chunk after parsing the entry, costing a
// round trip.
func entryDependencies(entry manifestItem) (css []string, preloads []string) {
	manifestData := parseManifestData()
	seen := map[string]bool{}

	var walk func(item manifestItem)
	walk = func(item manifestItem) {
		for _, cssFile := range item.Css {
			if !seen["css:"+cssFile] {
				seen["css:"+cssFile] = true
				css = append(css, cssFile)
			}
		}

		for _, key := range item.Imports {
			dep, ok := manifestData[key]

			if !ok || seen["js:"+key] {
				continue
			}

			seen["js:"+key] = true
			preloads = append(preloads, dep.File)
			walk(dep)
		}
	}

	walk(entry)

	return css, preloads
}

func parseManifestData() map[string]manifestItem {
	manifestOnce.Do(func() {
		if err := json.Unmarshal(config.ManifestData, &manifestCache); err != nil {
			log.Fatal("Unable to read the manifest.json\nPlease ensure you run `bunx vite build`")
			panic(err)
		}
	})

	return manifestCache
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
	if config.NodeEnv != string(config.AppEnvProduction) {
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
