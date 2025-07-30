import cookie from 'cookie'
import { $config } from 'utils/configuration'

interface IOptions {
  domain?: string
  encode?: boolean
  path?: string
  expires?: Date
  httpOnly?: boolean
  maxAge?: number
  priority?: 'low' | 'medium' | 'high'
  sameSite?: boolean | 'lax' | 'none' | 'strict'
  secure?: boolean
}

const defOptions: IOptions = {
  path: '/',
  expires: undefined,
  maxAge: undefined,
  sameSite: $config.isDEV ? undefined : 'lax',
  secure: false,
}

function mixOptions(options?) {
  return Object.assign(defOptions, options || {})
}

function get(name = ''): string | undefined {
  let str = {}
  // if ($global.isSSR) {
  //   str = cookie.parse(useSSRContext()?.req.headers.cookie)
  // }
  str = cookie.parse(document.cookie)
  return str[name] || null
}

function set(name: string, value: unknown, options?: IOptions) {
  // if ($global.isSSR) {
  //   // useSSRContext()?.res.setHeader('Set-Cookie', cookie.serialize(name, String(value), this.getOptions(options)))
  //   useSSRContext()?.res.cookie(name, String(value), this.getOptions(options))
  // }
  document.cookie = cookie.serialize(name, String(value), mixOptions(options))
}

function remove(value: string) {
  // if ($global.isSSR) {
  //   useSSRContext()?.res.cookie(name, '', options)
  // }
  document.cookie = cookie.serialize(value, '', mixOptions({ maxAge: -1 }))
}

const $cookies = {
  get,
  set,
  remove,
}

export { $cookies }
