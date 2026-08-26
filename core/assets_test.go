package core

import (
	"bytes"
	"encoding/json"
	stdhtml "html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/aaripurna/potash/config"
	"github.com/gofiber/template/html/v3"
)

// The manifest is parsed once per process, so every test that changes
// config.ManifestData has to clear the cache first.
func resetManifest(t *testing.T, data string) {
	t.Helper()

	manifestOnce = sync.Once{}
	manifestCache = nil
	manifestErr = nil

	config.ManifestData = []byte(data)
}

func setEnv(t *testing.T, nodeEnv string) {
	t.Helper()

	original := config.NodeEnv
	config.NodeEnv = nodeEnv
	config.ViteServerPort = "5173"

	t.Cleanup(func() { config.NodeEnv = original })
}

const singleEntryManifest = `{
  "assets/main.js": {
    "file": "assets/main-abc.js",
    "name": "main",
    "src": "assets/main.js",
    "isEntry": true,
    "css": ["assets/main-abc.css"]
  }
}`

// main -> shared -> deep, with css at two levels and shared reachable twice.
const splitManifest = `{
  "assets/main.js": {
    "file": "assets/main-abc.js",
    "src": "assets/main.js",
    "isEntry": true,
    "css": ["assets/main-abc.css"],
    "imports": ["_shared.js", "_deep.js"]
  },
  "_shared.js": {
    "file": "assets/shared-def.js",
    "css": ["assets/shared-def.css"],
    "imports": ["_deep.js"]
  },
  "_deep.js": {
    "file": "assets/deep-ghi.js"
  }
}`

// renderTemplate runs body through a real engine with the asset funcs
// registered, which is how these helpers are actually reached in production.
func renderTemplate(t *testing.T, body string, data any) (string, error) {
	t.Helper()

	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "page.html"), []byte(body), 0o600); err != nil {
		t.Fatalf("unable to write template: %v", err)
	}

	engine := html.New(dir, ".html")
	AssetHtml(engine)

	if err := engine.Load(); err != nil {
		t.Fatalf("unable to load templates: %v", err)
	}

	var out bytes.Buffer
	err := engine.Render(&out, "page", data)

	return out.String(), err
}

func TestViteAssetDevelopmentUsesViteServer(t *testing.T) {
	setEnv(t, "development")
	resetManifest(t, singleEntryManifest)

	got, err := renderTemplate(t, `{{ vite_asset "assets/main.js" }}`, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, `src="http://localhost:5173/assets/main.js"`) {
		t.Errorf("expected the vite dev server url, got %q", got)
	}

	if strings.Contains(got, "assets/main-abc.js") {
		t.Errorf("development should not read the manifest, got %q", got)
	}
}

func TestViteAssetProductionEmitsModuleScript(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, singleEntryManifest)

	got, err := renderTemplate(t, `{{ vite_asset "assets/main.js" }}`, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Without type="module" a code-split entry is a syntax error in the browser.
	if !strings.Contains(got, `<script type="module" crossorigin src="/assets/main-abc.js">`) {
		t.Errorf("expected a module script tag, got %q", got)
	}

	if !strings.Contains(got, `<link rel="stylesheet" href="/assets/main-abc.css">`) {
		t.Errorf("expected the entry stylesheet, got %q", got)
	}
}

// An entry with no css must still render its script tag.
func TestViteAssetProductionWithoutCss(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, `{"assets/page.js": {"file": "assets/page-xyz.js", "isEntry": true}}`)

	got, err := renderTemplate(t, `{{ vite_asset "assets/page.js" }}`, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, `src="/assets/page-xyz.js"`) {
		t.Errorf("expected a script tag for a css-less entry, got %q", got)
	}

	if strings.Contains(got, "stylesheet") {
		t.Errorf("did not expect a stylesheet link, got %q", got)
	}
}

func TestViteAssetProductionPreloadsSharedChunks(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, splitManifest)

	got, err := renderTemplate(t, `{{ vite_asset "assets/main.js" }}`, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{
		`<link rel="modulepreload" crossorigin href="/assets/shared-def.js">`,
		`<link rel="modulepreload" crossorigin href="/assets/deep-ghi.js">`,
		`<link rel="stylesheet" href="/assets/shared-def.css">`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %s in output, got %q", want, got)
		}
	}

	// _deep.js is reachable from both main and _shared; it must appear once.
	if n := strings.Count(got, "/assets/deep-ghi.js"); n != 1 {
		t.Errorf("expected the shared chunk preloaded once, got %d", n)
	}
}

// A broken manifest must fail the render, not take the process down.
func TestViteAssetBrokenManifestReturnsError(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, "")

	_, err := renderTemplate(t, `{{ vite_asset "assets/main.js" }}`, nil)

	if err == nil {
		t.Fatal("expected an error for an unreadable manifest")
	}

	if !strings.Contains(err.Error(), "bunx vite build") {
		t.Errorf("expected an actionable message, got %v", err)
	}
}

func TestViteAssetUnknownEntryReturnsError(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, singleEntryManifest)

	_, err := renderTemplate(t, `{{ vite_asset "assets/nope.js" }}`, nil)

	if err == nil {
		t.Fatal("expected an error for an entry missing from the manifest")
	}

	if !strings.Contains(err.Error(), "assets/nope.js") {
		t.Errorf("expected the entry name in the error, got %v", err)
	}
}

func TestPropsIsEscapedIntoTheAttribute(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, singleEntryManifest)

	heading := `He said "hi" & <left> </script>`
	data := map[string]any{"Island": map[string]any{"heading": heading}}

	got, err := renderTemplate(t, `<div data-props="{{ props .Island }}"></div>`, data)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No raw quote may survive, or the attribute could be broken out of.
	if strings.Contains(got, `data-props="{"`) {
		t.Errorf("raw quotes leaked into the attribute: %q", got)
	}

	// json.Marshal escapes the html-significant bytes before html/template ever
	// sees them, so < and & never reach the attribute literally.
	for _, unwanted := range []string{"<left>", "</script>"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("%s survived unescaped in %q", unwanted, got)
		}
	}

	// The invariant that matters: what a browser decodes must be the original.
	start := strings.Index(got, `data-props="`) + len(`data-props="`)
	end := strings.Index(got[start:], `"`) + start
	decoded := stdhtml.UnescapeString(got[start:end])

	var round map[string]string
	if err := json.Unmarshal([]byte(decoded), &round); err != nil {
		t.Fatalf("attribute did not decode to valid json: %v (%q)", err, decoded)
	}

	if round["heading"] != heading {
		t.Errorf("round trip changed the value:\n got %q\nwant %q", round["heading"], heading)
	}
}

func TestDevelopmentDependencies(t *testing.T) {
	tests := []struct {
		name     string
		nodeEnv  string
		wantHTML bool
	}{
		{name: "development injects the vite client", nodeEnv: "development", wantHTML: true},
		{name: "production injects nothing", nodeEnv: "production", wantHTML: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(t, tt.nodeEnv)
			resetManifest(t, singleEntryManifest)

			got, err := renderTemplate(t, `{{ development_dependencies }}`, nil)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if has := strings.Contains(got, "@vite/client"); has != tt.wantHTML {
				t.Errorf("vite client present = %v, want %v (got %q)", has, tt.wantHTML, got)
			}
		})
	}
}

func TestAssetsFinder(t *testing.T) {
	tests := []struct {
		name    string
		nodeEnv string
		entry   string
		want    string
		wantErr bool
	}{
		{name: "development points at the dev server", nodeEnv: "development", entry: "assets/main.js", want: "http://localhost:5173/assets/main.js"},
		{name: "production resolves through the manifest", nodeEnv: "production", entry: "assets/main.js", want: "/assets/main-abc.js"},
		{name: "surrounding whitespace is trimmed", nodeEnv: "production", entry: "  assets/main.js  ", want: "/assets/main-abc.js"},
		{name: "unknown entry errors", nodeEnv: "production", entry: "assets/nope.js", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(t, tt.nodeEnv)
			resetManifest(t, singleEntryManifest)

			got, err := assetsFinder(tt.entry)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEntryDependenciesOrderAndDedup(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, splitManifest)

	entry, err := manifestEntry("assets/main.js")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	css, preloads, err := entryDependencies(entry)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantCSS := []string{"assets/main-abc.css", "assets/shared-def.css"}
	if strings.Join(css, ",") != strings.Join(wantCSS, ",") {
		t.Errorf("css = %v, want %v", css, wantCSS)
	}

	wantPreloads := []string{"assets/shared-def.js", "assets/deep-ghi.js"}
	if strings.Join(preloads, ",") != strings.Join(wantPreloads, ",") {
		t.Errorf("preloads = %v, want %v", preloads, wantPreloads)
	}
}

// A chunk graph with a cycle must not hang or overflow the stack.
func TestEntryDependenciesHandlesCycles(t *testing.T) {
	setEnv(t, "production")
	resetManifest(t, `{
	  "assets/main.js": {"file": "assets/main-abc.js", "isEntry": true, "imports": ["_a.js"]},
	  "_a.js": {"file": "assets/a.js", "imports": ["_b.js"]},
	  "_b.js": {"file": "assets/b.js", "imports": ["_a.js"]}
	}`)

	entry, err := manifestEntry("assets/main.js")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, preloads, err := entryDependencies(entry)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(preloads) != 2 {
		t.Errorf("expected each chunk once, got %v", preloads)
	}
}

func TestParseManifestDataIsCached(t *testing.T) {
	resetManifest(t, singleEntryManifest)

	first, err := parseManifestData()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Changing the source without resetting must not be picked up: the parse is
	// behind a sync.Once because the manifest is embedded and immutable.
	config.ManifestData = []byte(`{"other.js": {"file": "other.js"}}`)

	second, err := parseManifestData()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := second["assets/main.js"]; !ok {
		t.Error("expected the cached manifest to be reused")
	}

	if len(first) != len(second) {
		t.Errorf("cache returned a different map: %d vs %d", len(first), len(second))
	}
}
