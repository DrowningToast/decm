// vite.config.ts

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import generouted from "@generouted/react-router/plugin";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
    plugins: [react(), generouted(), tailwindcss()],
    server: {
        port: 3000,
    },
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    build: {
        rollupOptions: {
            output: {
                manualChunks: undefined,
            },
        },
    },
    optimizeDeps: {
        exclude: ["**/*.test.tsx?", "**/*.spec.tsx?"],
    },
});
