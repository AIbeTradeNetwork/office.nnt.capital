export function $historBack() {
  const history = window.history
  if (history.length <= 1) return
  history.back()
}
