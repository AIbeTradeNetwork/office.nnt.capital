import { $config } from './configuration'

const versionKey = 'version'

export function $checkRefreshGuard() {
  const oldVer = localStorage.getItem('version')
  const curCer = $config.version

  if (oldVer !== curCer) {
    localStorage.setItem('version', curCer)
    if (localStorage.getItem(versionKey)) {
      window.location.reload()
    }
  }
}
