import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { router } from "sv-router/vite-plugin";

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), svelte(), router()],
  resolve: {
    tsconfigPaths: true,
  },
  // -----
  // workaround: https://github.com/MIERUNE/svelte-maplibre-gl/issues/206
  optimizeDeps: { exclude: ["maplibre-gl"] },
  // ----
});
