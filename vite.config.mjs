import { defineConfig, loadEnv } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig(({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return {
    plugins: [tailwindcss()],
    // JSX compiles straight to preact/jsx-runtime. The react aliases are the
    // escape hatch: `import { useState } from "react"` and most React libraries
    // resolve to preact/compat, and dropping these aliases swaps in real React
    // without touching component code.
    resolve: {
      alias: {
        react: "preact/compat",
        "react-dom": "preact/compat",
        "react-dom/test-utils": "preact/test-utils",
        "react/jsx-runtime": "preact/jsx-runtime",
      },
    },
    esbuild: {
      jsx: "automatic",
      jsxImportSource: "preact",
    },
    // Islands load via dynamic import, so vite cannot crawl to preact from the
    // entry at startup. Without this it discovers these mid-session, re-optimizes
    // and force-reloads the page - which shows up as random failures under test.
    optimizeDeps: {
      include: ["preact", "preact/hooks", "preact/jsx-runtime", "preact/jsx-dev-runtime"],
    },
    server: {
      port: process.env.VITE_SERVER_PORT || 5173,
    },
    build: {
      outDir: "dist",
      cssCodeSplit: true,
      manifest: true,
      rollupOptions: {
        input: ["assets/main.js"],
      },
    },
  };
});
