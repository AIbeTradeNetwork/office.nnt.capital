import { $config } from 'utils/configuration'
import { $localStorage } from 'utils/localStorage'
import { $ref } from './ref'
import { $authentication } from './authentication'
import { $modals } from './modals'
import { nextTick } from 'vue'
import { Router } from 'routes/index'

function clearQuery(query: string) {
  if (!query) return
  const queryCopy = Object.assign({}, Router.currentRoute.value.query)
  delete queryCopy[query]
  Router.push({ query: queryCopy })
}

function checkToken() {
  const params = new URLSearchParams(document.location.search)
  if (params.get('token')) {
    $authentication.setToken(params.get('token'))
    clearQuery('token')
  }
}

function checkModal() {
  const params = new URLSearchParams(document.location.search)
  if (params.get('modal') && $modals[params.get('modal')]) {
    setTimeout(() => {
      $modals[params.get('modal')].show({})
      clearQuery('modal')
    }, 2000)
  }
}

function checkRef() {
  const params = new URLSearchParams(document.location.search)
  if (params.get('ref')) $ref.saveToCookies(params.get('ref'))
}

function checkDevMode(code?: string) {
  const key = 'devmode'
  const params = new URLSearchParams(document.location.search)
  const devMode = code || params.get(key) || $localStorage.get(key)

  if (!devMode) return false

  if (devMode === '1') {
    $localStorage.set(key, '1')
    $config.isDevMode = true
    return
  }

  if (devMode === '0') {
    $localStorage.remove(key)
    $config.isDevMode = false
    return
  }

  $config.isDevMode = false
  $localStorage.remove(key)
}

function $queriesGuards() {
  checkToken()
  checkRef()
  checkDevMode()
  checkModal()
}

export { $queriesGuards, checkDevMode }
