const scrpt = document.getElementById('translations')

/**
 * Configuration options for locale settings.
 * @type {Object}
 * @property {string} fallbackLocale - The default language standard to return to. It appears locally index.html or autotrade.html.
 * @property {string} file - The initial file name for building a file name indicating the translation.
 * @property {string} current - The current language of the document.
 * @property {string} locale - The current language setting of the document.
 * @property {Array<string>} locales - List of supported locale codes.
 * @property {Object} map - Mapping of locale codes to their respective display information.
 */
const options = {
  fallbackLocale: 'en',
  file: scrpt.getAttribute('file'),
  current: scrpt.getAttribute('current'),
  locale: '',
  locales: scrpt.getAttribute('locales').split(','),
  selectors: {
    container: 'translation-selector',
    store: 'ecosystem_locale',
  },
  map: {
    // flags https://flagicons.lipis.dev/
    ru: {
      code: 'ru',
      flag: 'ru',
      name: 'Русский',
    },
    en: {
      code: 'en',
      flag: 'gb',
      name: 'English',
    },
    es: {
      code: 'es',
      flag: 'es',
      name: 'Español',
    },
  },
}

const store = {
  get locale() {
    return (localStorage.getItem(options.selectors.store) || '').split(/-|_/)[0]
  },
  set locale(locale) {
    ;(localStorage.setItem(options.selectors.store, locale) || '').split(/-|_/)[0]
  },
}

const browser = {
  get locale() {
    return (navigator.language || navigator.userLanguage || '').split(/-|_/)[0]
  },
}

function getLocaleMap(locale) {
  const lcl = locale || options.locale
  if (!lcl) return {}
  if (!options.map[lcl]) return {}
  return options.map[lcl]
}

function initLocale() {
  const locale = store.locale || browser.locale || options.fallbackLocale
  options.locale = locale
}

function checkLocale() {
  if (!options.locales.includes(options.locale)) {
    setLocale(options.fallbackLocale)
    return
  }

  if (options.current !== options.locale) {
    setLocale(options.locale)
    return
  }
}

async function setLocale(locale) {
  if (!locale) return
  const match = location.pathname.match(/(\/.*)\/\w+\.html$/)
  const path = match ? match[1] : location.pathname

  function fileName() {
    if (options.fallbackLocale === locale) return options.file
    return `${options.file}_${locale}`
  }

  const filePath = `${path}/${fileName()}.html`
  store.locale = locale
  location.pathname = filePath
}

function addTranslationSelector() {
  const containers = document.querySelectorAll(`.${options.selectors.container}`)
  containers.forEach((container) => {
    // container.style.position = "relative";

    const flag = document.createElement('div')
    // flag.style.width = "30px";
    // flag.style.height = "30px";
    flag.classList.add('fi', 'fi-' + getLocaleMap(options.locale).flag)
    container.appendChild(flag)

    const select = document.createElement('select')
    select.style.position = 'absolute'
    select.style.top = '0'
    select.style.left = '0'
    select.style.width = '100%'
    select.style.height = '100%'
    select.style.opacity = '0'
    select.style.cursor = 'pointer'
    container.appendChild(select)
    select.addEventListener('change', (event) => {
      setLocale(event.target.value)
    })

    for (const index in options.locales) {
      const locale = options.locales[index]
      const option = document.createElement('option')
      option.selected = getLocaleMap(locale).code === options.locale
      option.value = getLocaleMap(locale).code
      option.textContent = getLocaleMap(locale).name
      select.appendChild(option)
    }
  })
}

// INITS
initLocale()
checkLocale()

document.addEventListener('DOMContentLoaded', () => {
  addTranslationSelector()
})
