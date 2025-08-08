import { $i18nGlobal } from 'i18n/index'
import { ELocales } from 'types/enums'
import { computed } from 'vue'

const $documentsList = {
  privacyPolicy: {
    name: 'Privacy policy',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/108waGRyvAtljeP0-L-SmyZlbJws8hDfu/view?usp=drive_link'
        default:
          return 'https://drive.google.com/file/d/1zRZvBU0WdhFrRlmgcXl6Nr8LTDXR2HF5/view?usp=drive_link'
      }
    }),
  },
  servicesUserAgreement: {
    name: 'User service agreement',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/1boo3DTtoL-SoTPtax56Xtmz42x71kN1o/view?usp=drive_link'
        default:
          return 'https://drive.google.com/file/d/1Ujxaig7ayI19fetX-9NifPJgfcF6BHHl/view?usp=drive_link'
      }
    }),
  },
  riskWarning: {
    name: 'Risk warning',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://docs.google.com/https://drive.google.com/file/d/1Q1IibX3xBVp1eluANcinmgf7FZWAOppE/view?usp=drive_link'
        default:
          return 'https://drive.google.com/file/d/1QTPq7XZofuzP657NvJg_zemO10etd88C/view?usp=drive_link'
      }
    }),
  },
  disclaimer: {
    name: 'Disclaimer',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/18o8LRfkQhHQvk7jC0duJs311ODW3nsZP/view?usp=drive_link'
        default:
          return 'https://drive.google.com/file/d/1jBywceNE3pnuGMrIsCQhtLSCTxxZYZrk/view?usp=drive_link'
      }
    }),
  },
  termsOfUseABTFunction: {
    name: 'Terms of use',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/1ZbRJ3VGyelMf2SMPNWtlaU-ENCucveVq/view?usp=drive_link'
        default:
          return 'https://drive.google.com/file/d/1iwROSCnJvgKfKSFeRkHec0dV35hC9rXD/view?usp=drive_link'
      }
    }),
  },
}

export { $documentsList }
