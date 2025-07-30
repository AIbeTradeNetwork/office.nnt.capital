<script lang="ts">
  export default {
    name: 'ModalBot',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $t } from 'i18n/index'
  import { $requests } from 'queries/index'
  import { EBotTradeType, ECurrencies, EStrategyCategory } from 'types/enums'
  import { $riskLevelBadge } from 'utils/colors'
  import { $formatInt } from 'utils/formats'
  import { $modals } from 'utils/modals'
  import { $notify } from 'utils/notify'
  import { $store } from 'utils/store'
  import { computed, onBeforeMount, onMounted, reactive, ref } from 'vue'

  interface Form {
    category: StrategyCategory
    uid: string | null
    strategy_code: BotAddReq['strategy_code']
    key_uid: BotAddReq['key_uid'] | BotEditReq['key_uid']
    is_active: BotEditReq['is_active']
    trade_percent: Int
    trade_type: BotAddReq['trade_type'] | undefined
    trade_limit: BotAddReq['trade_limit']
    trade_reinvest: BotAddReq['trade_reinvest']
    leverage?: BotAddReq['leverage']
  }

  const botData = $modals.bot.data as Strategy & Bot
  const loading = ref(false)
  const checked = ref(false)
  const tradePercent = {
    min: 1,
    max: 100,
    step: 1,
  }

  // const exchanges = computed(() => $store.get('exchanges'))
  // const selectedExchange: Ref<Exchange['code']> = ref(botData.exchange_code)
  // const selectedExchange: Ref<Exchange['code']> = ref(botData.exchange_code || '')

  const modalType = ref<{
    title: string
    type: 'add' | 'edit'
  }>()

  const strategyDataByCode = computed(() =>
    $store.get('strategies').find((item) => item.code === (botData.strategy_code || botData.code)),
  )

  const ifShowLeverage = computed(() => {
    if (strategyDataByCode.value.fix_leverage > 0) return false
    if (!strategyDataByCode.value.min_leverage && !strategyDataByCode.value.max_leverage)
      return false
    if (strategyDataByCode.value.min_leverage > 0 || strategyDataByCode.value.max_leverage > 0)
      return true
    return false
  })

  const isTradeTypeLimitAndPercent = computed(() => {
    return strategyDataByCode.value.is_trade_percent && strategyDataByCode.value.is_trade_limit
  })

  const initialForm: Form = {
    category: strategyDataByCode.value.category,
    uid: botData.uid || null,
    strategy_code: botData.strategy_code || botData.code,
    key_uid: botData.key_uid || '',
    is_active: typeof botData.is_active === 'boolean' ? botData.is_active : true,
    trade_percent: botData.trade_percent || 1,
    trade_type: undefined,
    trade_limit: $formatInt(botData.trade_limit) || 0,
    trade_reinvest: botData.trade_reinvest || false,
    leverage:
      botData.leverage ||
      strategyDataByCode.value.fix_leverage ||
      strategyDataByCode.value.min_leverage ||
      1,
  }

  const form: Form = reactive({
    ...initialForm,
  })

  // const tradeTypeModel = computed({
  //   get: () => (form.trade_type === 'percent' ? false : true),
  //   set: (val) => {
  //     if (val) form.trade_type = EBotTradeType.limit
  //     else form.trade_type = EBotTradeType.percent
  //   },
  // })

  const keysByExchange = computed(() => {
    // .filter((item) => item.exchange_code === selectedExchange.value)
    return $store
      .get('keys')
      .filter((item) => strategyDataByCode.value.exchange_codes.includes(item.exchange_code))
    // item.exchange_code === selectedExchange.value
  })

  function setModalType() {
    if (!form.key_uid) {
      modalType.value = {
        title: $t('BotAdding'),
        type: 'add',
      }
      return
    }

    modalType.value = {
      title: $t('BotEditing'),
      type: 'edit',
    }
  }

  function setBotType() {
    if (botData.trade_type) {
      form.trade_type = botData.trade_type
      return
      // botData.trade_type || (isTradeTypeLimitAndPercent.value ? EBotTradeType.percent : undefined),
    }
    if (strategyDataByCode.value.is_trade_limit) {
      form.trade_type = EBotTradeType.limit
      form.trade_reinvest = false
      return
    }
    if (strategyDataByCode.value.is_trade_percent) {
      form.trade_type = EBotTradeType.percent
      return
    }
  }

  function checkFormFill() {
    if (!form.strategy_code) return
    if (!form.key_uid) return
    return true
  }

  function openModalDocuments() {
    $modals.documents.show()
  }

  async function doProcessQuestion(text: string) {
    return new Promise((resolve) => {
      $modals.question.onDeny = () => {
        return resolve(false)
      }

      $modals.question.onConfirm = () => {
        return resolve(true)
      }

      $modals.question.show({
        text,
      })
    })
  }

  async function botAddSent() {
    loading.value = true

    try {
      const response = await $requests.bots.add({
        ...form,
        ...{
          symbol_code: botData.symbol_code,
        },
      })

      if (response?.uid) {
        $notify.show({
          type: 'success',
          text: $t('info.botAdded'),
        })

        $store.updateBots()
      }
    } catch (error) {
      console.error(error)
    }
    loading.value = false
    $modals.bot.close()
  }

  async function botEditSent() {
    if (initialForm.is_active === true && initialForm.is_active !== form.is_active) {
      if (!(await doProcessQuestion($t('info.warning.botChange.text')))) return
    }

    loading.value = true
    try {
      const response = await $requests.bots.edit(form.uid, form)

      if (response?.uid) {
        $notify.show({
          type: 'success',
          text: $t('info.botAdited'),
        })

        $store.updateBots()
      }
    } catch (error) {
      console.error(error)
    }
    loading.value = false
    $modals.bot.close()
  }

  async function deleteBot() {
    if (!(await doProcessQuestion($t('info.warning.botDelete.text')))) return

    loading.value = true

    try {
      const response = await $requests.bots.del(form.uid)

      if (response?.uid) {
        $notify.show({
          type: 'success',
          text: $t('info.botDeleted'),
        })

        $store.updateBots()
      }
    } catch (error) {
      console.error(error)
    }

    loading.value = false

    $modals.bot.close()
  }

  function apply() {
    if (!checkFormFill()) {
      $notify.show({
        type: 'error',
        text: $t('error.fillAllFields'),
      })
      return
    }

    if (!checked.value) {
      $notify.show({
        type: 'error',
        text: $t('error.requiredAgreement'),
      })
      return
    }

    if (modalType.value.type === 'add') botAddSent()

    if (modalType.value.type === 'edit') botEditSent()
  }

  function close() {
    $modals.bot.close()
  }

  function checkTradePercent(el: HTMLInputElement) {
    if (!strategyDataByCode.value.is_trade_percent) return

    const _value = Number(el.value)

    if (_value > tradePercent.max) {
      form.trade_percent = tradePercent.max
      return
    }
    if (_value < tradePercent.min) {
      form.trade_percent = tradePercent.min
      return
    }
  }

  function checkMinMaxLimit(el?: HTMLInputElement) {
    if (!strategyDataByCode.value.is_trade_limit) return

    const _value = el ? Number(el.value) : form.trade_limit
    const min = $formatInt(strategyDataByCode.value.min_deposit)
    const max = $formatInt(strategyDataByCode.value.max_deposit)

    if (min && _value < min) {
      form.trade_limit = min
      return
    }
    if (max && _value > max) {
      form.trade_limit = max
      return
    }
  }

  function getStrategyByCode() {
    return $store.get('strategies').find((item) => item.code === botData.strategy_code)
  }

  function showInfo(type: EBotTradeType | EStrategyCategory) {
    if (type === EBotTradeType.limit) {
      $modals.info.show({
        title: `${$t('TradeType')} - ${$t('Limit')}`,
        text: $t('info.tradeType.limit'),
      })
    }

    if (type === EBotTradeType.percent) {
      $modals.info.show({
        title: `${$t('TradeType')} - ${$t('Percent')}`,
        text: $t('info.tradeType.percent'),
      })
    }

    if (type === EStrategyCategory.futures) {
      $modals.info.show({
        title: `${$t('Leverage')}`,
        text: $t('info.bot.leverage'),
      })
    }
  }

  function showReinvestInfo() {
    $modals.info.show({
      title: $t('Reinvest'),
      text: $t('info.tradeReinvest'),
    })
  }

  onBeforeMount(() => {
    setModalType()
    setBotType()
    checkMinMaxLimit()
  })

  onMounted(() => {
    $store.updateKeys()
  })
</script>

<template>
  <Modal :modal="$modals.bot" :title="modalType.title">
    <div class="form">
      <div class="form-control">
        <label class="label">
          <span class="label-text">
            {{ $t('Strategy') }}
          </span>
        </label>
        <div class="ml-1">
          <div class="flex justify-between">
            <p class="font-bold">{{ botData.name || getStrategyByCode()?.name || '-' }}</p>
            <!-- <p class="badge badge-outline badge-warning badge-sm">{{ botData.symbol_code }}</p> -->
          </div>

          <div class="space-y-2 my-3">
            <div class="text-sm flex flex-wrap gap-2 my-2">
              <template v-if="(strategyDataByCode?.symbol_codes || []).length">
                <p
                  v-for="symbol_code in strategyDataByCode.symbol_codes"
                  :key="symbol_code"
                  class="badge badge-warning badge-md"
                >
                  {{ symbol_code }}
                </p>
              </template>

              <template v-else>
                <p class="badge badge-warning badge-md">
                  {{ 'All coins' }}
                </p>
              </template>

              <!-- <p class="badge badge-warning badge-md">{{ botData.symbol_code || 'All coins' }}</p> -->

              <p v-if="strategyDataByCode?.exchange_code" class="badge badge-outline badge-md">
                <b>{{ strategyDataByCode.exchange_code }}</b>
              </p>

              <p v-if="strategyDataByCode?.category" class="badge badge-outline badge-md">
                <!-- {{ $t('Category') }}:&nbsp; -->
                <b>{{ $t(strategyDataByCode.category) }}</b>
              </p>
            </div>

            <div class="text-sm flex flex-wrap gap-2 my-2">
              <p
                v-if="strategyDataByCode?.risk_level"
                class="badge badge-outline badge-md"
                :class="$riskLevelBadge[strategyDataByCode.risk_level]"
              >
                {{ $t('RiskLevel') }}:&nbsp;
                <b>{{ strategyDataByCode.risk_level }}</b>
              </p>

              <p v-if="strategyDataByCode?.pos_profit > 0" class="badge badge-outline badge-md">
                {{ $t('before') }}&nbsp;+
                <b>{{ $formatInt(strategyDataByCode.pos_profit) }}%</b>
              </p>

              <p v-if="strategyDataByCode?.min_deposit" class="badge badge-outline badge-md">
                <!-- {{ $t('PotentialProfitability') }}:&nbsp; -->
                {{ $t('MinDeposit') }}:&nbsp;
                <b>{{ $formatInt(strategyDataByCode.min_deposit) }}</b>
              </p>

              <p v-if="strategyDataByCode?.max_deposit" class="badge badge-outline badge-md">
                <!-- {{ $t('PotentialProfitability') }}:&nbsp; -->
                {{ $t('MaxDeposit') }}:&nbsp;
                <b>{{ $formatInt(strategyDataByCode.max_deposit) }}</b>
              </p>

              <p v-if="strategyDataByCode?.share_profit > 0" class="badge badge-outline badge-md">
                {{ $t('ShareProfit') }}:&nbsp;
                <b>{{ $formatInt(strategyDataByCode.share_profit) }}%</b>
              </p>
            </div>
          </div>

          <div class="space-y-4">
            <p
              v-for="par in ($t(`strategies.${botData.code || botData.strategy_code}`) || '').split(
                '&n&',
              )"
              :key="par"
              class="text-xs italic"
            >
              {{ par }}
            </p>
          </div>
        </div>
      </div>

      <!-- <div class="form-control">
        <label class="label">
          <span class="label-text">{{ $t('ChoosingExchange') }}</span>
        </label>
        <select v-model="selectedExchange" class="select select-bordered w-full">
          <option disabled value=""></option>
          <option v-for="item in exchanges" :key="item.code" :value="item.code">
            {{ item.name }}
          </option>
        </select>
      </div> -->

      <div class="form-control">
        <label class="label">
          <span class="label-text">
            {{ $t('Key') }}
          </span>
        </label>
        <select v-model="form.key_uid" class="select select-bordered w-full">
          <option value="" :disabled="true"></option>
          <option v-for="item in keysByExchange" :key="item.uid" :value="item.uid">
            {{ item.name }}: {{ item.uid }}
          </option>
        </select>
      </div>

      <div v-if="isTradeTypeLimitAndPercent" class="form-control flex justify-start -mb-2">
        <span class="label-text">{{ $t('TradeType') }}</span>
        <div class="inline-flex items-center gap-2">
          <label class="label cursor-pointer gap-2">
            <input
              v-model="form.trade_type"
              type="radio"
              name="trade_type"
              class="radio"
              :value="EBotTradeType.percent"
              :checked="form.trade_type === EBotTradeType.percent"
            />
            <span class="label-text">{{ $t('Percent') }}</span>
          </label>
          <span class="inline-block cursor-pointer" @click="showInfo(EBotTradeType.percent)">
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
                d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"
              />
            </svg>
          </span>
        </div>

        <div class="inline-flex items-center gap-2">
          <label class="label cursor-pointer gap-2">
            <input
              v-model="form.trade_type"
              type="radio"
              name="trade_type"
              class="radio"
              :value="EBotTradeType.limit"
              :checked="form.trade_type === EBotTradeType.limit"
            />
            <span class="label-text">{{ $t('Limit') }}</span>
          </label>
          <span class="inline-block cursor-pointer" @click="showInfo(EBotTradeType.limit)">
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
                d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"
              />
            </svg>
          </span>
        </div>
      </div>

      <div
        v-if="
          strategyDataByCode.is_trade_percent &&
          (form.trade_type ? form.trade_type === EBotTradeType.percent : true)
        "
        class="form-control"
      >
        <label class="label">
          <span class="label-text">{{ $t('TradePercentOnDeal') }}</span>
        </label>
        <label class="input input-bordered flex items-center gap-2">
          <input
            v-model="form.trade_percent"
            type="number"
            :step="tradePercent.step"
            :min="tradePercent.min"
            :max="tradePercent.max"
            class="grow bg-transparent"
            name="trade_percent"
            list="trade_percent_list"
            @change="checkTradePercent($event.target as HTMLInputElement)"
          />
          %
        </label>
        <datalist id="trade_percent_list">
          <option v-for="(_, index) in new Array(100 / 5)" :key="index" :value="(index + 1) * 5" />
        </datalist>
      </div>

      <div
        v-if="
          strategyDataByCode.is_trade_limit &&
          (form.trade_type ? form.trade_type === EBotTradeType.limit : true)
        "
        class="form-control"
      >
        <label class="label">
          <span class="label-text">{{ $t('Limit') }} {{ ECurrencies.usdt.toUpperCase() }}</span>
        </label>
        <label class="input input-bordered flex items-center gap-2">
          <input
            v-model="form.trade_limit"
            type="number"
            class="grow bg-transparent"
            name="trade_limit"
            @change="checkMinMaxLimit($event.target as HTMLInputElement)"
          />
        </label>
      </div>

      <div v-if="strategyDataByCode.is_reinvest" class="form-control">
        <div class="inline-flex items-center gap-2">
          <label class="label justify-start cursor-pointer">
            <input v-model="form.trade_reinvest" type="checkbox" class="checkbox mr-2" />
            <span class="label-text">{{ $t('Reinvest') }}</span>
          </label>
          <span class="inline-block cursor-pointer" @click="showReinvestInfo">
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
                d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"
              />
            </svg>
          </span>
        </div>
      </div>

      <div v-if="ifShowLeverage" class="form-control">
        <!-- v-if="strategyDataByCode.category === EStrategyCategory.futures" -->
        <label class="label">
          <span class="label-text">
            {{ $t('Leverage') }}:
            <b>{{ form.leverage }}</b>
          </span>
          <span class="inline-block cursor-pointer" @click="showInfo(EStrategyCategory.futures)">
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
                d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"
              />
            </svg>
          </span>
        </label>
        <input
          v-model="form.leverage"
          type="range"
          :min="strategyDataByCode.min_leverage"
          :max="strategyDataByCode.max_leverage"
          class="range"
        />
        <div class="w-full flex justify-between text-xs px-2">
          <span>
            <b>{{ strategyDataByCode.min_leverage }}</b>
          </span>
          <span></span>
          <span>
            <b>{{ strategyDataByCode.max_leverage }}</b>
          </span>
        </div>
      </div>

      <div class="form-control inline-block">
        <label class="inline-flex cursor-pointer label gap-2">
          <span class="label-text">{{ $t('ActivityBot') }}:</span>
          <input
            v-model="form.is_active"
            type="checkbox"
            class="toggle toggle-primary"
            :checked="form.is_active"
          />
          <span class="label-text" :class="`${form.is_active ? 'text-success' : 'text-error'}`">
            {{ form.is_active ? $t('On') : $t('Off') }}
          </span>
        </label>
      </div>

      <!-- <div class="form-control">
        <label class="label">
          <span class="label-text">
            {{ $t('Cryptocurrency') }}
          </span>
        </label>
        <select v-model="form.currency" class="select select-bordered w-full">
          <option>Homer</option>
          <option>Marge</option>
          <option>Bart</option>
          <option>Lisa</option>
          <option>Maggie</option>
        </select>
      </div> -->

      <!-- <div class="form-control">
        <label class="label">
          <span class="label-text">
            {{ $t('Bid', { currency: $config.currency }) }}
          </span>
        </label>
        <input v-model="form.bid" type="number" class="w-full input input-bordered" name="bid" />
      </div> -->
    </div>

    <div class="form-control">
      <label class="label justify-start cursor-pointer">
        <input v-model="checked" type="checkbox" class="checkbox mr-2" />
        <span class="label-text">
          {{ $t('ApplyAgreement') }}.
          <a class="link link-info" href="#" @click="openModalDocuments">
            {{ $t('Show') }} {{ $t('Documents').toLowerCase() }}
          </a>
        </span>
      </label>
    </div>

    <template #actions>
      <div class="flex flex-1 justify-between">
        <div>
          <button
            v-if="modalType.type === 'edit'"
            class="btn btn-error btn-outline"
            @click="deleteBot"
          >
            {{ $t('Delete') }}
          </button>
        </div>

        <div class="flex gap-4">
          <button class="btn btn-primary" @click="apply">
            {{ $t('Apply') }}
          </button>
          <button class="btn" @click="close">
            {{ $t('Close') }}
          </button>
        </div>
      </div>
    </template>
  </Modal>
</template>
