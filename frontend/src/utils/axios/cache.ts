import generateID from 'uuid-by-string'

const cacheTime = 1000 * 2
const cache: {
  [id: string]: { data: { [key: string]: any }; timer: number }
} = {}

function axiosGetCache(requets: { [key: string]: any }): undefined | { [key: string]: any } {
  if (!requets) return undefined

  const now = new Date().getTime()
  const id = generateID(JSON.stringify(requets))

  if (!id || !cache[id]) return undefined
  if (cache[id].timer < now) {
    delete cache[id]
    return undefined
  }

  const cacheData = cache[id].data

  return cacheData ? { ...cacheData } : undefined
}

function axiosSetCache(requets: { [key: string]: any }, data: any): undefined | boolean {
  if (!requets || !data) return

  const id = generateID(JSON.stringify(requets))

  if (!id) return

  cache[id] = {
    data: {
      data,
    },
    timer: new Date().getTime() + cacheTime,
  }
}

export { axiosGetCache, axiosSetCache }
