import axios, { AxiosResponse, type AxiosRequestConfig, type AxiosInstance } from 'axios'
import { setInterceptorsGuard } from './axios/interceptorsGuard'
import { axiosGetCache, axiosSetCache } from './axios/cache'

type TCache = boolean | 'update'

// let isRefreshState = false
const $axiosInstance = axios.create() as AxiosInstance

setInterceptorsGuard($axiosInstance)

/**
 * Asynchronous function for making a request.
 *
 * @param {Object} data - Object containing requestConfig, query, variables, and cache
 * @return {Promise<T | AxiosResponse<T>>} Promise of type T or AxiosResponse of type T
 */
async function $request<T>(data: {
  requestConfig?: AxiosRequestConfig
  query: string
  variables?: { [key: string]: any }
  cache?: TCache
}) {
  const config: AxiosRequestConfig & { cache?: TCache | { time: number } } = {
    method: 'POST',
    url: import.meta.env.VITE_APP_URL_GQL,
    data: {
      query: data.query,
      variables: data.variables,
    },
    ...(data.requestConfig || {}),
  }

  if (data.cache && data.cache !== 'update') {
    const cacheData = axiosGetCache(config.data)
    if (cacheData) return Promise.resolve(cacheData)
  }

  try {
    const response = (await $axiosInstance(config)) as AxiosResponse<T>

    if (data.cache) axiosSetCache(config.data, response.data)

    return response
  } catch (error) {
    return Promise.reject(error)
  }
}

export { $request }
