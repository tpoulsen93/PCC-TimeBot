import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// The SPA is served same-origin by the Go binary in production, so API calls
// use relative paths. In local development, proxy /api to the Go server.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api": "http://localhost:3000",
    },
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
