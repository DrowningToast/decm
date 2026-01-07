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
        proxy: {
            "/api": {
                target: process.env.VITE_CORE_BACKEND_API || "http://localhost:8080",
                changeOrigin: true,
                secure: false,
            },
        },
    },
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    build: {
        rollupOptions: {
            external: [/.*\.test\.tsx?$/, /.*\.spec\.tsx?$/],
        },
    },
});
