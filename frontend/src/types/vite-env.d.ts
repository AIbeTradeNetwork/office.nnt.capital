/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly NODE_ENV: 'development' | 'production'
  readonly ENVIRONMENT: 'staging' | 'production'
  readonly VITE_APP_URL_GQL: string
  readonly VITE_APP_URL_GQL_BOTS: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
