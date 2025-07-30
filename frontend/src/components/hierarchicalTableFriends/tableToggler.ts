function getAllSibling(el: HTMLElement, deep?: boolean): HTMLElement[] {
  const siblings = []

  if (!el) return siblings
  let sibling = el.nextElementSibling
  while (sibling) {
    if (sibling.nodeType === 1 && sibling !== el) {
      const siblingId = Number(sibling.id)
      const levelId = Number(el.id)
      const nextLevelId = Number(el.id) + 1

      if (levelId === siblingId) break

      if (deep) {
        if (siblingId >= nextLevelId) siblings.push(sibling)
      } else {
        if (siblingId === nextLevelId) {
          siblings.push(sibling)
        }
      }
    }

    sibling = sibling.nextElementSibling
  }

  return siblings || []
}

function toggle(event) {
  if (!event.target) return
  const t = event.target as HTMLElement
  const tr = t.closest('tr')
  const isOpen = tr.classList.contains('open')

  getAllSibling(tr, isOpen).forEach((item) => {
    if (isOpen) {
      item.classList.add('hidden')
      item.classList.remove('open')
    } else {
      item.classList.remove('hidden')
    }
  })

  tr.classList.toggle('open')
}

export { toggle }
