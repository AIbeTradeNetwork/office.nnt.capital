declare interface IApiError {
  code: number
  error: string
  message: string
}

function $isApiError(x: unknown): x is IApiError {
  if (x && typeof x === 'object' && 'code' in x) {
    return true
  }
  return false
}

export { $isApiError }
