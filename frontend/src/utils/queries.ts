import { Router } from 'routes/index'

function get(key: string): string | undefined {
  if (!key) throw new Error('queries:get: key is required')
  return Router.currentRoute.value.query[key] as string
}

function set(key: string, value: string): void {
  if (!key || !value) throw new Error(`queries:get: key(${key}) or ${value} is required`)
  const queries = Object.assign({}, Router.currentRoute.value.query)
  queries[key] = value
  Router.push({ query: queries })
}

function remove(key) {
  if (!key) throw new Error(`queries:get: key is required`)
  const queries = Object.assign({}, Router.currentRoute.value.query)
  delete queries[key]
  Router.push({ query: queries })
}

const $queries = {
  get,
  set,
  remove,
}

export { $queries }
