import copy from 'copy-to-clipboard'
import { $config } from './configuration'
import { $notify } from './notify'
import { $t } from 'i18n/index'

function $copyToClipboard(text: string, callback?: () => void) {
  copy(text || '', {
    debug: $config.isDEV,
    format: 'text/plain',
    onCopy: () => {
      $notify.show({
        type: 'success',
        text: $t('Copied'),
      })

      if (callback) callback()
    },
  })
}

export { $copyToClipboard }
