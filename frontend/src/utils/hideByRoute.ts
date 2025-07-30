import { Router, routes } from 'routes/index'
import { EHideClasses } from 'types/enums'

export function hideByRoute(type: EHideClasses) {
  const routesMapForAutoTrade = routes
    .filter((item) => {
      if ([
        'FarmUDEX', 'Payments', 'Purchases', 'Friends'
      ].includes(item.name as string))
        return false
      return true
    })
    .map((item) => item.name)

  if (
    type === EHideClasses.autotrade &&
    routesMapForAutoTrade.includes(Router.currentRoute.value.name)
  ) {
    return true
  }

  return false
}
