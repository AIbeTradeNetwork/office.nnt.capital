import {
  TonConnectUI,
  THEME,
  WalletsModalState,
  TonConnectUiOptionsWithManifest,
  Locales,
  ConnectedWallet,
  TonConnectError,
} from '@tonconnect/ui'
import { Address } from '@ton/core'
import { $getLocaleOptions, $t } from 'i18n/index'
import { getTonPayload, setTonWallet } from 'queries/ton'
import { $config } from 'utils/configuration'
import { $me } from 'utils/me'
import { $notify } from 'utils/notify'
import { ref } from 'vue'
import { $screen } from 'utils/screen'
import { $formatToInt } from 'utils/formats'

let tonConnectUI: TonConnectUI
const isLoading = ref(true)
const isConnected = ref(false)

export function useTonConnect() {
  const uiOptions: TonConnectUiOptionsWithManifest = {
    language: $getLocaleOptions().locale as unknown as Locales,
    enableAndroidBackHandler: false,
    uiPreferences: {
      theme: THEME.DARK,
    },
    manifestUrl: 'https://office.aibetrade.com/tonconnect-manifest.json',
  }

  function checkErrors() {
    if (!tonConnectUI) {
      $notify.show({
        type: 'error',
        text: 'TonConnectUI is not initialized',
      })
      return true
    }
    return false
  }

  const storage = {
    getAll() {
      return localStorage
    },

    clientTonInformation() {
      Object.keys(localStorage).forEach((key) => {
        if (key.match(/ton-/)) {
          return localStorage.removeItem(key)
        }
      })
    },
  }

  function initTonConnectUI() {
    tonConnectUI = new TonConnectUI(uiOptions)
  }

  const modal = {
    _subscribed: () => {},
    async open() {
      if (checkErrors()) return
      tonConnectUI.openModal()
    },

    close() {
      if (checkErrors()) return
      tonConnectUI.closeModal()
    },

    getState() {
      if (checkErrors()) return
      return tonConnectUI.modalState
    },
    unsubscribe() {
      if (this._subscribed) this._subscribed()
    },
    subscribe(fn?) {
      this.unsubscribe()
      if (checkErrors()) return
      this._subscribed = tonConnectUI.onModalStateChange((state: WalletsModalState) => {
        if (fn) fn(state)
      })
    },
  }

  const wallet = {
    _subscribed: () => {},
    subscribedFunctions: [],
    current() {
      if (checkErrors()) return
      if (!isConnected.value) {
        return
      }
      return tonConnectUI.wallet
    },
    info() {
      if (checkErrors()) return
      if (!isConnected.value) {
        return
      }
      // @ts-expect-error: ts(2339)
      return tonConnectUI.walletInfo
    },
    account() {
      if (checkErrors()) return
      if (!isConnected.value) {
        return
      }
      return tonConnectUI.account
    },

    isConnectionRestored: async () => {
      return await tonConnectUI.connectionRestored
    },

    isConnected() {
      if (tonConnectUI && tonConnectUI.connected) {
        isConnected.value = true
      } else {
        isConnected.value = false
      }

      return isConnected
    },
    open() {
      const w = tonConnectUI.wallet
      // @ts-expect-error: ts(2339)
      const link = $screen.isMobile ? w.deepLink : w.universalLink
      if (link) window.open(link, '_blank')
    },
    getBounceableAddress() {
      if (checkErrors()) return
      if (!isConnected.value) {
        return
      }
      return Address.parse(wallet.account().address).toString({ bounceable: true })
    },
    getNonBounceableAddress() {
      if (checkErrors()) return
      if (!isConnected.value) {
        return
      }
      return Address.parse(wallet.account().address).toString({ bounceable: false })
    },
    async disconnect() {
      if (checkErrors()) return
      storage.clientTonInformation()
      await tonConnectUI.disconnect()
    },
    unsubscribe() {
      if (this._subscribed) this._subscribed()
    },
    async transferToAPersonalAccount(o: { amount: number }) {
      if (!wallet.isConnected().value) {
        return $notify.show({
          text: $t('error.walletNotConnected'),
          type: 'error',
        })
      }

      if (!o.amount) {
        return $notify.show({
          text: 'Error: Amount not defined',
          type: 'error',
        })
      }

      const transaction = {
        validUntil: Math.floor(Date.now() / 1000) + 60, // 60 sec
        messages: [
          {
            address: $config.tonAddress,
            amount: $formatToInt(o.amount, { precision: 9 }) + '',
          },
        ],
      }

      try {
        const result = await tonConnectUI.sendTransaction(transaction)

        // const myHeaders = new Headers()
        // myHeaders.append('Content-Type', 'application/json')
        // myHeaders.append('Accept', 'application/json')

        // const requestOptions = {
        //   method: 'POST',
        //   headers: myHeaders,
        //   body: result.boc,
        //   redirect: 'follow',
        // }

        // // @ts-expect-error: ts(2322)
        // fetch('https://test.ton.org/api/v2/sendBoc', requestOptions)
        //   .then((response) => response.text())
        //   .then((result) => console.log('fetch', result))
        //   .catch((error) => console.error('fetch', error))

        return result
      } catch (error) {
        let err = 'Error: ton pay'

        if (error instanceof Error) {
          err = error.message
        }

        $notify.show({
          text: err,
          type: 'error',
        })

        console.error(error)

        return Promise.reject(error)
      }
    },
    subscribe() {
      this.unsubscribe()
      if (checkErrors()) return
      this._subscribed = tonConnectUI.onStatusChange(
        (connectedWallet: ConnectedWallet) => {
          ;(async () => {
            wallet.isConnected()

            for (const fn of this.subscribedFunctions) {
              if (fn && fn.constructor.name === 'AsyncFunction') {
                await fn(connectedWallet)
              } else {
                throw new Error('async function expected')
              }
            }
          })()
        },
        (error: TonConnectError) => {
          $notify.show({
            type: 'warning',
            text: error.message,
          })
        },
      )
    },
  }

  const events = {
    _prefix: 'ton-connect-ui-',
    functions: [],
    list: [
      'connection-started',
      'connection-completed',
      'connection-error',
      'connection-restoring-started',
      'connection-restoring-completed',
      'connection-restoring-error',
      'disconnection',
      'transaction-sent-for-signature',
      'transaction-signed',
      'transaction-signing-failed',
    ],
    _eventMethod(event) {
      events.functions.forEach((fn) => {
        fn(event)
      })
    },
    removeEvents() {
      this.list.forEach((name) => {
        window.removeEventListener(`${this._prefix}${name}`, this._eventMethod)
      })
    },
    addEvents() {
      this.removeEvents()
      this.list.forEach((name) => {
        window.addEventListener(`${this._prefix}${name}`, this._eventMethod)
      })
    },
  }

  function checkMatchMeWalletAndConnectedWallet() {
    function disconnect() {
      isConnected.value = false
      // wallet.unsubscribe()
      wallet.disconnect()
    }

    const connectedWallet = wallet.getBounceableAddress()

    if (!$me.data.ton_wallet) return false

    if ($me.data.ton_wallet !== connectedWallet) {
      disconnect()
      $notify.show({
        type: 'error',
        text: $t('error.walletExists'),
        see: () => {},
      })
      return true
    }

    if ($me.data.ton_wallet === connectedWallet) return true

    return false
  }

  async function registerWallet(walletData) {
    if (!wallet.isConnected().value) return 1

    try {
      await $me.update()

      if (checkMatchMeWalletAndConnectedWallet()) return 1

      if (
        !walletData ||
        !walletData.connectItems ||
        !('proof' in walletData.connectItems.tonProof)
      ) {
        return 1
      }

      const tonProof = Object.assign({}, walletData.connectItems.tonProof)
      tonProof.address = walletData.account.address
      tonProof.proof.domain = tonProof.proof.domain.value
      tonProof.network = walletData.account.chain
      tonProof.proof.state_init = walletData.account.walletStateInit

      const tonProofPayload = await setTonWallet(JSON.stringify(tonProof))

      if (!tonProofPayload) console.error('registerWallet tonProofPayload empty')
    } catch (error) {
      $notify.show({
        error: error,
      })
      console.error(error)
    }

    return 1
  }

  async function walletAuth() {
    tonConnectUI.setConnectRequestParameters({ state: 'loading' })

    const tonProofPayload = await getTonPayload()

    if (!tonProofPayload) {
      tonConnectUI.setConnectRequestParameters(null)
      // tonProofPayloadStatusChange()
      throw new Error('tonProofPayload is empty')
    }

    tonConnectUI.setConnectRequestParameters({
      state: 'ready',
      value: { tonProof: tonProofPayload },
    })
  }

  ;(async () => {
    if (tonConnectUI) return

    isLoading.value = true

    wallet.subscribedFunctions.push(registerWallet)
    storage.clientTonInformation()

    initTonConnectUI()

    // wallet.disconnect()
    events.addEvents()
    modal.subscribe()
    wallet.subscribe()

    await walletAuth()

    isLoading.value = false
  })()

  return {
    tonConnectUI,
    uiOptions,
    modal,
    wallet,
    events,
  }
}
