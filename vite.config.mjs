import { defineConfig, loadEnv } from 'vite'
import tailwindcss from "@tailwindcss/vite";
import tsconfigPaths from "vite-tsconfig-paths";


export default defineConfig(({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return {
    plugins: [tailwindcss(), tsconfigPaths()],
    server: {
      port: process.env.VITE_SERVER_PORT || 5173
    },
    build: {
      outDir: "dist",
      cssCodeSplit: true,
      manifest: true,
      minify: true,
      rollupOptions: {
        input: [
          "assets/main.ts"
        ],
      }
    }
  }
})