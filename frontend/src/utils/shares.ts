import { $me } from './me'

export const tgAppLink = 'https://t.me/nntcapitalbot'

function createShareLink() {
  // const WebApp = (window as any)?.Telegram?.WebApp
  // if (WebApp && WebApp.initData) {
  //   return `tg://msg_url?url=https://t.me/nntcapitalbot/app?startapp=${$me.data.uid}`
  // } else {

  // }
  return `https://t.me/share/url?url=${tgAppLink}/nnt?startapp=${$me.data.uid}`
}

export function openShareLink(only?: 'web' | 'telegramm') {
  const WebApp = (window as any)?.Telegram?.WebApp

  function openWeb() {
    const a = document.createElement('a')
    a.href = createShareLink()
    a.target = '_blank'
    a.click()
  }

  function openTelegramm() {
    WebApp.openTelegramLink(createShareLink())
  }

  if (only && only === 'web') {
    openWeb()
    return
  }

  if (only && only === 'telegramm') {
    openTelegramm()
    return
  }

  if (WebApp && WebApp.initData) {
    openTelegramm()
  } else {
    openWeb()
  }
}

export function createTelegrammRefLink(ref: string) {
  return `${tgAppLink}/nnt?startapp=${ref}`
}

export function createRefLink(only?: 'web' | 'telegramm') {
  const WebApp = (window as any)?.Telegram?.WebApp

  function createWeb() {
    return `${location.origin}/?ref=${$me.data.uid}`
  }

  if (only && only === 'web') return createWeb()
  if (only && only === 'telegramm') return createTelegrammRefLink($me.data.uid)
  if (WebApp && WebApp.initData) return `${tgAppLink}/app?startapp=${$me.data.uid}`
  return `${location.origin}/?ref=${$me.data.uid}`
}
