import Hammer from 'hammerjs'
import { $screen } from 'utils/screen'
let hammer = null
let container = null
let blockInner = null

let minScale = 0
let maxScale = 0
let containerWidth
let containerHeight
let blockInnerX = 0
let blockInnerY = 0
let blockInnerScale = 0

let displayDefaultWidth
let displayDefaultHeight

let rangeX = 0
let rangeMaxX = 0
let rangeMinX = 0

let rangeY = 0
let rangeMaxY = 0
let rangeMinY = 0

let blockInnerCurrentX = 0
let blockInnerCurrentY = 0
let blockInnerCurrentScale = 0

function clamp(value, min, max) {
  return Math.min(Math.max(min, value), max)
}

function clampScale(newScale) {
  return clamp(newScale, minScale, maxScale)
}

function resizeContainer() {
  setScales()

  containerWidth = container.offsetWidth
  containerHeight = container.offsetHeight
  if (displayDefaultWidth !== undefined && displayDefaultHeight !== undefined) {
    displayDefaultWidth = blockInner.offsetWidth
    displayDefaultHeight = blockInner.offsetHeight
    updateRange()
    blockInnerCurrentX = clamp(blockInnerX, rangeMinX, rangeMaxX)
    blockInnerCurrentY = clamp(blockInnerY, rangeMinY, rangeMaxY)
    updateContainerPositions(blockInnerCurrentX, blockInnerCurrentY, blockInnerCurrentScale)
  }
}

function updateContainerPositions(x, y, scale) {
  const transform = `translateX(${x}px) translateY(${y}px) translateZ(0px) scale(${scale},${scale})`
  blockInner.style.transform = transform
}

function updateRange() {
  rangeX = Math.max(0, Math.round(displayDefaultWidth * blockInnerCurrentScale) - containerWidth)
  rangeY = Math.max(0, Math.round(displayDefaultHeight * blockInnerCurrentScale) - containerHeight)

  rangeMaxX = Math.round(rangeX / 2)
  rangeMinX = 0 - rangeMaxX

  rangeMaxY = Math.round(rangeY / 2)
  rangeMinY = 0 - rangeMaxY
}

function setBlockInner() {
  displayDefaultWidth = blockInner.offsetWidth
  displayDefaultHeight = blockInner.offsetHeight
  rangeX = Math.max(0, displayDefaultWidth - containerWidth)
  rangeY = Math.max(0, displayDefaultHeight - containerHeight)
}

function preventDefault(e) {
  e.preventDefault()
}

function eventWheel(e) {
  blockInnerScale = blockInnerCurrentScale = clampScale(blockInnerScale + e.wheelDelta / 800)
  updateRange()
  blockInnerCurrentX = clamp(blockInnerCurrentX, rangeMinX, rangeMaxX)
  blockInnerCurrentY = clamp(blockInnerCurrentY, rangeMinY, rangeMaxY)
  updateContainerPositions(blockInnerCurrentX, blockInnerCurrentY, blockInnerScale)
}
function addListeners() {
  window.addEventListener('resize', resizeContainer, true)
  container.addEventListener('wheel', eventWheel, false)
  blockInner.addEventListener('mousedown', preventDefault, false)

  hammer.get('pinch').set({ enable: true })
  hammer.get('pan').set({ direction: Hammer.DIRECTION_ALL })
  hammer.on('pan', (ev) => {
    blockInnerCurrentX = clamp(blockInnerX + ev.deltaX, rangeMinX, rangeMaxX)
    blockInnerCurrentY = clamp(blockInnerY + ev.deltaY, rangeMinY, rangeMaxY)
    updateContainerPositions(blockInnerCurrentX, blockInnerCurrentY, blockInnerScale)
  })
  hammer.on('pinch pinchmove', (ev) => {
    blockInnerCurrentScale = clampScale(ev.scale * blockInnerScale)
    updateRange()
    blockInnerCurrentX = clamp(blockInnerX + ev.deltaX, rangeMinX, rangeMaxX)
    blockInnerCurrentY = clamp(blockInnerY + ev.deltaY, rangeMinY, rangeMaxY)
    updateContainerPositions(blockInnerCurrentX, blockInnerCurrentY, blockInnerCurrentScale)
  })
  hammer.on('panend pancancel pinchend pinchcancel', () => {
    blockInnerScale = blockInnerCurrentScale
    blockInnerX = blockInnerCurrentX
    blockInnerY = blockInnerCurrentY
  })
}

function setScales() {
  if ($screen.isMobile.value) {
    blockInnerScale = 0.3
    blockInnerCurrentScale = 0.3
    minScale = 0.2
    maxScale = 1.5
  } else {
    blockInnerScale = 1
    blockInnerCurrentScale = 1
    minScale = 0.4
    maxScale = 2
  }
}

function initEvents(container_: HTMLElement, inner: HTMLElement) {
  container = container_
  blockInner = inner.querySelector('div')

  hammer = new Hammer(container)

  setScales()
  setBlockInner()
  resizeContainer()
  addListeners()
}

function destroyEvents() {
  if (!hammer) return
  hammer.destroy()
  hammer = null
  window.removeEventListener('resize', resizeContainer, true)
  container.removeEventListener('wheel', eventWheel, false)
  blockInner.removeEventListener('mousedown', preventDefault, false)
  container = null
  blockInner = null
}

function reloadEvents() {
  if (!container || !blockInner) return

  const container_ = container
  const blockInner_ = blockInner

  destroyEvents()
  initEvents(container_, blockInner_)
}

export { initEvents, destroyEvents, reloadEvents }
