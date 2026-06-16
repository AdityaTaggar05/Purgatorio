import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// Backend routes used by the frontend live under these prefixes.
// In dev, proxy them straight to the locally-running backend (e.g. via
// `docker compose up purg-backend`, which publishes 8080:8080).
const BACKEND_URL = "http://localhost:8080";
const API_PREFIXES = ["/auth", "/user", "/base", "/shop", "/army"];

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    proxy: Object.fromEntries(
      API_PREFIXES.map((prefix) => [
        prefix,
        { target: BACKEND_URL, changeOrigin: true, ws: true },
      ])
    ),
  },
});
