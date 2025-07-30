import { computed } from 'vue'
import { $me } from './me'
import { $store } from './store'
import { $config } from './configuration'
import { ERoles } from 'types/enums'

const cfg = computed(() => {
  return $store.get('cfg')
})

export const $isNotTMeMail = computed(() => !$me.data.email.match('@t.me'))

export const $isNotOurCompany = computed(() => {
  return $store.get('cfg').default_ref_uid !== $me.data._ref_uid
})

export const $isRefInviteEnds = computed(() => $me.data.lim_ref_uid)

export const $isAUnlimitedInvite = computed(
  () => cfg.value.unlim_invite || $me.data.unlim_invite || $config.role === ERoles.distributor,
)
