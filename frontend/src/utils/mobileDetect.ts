const isMobileBrowser = /Mobi|Android|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
  navigator.userAgent,
)

const $device = {
  isMobile: isMobileBrowser,
  isPc: !isMobileBrowser,
}

export { $device }
