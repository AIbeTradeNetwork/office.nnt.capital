import { interpolate } from 'culori'

export function $oklchToLch(color) {
  if (!color) return ''
  if (!color.match(/\d/)) return ''

  const q = interpolate([color], 'lch')(0)
  // oklch(0.829011 0.031335 222.959324)
  // rgb(178 204 214)
  // lch(80 11.03 227.26);
  // oklch(0.83 0.03 222.95)
  return `${q.mode}(${q.l} ${q.c} ${q.h})`
}

// export function $lchToHex(color) {
//   if (!color) return ''
//   return formatHex(color)
// }
