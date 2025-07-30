import { $i18nGlobal } from 'i18n/index'
import { ELocales } from 'types/enums'
import { computed } from 'vue'

const $documentsList = {
  privacyPolicy: {
    name: 'Privacy policy',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/1dpIF9jTSnOxdvu326kfpdZcKMY4yA9Pd/view'
        default:
          return 'https://drive.google.com/file/d/1LxknapO7htIvlZdEtLXImiRbEq6WtBiI/view'
      }
    }),
  },
  servicesUserAgreement: {
    name: 'User service agreement',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/1JOqJCJWhz7bQ4z2lWMos-3BE2b8rJKga/view'
        default:
          return 'https://drive.google.com/file/d/1TQVPZcDQkWXJ1Ksb1U49nFdIiuuZ-9vr/view'
      }
    }),
  },
  riskWarning: {
    name: 'Risk warning',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/1W86aU9FOgdr2-jbzsHT29GDmIOeyOTrn/view'
        default:
          return 'https://drive.google.com/file/d/1tzYCdxnOJhZoT-P7GOb0rdUIW3jVIRyN/view'
      }
    }),
  },
  disclaimer: {
    name: 'Disclaimer',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/1gCUDKcmpvV6mFUwcKmvamOK3CDy0HEkG/view'
        default:
          return 'https://drive.google.com/file/d/1KCcIaak73BBK9kWdXNTe4ZaYNFnDrx4r/view'
      }
    }),
  },
  termsOfUseABTFunction: {
    name: 'Terms of use of the ABT function',
    link: computed(() => {
      switch ($i18nGlobal.locale.value) {
        case ELocales.ru_RU:
          return 'https://drive.google.com/file/d/1bbMhEn_kHJv1mZj0953axSaLuLfkSLf4/view'
        default:
          return 'https://drive.google.com/file/d/13S2UhPlQNw4C-4qKN6EdIKblVtFGtdu1/view'
      }
    }),
  },
}

export { $documentsList }
