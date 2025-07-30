import { $cookies } from './cookies'

enum Keys {
  ref = 'ref',
}

function getFromCookies() {
  return $cookies.get(Keys.ref)
}

function saveToCookies(ref: string) {
  if (!ref) $cookies.remove(Keys.ref)
  $cookies.set(Keys.ref, ref, {
    expires: new Date(new Date().getTime() + 60000 * 60 * 24 * 7),
  })
}

const $ref = {
  getFromCookies,
  saveToCookies,
}

export { $ref }
