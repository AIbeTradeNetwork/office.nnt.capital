import('./confetti.css')

let container = undefined
const colors = ['orange', 'yellow']

function rnd(min, max) {
  min = parseInt(min)
  max = parseInt(max)
  return Math.floor(Math.random() * (max - min + 1)) + min
}

export function conffetiStart(htmlElement: HTMLElement) {
  if (!htmlElement) return
  container = htmlElement
  const count = (htmlElement.getBoundingClientRect().width / 100) * 5
  const particle = document.createElement('span')
  particle.classList.add('coin-confetti')

  for (let i = 0; i < count; i++) {
    const n = particle.cloneNode() as HTMLSpanElement
    n.style.backgroundColor = colors[rnd(0, colors.length - 1)]
    n.style.top = `${rnd(1, 1)}%`
    n.style.left = `${rnd(0, 100)}%`
    n.style.width = `${rnd(5, 9)}px`
    n.style.height = n.style.width
    n.style.animationDelay = `${rnd(0, 30) / 10}s`
    htmlElement.appendChild(n)
  }
}

export function conffetiStop() {
  if (!container) return
  container.innerHTML = ''
  container = undefined
}
