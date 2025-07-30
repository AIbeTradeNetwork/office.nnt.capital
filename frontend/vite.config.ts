import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { viteStaticCopy } from 'vite-plugin-static-copy'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import { resolve } from 'path'
import { checker } from 'vite-plugin-checker'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

// import { obfuscator, obfuscatorConfig } from './obfuscator.js'

const oklchToRGBA = () => {
  function replace(code: string) {
    return code
      .replace(/(@supports\snot\s\(color:\soklch\(0%\s0\s0\)\))/gm, '@media screen')
      .replace(/(@supports\s\(color:\soklch\(0%\s0\s0\)\))/gm, '@supports (rem:rem)')
      .replace(/oklch\((var\(--.+\))\/var\(--.+?\)\)/gm, '$1')
      .replace(/oklch\((var\(--.+\))\/1\)/gm, '$1')
      .replace(/oklch\((var\(--.+,var\(--.+\)\))\/var\(--.+?\)\)/gm, '$1')
      .replace(/oklch\((var\(--.+,\svar\(--.+\)\))\s\/\svar\(--.+?\)\)/gm, '$1')
      .replace(/oklch\((var\(--.+\))\/.+?\)/gm, '$1')
    // .replace(/rgb\((\d{1,}\s\d{1,}\s\d{1,})\s\/\s.*?\){1,2}/gm, 'rgb($1)')
  }
  return {
    name: 'oklchToRGBA',
    transform(code, id) {
      if (id.endsWith('.css')) {
        return {
          code: replace(code),
          map: null,
        }
      }
    },
    generateBundle(options, bundle) {
      Object.keys(bundle).forEach((key) => {
        if (key && key.match(/.css$/g)) {
          bundle[key].source = replace(bundle[key].source)
        }
      })
    },
  }
}

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  // Set the third parameter to '' to load all env regardless of the `VITE_` prefix.
  // const env = loadEnv(mode, process.cwd(), '')

  return {
    logLevel: 'info',
    plugins: [
      nodePolyfills(),
      checker({
        typescript: true,
      }),
      viteStaticCopy({
        targets: [
          {
            src: resolve(__dirname, './src') + '/i18n/locales',
            dest: '',
          },
          // {
          //   src: resolve(__dirname, './src') + '/index.html',
          //   dest: '',
          // },
          // {
          //   src: resolve(__dirname, './src') + '/landing/coin/index.html',
          //   dest: '',
          // },
        ],
      }),
      vue({
        template: {
          // transformAssetUrls
        },
      }),
      VueI18nPlugin({
        /* options */
        // locale messages resource pre-compile option
        runtimeOnly: false,
        include: [resolve(__dirname, './src/i18n/locales/**')],
      }),
      // obfuscator({
      //   global: true,
      //   exclude: ['./node_modules/**/*'],
      //   ...obfuscatorConfig,
      // }),
      // splitVendorChunkPlugin(),
      oklchToRGBA(),
    ],
    // css: {
    //   postcss: {
    //     plugins: [tailwindcss()],
    //   },
    // },

    css: {
      devSourcemap: mode === 'development',
    },

    resolve: {
      alias: {
        public: resolve(__dirname, './public'),
        root: resolve(__dirname, './'),
        src: resolve(__dirname, './src'),
        assets: resolve(__dirname, './src/assets'),
        components: resolve(__dirname, './src/components'),
        config: resolve(__dirname, './src/config'),
        i18n: resolve(__dirname, './src/i18n'),
        layouts: resolve(__dirname, './src/layouts'),
        pages: resolve(__dirname, './src/pages'),
        routes: resolve(__dirname, './src/routes'),
        types: resolve(__dirname, './src/types'),
        utils: resolve(__dirname, './src/utils'),
        queries: resolve(__dirname, './src/queries'),
        coin: resolve(__dirname, './src/landing/coin'),
      },
    },
    optimizeDeps: {
      assetsInlineLimit: 0,
      // sourcemap: true,
      esbuildOptions: {
        target: ['esnext'],
        define: {
          global: 'globalThis',
          __VUE_I18N_FULL_INSTALL__: JSON.stringify(false),
          __VUE_I18N_LEGACY_API__: JSON.stringify(false),
          __INTLIFY_PROD_DEVTOOLS__: JSON.stringify(false),
        },
      },
    },
    server: {
      open: true,
      port: 8080,
    },
    build: {
      sourcemap: mode === 'development',
      ssr: false,
      ssrManifest: false,
      emptyOutDir: true,
      manifest: true,
      cssCodeSplit: true,
      minify: mode !== 'development' ? 'esbuild' : false,
      target: ['es2020'],
      assetsInlineLimit: 0,
      commonjsOptions: {
        transformMixedEsModules: true,
      },
      reportCompressedSize: true,
      chunkSizeWarningLimit: 1024,
      modulePreload: {
        polyfill: true,
      },
      rollupOptions: {
        // input: {
        //   index: resolve(__dirname, './index.html'),
        //   // coint: resolve(__dirname, './src/landing/coin/coin.html'),
        // },
        output: {
          manualChunks(id, { getModuleInfo }) {
            // const mainEntry = './src/main.ts'
            // const cointEntry = './src/landing/coin/main.js'

            if (id.includes('vue')) {
              return 'vue-vendor'
            }

            if (id.includes('@tonconnect') || id.includes('@ton')) {
              return 'ton-venfor'
            }

            if (id.includes('decimal')) {
              return 'decimal-venfor'
            }

            // Получаем информацию о модуле
            const moduleInfo = getModuleInfo(id)

            // if (moduleInfo.importers.every((importer) => importer.includes(mainEntry))) {
            //   return 'main'
            // }

            // // // Проверяем, импортируется ли модуль исключительно в secondary.js
            // if (moduleInfo.importers.every((importer) => importer.includes(cointEntry))) {
            //   return 'coin'
            // }

            // Создаем отдельный чанк для модулей, которые импортируются динамически
            if (moduleInfo.dynamicImporters.length > 0) {
              return 'async-imports'
            }
          },
        },
      },
    },
  }
})
