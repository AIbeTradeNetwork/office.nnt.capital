<script lang="ts">
  export default {
    name: 'HierarchicalTree',
  }
</script>

<script setup lang="ts">
  import { ref, Ref, onMounted, onBeforeUnmount, computed, nextTick, provide } from 'vue'
  import Branch from './branch.vue'
  import UsersList from './usersList.vue'
  import { $requests } from 'queries/index'
  import { $me } from 'utils/me'
  import { initEvents, destroyEvents, reloadEvents } from './events'

  const tree = ref([])
  const matrix = []
  const maxRows = 3
  const maxCols = 2
  const serverData = ref<TeamUser[]>([])
  const container: Ref<HTMLElement | null> = ref(null)
  const containerInner: Ref<HTMLElement | null> = ref(null)
  const usersList = ref()

  const createNodeData = (data = null) => {
    return { data: data, childrens: [] }
  }

  function createRows() {
    const queue = []
    let row = 1
    let col = 1

    queue.push(tree.value[0])

    while (queue.length && row <= maxRows) {
      const matrixRow = []
      queue.forEach(() => {
        const current = queue.shift()

        for (let i = 1; i <= maxCols; i++) {
          const node = createNodeData()
          matrixRow.push(node)
          current.childrens.push(node)
        }
        queue.push(...current.childrens)
      })

      matrix.push(matrixRow)

      col = 1
      row++
    }
  }

  function getFromForRow(place, row) {
    return (Number(place.col) - 1) * Math.pow(2, Number(row) - Number(place.row))
  }

  function getArrayIndex(placeRoot, place) {
    return {
      row: place.row - placeRoot.row,
      col: place.col - getFromForRow(placeRoot, place.row),
    }
  }

  function fillRows() {
    const rootElPlace = { row: Number(matrix[0][0].place.row), col: Number(matrix[0][0].place.col) }

    serverData.value.forEach((el) => {
      const place = { row: el.place.row, col: el.place.col }
      const elPlace = getArrayIndex(rootElPlace, place)
      matrix[elPlace.row][elPlace.col - 1].data = el
    })
  }

  function createTree(user: Partial<User>) {
    tree.value.push(createNodeData(user))
    matrix.push([user])
    createRows()
    fillRows()
  }

  async function loadData(user: Partial<User>) {
    tree.value = []
    matrix.length = 0

    try {
      serverData.value = await $requests.teams.getBin({
        user_uid: user.uid,
        rows: maxRows,
      })

      createTree(user)

      usersList.value.addUser(user)
    } catch (error) {
      console.error(error)
    }
  }

  provide('hierarchicalTreeLoadData', loadData)

  onMounted(async () => {
    await loadData($me.data)
    nextTick(() => {
      initEvents(container.value, containerInner.value)
    })
    // zoomInit(container.value)
  })

  onBeforeUnmount(() => {
    destroyEvents()
    // zoomDestroy()
    // dragDestroy()
  })
</script>

<template>
  <div ref="container" class="mobile:fixed pc:absolute hierarchical_tree">
    <div ref="containerInner" class="hierarchical_tree-content">
      <div>
        <ul>
          <li>
            <Branch
              v-for="(item, index) in tree"
              :key="item.data?.uid || index"
              :branch-data="item"
            />
          </li>
        </ul>
      </div>
    </div>

    <!-- <button class="absolute bottom-2 right-2 btn btn-circle cursor-pointer" @click="reset">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"
        />
      </svg>
    </button> -->
  </div>

  <UsersList ref="usersList" @selected="loadData($event)" />
</template>

<style lang="sass" scoped>
  @import 'index'
</style>
