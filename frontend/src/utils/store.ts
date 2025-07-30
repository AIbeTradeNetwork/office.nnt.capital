import { $requests } from 'queries/index'
import { EStrategyType } from 'types/enums'
import { Ref, reactive, ref } from 'vue'

const store = reactive({
  exchanges: <Exchange[]>[],
  keys: <Key[]>[],
  strategies: <Strategy[]>[],
  currencies: <Currency[]>[],
  notifications: <INotification[]>[],
  bots: <Bot[]>[],
  cfg: <Cfg>{},
  brokerBalance: <Broker.Balance>{},
  products: <Product[]>[],
})

function set<K extends keyof typeof store>(key: K, value: (typeof store)[K]): void {
  store[key] = value
}

function get<K extends keyof typeof store>(key: K): (typeof store)[K] {
  return store[key]
}

async function updateExchanges() {
  return await $requests.exchanges.get()
}

async function updateKeys() {
  return await $requests.keys.get()
}

async function updateStrategies(type?: EStrategyType) {
  if (type === EStrategyType.classic) {
    return await $requests.strategies.get({ type: EStrategyType.classic })
  }
  if (type === EStrategyType.ico) {
    return await $requests.strategies.get({ type: EStrategyType.ico })
  }
  return Promise.allSettled([
    $requests.strategies.get({ type: EStrategyType.classic }),
    $requests.strategies.get({ type: EStrategyType.ico }),
  ])
}

async function updateCurrencies() {
  return await $requests.currencies.get()
}

async function updateBots() {
  return await $requests.bots.get()
}

async function updateNotifications() {
  return await $requests.notifications.get()
}

async function updateCfg() {
  return await $requests.cfg.get()
}

async function updateProducts() {
  return await $requests.products.get({})
}

async function clearProducts() {
  set('products', [])
}

async function updateBrokerBalance() {
  return await $requests.broker.balance()
}

const $store = {
  store,
  set,
  get,
  updateExchanges,
  updateKeys,
  updateStrategies,
  updateCurrencies,
  updateBots,
  updateCfg,
  updateNotifications,
  updateProducts,
  clearProducts,
  updateBrokerBalance,
}

export { $store }
