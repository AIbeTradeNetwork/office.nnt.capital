<script lang="ts">
  export default {
    name: 'TaskMarafon',
  }
</script>

<script setup lang="ts">
  import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
  import { addDays, endOfDay, format, startOfDay } from 'date-fns'
  import { UTCDate } from '@date-fns/utc'
  import { $textSplit } from 'utils/text'
  import { $requests } from 'queries/index'
  import { $me } from 'utils/me'

  interface Step {
    id: number
    start: number
    end: number
    invited: number
    done: boolean
  }

  const steps = ref<Step[]>([])
  const refs = ref<TeamUser[]>([])

  const days = 7
  const maxInvite = 2

  const globalStart = startOfDay(new UTCDate(2024, 5, 25, 0, 0)).getTime()
  const globalEnd = endOfDay(addDays(new UTCDate(globalStart), days - 1)).getTime()
  const endDayTime = ref('')

  const loading = ref(true)
  const isFail = ref(false)
  const isSuccess = ref(false)
  const currentStep = computed(() =>
    steps.value.find((step) => UTCDate.now() >= step.start && UTCDate.now() <= step.end),
  )

  function createArray() {
    // const nowStart = startOfDay(now).getTime()
    // const nowEnd = endOfDay(now).getTime()

    steps.value = []

    for (let i = 0; i < days; i++) {
      const id = i
      const stepStart = addDays(new UTCDate(globalStart), i).getTime()
      const stepEnd = endOfDay(addDays(new UTCDate(globalStart), i)).getTime()
      const invited = (() => {
        // if (i <= 4) return 2
        // if (i === 5) return 1
        return refs.value.filter((ref) => {
          return ref.created_at >= stepStart && ref.created_at <= stepEnd
        }).length
      })()

      const done = (() => {
        const now = UTCDate.now()
        if (isFail.value) return false

        if (stepEnd < now) {
          if (invited >= maxInvite) return true
          isFail.value = true
          return false
        }

        if (now >= stepStart && now <= stepEnd) {
          return undefined
        }

        return null
      })()

      steps.value.push({
        id,
        start: stepStart,
        end: stepEnd,
        invited,
        done,
      })
    }
  }

  function checkDone() {
    const successSteps = steps.value.filter((step) => step.done === true).length >= days
    if (successSteps) isSuccess.value = true
  }

  const timer = {
    time: null as ReturnType<typeof setTimeout>,
    tick() {
      if (currentStep.value) {
        endDayTime.value = format(new UTCDate(currentStep.value.end - UTCDate.now()), 'HH:mm:ss')
      } else {
        return '00:00:00'
      }
    },
    start() {
      this.tick()
      this.time = setTimeout(() => {
        this.start()
      }, 1000)
    },
    stop() {
      if (this.time) clearTimeout(this.time)
    },
  }

  async function init() {
    loading.value = true

    refs.value = await $requests.teams.getRef({
      user_uid: $me.data.uid,
    })

    createArray()
    checkDone()
    timer.start()

    loading.value = false
  }

  async function reload() {
    if (loading.value) return
    init()
  }

  onMounted(async () => {
    init()
  })

  onBeforeUnmount(() => {
    timer.stop()
  })
</script>

<template>
  <div class="mt-6 space-y-4">
    <div class="flex items-center font-bold text-lg text-gradient whitespace-nowrap">
      {{ $t('tasks.task.marafon.title') }}
      &nbsp;
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-6 cursor-pointer hover:text-primary"
        :class="loading ? 'animate-spin' : ''"
        @click="reload"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"
        />
      </svg>
    </div>

    <p class="text-info">{{ $t('tasks.task.marafon.requirement') }}</p>

    <p>
      {{ $t('tasks.task.marafon.time_marathon') }}:
      <b>{{ format(new UTCDate(globalStart), 'P') }}-{{ format(new UTCDate(globalEnd), 'P') }}</b>
    </p>

    <div v-if="loading" class="text-center font-bold my-10 text-info">{{ $t('Loading') }}</div>
    <div v-else class="overflow-x-auto">
      <ul class="steps w-full [&_.step]:min-w-[3rem]">
        <li
          v-for="step in steps"
          :key="step.id"
          class="step"
          :class="`
            ${typeof step.done === 'undefined' && 'step-primary'}
            ${step.done === null && 'step-neutral'}
            ${step.done === false && 'step-error'}
            ${step.done === true && 'step-success'}
          `"
        >
          {{ step.invited }}/{{ maxInvite }}
        </li>
      </ul>
    </div>

    <div>
      <p v-for="(t, index) in $textSplit($t('tasks.task.marafon.reward'))" :key="index">{{ t }}</p>

      <div class="pb-4"></div>

      <p v-if="isFail" class="text-base font-bold text-error">
        {{ $t('tasks.task.marafon.fail') }}
        <span class="text-2xl">&#129402;</span>
      </p>
      <p v-else-if="isSuccess" class="text-base font-bold text-success">
        {{ $t('tasks.task.marafon.success') }}
        <span class="text-2xl">&#129395;</span>
      </p>

      <div v-else>
        <div class="pb-4"></div>

        <p class="text-warning">
          {{ $t('tasks.task.marafon.invited_today') }}: {{ currentStep?.invited || 0 }}/{{
            maxInvite
          }}
        </p>
        <p class="text-warning">
          {{ $t('tasks.task.marafon.remain') }}: {{ endDayTime || '00:00:00' }}
        </p>
      </div>
    </div>
  </div>
</template>
