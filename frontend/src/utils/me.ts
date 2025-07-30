import { reactive } from 'vue'
import { $config } from './configuration'
import { $store } from './store'
import { $requests } from 'queries/index'

const _me: Partial<User> = reactive({})

const $me = {
  data: _me,

  set: (data: Partial<User>) => {
    if (!data) return
    // TODO Change _ref_uid to origin_ref_uid
    if (data.ref_uid === $store.get('cfg').default_ref_uid) {
      data._ref_uid = data.ref_uid + ''
      data.ref_uid = ''
    } else {
      data._ref_uid = data.ref_uid
    }

    Object.assign(_me, data)
  },

  update: async () => {
    await $requests.me.get()
  },

  clear: () => {
    Object.keys(_me).forEach((key) => {
      delete _me[key]
    })
  },
}
export { $me }
