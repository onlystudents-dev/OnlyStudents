import preact from '@preact/preset-vite'
import { defineConfig } from 'vite'
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    preact(),
    tailwindcss(),
  ],
  build: {
    rolldownOptions: {
      output: {
        codeSplitting: {
          groups: [
            {
              name: 'libs',
              test: /node_modules/,
              minSize: 100000,
              maxSize: 500000,
              priority: 10,
            },
          ],
        },
      }
    },
    target: 'es2022',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
        passes: 3,
        pure_funcs: ['console.debug'],
      },
      format: { comments: false },
      mangle: { safari10: true },
    },
    sourcemap: false,
    modulePreload: { polyfill: false },
    reportCompressedSize: true,
  },
})