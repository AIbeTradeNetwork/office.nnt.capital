import { createI18n } from 'vue-i18n'
import { US as iconUS, RU as iconRU } from 'country-flag-icons/string/3x2'
import { enUS as fnsLocaleEN, ru as fnsLocaleRU, Locale as IDateFnsLocale } from 'date-fns/locale'
import { ELocales } from 'types/enums'
import { $config } from 'utils/configuration'
import { $localStorage } from 'utils/localStorage'

const defaultLocale: ELocales = ELocales.en_US

const localesOptions: {
  locale: ELocales
  name: string
  icon: string
  map: string
  dateFnsLocale: IDateFnsLocale
}[] = [
  {
    locale: ELocales.en_US,
    name: 'English',
    icon: iconUS,
    map: 'en',
    dateFnsLocale: fnsLocaleEN,
  },
  {
    locale: ELocales.ru_RU,
    name: 'Русский',
    icon: iconRU,
    map: 'ru',
    dateFnsLocale: fnsLocaleRU,
  },
] as const

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: defaultLocale,
  fallbackLocale: ELocales.en_US,
  formatFallbackMessages: true,
  messages: {},
  pluralRules: {
    // https://docs.translatehouse.org/projects/localization-guide/en/latest/l10n/pluralforms.html
    // [ELocales.fr_FR]: function (choice) {
    //   // 1 яблоко / 99 яблок
    //   if (choice > 1) return 0
    //   return 1
    // },
    [ELocales.en_US]: function (choice) {
      // 1 яблоко / 99 яблок
      if (choice > 1) return 1
      return 0
    },
    [ELocales.ru_RU]: function (choice) {
      // 1 яблоко / 52 яблока / 99 яблок
      if (choice % 10 === 1 && choice % 100 !== 11) return 0
      if (choice % 10 >= 2 && choice % 10 <= 4 && (choice % 100 < 10 || choice % 100 >= 20))
        return 1
      return 2
    },
  },
})

const $i18nGlobal = i18n.global
const $t = i18n.global.t

function getStorageData() {
  const locale = $localStorage.get('locale')
  return (locale as ELocales) || undefined
}

function setStorageData(locale: string) {
  $localStorage.set('locale', locale)
}

function $getLocaleOptions(locale?: string, key?: keyof (typeof localesOptions)[0]) {
  const finded = localesOptions.find(
    (item) => item[key || 'locale'] === (locale || i18n.global.locale.value),
  )
  return finded || localesOptions[0]
}

function getLocale(): ELocales {
  const locale = (getStorageData() || navigator.language || defaultLocale || '').split('-')[0]

  const filteredLocaleOptions = $getLocaleOptions(locale, 'map')

  if (!filteredLocaleOptions) return defaultLocale
  return filteredLocaleOptions.locale
}

function $loadLocale(locale: string) {
  const path = `/locales/${locale}.json?${$config.version}`

  return fetch(path)
    .then((response) => {
      if (response.ok) return response.json()
      throw new Error(`error:$loadLocale:locale=${locale}:path=${path}`)
    })
    .catch((error) => {
      return Promise.reject(error)
    })
}

async function $setLocale(locale: string) {
  const lcl = locale || defaultLocale

  if (!i18n.global.availableLocales.includes(lcl as any)) {
    const messages = await $loadLocale(lcl)
    if (messages) i18n.global.setLocaleMessage(lcl, messages)
  }

  i18n.global.locale.value = lcl as any
  setStorageData(lcl)
  ;(document as any).querySelector('html').setAttribute('lang', lcl.split('-')[0])

  return Promise.resolve(true)
}

async function $i18nGuard() {
  const locale = getLocale()

  await $setLocale(locale)

  return Promise.resolve(true)
}

const $i18nOptions = {
  locales: ELocales,
  localesOptions,
  defaultLocale,
}

export {
  i18n,
  $i18nGuard,
  $setLocale,
  $loadLocale,
  $getLocaleOptions,
  $i18nGlobal,
  $i18nOptions,
  $t,
}
