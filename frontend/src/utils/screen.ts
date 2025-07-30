import { computed, ref } from 'vue'
import { $throttle } from './throttle'

const defaults = {
  zr: 0,
  sm: 640,
  md: 768,
  lg: 1279,
  xl: 1280,
}

const currentSize = ref(0)

const isMobile = computed(() => currentSize.value >= 0 && currentSize.value < defaults.lg)
const isPc = computed(() => currentSize.value >= defaults.lg)

function resizeHandler(event: Event | Window) {
  if (!event) return
  if ('target' in event) {
    currentSize.value = (event.target as Window).innerWidth
  } else {
    currentSize.value = event.innerWidth
  }
}

const throttleResize = $throttle(resizeHandler, 600)

function $screenGuard() {
  window.addEventListener('resize', throttleResize, true)
  resizeHandler(window)
}

function $defaultsScreenToTilwind() {
  const defaultsCopy = Object.assign({}, defaults)

  for (const item in defaultsCopy) {
    defaultsCopy[item] = defaultsCopy[item] + 'px'
  }

  return defaultsCopy
}

const $screen = {
  isMobile,
  isPc,
}

export { $screen, $defaultsScreenToTilwind, $screenGuard }
