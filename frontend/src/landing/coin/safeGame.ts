import { $t } from 'i18n/index'
import JSConfetti from 'js-confetti'
import { $requests } from 'queries/index'
import { ECurrencies, Precision } from 'types/enums'
import { $modals } from 'utils/modals'
import { ComponentInternalInstance, computed, ref, watch } from 'vue'

export function useSafeGame() {
  let safeZone = undefined as HTMLElement | undefined
  let coinCanvas = undefined as HTMLElement | undefined
  let safeZoneKey = undefined as HTMLElement | undefined
  let safeZoneInfo = undefined as HTMLElement | undefined
  let coinImg = undefined as HTMLImageElement | undefined
  let safeZoneWarning = undefined as HTMLElement | undefined
  let currentNumberZone = undefined
  let numbersZone = undefined

  const numbersLength = 12
  const numbersAngle = 360 / numbersLength
  const codeLength = 8

  function CreateNumbers() {
    safeZone.classList.remove('hidden')

    numbersZone = document.createElement('div')
    numbersZone.classList.add(
      'absolute',
      'w-full',
      'h-full',
      'flex',
      'outline-none',
      'select-none',
      'cursor-pointer',
    )
    safeZone.appendChild(numbersZone)

    const boxD = numbersZone.getBoundingClientRect().width / 2 - 11
    for (let i = 0; i < numbersLength; i++) {
      const number = document.createElement('div')
      number.classList.add(
        'absolute',
        'flex',
        'items-center',
        'justify-center',
        'text-white',
        'w-6',
        'h-6',
        'rounded-full',
        'bg-black',
        'pointer-events-none',
        'outline-none',
        'select-none',
      )

      number.innerText = String(i + 1)
      number.style.top = -boxD * Math.cos((360 / numbersLength / 180) * i * Math.PI) + boxD + 'px'
      number.style.left = boxD * Math.sin((360 / numbersLength / 180) * i * Math.PI) + boxD + 'px'
      number.style.transform = `rotate(${(360 / numbersLength) * i}deg)`
      numbersZone.appendChild(number)
    }
  }

  function CreateCurrentNumber() {
    currentNumberZone = document.createElement('div')
    currentNumberZone.classList.add(
      'w-6',
      'h-6',
      'absolute',
      'top-0',
      'left-2/4',
      '-translate-x-2/4',
      'flex',
      'items-center',
      'justify-center',
      'border',
      'border-1',
      'rounded-full',
      'pointer-events-none',
      'outline-none',
      'select-none',
      'bg-white',
      'text-black',
    )

    const text = document.createElement('div')
    text.innerText = '1'
    text.classList.add('text-black')
    currentNumberZone.appendChild(text)

    const radialPercent = document.createElement('div')
    radialPercent.classList.add(
      'absolute',
      '-top-[2px]',
      '-left-[2px]',
      'w-[26px]',
      'h-[26px]',
      'radial-progress',
      'text-green-500',
      'overflow-hidden',
    )
    radialPercent.role = 'progressbar'
    // radialPercent.style.setProperty('--size', '20rem')
    radialPercent.style.setProperty('--thickness', '4px')
    radialPercent.style.setProperty('--value', '0')

    CreateCurrentNumber.prototype.setPercent = (percent: string) => {
      radialPercent.style.setProperty('--value', percent)
    }

    currentNumberZone.appendChild(radialPercent)

    safeZone.appendChild(currentNumberZone)
  }

  const ConfetiTimer = {
    JSConfetti: new JSConfetti(),
    time: null as ReturnType<typeof setTimeout>,
    async tick() {
      return await this.JSConfetti.addConfetti()
    },
    async start() {
      await this.tick()
      this.time = setTimeout(() => {
        this.start()
      }, 100)
    },
    stop() {
      this.JSConfetti.clearCanvas()
      if (this.time) clearTimeout(this.time)
    },
  }

  const Combo = {
    el: undefined as HTMLElement | undefined,
    hintEl: undefined as HTMLElement | undefined,
    current: null,
    data: ref([]),
    dataHint: ref([]),
    keysForBuy: [6, 7],
    keyCurrentPosition: computed(() => Combo.data.value.length),
    hintOpenKeysLength: computed(
      () => Combo.dataHint.value.filter((item) => !isNaN(Number(item))).length + 1,
    ),
    isKeyForBuy: computed(() => {
      if (Combo.dataHint.value.length === 0) return false
      // const isOpen = !!Number(Combo.dataHint.value[Combo.keyCurrentPosition.value - 1])
      // const isKeyInBuy = Combo.keysForBuy.includes(Combo.keyCurrentPosition.value)
      // if (!isOpen && isKeyInBuy) return true
      if (Combo.keysForBuy.includes(Combo.hintOpenKeysLength.value)) return true
      return false
    }),

    addHint() {
      const div = document.createElement('div')

      div.innerHTML = ''

      div.classList.add(
        'hint',
        'hidden',
        'absolute',
        'top-0',
        'left-2/4',
        '-translate-x-2/4',
        '-mt-[10px]',
        'bg-base-300',
        'flex',
        'justify-center',
        'items-center',
        'rounded-full',
        'px-2',
        'h-[20px]',
        'border-[1px]',
        'border-white',
        'text-[11px]',
        'tracking-wider',
      )

      Combo.hintEl = div
      Combo.el.appendChild(div)
      Combo.loadHint()
    },

    removeHint() {
      Combo.hintEl.parentNode?.removeChild(Combo.hintEl)
    },

    fillHint(text: string) {
      const stringToArr = (text || '').split(',').filter((item) => item)
      Combo.dataHint.value = stringToArr
      Combo.hintEl.innerText = stringToArr.join(' ')
      if (!stringToArr.length) InviteWarning.open()
    },

    async loadHint() {
      let response
      try {
        response = await $requests.games.user_safe()
        if (response && response.secret) {
          Combo.hintEl.classList.remove('hidden')
        }
      } catch (error) {
        console.error(error)
      }
      Combo.fillHint(response && response.secret ? response.secret : '')
      // Combo.fillHint('1,2,*,4,5,*,7,*')
    },

    addReset() {
      const btn = document.createElement('div')
      btn.innerHTML = `<svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M19.5 12c0-1.232-.046-2.453-.138-3.662a4.006 4.006 0 0 0-3.7-3.7 48.678 48.678 0 0 0-7.324 0 4.006 4.006 0 0 0-3.7 3.7c-.017.22-.032.441-.046.662M19.5 12l3-3m-3 3-3-3m-12 3c0 1.232.046 2.453.138 3.662a4.006 4.006 0 0 0 3.7 3.7 48.656 48.656 0 0 0 7.324 0 4.006 4.006 0 0 0 3.7-3.7c.017-.22.032-.441.046-.662M4.5 12l3 3m-3-3-3 3"
        />
      </svg>`

      btn.classList.add(
        'reset',
        'absolute',
        'top-2/4',
        '-right-[12px]',
        '-mt-[17px]',
        'bg-base-300',
        'pointer-events-auto',
        'cursor-pointer',
        'flex',
        'justify-center',
        'items-center',
        'rounded-full',
        'px-1',
        'py-1',
        'border-[1px]',
        'border-white',
      )

      btn.onclick = () => this.clear()

      this.el.appendChild(btn)
    },

    removeReset() {
      this.el.removeChild(this.el.querySelector('.reset'))
    },

    fill() {
      this.el.querySelectorAll('div:not([class])').forEach((el) => el.remove())
      this.data.value.forEach((n: string) => {
        const div = document.createElement('div')
        div.innerText = n
        this.el.appendChild(div)
      })
    },

    check() {
      if (InviteWarning.isOpen) return

      this.data.value.push(this.current)

      this.fill()

      if (this.data.value.length < codeLength) return

      $requests.games
        .hack_safe(Combo.data.value.join(','))
        .then((response: UserSafe) => {
          if (response && response.uid) {
            Combo.clear()
            ConfetiTimer.start()
            Combo.loadHint()

            $modals.prize.onClose = () => {
              ConfetiTimer.stop()
            }

            $modals.prize.show({
              amount: response.variant.amount,
              currency: response.variant.type,
              precision: Precision[response.variant.type],
            })
          }
        })
        .finally(() => {
          if (this.data.value.length === codeLength) {
            setTimeout(() => {
              this.clear()
            }, 300)
          }
        })
    },

    setCurrent(n: number) {
      Combo.current = n
      currentNumberZone.querySelector('div').innerText = String(n)
    },

    show() {
      this.el.classList.remove('hidden')
      this.el.querySelector('.reset')?.classList.remove('hidden')
    },
    hide() {
      this.el.classList.add('hidden')
      this.el.querySelector('.reset')?.classList.add('hidden')
    },
    clear() {
      this.data.value = []
      this.fill()
    },
  }

  const safeBuyKeys = {
    unwatch: watch(Combo.isKeyForBuy, (value) => {
      if (typeof value === 'boolean' && value) {
        safeZoneKey.classList.remove('hidden')
      } else {
        safeZoneKey.classList.add('hidden')
      }
    }),

    openBuy() {
      if (!Combo.keysForBuy.includes(Combo.hintOpenKeysLength.value)) return

      $modals.farmingSubscription.show(`safe_${Combo.hintOpenKeysLength.value}th_digit`, {
        onSuccess: () => {
          Combo.loadHint()
        },
      })

      // $modals.paySystem.show({
      //   type: 'product',
      //   what: $t('safe_game.buy_key_text', { key: Combo.hintOpenKeysLength.value }),
      //   code: `safe_${Combo.hintOpenKeysLength.value}th_digit`,
      //   currency_code: ECurrencies.usdHidden,
      //   onSuccess: () => {
      //     Combo.loadHint()
      //   },
      // })
    },
  }

  const InviteWarning = {
    isOpen: false,
    open() {
      this.isOpen = true
      safeZoneWarning.classList.remove('hidden')
    },
    close() {
      this.isOpen = false
      safeZoneWarning.classList.add('hidden')
    },
  }

  const SelectNumber = {
    timerMove: <ReturnType<typeof setTimeout>>undefined,
    timerSelect: <ReturnType<typeof setTimeout>>undefined,
    timerPercent: <ReturnType<typeof setTimeout>>undefined,
    timeWait: 1000,
    percent: 0,
    percentPart: 100 / 4,

    setPercent() {
      CreateCurrentNumber.prototype.setPercent(this.percent)
      this.percent = this.percent + this.percentPart
      this.timerPercent = setTimeout(() => {
        if (this.percent <= 100) {
          SelectNumber.setPercent()
        } else {
          this.percent = 0
          CreateCurrentNumber.prototype.setPercent(this.percent)
        }
      }, this.timeWait / 4)
    },

    clearPercent() {
      this.percent = 0
      CreateCurrentNumber.prototype.setPercent(this.percent)
    },

    startTimer() {
      this.timerMove = setTimeout(() => {
        this.setPercent()
        this.timerSelect = setTimeout(() => {
          Combo.check()
        }, this.timeWait)
      }, 300)
    },

    clearTimer() {
      // SetBtn.hide()
      this.clearPercent()
      clearTimeout(this.timerPercent)
      clearTimeout(this.timerSelect)
      clearTimeout(this.timerMove)
    },
  }

  function showInfo() {
    $modals.info.show({
      title: $t('info.safe_game.title'),
      text: $t('info.safe_game.text'),
    })
  }

  const Events = {
    overflowClass: ['overscroll-contain'],
    work: false,
    startClick: 0,
    startRotated: 0,
    rotated: 0,
    scrollBlockStart() {
      document.documentElement.classList.add(...Events.overflowClass)
      document.body.classList.add(...Events.overflowClass)
    },
    scrollBlocStop() {
      document.documentElement.classList.remove(...Events.overflowClass)
      document.body.classList.remove(...Events.overflowClass)
    },

    down(event: MouseEvent | TouchEvent) {
      event.preventDefault()
      Events.scrollBlockStart()
      Events.work = true
      Events.startClick = Events.getAngle(event)
      Events.startRotated = Events.rotated
    },

    up(event: MouseEvent | TouchEvent) {
      Events.work = false
      Events.scrollBlocStop()
    },

    getAngle(event: MouseEvent | TouchEvent) {
      const ev = event instanceof TouchEvent ? event.touches[0] : event
      const rect = safeZone.getBoundingClientRect()
      const dX = ev.clientX - rect.left - rect.width / 2
      const dY = ev.clientY - rect.top - rect.height / 2
      const angle = Math.trunc((Math.atan2(dY, dX) * 180) / Math.PI)
      return angle
    },

    move(event: MouseEvent | TouchEvent) {
      event.preventDefault()

      if (!Events.work) return

      const delta = Events.startClick - Events.startRotated

      Events.rotated = Events.getAngle(event) - delta

      if (Events.rotated < 0) Events.rotated += 360
      if (Events.rotated > 360) Events.rotated -= 360

      numbersZone.style.transform = `rotate(${Events.rotated}deg)`
      coinCanvas.style.transform = `rotate(${Events.rotated}deg)`
      coinImg.style.transform = `rotate(${Events.rotated}deg)`

      const reverse = 360 - Events.rotated
      const number = Math.trunc(Math.round(reverse / numbersAngle)) + 1
      const deltaNumber = Math.abs(number > numbersLength ? 1 : number)

      Combo.setCurrent(deltaNumber)

      SelectNumber.clearTimer()
      SelectNumber.startTimer()
    },

    add() {
      safeZone.addEventListener('mouseup', this.up, true)
      safeZone.addEventListener('mousedown', this.down, true)
      safeZone.addEventListener('mouseout', this.up, true)
      safeZone.addEventListener('mousemove', this.move, true)

      safeZone.addEventListener('touchend', this.up, true)
      safeZone.addEventListener('touchstart', this.down, true)
      safeZone.addEventListener('touchmove', this.move, true)
      safeZoneKey.addEventListener('click', safeBuyKeys.openBuy)
      safeZoneInfo.addEventListener('click', showInfo)
    },

    remove() {
      safeZone.removeEventListener('mouseup', this.up, true)
      safeZone.removeEventListener('mousedown', this.down, true)
      safeZone.removeEventListener('mouseout', this.up, true)
      safeZone.removeEventListener('mousemove', this.move, true)

      safeZone.removeEventListener('touchend', this.up, true)
      safeZone.removeEventListener('touchstart', this.down, true)
      safeZone.removeEventListener('touchmove', this.move, true)
      safeZoneKey.removeEventListener('click', safeBuyKeys.openBuy)
      safeZoneInfo.removeEventListener('click', showInfo)
    },
  }

  function startGame(instance: ComponentInternalInstance) {
    safeZoneInfo = instance.refs.safeZoneInfo as HTMLElement
    safeZoneKey = instance.refs.safeZoneKey as HTMLElement
    safeZone = instance.refs.safeZone as HTMLElement
    coinCanvas = instance.refs.coinCanvas as HTMLElement
    coinImg = instance.refs.coinImg as HTMLImageElement
    safeZoneWarning = instance.refs.safeZoneWarning as HTMLElement
    Combo.el = instance.refs.arrowsCombo as HTMLElement
    safeZoneInfo.classList.remove('hidden')
    CreateNumbers()
    CreateCurrentNumber()
    Combo.clear()
    Combo.addReset()
    Combo.show()
    Combo.addHint()

    Events.add()
  }

  function stopGame() {
    numbersZone.style.transform = null
    coinCanvas.style.transform = null
    coinImg.style.transform = null
    Events.remove()
    SelectNumber.clearTimer()
    Combo.clear()
    Combo.hide()
    Combo.removeHint()
    Combo.removeReset()
    ConfetiTimer.stop()
    safeBuyKeys.unwatch()
    InviteWarning.close()
    safeZone.classList.add('hidden')
    safeZoneKey.classList.add('hidden')
    safeZoneInfo.classList.add('hidden')
    safeZone.innerHTML = ''
    safeZone = undefined
    coinCanvas = undefined
    coinImg = undefined
    currentNumberZone = undefined
    numbersZone = undefined
    safeZoneKey = undefined
    safeZoneInfo = undefined
    safeZoneWarning = undefined
  }

  return { startGame, stopGame }
}
