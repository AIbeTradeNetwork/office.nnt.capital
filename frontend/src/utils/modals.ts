import { ECoinFarmingActions } from 'types/enums'
import { reactive } from 'vue'

const modalsData = {
  active: <boolean>false,
  onShow: (...args: any[]) => {},
  onClose: (...args: any[]) => {},
}

// const coinShop = reactive({
//   ...modalsData,
//   show() {
//     this.onShow()
//     this.active = true
//   },
//   close() {
//     this.onClose()
//     this.active = false
//   },
// })

const any = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const brokerPayout = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const brokerAddingFunds = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const balance = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const addingFunds = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const farmingSubscription = reactive({
  ...modalsData,
  type: '',
  data: null,
  onSuccess: <(...args: any[]) => void | null>null,
  show(
    type: typeof farmingSubscription.type,
    o?: {
      data?: typeof farmingSubscription.data
      onSuccess?: typeof farmingSubscription.onSuccess
    },
  ) {
    this.type = type
    this.data = o ? o.data : null
    this.onShow()
    this.onSuccess = o ? o.onSuccess : null
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.type = ''
    this.data = null
    this.onSuccess = null
  },
})

const farmingInfo = reactive({
  type: <ECoinFarmingActions | undefined>undefined,
  ...modalsData,
  show(type: ECoinFarmingActions) {
    this.onShow()
    this.type = type
    this.active = true
  },
  close() {
    this.onClose()
    this.type = undefined
    this.active = false
  },
})

const info = reactive({
  title: '',
  text: '',
  ...modalsData,
  show({ title = '', text = '' }) {
    this.onShow()
    this.title = title
    this.text = text
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.title = ''
    this.text = ''
  },
})

const prize = reactive({
  title: '',
  amount: 0,
  currency: '',
  precision: 0,
  ...modalsData,
  show({ title = '', amount = 0, currency = '', precision = 0 }) {
    this.onShow()
    this.title = title
    this.amount = amount
    this.currency = currency
    this.precision = precision
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.amount = 0
    this.currency = ''
    this.precision = 0
    this.title = ''
  },
})

const warning = reactive({
  title: '',
  text: '',
  ...modalsData,
  show({ title = '', text = '' }) {
    this.onShow()
    this.title = title
    this.text = text
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.title = ''
    this.text = ''
  },
})

const error = reactive({
  title: '',
  text: '',
  ...modalsData,
  show({ title = '', text = '' }) {
    this.onShow()
    this.title = title
    this.text = text
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.title = ''
    this.text = ''
  },
})

const PDF = reactive({
  title: '',
  src: '',
  ...modalsData,
  show({ title = '', src = '' }) {
    this.onShow()
    this.title = title
    this.src = src
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.title = ''
    this.src = ''
  },
})

const settings = reactive({
  title: '',
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.title = ''
  },
})

const exchange = reactive({
  data: <Exchange>null,
  ...modalsData,
  show(data: Exchange) {
    this.onShow()
    this.data = data
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.data = null
  },
})

const exchanges = reactive({
  data: <Exchange>null,
  ...modalsData,
  show(data: Exchange) {
    this.onShow()
    this.data = data
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.data = null
  },
})

const bot = reactive({
  data: <Strategy | Bot>null,
  ...modalsData,
  show(data: Strategy | Bot) {
    this.onShow()
    this.data = data
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.data = null
  },
})

const strategies = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const documents = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const userSocial = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const paySystem = reactive({
  data: <ITransaction>null,
  ...modalsData,
  show(data: ITransaction & { onSuccess?: (...args: any[]) => void }) {
    this.onShow()
    this.data = data
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.data = null
  },
})

const keys = reactive({
  data: <Exchange>null,
  ...modalsData,
  show(data: Exchange) {
    this.onShow()
    this.data = data
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.data = null
  },
})

const question = reactive({
  title: '',
  text: '',
  ...modalsData,
  onConfirm: () => {},
  onDeny: () => {},
  show({ title = '', text = '' }) {
    this.onShow()
    this.title = title
    this.text = text
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
    this.title = ''
    this.text = ''
  },
})

const payout = reactive({
  ...modalsData,
  onSucess: () => {},
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const promocode = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const ranks = reactive({
  ...modalsData,
  data: <{ [key: string]: number }>{},
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    ranks.data = {}
    this.onClose()
    this.active = false
  },
})

const goToPro = reactive({
  ...modalsData,
  show() {
    this.onShow()
    this.active = true
  },
  close() {
    this.onClose()
    this.active = false
  },
})

const $modals = {
  any,
  info,
  warning,
  error,
  PDF,
  settings,
  exchange,
  exchanges,
  bot,
  strategies,
  documents,
  userSocial,
  paySystem,
  keys,
  question,
  payout,
  promocode,
  farmingInfo,
  ranks,
  farmingSubscription,
  addingFunds,
  balance,
  brokerPayout,
  brokerAddingFunds,
  prize,
  goToPro,
  // coinShop,
}

export { $modals }
