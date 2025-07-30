import { reactive, watch } from 'vue'

let itemID = 0
const time = 3000

const list: INotifyForList[] = reactive([])

function remove(id: number) {
  const index = list.findIndex((item) => item.id === id)
  if (index > -1) {
    clearTimeout(list[index].timeout)
    list.splice(index, 1)
  }
}

function show(item: INotify): void {
  const newId = itemID++

  const newItem: INotifyForList = {
    ...item,
    id: newId,
    timeout: (() => {
      if (typeof item.autohide === 'boolean' && item.autohide === false) return
      if (typeof item.autohide === 'number' && item.autohide <= 0) return
      if (item.see) return
      if (item.accept) return
      if (item.decline) return
      return setTimeout(
        () => remove(newId),
        typeof item.autohide === 'number' ? item.autohide : time,
      )
    })(),
  }

  if (newItem.see) {
    newItem.see = function () {
      item.see()
      remove(newId)
    }
  }

  if (newItem.accept) {
    newItem.accept = function () {
      item.accept()
      remove(newId)
    }
  }

  if (newItem.decline) {
    newItem.decline = function () {
      item.decline()
      remove(newId)
    }
  }

  if (newItem.error) {
    if (typeof newItem.error === 'string') {
      newItem.text = newItem.error
      newItem.type = 'error'
    }

    if (newItem.error instanceof Error) {
      newItem.text = newItem.error.message
      newItem.type = 'error'
    }
  }

  list.push(newItem)
}
const $notify = {
  show,
  list,
}

export { $notify }
