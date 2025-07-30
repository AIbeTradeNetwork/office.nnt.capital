import { ref, Ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue'
// import { zoomInit, zoomDestroy, resetZoom } from './zoom'
import { dragInit, dragDestroy } from './drag'
import { $requests } from 'queries/index'
// import { $me } from 'utils/me'
import { $me } from 'utils/me'

const matrixData: Ref<Partial<TeamUser | User>[][]> = ref([[$me.data]])

const container: Ref<HTMLElement | null> = ref(null)

const userBranchDataForRequest = (name) => {
  return {
    user_uid: name,
    rows: 5,
  }
}

async function loadData() {
  try {
    const response = await $requests.teams.getBin(userBranchDataForRequest('SYSTEM'))
    if (!response) return

    response.forEach((item) => {
      const row = BigInt(item.place.row)
      const col = BigInt(item.place.col)
      if (!matrixData.value[row.toString()]) matrixData.value[row.toString()] = []
      matrixData.value[row.toString()][(col - BigInt(1)).toString()] = item
    })

    for (let i = 0; i < matrixData.value[3].length; i++) {
      if (!matrixData.value[3][i]) matrixData.value[3][i] = null
    }

    nextTick(() => {
      dragInit(container.value)
    })
  } catch (error) {
    console.error(error)
  }
}

onMounted(async () => {
  await loadData()
  // zoomInit(container.value)
})

onBeforeUnmount(() => {
  // zoomDestroy()
  dragDestroy()
})
