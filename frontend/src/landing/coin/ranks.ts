import { Decimal } from 'decimal.js'
import { CABTCounts, EABTRanks } from 'types/enums'
import { ComputedRef, computed, ref } from 'vue'

const farmedCoin = ref(new Decimal(0))

const currentRank = computed(() => {
  // return CABTCounts[EABTRanks['archon-I']]
  return Object.values(CABTCounts).find((item) => {
    if (item.min > -1 && item.max) {
      if (
        farmedCoin.value.toNumber() >= new Decimal(item.min).toNumber() &&
        farmedCoin.value.toNumber() < new Decimal(item.max).toNumber()
      ) {
        return item
      }
    } else {
      if (farmedCoin.value.toNumber() >= new Decimal(item.min).toNumber()) {
        return item
      }
    }
  })
})

const nextRank: ComputedRef<EABTRanksCountItem | undefined> = computed(() => {
  const arr = Object.keys(CABTCounts)
  const currentRankIndex = arr.findIndex((item) => item === currentRank.value.code)
  return CABTCounts[arr[currentRankIndex + 1]] || CABTCounts[arr[currentRankIndex]]
})

const getPercent = computed(() => {
  if (!currentRank.value.max) return 100

  return farmedCoin.value
    .minus(currentRank.value.min)
    .mul(100)
    .div(new Decimal(currentRank.value.max).minus(currentRank.value.min))
    .toNumber()
})

function setClaimed(count: string, precision: number) {
  if (!count || isNaN(parseFloat(count))) count = '0'
  if (!precision) precision = 9
  farmedCoin.value = new Decimal(count)
}

const ABTRank = {
  currentRank,
  nextRank,
  setClaimed,
  farmedCoin,
  getPercent,
}

export default ABTRank
