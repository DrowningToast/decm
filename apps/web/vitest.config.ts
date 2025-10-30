import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'happy-dom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    env: {
      VITE_CORE_BACKEND_API: 'http://localhost:8080/api/v1',
      VITE_WALLETCONNECT_PROJECT_ID: 'test-project-id',
    },
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/test/',
        'src/**/*.test.{ts,tsx}',
        'src/**/*.spec.{ts,tsx}',
        'src/router.ts',
        'src/vite-env.d.ts'
      ]
    },
    include: ['src/**/*.{test,spec}.{ts,tsx}'],
    exclude: ['node_modules', 'dist', 'cypress', 'playwright']
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@decm/api': path.resolve(__dirname, '../../packages/api/src')
    }
  },
  define: {
    'import.meta.env.VITE_CORE_BACKEND_API': JSON.stringify('http://localhost:8080/api/v1'),
    'import.meta.env.VITE_WALLETCONNECT_PROJECT_ID': JSON.stringify('test-project-id'),
  }
})