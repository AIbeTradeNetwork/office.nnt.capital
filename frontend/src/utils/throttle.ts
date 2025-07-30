const $throttle = function (callee, time?) {
  let timer = null

  return function (...args) {
    clearTimeout(timer)

    timer = setTimeout(() => {
      callee(...args)
      clearTimeout(timer)
      timer = null
    }, time || 1000)
  }
}

export { $throttle }
