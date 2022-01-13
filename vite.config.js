import { defineConfig } from "vite";
import { resolve } from "path";

export default defineConfig({
  base: process.env.NODE_ENV === "production" ? "/static/" : "/",
  build: {
    outDir: "static",
    manifest: true,
    sourcemap: false,
    brotliSize: false,
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      input: {
        main: resolve(__dirname, "templates/index.html"),
        notfound: resolve(__dirname, "templates/404.html"),
      },
    },
  },
  plugins: [],
});
