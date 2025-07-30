import { ComponentInternalInstance, ref } from 'vue'
import { EMenu } from './enums'

import { CountUp, CountUpOptions } from 'countup.js'
import { conffetiStart, conffetiStop } from './confetti/confetti'

export const menu = {
  active: ref(EMenu.farm),
  set(menu: EMenu) {
    this.active.value = menu
    window.scrollTo({ top: 0, behavior: 'smooth' })
  },
}

export function useCounter(instance: ComponentInternalInstance) {
  const counter = {
    countUp: undefined,
    element: undefined,
    currentValue: ref(0),
    options: {
      start: 1,
      end: 100,
      decimalPlaces: 4,
      duration: 5,
      separator: '.',
    },
    start(duration: number, precision: number, start: number, end: number) {
      if (!start || !end) return
      this.options.duration = duration
      this.options.decimalPlaces = precision
      this.options.start = start
      this.options.end = end
      this.element = instance.refs.counterElement as HTMLElement
      this.countUp = new CountUp(this.element, this.options.end, {
        startVal: this.options.start,
        decimalPlaces: this.options.decimalPlaces,
        duration: this.options.duration,
        separator: this.options.separator,
        useEasing: false,
      })
      this.countUp.start()
    },
    getValue() {
      if (!this.countUp) return
      this.currentValue.value = this.countUp.frameVal
    },
    update(end: number) {
      if (!this.countUp) return
      this.options.end = end
      this.countUp.update(this.options.end)
    },
    stop() {
      if (!this.countUp) return
      this.countUp.pauseResume()
      this.countUp = null
    },
  }

  return counter
}

export function useConffeti(instance: ComponentInternalInstance) {
  return {
    element: null,
    working: false,
    start() {
      if (this.working) return
      this.element = instance.refs.conffetiElement
      if (!this.element) return
      conffetiStart(this.element)
      this.working = true
    },
    stop() {
      if (!this.element) return
      this.element = null
      this.working = false
      conffetiStop()
    },
  }
}

export function useTimeout(instance: ComponentInternalInstance) {
  return {
    timerId: null,
    timeout: 1000,
    functions: [],
    initFuctions() {
      this.functions.forEach((fn) => fn())
    },
    tick() {
      this.timerId = setTimeout(() => {
        this.initFuctions()
        if (this.timerId) this.tick()
      }, this.timeout)
    },
    start(timeout = 1000) {
      this.timeout = timeout
      this.initFuctions()
      this.tick()
    },
    stop() {
      clearTimeout(this.timerId)
      this.functions = []
    },
  }
}

const coinCanvasOptions = {
  width: 250,
  height: 250,
  radius: 104,
  stepLength: 18,
}

export function drawCoinCanvas(instance, percents = 0) {
  const coinCanvas = instance.refs.coinCanvas as HTMLCanvasElement
  coinCanvas.width = coinCanvasOptions.width
  coinCanvas.height = coinCanvasOptions.height

  const ctx = coinCanvas.getContext('2d')
  const x = coinCanvas.width / 2
  const y = coinCanvas.height / 2
  const stepLength = coinCanvasOptions.stepLength
  const step = 10
  const rotateAngle = Math.PI * 90

  ctx.clearRect(0, 0, coinCanvas.width, coinCanvas.height)

  for (let i = 0; i < 360; i += step) {
    const x1 = x + coinCanvasOptions.radius * Math.cos((i * Math.PI - rotateAngle) / 180)
    const y1 = y + coinCanvasOptions.radius * Math.sin((i * Math.PI - rotateAngle) / 180)
    const x2 =
      x + (coinCanvasOptions.radius + stepLength) * Math.cos((i * Math.PI - rotateAngle) / 180)
    const y2 =
      y + (coinCanvasOptions.radius + stepLength) * Math.sin((i * Math.PI - rotateAngle) / 180)

    // const opacity = ((i + step) / 360).toFixed(1)
    // const color = `rgba(255,255,255, ${opacity === '0.0' ? 0.1 : opacity})`

    ctx.stroke

    if (i <= 360 * (percents / 100)) {
      ctx.strokeStyle = 'gold'
    } else {
      ctx.strokeStyle = 'rgba(255,255,255,0.3)'
    }

    ctx.beginPath()
    ctx.moveTo(x1, y1)
    ctx.lineTo(x2, y2)
    ctx.lineWidth = 3
    ctx.shadowColor = 'rgba(0,0,0,0.3)'
    ctx.shadowOffsetX = 1
    ctx.shadowOffsetY = 1
    // getComputedStyle(document.body).getPropertyValue('--coin-color-1')
    ctx.lineCap = 'round'
    ctx.stroke()
  }
}

// export function drawDonutCanvas(instance, percentage) {
//   if (!percentage) percentage = 0
//   if (percentage > 100) percentage = 100

//   const coinCanvas = instance.refs.donutCanvas as HTMLCanvasElement
//   coinCanvas.width = coinCanvasOptions.width
//   coinCanvas.height = coinCanvasOptions.height

//   const ctx = coinCanvas.getContext('2d')
//   const centerX = coinCanvas.width / 2
//   const centerY = coinCanvas.height / 2
//   const radiusInner = centerX - 11

//   const startAngle = Math.PI * -0.5
//   const endAngle = startAngle + Math.PI * 2 * (percentage / 100)

//   ctx.clearRect(0, 0, coinCanvas.width, coinCanvas.height)
//   ctx.beginPath()
//   ctx.arc(centerX, centerY, radiusInner, startAngle, endAngle)
//   ctx.lineWidth = coinCanvasOptions.stepLength + 3
//   ctx.strokeStyle = getComputedStyle(document.body).getPropertyValue('--coin-color-2')
//   ctx.lineCap = 'round'
//   ctx.stroke()
// }
