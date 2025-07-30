import { $i18nGlobal } from 'i18n/index'
import { $config } from './configuration'
import { ECurrencies } from 'types/enums'
import { Decimal } from 'decimal.js'

/**
 * Formats a price by locale with optional currency symbol and mathematical sign.
 *
 * @param {Object} options - The options for formatting the price.
 * @param {number | 'string'} options.count - The price to be formatted.
 * @param {ECurrencies} [options.currency] - The currency to be used for formatting. Defaults to the configured currency.
 * @param {boolean} [options.negative] - Whether to display a negative sign.
 * @param {boolean} [options.positive] - Whether to display a positive sign.
 * @return {string} The formatted price with optional currency symbol and mathematical sign.
 */
export const $formatPriceByLocale = function ({
  count,
  currency,
  negative,
  positive,
}: {
  count: number | 'string'
  currency?: ECurrencies | string
  negative?: boolean
  positive?: boolean
}) {
  if (typeof count === 'string' && count.match(/\d/g)) count = JSON.parse(count)
  if (typeof count !== 'number' || !count) count = 0
  if (!currency) currency = ECurrencies['']

  const currencyCode = (() => {
    const c = (currency || $config.currency || ECurrencies.usdHidden).toLowerCase()
    if (c.toLowerCase() === ECurrencies.usdHidden.toLowerCase()) {
      return ECurrencies.usd.toUpperCase()
    }

    return c.toUpperCase()
  })()

  // const getSign = () => {
  //   return String(count).match(/\+|-/) || ['']
  // }

  const options: Intl.NumberFormatOptions = {
    style: 'currency',
    currency: ECurrencies.usdHidden.toUpperCase(),
    // code, symbol, name
    currencyDisplay: 'code',
  }

  const formatted = new Intl.NumberFormat($i18nGlobal.locale.value, options)
    .format(count)
    .replace(ECurrencies.usdHidden.toUpperCase(), currencyCode)
  // .replace(/\+|-/, '')
  // .replace(/(\d{1,})/g, `${getSign()}$1`)

  const withMathSymbol = (() => {
    let out = ''
    if (positive) out = '+'
    if (negative) out = '-'
    return out
  })()

  if (withMathSymbol) {
    return withMathSymbol + ' ' + formatted
  }

  return formatted
}

/**
 * Formats the given amount to Ether based on the provided currency.
 *
 * @param {number} amount - The amount to format.
 * @param {ECurrencies} [currency] - The currency to use for formatting. Defaults to the value of `$config.currency`.
 * @throws {Error} Throws an error if `amount` is falsy or if `currency` is not a valid currency.
 * @return {number} The formatted amount in Ether.
 */
export function $formatAmountToEther(amount: number, currency?: ECurrencies): number {
  if (!amount) {
    throw new Error('formatAmountToEther: empty amount/currency')
  }

  const _currency = currency || $config.currency

  if (!(_currency in ECurrencies)) {
    throw new Error('formatAmountToEther: empty amount/currency')
  }

  if (typeof amount !== 'number') amount = Number(amount)
  if (_currency === ECurrencies.cv) return amount * 10 ** 16
  // TODO: ECurrencies.usd changed to udex and should be removed later
  if (_currency === ('usd' as ECurrencies)) return amount * 10 ** 16
  return amount
}

/**
 * Formats the given number by dividing it by 10 to the power of the specified shift.
 *
 * @param {Int} number - The number to be formatted.
 * @return {number} The formatted number.
 */
export function $formatInt(
  amount: Int,
  o?: { precision?: number; currency_code?: ECurrencies },
): number {
  if (!amount) amount = 0
  const precision =
    o && o.precision ?
      o.precision
    : (() => {
        if (o && o.currency_code === ECurrencies.udex) return 9
        if (o && o.currency_code === ECurrencies.usdHidden) return 2
        return 2
      })()
  return new Decimal(amount).div(10 ** precision).toNumber()
}

/**
 * Formats a given number by multiplying it with a power of 10 and then
 * rounding it to a specified number of decimal places.
 *
 * @param {number} number - The number to be formatted.
 * @return {number} The formatted number.
 */
export function $formatToInt(amount: number, o?: { precision: number }): number {
  if (!amount) amount = 0
  const precision = o ? o.precision : 2
  return new Decimal(amount)
    .toDecimalPlaces(2)
    .mul(10 ** precision)
    .toNumber()
}

/**
 * Formats a number into a string with a 'k' or 'kk' suffix to represent thousands or millions.
 *
 * @param {number} num - The number to be formatted.
 * @return {string|number} The formatted number as a string with a 'k' or 'kk' suffix, or the original number if it is not a number or less than 1000.
 */
export function $formatK(num: number) {
  if (isNaN(num)) return num

  if (num >= 1000) {
    return (num / 1000).toFixed(0) + 'k'
  } else if (num >= 1000000) {
    return (num / 1000000).toFixed(0) + 'kk'
  }

  return num
}
