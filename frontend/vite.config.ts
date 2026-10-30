import preact from '@preact/preset-vite'
import {defineConfig, type Plugin} from 'vite'
import tailwindcss from "@tailwindcss/vite";
import * as path from "node:path";
import * as fs from "node:fs";

// yes this is made by ai but i thought of the idea
function minifyPublicJSON(): Plugin {
    return {
        name: "minify-public-json",
        closeBundle() {
            const dir = path.resolve("dist/assets/languages");

            if (!fs.existsSync(dir)) return;

            for (const file of fs.readdirSync(dir)) {
                if (!file.endsWith(".json")) continue;

                const filePath = path.join(dir, file);
                // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
                const data = JSON.parse(fs.readFileSync(filePath, "utf8"));

                fs.writeFileSync(filePath, JSON.stringify(data));
            }
        },
    };
}

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        preact(),
        tailwindcss(),
        minifyPublicJSON(),
    ],
    build: {
        rolldownOptions: {
            input: {
                main: path.resolve(__dirname, "index.html"),
                nf: path.resolve(__dirname, "404.html"),
            },
            output: {
                codeSplitting: {
                  groups: [
                      {
                          name: 'opaque',
                          test: /node_modules[\\/]@serenity-kit[\\/]opaque/,
                          priority: 30
                      },
                      {
                          name: 'crypto',
                          test: /src[\\/]util[\\/]crypto[\\/]/,
                          maxSize: 250000,
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
