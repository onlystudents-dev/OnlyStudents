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
              name: 'opaque',
              test: /node_modules[\\/]@serenity-kit[\\/]opaque/,
              priority: 30
            },
            {
              name: 'shared',
              test: /src[\\/](navbar|util|types)[\\/]/,
              priority: 25
            },
            {
              name: 'libs',
              test: /node_modules[\\/]/,
              maxSize: 250000,
              priority: 20,
            },
            {
              name: 'auth',
              test: /src[\\/]auth[\\/]/,
              maxSize: 50000,
              priority: 10,
            },
            {
              name: 'guardian',
              test: /src[\\/]ui[\\/]guardian[\\/]/,
              maxSize: 50000,
              priority: 10,
            },
            {
              name: 'me',
              test: /src[\\/]ui[\\/]me[\\/]/,
              maxSize: 50000,
              priority: 10,
            },
            {
              name: 'student',
              test: /src[\\/]ui[\\/]student[\\/]/,
              maxSize: 50000,
              priority: 10,
            },
            {
              name: 'teacher',
              test: /src[\\/]ui[\\/]teacher[\\/]/,
              maxSize: 50000,
              priority: 10,
            },
            {
              name: 'admin',
              test: /src[\\/]admin[\\/]/,
              maxSize: 50000,
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
