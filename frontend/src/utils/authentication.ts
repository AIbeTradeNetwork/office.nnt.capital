import { ref } from 'vue'
import { $cookies } from 'utils/cookies'
import { Router } from 'routes/index'
import { $jwtDecode } from './jwtDecode'

enum Keys {
  token = 'token',
  refreshToken = 'refresh_token',
}

const isAuthenticated = ref(false)

function setToken(token) {
  if (!token) return
  $cookies.set(Keys.token, token, {
    expires: new Date($jwtDecode(token).exp * 1000),
  })
  if (token) isAuthenticated.value = true
  return
}

function getToken() {
  return $cookies.get(Keys.token)
}

function getRefreshToken() {
  return $cookies.get(Keys.refreshToken)
}

function setRefreshToken(token) {
  if (!token) return
  $cookies.set(Keys.refreshToken, token, {
    expires: new Date($jwtDecode(token).exp * 1000),
  })
  return
}

function $authenticationGuard() {
  getToken() ? (isAuthenticated.value = true) : (isAuthenticated.value = false)
}

function logIn(data: AuthRes) {
  if (!data) return

  if (data.refresh_token) {
    setRefreshToken(data.refresh_token)
  }

  if (data.auth_token) {
    setToken(data.auth_token)
    isAuthenticated.value = true
    Router.push({ name: 'Index' })
  }
}

function logOut() {
  $cookies.remove(Keys.token)
  $cookies.remove(Keys.refreshToken)
  isAuthenticated.value = false
  location.reload()
}

const $authentication = {
  isAuthenticated,
  logIn,
  logOut,
  getToken,
  setToken,
  getRefreshToken,
  setRefreshToken,
}

export { $authentication, $authenticationGuard }
