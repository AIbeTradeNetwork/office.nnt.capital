import { jwtDecode as decode } from 'jwt-decode'

interface ITokenData {
  exp: number
  iat: number
  iss: string
  uid: string
}

export function $jwtDecode(token: string): ITokenData {
  return decode(token)
}
