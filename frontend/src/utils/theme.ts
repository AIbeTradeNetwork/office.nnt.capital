import { Ref, ref, watch } from 'vue'
import { $oklchToLch } from 'utils/oklch'
import { $localStorage } from './localStorage'

interface IThemeColors {
  primary: string
  primaryFocus: string
  primaryContent: string
  secondary: string
  secondaryFocus: string
  secondaryContent: string
  accent: string
  accentFocus: string
  accentContent: string
  neutral: string
  neutralFocus: string
  neutralContent: string
  base100: string
  base200: string
  base300: string
  baseContent: string
  info: string
  infoContent: string
  success: string
  successContent: string
  warning: string
  warningContent: string
  error: string
  errorContent: string
}

enum EThemes {
  light = 'light',
  dark = 'dark',
  auto = 'auto',
}

const current = ref<EThemes>(EThemes.dark)
// Ref<EThemes> = ref(($localStorage.get('theme') as EThemes) ||

// watch(current, () => {
//   set()
//   $localStorage.set('theme', current.value)
// })

const themeMqLight =
  (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)')) || false

function activateTheme(theme: EThemes) {
  // current.value = theme
  // document.querySelector('html').setAttribute('data-theme', theme)
}

function set() {
  // if (!(current.value in EThemes)) {
  //   current.value === EThemes.auto
  // }
  // if (current.value === EThemes.auto) {
  //   if (themeMqLight && themeMqLight.matches) {
  //     activateTheme(EThemes.light)
  //     return
  //   }
  //   activateTheme(EThemes.dark)
  //   return
  // }
  // activateTheme(current.value)
}

function watchThemeEventChange() {
  // if (!themeMqLight) return
  // themeMqLight.onchange = function () {
  //   set()
  // }
}

function colors() {
  const computedStyles = getComputedStyle(document.querySelector(':root'))
  const conver = (value) => $oklchToLch(`oklch(${computedStyles.getPropertyValue(value)})`)
  return {
    primary: conver('--p'),
    primaryFocus: conver('--pf'),
    primaryContent: conver('--pc'),
    secondary: conver('--s'),
    secondaryFocus: conver('--sf'),
    secondaryContent: conver('--sc'),
    accent: conver('--a'),
    accentFocus: conver('--af'),
    accentContent: conver('--ac'),
    neutral: conver('--n'),
    neutralFocus: conver('--nf'),
    neutralContent: conver('--nc'),
    base100: conver('--b1'),
    base200: conver('--b2'),
    base300: conver('--b3'),
    baseContent: conver('--bc'),
    info: conver('--in'),
    infoContent: conver('--inc'),
    success: conver('--su'),
    successContent: conver('--suc'),
    warning: conver('--wa'),
    warningContent: conver('--wac'),
    error: conver('--er'),
    errorContent: conver('--erc'),
  }
}

function guard() {
  // set()
  // watchThemeEventChange()
}

const $theme = {
  current,
  set,
  colors,
  guard,
  themes: EThemes,
}

export { $theme }
