<!-- https://tradingview.github.io/lightweight-charts/ -->

<script lang="ts">
  export default {
    name: 'ChartDefault',
  }
</script>

<script setup lang="ts">
  import { ref, onMounted, Ref, onUnmounted } from 'vue'
  import {
    createChart,
    CrosshairMode,
    IChartApi,
    ISeriesApi,
    LineStyle,
    LineWidth,
    SeriesType,
  } from 'lightweight-charts'
  import { $theme } from 'utils/theme'

  const root: Ref<HTMLElement> = ref()
  let chart: IChartApi
  let mianSeriesType: ISeriesApi<SeriesType>

  const data = [
    { time: '2019-04-11', value: 80.01 },
    { time: '2019-04-12', value: 96.63 },
    { time: '2019-04-13', value: 76.64 },
    { time: '2019-04-14', value: 81.89 },
    { time: '2019-04-15', value: 74.43 },
    { time: '2019-04-16', value: 80.01 },
    { time: '2019-04-17', value: 96.63 },
    { time: '2019-04-18', value: 76.64 },
    { time: '2019-04-19', value: 81.89 },
    { time: '2019-04-20', value: 74.43 },
  ]

  const colors = function () {
    return {
      text: $theme.colors().baseContent,
      vertLines: 'transparent',
      horzLines: $theme.colors().baseContent,
      font: getComputedStyle(document.querySelector('html')).getPropertyValue('fontFamily'),
    }
  }

  const chartOptions = () => {
    return {
      autoSize: true,
    }
  }

  function resizing(remove?) {
    function resize() {
      if (!chart) return
      chart.resize(
        root.value.getBoundingClientRect().width,
        root.value.getBoundingClientRect().height,
      )
    }

    if (remove) return window.removeEventListener('resize', () => resize)
    window.addEventListener('resize', () => resize)
  }

  function createTheme() {
    chart.applyOptions({
      layout: {
        textColor: colors().text,
        background: { color: 'transparent' },
      },
      grid: {
        vertLines: { color: colors().vertLines },
        horzLines: { color: colors().horzLines },
      },
      crosshair: {
        // Change mode from default 'magnet' to 'normal'.
        // Allows the crosshair to move freely without snapping to datapoints
        mode: CrosshairMode.Normal,

        // Vertical crosshair line (showing Date in Label)
        vertLine: {
          width: 100 as LineWidth,
          color: '#C3BCDB44',
          style: LineStyle.Solid,
          labelBackgroundColor: '#9B7DFF',
        },

        // Horizontal crosshair line (showing Price in Label)
        horzLine: {
          color: '#9B7DFF',
          labelBackgroundColor: '#9B7DFF',
        },
      },
    })

    mianSeriesType.applyOptions({
      color: 'green',
      wickUpColor: 'rgb(54, 116, 217)',
      upColor: 'rgb(54, 116, 217)',
      wickDownColor: 'rgb(225, 50, 85)',
      downColor: 'rgb(225, 50, 85)',
      borderVisible: true,
      title: 'title',
    })
  }

  function createChar() {
    chart = createChart(root.value, chartOptions())
    mianSeriesType = chart.addLineSeries()
    mianSeriesType.setData(data)
    chart.timeScale().applyOptions({
      borderColor: 'red',
    })

    chart.timeScale().fitContent()
  }

  onMounted(() => {
    createChar()
    createTheme()
    resizing()
  })

  onUnmounted(() => {
    if (chart) {
      chart.remove()
      chart = null
    }
    resizing('clear')
  })
</script>

<template>
  <div ref="root" class="box w-[500px] h-[500px]"></div>
</template>
