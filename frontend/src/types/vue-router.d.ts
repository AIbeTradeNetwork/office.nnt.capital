export {}

declare module 'vue-router' {
  interface RouteMeta {
    layout: DefineComponent
    keepAlive?: boolean
    redirect?: RouteRecordRedirect
  }
}
