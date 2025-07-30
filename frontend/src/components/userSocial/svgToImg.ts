// https://github.com/tsayen/dom-to-image#usage
// https://github.com/tatsuyasusukida/svg-to-png/blob/main/src/app/original/page.tsx

import domtoimage from 'dom-to-image'
import { $i18nGlobal } from 'i18n/index'

const fileName = 'social'
const fileExt = 'jpg'
const dateTime = () => {
  const date = new Date()

  const formatedDate = Intl.DateTimeFormat($i18nGlobal.locale.value, {
    year: 'numeric',
    month: 'numeric',
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
    hour12: false,
  }).format(new Date())

  return formatedDate.replace(/,\s/g, '_').replace(/:/g, '.')
}

function triggerDownload(imgURI: string): Promise<true> {
  return new Promise((resolve, reject) => {
    try {
      const evt = new MouseEvent('click', {
        view: window,
        bubbles: false,
        cancelable: true,
      })
      const a = document.createElement('a')
      a.setAttribute('download', `${fileName}_${dateTime()}.${fileExt}`)
      a.setAttribute('href', imgURI)
      a.setAttribute('target', '_blank')
      a.dispatchEvent(evt)
      resolve(true)
    } catch (err) {
      reject(err)
    }
  })
}

function downloadSvg(svg: SVGElement): Promise<true> {
  return new Promise((resolve, reject) => {
    const svgData = new XMLSerializer().serializeToString(svg)
    const svgDataBase64 = btoa(unescape(encodeURIComponent(svgData)))
    const svgDataUrl = `data:image/svg+xml;charset=utf-8;base64,${svgDataBase64}`

    let image = new Image()

    async function load() {
      const width = '1080px'
      const height = '1920px'
      const canvas = document.createElement('canvas')

      canvas.setAttribute('width', width)
      canvas.setAttribute('height', height)

      const context = canvas.getContext('2d') as CanvasRenderingContext2D
      context.drawImage(image, 0, 0, parseInt(width, 10), parseInt(height, 10))

      const dataUrl = canvas.toDataURL(`image/${fileName}.${fileExt}`)

      try {
        await triggerDownload(dataUrl)
        resolve(true)
      } catch (err) {
        reject(err)
      }

      image.removeEventListener('load', load)
      image = undefined
    }

    image.addEventListener('load', load)
    image.src = svgDataUrl
  })
}

function saveImage(id: string): Promise<boolean> {
  return new Promise((resolve, reject) => {
    try {
      if (!id) {
        reject(false)
        return
      }

      domtoimage.toSvg(document.getElementById(id)).then(async (dataUrl) => {
        let div = document.createElement('div')
        div.innerHTML = dataUrl

        await downloadSvg(div.querySelector('svg'))

        resolve(true)

        div = undefined
      })
    } catch (err) {
      reportError(err)
      reject(err)
    }
  })
}

export { saveImage }
