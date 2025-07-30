function set(key: string, value: string): void {
  if (!key || !value) throw new Error(`error:storage:set key=${key}, value=${value}`)
  localStorage.setItem(key, value)
}

function get(key: string): string | null {
  if (!key) throw new Error(`error:storage:get key=${key}`)
  return localStorage.getItem(key)
}

function remove(key: string): void {
  if (!key) throw new Error(`storage:error:remove key=${key}`)
  localStorage.removeItem(key)
}

const $localStorage = {
  set,
  get,
  remove,
}

export { $localStorage }
