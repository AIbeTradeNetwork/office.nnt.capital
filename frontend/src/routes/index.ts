/** @type {import('vue-router').Config} */

import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { $notify } from 'utils/notify'
import AppLaout from 'layouts/app/index.vue'
// import LandingLaout from 'layouts/landing/index.vue'

import AuthorizationLaout from 'layouts/authorization/index.vue'
import { $authentication } from 'utils/authentication'
import { $config } from 'utils/configuration'
// import { $isNotOurCompany } from 'utils/checks'
// import { EHideClasses } from 'types/enums'
// const redirectToHome = (to: any, from: any, next: any) => {
//     const etherStore = useEtherStore()
//     if (etherStore.web3Provider && etherStore.netId) {
//         next({ name: 'Game' })
//     } else {
//         next()
//     }
// }

async function checkAuth() {
  if ($authentication.isAuthenticated.value) return true
  return false
}

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Index',
    // component: import('pages/dashboard/index.vue'),
    // meta: { layout: AppLaout, keepAlive: true },
    redirect: () => {
      return { name: 'FarmUDEX' }
    },
  },
  // {
  //   path: '/dashboard',
  //   name: 'Dashboard',
  //   component: () => import('pages/dashboard/index.vue'),
  //   meta: { layout: AppLaout, keepAlive: false },
  // },
  {
    path: '/farmUDEX',
    name: 'FarmUDEX',
    component: () => import('coin/index.vue'),
    meta: { layout: AppLaout, keepAlive: false },
    children: [
      {
        path: 'autotrade',
        name: 'Autotrade',
        component: () => import('coin/components/autotrade.vue'),
      },
    ],
  },
  {
    path: '/signin',
    name: 'Authorization',
    component: () => import('pages/authorization/index.vue'),
    meta: { layout: AuthorizationLaout, keepAlive: false },
  },
  {
    path: '/signup',
    name: 'Registration',
    component: () => import('pages/registration/index.vue'),
    meta: { layout: AuthorizationLaout, keepAlive: false },
  },
  // {
  //   path: '/friends',
  //   name: 'Friends',
  //   component: () => import('pages/referral_network/index.vue'),
  //   meta: { layout: AppLaout, keepAlive: true },
  // },
  {
    path: '/friends',
    name: 'Friends',
    component: () => import('pages/referral_network_friends/index.vue'),
    meta: { layout: AppLaout, keepAlive: true },
  },
  {
    path: '/partners',
    name: 'Partners',
    component: () => import('pages/referral_network_partners/index.vue'),
    meta: { layout: AppLaout, keepAlive: true },
  },
  {
    path: '/purchases',
    name: 'Purchases',
    component: () => import('pages/purchases/index.vue'),
    meta: { layout: AppLaout, keepAlive: true },
  },
  {
    path: '/tariffs',
    name: 'Tariffs',
    // component: () => import('pages/tariffs/index.vue'),
    meta: { layout: AppLaout, keepAlive: true },
    redirect: () => {
      return { name: 'Index' }
    },
  },
  {
    path: '/shop',
    name: 'Shop',
    component: () => import('pages/shop/index.vue'),
    meta: { layout: AppLaout, keepAlive: true },
  },
  // {
  //   path: '/tariffs_customer',
  //   name: 'TariffsCustomer',
  //   component: () => import('pages/tariffs_customer/index.vue'),
  //   meta: { layout: AppLaout, keepAlive: true },
  // },
  // {
  //   path: '/tariffs_parntners',
  //   name: 'TariffsParntners',
  //   component: () => import('pages/tariffs_parntners/index.vue'),
  //   meta: { layout: AppLaout, keepAlive: true },
  // },
  {
    path: '/profile',
    name: 'Profile',
    // component: () => import('pages/profile/index.vue'),
    meta: { layout: AppLaout, keepAlive: true },
    redirect: () => {
      return { name: 'Index' }
    },
  },
  {
    path: '/payments',
    name: 'Payments',
    component: () => import('pages/payments/index.vue'),
    meta: { layout: AppLaout, keepAlive: true },
  },
  {
    path: '/autotrade_pro',
    name: 'AutotradePRO',
    meta: { layout: AppLaout, keepAlive: true },
    component: () => import('pages/autotrade_pro/index.vue'),
  },

  {
    path: '/:catchAll(.*)',
    name: 'NotFound',
    meta: { layout: AppLaout, keepAlive: false },
    redirect: () => {
      return { name: 'Index' }
    },
  },

  // {
  //   path: '/landing',
  //   meta: { layout: LandingLaout, keepAlive: true },
  //   children: [
  //     {
  //       path: 'coin',
  //       name: 'CoinLanding',
  //       component: () => import('coin/index.vue'),
  //     },
  //   ],
  // },

  // {
  //   path: '/coin',
  //   name: 'Coin',
  //   meta: { layout: AppLaout, keepAlive: true },
  //   component: () => import('coin/index.vue'),
  // },
]

const Router = createRouter({
  history: createWebHistory(),
  routes,
  linkActiveClass: 'route-active-link',
  linkExactActiveClass: 'route-exact-active-link',
})

Router.beforeEach(async (to, from, next) => {
  if (to?.meta?.layout?.name !== 'LayoutAuth' && !(await checkAuth())) {
    $notify.show({
      type: 'error',
      title: 'Ошибка авторизации',
      text: 'Ваша сессия истекла или вы не авторизованы. Пожалуйста, войдите заново.'
    })
    setTimeout(() => {
      next({ name: 'Authorization' })
    }, 2000) // 2 секунды на прочтение уведомления
    return
  }
  next()
})

function $routesCustomGuard() {
  if (!$config.isDevMode) return

  const customRoutes = [
    {
      path: '/exchanges',
      name: 'Exchanges',
      component: () => import('pages/exchanges/index.vue'),
      meta: { layout: AppLaout, keepAlive: true },
    },
    {
      path: '/strategies',
      name: 'Strategies',
      meta: { layout: AppLaout, keepAlive: true },
      component: () => import('pages/strategies/index.vue'),
    },
    {
      path: '/strategies_ico',
      name: 'StrategiesIco',
      meta: { layout: AppLaout, keepAlive: true },
      component: () => import('pages/strategies_ico/index.vue'),
    },
    {
      path: '/trading_bots',
      name: 'TradingBots',
      meta: { layout: AppLaout, keepAlive: true },
      component: () => import('pages/trading_bots/index.vue'),
    },
  ]

  customRoutes.forEach((r) => Router.addRoute(r))
}

export { Router, routes, $routesCustomGuard }
