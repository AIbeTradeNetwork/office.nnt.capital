// https://date-fns.org/v2.23.0/docs/format
import { format, setDefaultOptions } from 'date-fns'
import { $getLocaleOptions, $i18nGlobal } from 'i18n/index'
import { watch } from 'vue'

function $formatDateGuard() {
  watch(
    $i18nGlobal.locale,
    () => {
      setDefaultOptions({
        locale: $getLocaleOptions().dateFnsLocale,
      })
    },
    { immediate: true },
  )
}

/**
 * Formats a date according to the specified format.
 *
 * @param {any} date - the date to be formatted
 * @param {string} format - optional format string
 * @template format date, time, date|time, YYYY-MM-DDTHH:mm:ss.SSSZ, Unicode, Ii, R, P, Pp, PPpp, t
 * - These patterns are not in the Unicode Technical Standard
 * - `i`: ISO day of week
 * - `I`: ISO week of year
 * - `R`: ISO week-numbering year
 * - `t`: seconds timestamp
 * - `T`: milliseconds timestamp
 * - `o`: ordinal number modifier
 * - `P`: long localized date
 * - `p`: long localized time
 * @return {string} the formatted date
 */
function $formatDate(date, _format = 'date') {
  if (!date) return ''

  const getFormat = () => {
    let out = ''

    if (_format.match('date')) out += 'P'
    if (_format.match('time')) out += 'p'

    return out || _format
  }
  return format(date, getFormat(), {
    locale: $getLocaleOptions().dateFnsLocale,
  })
}

function $formatDay(period: number, format?: 'day') {
  if (!period) return 0

  return period / (1000 * 1000 * 1000 * 60 * 60 * 24)
}

export { $formatDate, $formatDay, $formatDateGuard }
