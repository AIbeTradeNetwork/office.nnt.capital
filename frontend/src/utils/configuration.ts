import { computed, reactive, ref, Ref } from 'vue'
import { ECurrencies, ERoles } from 'types/enums'
import { $store } from './store'

const $config = reactive({
  version: import.meta.env.VITE_BUILD_VERSION || '---',
  baseUrl: import.meta.env.BASE_URL,
  appDomain: 'aibetrade.com',
  appName: 'AiBeTrade',
  supportLink: 'abtsupportbot',
  mode: import.meta.env.MODE,
  isDEV: import.meta.env.DEV,
  isPROD: import.meta.env.PROD,
  role: ERoles.client,
  currency: ECurrencies.udex,
  isDevMode: true,
  tonAddress: '',
})

export { $config }
