import './base.css'
import './App.css'
import './theme.css'

import { createApp } from 'vue'

import { $queriesGuards } from 'utils/queriesGuards'
import { $theme } from 'utils/theme'
import { $screenGuard } from 'utils/screen'
import { $authenticationGuard } from 'utils/authentication'
import { i18n, $i18nGuard } from 'i18n/index'
import { $formatDateGuard } from 'utils/date'
import { $routesCustomGuard, Router } from 'routes/index'
import App from './App.vue'
import { $initTelegramAuthFrameAppGuard } from 'utils/auth/telegram'
import { $checkRefreshGuard } from 'utils/refresh'
;(async () => {
  try {
    $checkRefreshGuard()

    await $initTelegramAuthFrameAppGuard()

    $theme.guard()
    $screenGuard()
    $authenticationGuard()
    await $i18nGuard()
    $formatDateGuard()

    createApp(App).use(i18n).use(Router).mount('#app')

    $queriesGuards()

    $routesCustomGuard()
  } catch (error) {
    console.error(error)
  }
  // app.component('App', AppLaout).component('Authorization', AuthorizationLaout)
})()
