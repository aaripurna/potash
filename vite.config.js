import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return {
    server: {
      port: process.env.VITE_SERVER_PORT || 5173
    },
    build: {
      cssCodeSplit: true,
      outDir: "./public/vite/",
      manifest: true,
      rollupOptions: {
        input: [
          "assets/main.js"
        ],
      }
    }
  }
})