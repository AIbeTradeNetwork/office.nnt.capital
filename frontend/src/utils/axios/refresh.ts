import { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'
import { $requests } from 'queries/index'

export async function refreshToken(
  o: AxiosResponse | AxiosError,
): Promise<AxiosRequestConfig | null> {
  const originalRequest: AxiosRequestConfig & { _retry? } = { ...o.config }

  const tryRefreshToken = async function () {
    originalRequest._retry = true

    const tokens = await $requests.auth.refresh()

    if (tokens?.auth_token) {
      originalRequest.headers['Authorization'] = 'Bearer ' + tokens.auth_token

      return originalRequest
    }

    return null
  }

  if ('response' in o && 'status' in o.response && !originalRequest._retry) {
    if ([403, 401].includes(o.response.status)) {
      return await tryRefreshToken()
    }
  }

  if ('data' in o && 'errors' in o.data && !originalRequest._retry) {
    const isNeedRefresh = o.data.errors.find((item) => {
      if (item.message.match(/AccessDenied/gi)) return item
    })

    if (isNeedRefresh) return await tryRefreshToken()
  }

  return null
}
