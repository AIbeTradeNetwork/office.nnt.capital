export function getColorSign(val) {
  if (val > 0) {
    return 'text-success'
  } else if (val < 0) {
    return 'text-error'
  }
  return ''
}

export function getSign(val) {
  if (val > 0) {
    return '+'
  } else if (val < 0) {
    return ''
  }
}
