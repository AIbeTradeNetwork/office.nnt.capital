function calculation(trade: Trade) {
  return ((trade.price_close - trade.price_open) / trade.price_open) * 100
}

function styles(trade: Trade) {
  if (!trade.price_close) return ''
  const val = calculation(trade)

  if (val > 0) {
    return 'text-success'
  }

  if (val < 0) {
    return 'text-error'
  }
}

function value(trade: Trade) {
  if (!trade.price_close) return '-'

  const val = calculation(trade)

  if (val > 0) {
    return `+${val.toFixed(2)}`
  }

  if (val < 0) {
    return `${val.toFixed(2)}`
  }

  return '0.00'
}

const profitability = {
  styles,
  value,
}

export { profitability }
