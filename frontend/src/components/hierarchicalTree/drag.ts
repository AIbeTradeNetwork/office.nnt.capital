let el: HTMLElement | null = null
let parent: HTMLElement | null = null
let isWorking = false
let resizeTimer: ReturnType<typeof setTimeout>
const elStartPosition = { x: 0, y: 0 }
const startPositionMouse = { x: 0, y: 0 }
const mousePosition = { x: 0, y: 0 }
const isMobile = () => navigator.maxTouchPoints > 1

function resize() {
  if (!el) return
  clearTimeout(resizeTimer)
  resizeTimer = setTimeout(() => {
    resetPosition()
  }, 1000)
}

function pointerDown(ev) {
  if (isMobile() && ev.touches.length > 1) return

  startPositionMouse.x =
    isMobile() ? ev.touches[0].clientX : Math.abs(parent.offsetLeft - ev.clientX)
  startPositionMouse.y =
    isMobile() ? ev.touches[0].clientY : Math.abs(parent.offsetTop - ev.clientY)
  elStartPosition.x = el.getBoundingClientRect().left - parent.getBoundingClientRect().left
  elStartPosition.y = el.getBoundingClientRect().top - parent.getBoundingClientRect().top
  isWorking = true
}

function pointerUp() {
  startPositionMouse.x = 0
  startPositionMouse.y = 0
  isWorking = false
}

function pointerMove(ev) {
  if (!el) return
  if (isMobile() && ev.touches && ev.touches.length > 1) return

  ev.preventDefault()

  mousePosition.x = isMobile() ? ev.touches[0].clientX : Math.abs(parent.offsetLeft - ev.clientX)
  mousePosition.y = isMobile() ? ev.touches[0].clientY : Math.abs(parent.offsetTop - ev.clientY)

  if (!isMobile() && !isWorking) return

  const offsetSize = {
    x: parent.getBoundingClientRect().width / 3,
    y: parent.getBoundingClientRect().height / 3,
  }

  const moveX = mousePosition.x - startPositionMouse.x
  const moveY = mousePosition.y - startPositionMouse.y

  const elPositionX = elStartPosition.x + moveX
  const elPositionY = elStartPosition.y + moveY

  const bountLeft = elPositionX + el.getBoundingClientRect().width > offsetSize.x
  const bountRight = elPositionX < parent.getBoundingClientRect().width - offsetSize.x

  const boundTop = elPositionY + el.getBoundingClientRect().height > offsetSize.y
  const boundBottom = elPositionY < parent.getBoundingClientRect().height - offsetSize.y

  if (bountLeft && bountRight) el.style.left = `${elPositionX}px`
  if (boundTop && boundBottom) el.style.top = `${elPositionY}px`
  el.style.transformOrigin = '0 0'
}

function parentPointerleave() {
  pointerUp()
}

function addEvents() {
  if (!el) return
  if (!parent) return

  window.addEventListener('resize', resize, true)

  if (isMobile()) {
    parent.addEventListener('touchstart', pointerDown, true)
    parent.addEventListener('touchmove', pointerMove, true)
  } else {
    parent.addEventListener('mouseleave', parentPointerleave, true)
    parent.addEventListener('mousedown', pointerDown, true)
    parent.addEventListener('mouseup', pointerUp, true)
    parent.addEventListener('mousemove', pointerMove, true)
  }
}

function removeEvents() {
  window.removeEventListener('resize', resize, true)

  parent.removeEventListener('mouseleave', parentPointerleave, true)
  parent.removeEventListener('mousedown', pointerDown, true)
  parent.removeEventListener('mouseup', pointerUp, true)
  parent.removeEventListener('mousemove', pointerMove, true)

  parent.addEventListener('touchsend', pointerDown, true)
  parent.addEventListener('touchstart', pointerUp, true)
  parent.addEventListener('touchmove', pointerMove, true)
}

function initStartPosition() {
  if (!el) return
  const left = parent.getBoundingClientRect().width / 2 - el.getBoundingClientRect().width / 2
  elStartPosition.x = left
  el.style.left = left + 'px'
  el.style.top = '0px'
}

function dragInit(element: HTMLElement) {
  if (!element) return
  el = element
  parent = element.parentElement
  addEvents()
  initStartPosition()
}

function dragDestroy() {
  removeEvents()
}

function resetPosition() {
  if (!el) return
  elStartPosition.x = 0
  elStartPosition.y = 0
  mousePosition.x = 0
  mousePosition.y = 0
  startPositionMouse.x = 0
  startPositionMouse.y = 0
  initStartPosition()
}

export { dragInit, dragDestroy, resetPosition }
