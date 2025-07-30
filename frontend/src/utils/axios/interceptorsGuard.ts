import { $notify } from 'utils/notify'
import { createHeaders } from './createHeaders'
import { refreshToken } from './refresh'
import { AxiosRequestConfig } from 'axios'
import { $t } from 'i18n/index'

export function setInterceptorsGuard(axiosInstance) {
  axiosInstance.interceptors.request.use(
    (request) => {
      request.headers = Object.assign(request.headers, createHeaders())
      return request
    },

    (error) => {
      return Promise.reject(error)
    },
  )

  axiosInstance.interceptors.response.use(
    async (response) => {
      const refreshTokenOriginalRequest = await refreshToken(response)
      if (refreshTokenOriginalRequest) {
        return axiosInstance(refreshTokenOriginalRequest)
      }

      if ('errors' in response.data && response.data.errors[0]) {
        if (response.data.errors[0].message === 'comboNotFound') return response

        const extensions = response.data.errors[0].extensions
        if (extensions && extensions.errors && extensions.errors[0]) {
          $notify.show({
            type: 'error',
            title: $t('error.validationError'),
            text: $t(`error.${extensions.errors[0].code}`),
          })
        } else {
          $notify.show({
            type: 'error',
            title: $t('Error'),
            text: $t(`error.${response.data.errors[0].message}`),
          })
        }
      }

      return response
    },

    async (error) => {
      const refreshTokenOriginalRequest = await refreshToken(error)
      if (refreshTokenOriginalRequest) {
        return axiosInstance(refreshTokenOriginalRequest)
      }

      if ('message' in error) {
        $notify.show({
          type: 'error',
          title: $t('Error'),
          text: $t(`error.${error.message}`),
        })
      }

      return Promise.reject(error)
    },
  )
}
