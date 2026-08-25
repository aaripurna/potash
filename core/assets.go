package core

import (
	"encoding/json"
	"fmt"
	"html/template"
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
	manifestErr   error
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
		"vite_asset", func(name string) (template.HTML, error) {
			if config.NodeEnv != string(config.AppEnvProduction) {
				return template.HTML(fmt.Sprintf(`
					<script type="module" src="http://localhost:%s/%s"></script>
				`, config.ViteServerPort, strings.TrimSpace(name))), nil
			}

			entry, err := manifestEntry(name)

			if err != nil {
				return "", err
			}

			cssFiles, preloads, err := entryDependencies(entry)

			if err != nil {
				return "", err
			}

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

			return template.HTML(out.String()), nil
		},
	)

	engine.AddFunc(
		"asset_path", func(name string) (string, error) {
			return assetsFinder(name)
		},
	)

	// props serializes a value for an island's data-props attribute. It returns
	// a plain string so html/template applies its own attribute escaping.
	engine.AddFunc(
		"props", func(value any) (string, error) {
			encoded, err := json.Marshal(value)

			if err != nil {
				return "", fmt.Errorf("unable to encode island props: %w", err)
			}

			return string(encoded), nil
		},
	)
}

// entryDependencies walks an entry and its static imports, collecting the
// stylesheets to link and the shared chunks to preload. Without the preloads
// the browser only discovers a shared chunk after parsing the entry, costing a
// round trip.
func entryDependencies(entry manifestItem) (css []string, preloads []string, err error) {
	manifestData, err := parseManifestData()

	if err != nil {
		return nil, nil, err
	}

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

	return css, preloads, nil
}

// parseManifestData returns an error rather than exiting: a missing or broken
// manifest should fail the request being rendered, not take the whole server
// down with it.
func parseManifestData() (map[string]manifestItem, error) {
	manifestOnce.Do(func() {
		if err := json.Unmarshal(config.ManifestData, &manifestCache); err != nil {
			manifestErr = fmt.Errorf("unable to read manifest.json, please run `bunx vite build`: %w", err)
		}
	})

	return manifestCache, manifestErr
}

func manifestEntry(name string) (manifestItem, error) {
	manifestData, err := parseManifestData()

	if err != nil {
		return manifestItem{}, err
	}

	item, ok := manifestData[strings.TrimSpace(name)]

	if !ok {
		return manifestItem{}, fmt.Errorf("unable to find %s in your assets list", strings.TrimSpace(name))
	}

	return item, nil
}

func assetsFinder(name string) (string, error) {
	if config.NodeEnv != string(config.AppEnvProduction) {
		return fmt.Sprintf("http://localhost:%s/%s", config.ViteServerPort, strings.TrimSpace(name)), nil
	}

	item, err := manifestEntry(name)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s", item.File), nil
}
