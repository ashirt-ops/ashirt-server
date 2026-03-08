import { defineConfig, type Plugin } from 'vitest/config'
import path from 'path'

const mockStylesheetsPlugin: Plugin = {
  name: 'mock-stylesheets',
  enforce: 'pre',
  // Intercept ESM-style stylesheet imports
  resolveId(id) {
    if (id.endsWith('.styl') || id === './stylesheet') {
      return '\0virtual:stylesheet'
    }
  },
  load(id) {
    if (id === '\0virtual:stylesheet') {
      return 'export default {}'
    }
  },
  // Rewrite CJS-style require('./stylesheet') calls before execution
  transform(code, id) {
    if (/\.[tj]sx?$/.test(id) && code.includes("require('./stylesheet')")) {
      return { code: code.replace(/require\('\.\/stylesheet'\)/g, '{}'), map: null }
    }
  },
}

export default defineConfig({
  plugins: [mockStylesheetsPlugin],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['src/test_helpers/setup.ts'],
    include: ['src/**/*.test.{ts,tsx}'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      thresholds: {
        lines: 10,
        functions: 10,
      },
    },
  },
  resolve: {
    alias: {
      src: path.resolve(__dirname, './src'),
    },
  },
})
