import { $authentication } from 'utils/authentication'

export function createHeaders() {
  const headers = {}
  if ($authentication.getToken()) {
    headers['Authorization'] = 'Bearer ' + $authentication.getToken()
  }
  return headers
}
