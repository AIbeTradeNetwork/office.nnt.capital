declare type Int = number

declare interface INotify {
  type?: 'default' | 'info' | 'success' | 'warning' | 'error'
  title?: string
  text?: string
  autohide?: boolean | number
  accept?: () => void
  decline?: () => void
  see?: () => void
  error?: unknown
}

declare interface INotifyForList extends INotify {
  id: number
  timeout: ReturnType<typeof setTimeout> | undefined
}

declare interface ITransaction {
  type: 'product' | 'tariff' | 'distributor'
  what: string
  code: string
  price?: number
  category?: string
  precision?: number
  currency_code?: string
}

declare interface Paging {
  limit: number
  skip: number
}
