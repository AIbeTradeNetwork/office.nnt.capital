export function $useTelegram() {
  const telegram = (window as any)?.Telegram

  return {
    telegram,
    inTelegramApp: telegram?.WebApp?.initData,
    inTelegramAppIOS: false, //telegram?.WebApp?.initData && telegram?.WebApp.platform === 'ios',
    webApp: telegram?.WebApp,
  }
}
