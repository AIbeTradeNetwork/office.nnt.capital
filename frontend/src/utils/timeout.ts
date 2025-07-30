import { Ref, ref } from 'vue'

const defValue = '00:00:00'

function $startTimer(end: number, value: Ref<string>) {
  const start = new Date().getTime()

  if (!start || !end || !value.value || !value['dep']) return

  if (start > end) {
    value.value = defValue
    return
  }

  const time = new Date(end - start).getTime()
  const hours = ('0' + Math.floor((time % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))).slice(-2)
  const minutes = ('0' + Math.floor((time % (1000 * 60 * 60)) / (1000 * 60))).slice(-2)
  const seconds = ('0' + Math.floor((time % (1000 * 60)) / 1000)).slice(-2)

  value.value = `${hours}:${minutes}:${seconds}`

  return setTimeout(() => {
    $startTimer(end, value)
  }, 1000)
}

export { $startTimer }
