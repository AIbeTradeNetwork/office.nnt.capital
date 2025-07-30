let zoom = 1
let el: HTMLElement | null = null
let parent: HTMLElement | null = null
const startPositionPinch = { x: 0, y: 0, distance: 0 }

const options = {
  maxZoom: 2,
  minZoom: 0.4,
  startZoom: 1,
  zoomFactor: {
    pc: 0.02,
    mobile: 0.01,
  },
}

function addStyles() {
  if (!el) return
  el.style.transform = `scale(${zoom})`
  el.style.transformOrigin = '0 0'
}

function addStartZoom() {
  if (+el.style.scale === zoom) return
  zoom = options.startZoom
  addStyles()
}

function selectZoomAndSave({ zoomIn, zoomOut, factor }) {
  if (zoomIn) zoom = zoom + factor
  if (zoomOut) zoom = zoom - factor
}

function wellevent(ev) {
  ev.preventDefault()

  const wheelDelta = ev.wheelDelta
  const isScrollUp = wheelDelta < 0
  const isScrollDown = wheelDelta > 0

  if (isScrollUp && zoom >= options.maxZoom) {
    zoom = options.maxZoom
    return
  }

  if (isScrollDown && zoom <= options.minZoom) {
    zoom = options.minZoom
    return
  }

  selectZoomAndSave({ zoomIn: isScrollUp, zoomOut: isScrollDown, factor: options.zoomFactor.pc })
  addStyles()
}

function distancePinch(ev) {
  return Math.hypot(
    ev.touches[0].pageX - ev.touches[1].pageX,
    ev.touches[0].pageY - ev.touches[1].pageY,
  )
}

function touchStart(ev) {
  if (ev.touches.length === 2) {
    ev.preventDefault()
    startPositionPinch.x = (ev.touches[0].pageX + ev.touches[1].pageX) / 2
    startPositionPinch.y = (ev.touches[0].pageY + ev.touches[1].pageY) / 2
    startPositionPinch.distance = distancePinch(ev)
  }
}

function touchMove(ev) {
  if (ev.touches.length === 2) {
    ev.preventDefault()

    const zoomIn = distancePinch(ev) > startPositionPinch.distance
    const zoomOut = distancePinch(ev) < startPositionPinch.distance

    if (zoomIn && zoom >= options.maxZoom) {
      zoom = options.maxZoom
      return
    }

    if (zoomOut && zoom <= options.minZoom) {
      zoom = options.minZoom
      return
    }

    selectZoomAndSave({ zoomIn, zoomOut, factor: options.zoomFactor.mobile })

    addStyles()
  }
}

function addEvents() {
  if (!el) return
  if (navigator.maxTouchPoints <= 1) {
    parent.addEventListener('mousewheel', wellevent, true)
  } else {
    parent.addEventListener('touchstart', touchStart, true)
    parent.addEventListener('touchmove', touchMove, true)
  }
}

function removeEvents() {
  if (!el) return
  parent.addEventListener('mousewheel', wellevent, true)
  parent.addEventListener('touchstart', touchStart, true)
  parent.addEventListener('touchmove', touchMove, true)
}

function zoomInit(element: HTMLElement) {
  if (!element) return
  el = element
  parent = element.parentElement
  addStartZoom()
  addEvents()
}

function zoomDestroy() {
  removeEvents()
}

function resetZoom() {
  zoom = options.startZoom
  addStyles()
}

export { zoomInit, zoomDestroy, resetZoom }
