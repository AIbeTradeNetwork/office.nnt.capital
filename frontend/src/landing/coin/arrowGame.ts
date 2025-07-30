import { ComponentInternalInstance, ref } from 'vue'
import JSConfetti from 'js-confetti'
import { $requests } from 'queries/index'
import { $modals } from 'utils/modals'
import { $t } from 'i18n/index'
import { $formatInt, $formatPriceByLocale } from 'utils/formats'

enum ArrowsMap {
  'up' = 1,
  'down' = 3,
  'left' = 2,
  'right' = 4
}

type ArrowDirections = keyof typeof ArrowsMap

export function useArrowGame() {
  const jsConfetti = new JSConfetti()

  const testCombinations = ['up,down']

  const options = {
    comboLength: 8
  }

  let coinBlock = undefined as HTMLElement | undefined
  let arrowsZoneEl = undefined as HTMLElement | undefined
  let arrowsComboEl = undefined as HTMLElement | undefined

  const arrowsCombo = {
    combo: ref(''),
    createArrow(dir: ArrowDirections) {
      const cssRotate = (d: number) => `transform: rotate(${d}deg);`
      const rotate = dir === 'up' ? cssRotate(0) : dir === 'right' ? cssRotate(90) : dir === 'down' ? cssRotate(180) : dir === 'left' ? cssRotate(-90) : ''
      const arrow = `
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-[19px]" style="${rotate}">>
          <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 15.75 7.5-7.5 7.5 7.5" />
        </svg>
      `

      var div = document.createElement('div')
      div.innerHTML = arrow.trim()
      return div.firstChild
    },
    encodeCombo() {
      arrowsComboEl.innerHTML = ''

      return (this.combo.value.split('')).forEach((a: string) => {
        arrowsComboEl.appendChild(this.createArrow(ArrowsMap[a] as ArrowDirections))
      })
    },
    show () {
      arrowsComboEl.classList.remove('hidden')
    },
    hide () {
      arrowsComboEl.classList.add('hidden')
    },
    clear() {
      arrowsCombo.combo.value = ''
      arrowsCombo.hide()
    }
  }

  const confetiTimer = {
    time: null as ReturnType<typeof setTimeout>,
    async tick() {
      return await jsConfetti.addConfetti()
    },
    async start() {
      await this.tick()
      this.time = setTimeout(() => {
        this.start()
      }, 100)
    },
    stop() {
      jsConfetti.clearCanvas()
      if (this.time) clearTimeout(this.time)
    },
  }

  async function checkCombo() {
    if (!arrowsComboEl) return
    if (!arrowsZoneEl) return
    
    const combo = arrowsCombo.combo.value.split('').map((a: string) => ArrowsMap[a] as ArrowDirections).join(',')

    $requests.games.set_combo(combo).then((response: Combo) => {
      if (response && response.uid) {
        arrowsCombo.clear()
        confetiTimer.start()
    
        $modals.prize.onClose = () => {
          confetiTimer.stop()
        }
    
        $modals.prize.show({
          amount: response.amount,
          currency: response.currency_code,
          precision: response.precision
        })
      }

      
    }).finally(() => {
      if (arrowsCombo.combo.value.length >= options.comboLength) {
        setTimeout(() => {arrowsCombo.clear()}, 300)
      }
    })
  }

  function doStep() {
    if (!arrowsComboEl) return
    if (!arrowsZoneEl) return

    if (arrowsCombo.combo.value.length) {
      arrowsCombo.show()
      arrowsCombo.encodeCombo()
      checkCombo()
      return
    }
  }

  const arrowClick = {
    zoneBoxCLick(event: Event) {
      if (!coinBlock) return
      const e = event instanceof TouchEvent ? event.touches[0] : <MouseEvent>event
      const rect = coinBlock.getBoundingClientRect()
      const d = 40
      const dX = e.clientX - rect.left - rect.width / 2
      const dY = e.clientY - rect.top - rect.height / 2
      const tY = ((dX / rect.width) * -d)
      const tX = ((dY / rect.height) * -d)
      coinBlock.style.setProperty('--x', `${tX}deg`)
      coinBlock.style.setProperty('--y', `${tY}deg`)
  
      setTimeout(() => {
        coinBlock.style.setProperty('--x', `0deg`)
        coinBlock.style.setProperty('--y', `0deg`)
      }, 100)
    },
  
    event (event: Event) {
      const dir = (event.target as HTMLElement).getAttribute('data-dir') as ArrowDirections
      arrowsCombo.combo.value = arrowsCombo.combo.value + ArrowsMap[dir]
      doStep()
    },
    add() {
      coinBlock.addEventListener('click', arrowClick.zoneBoxCLick, true)
      arrowsZoneEl.querySelectorAll('button').forEach((b: HTMLElement) => {
        b.classList.remove('hidden')
        b.addEventListener('click', arrowClick.event)
      })
    },
    remove() {
      coinBlock.removeEventListener('click', arrowClick.zoneBoxCLick, true)
      arrowsZoneEl.querySelectorAll('button').forEach((b: HTMLElement) => {
        b.removeEventListener('click', arrowClick.event)
      })
    }
  }

  function startGame(instance: ComponentInternalInstance) {
    coinBlock = instance.refs.coinBlock as HTMLElement
    arrowsZoneEl = instance.refs.arrowsZone as HTMLElement
    arrowsComboEl = instance.refs.arrowsCombo as HTMLElement

    arrowsZoneEl.classList.remove('hidden')

    arrowsCombo.clear()
    arrowClick.add()
  }

  function stopGame() {
    jsConfetti.clearCanvas()
    arrowsCombo.clear()
    arrowClick.remove()
    arrowsZoneEl.classList.add('hidden')
    arrowsComboEl.classList.add('hidden')
    arrowsComboEl.innerHTML = ''
    coinBlock = undefined
    arrowsZoneEl = undefined
    arrowsComboEl = undefined

    if (arrowsCombo.combo.value.length >= options.comboLength) {
      arrowsCombo.clear()
    }
  }

  return { startGame, stopGame }
}
