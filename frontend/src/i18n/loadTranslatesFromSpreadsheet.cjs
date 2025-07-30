// https://theoephraim.github.io/node-google-spreadsheet/#/guides/exports

'use strict'

const fs = require('fs')
const fsp = require('fs/promises')
const { resolve } = require('path')
const { GoogleSpreadsheet } = require('google-spreadsheet')
const fileExt = 'tsv'
const spreadsheet = '1DTnWJ4fxZxVcZSKZJwA-B9ipkdwMiJKXHL1H48ehyVo'
const key = 'AIzaSyDU3WU3WZaLVqvj7Dff27irrMrZFIaSJZY'
const translateObject = {}
const tabs = require('./translateTabs.json')

function getUrlToFile(name, ext, path = '') {
  return resolve(__dirname, `${path ? path + '/' : ''}${name}.${ext || fileExt}`)
}

function clearText(text) {
  return (text || '').replace('\r', '')
}

async function transformFile(name) {
  const rawdata = fs.readFileSync(getUrlToFile(name))
  const content = rawdata.toString().replace(/\r/, '')
  const rows = content.split(/\n/)
  const columns = rows[0].split(/\t/)

  for (let i = 1; i < columns.length; i++) {
    const lang = columns[i]
    const data = {}

    rows.forEach((line) => {
      const row = line.split('\t')
      const keys = row[0].split('.')

      if (keys.length > 1) {
        let obj = data

        keys.forEach((val, index) => {
          if (index + 1 === keys.length) {
            obj[val] = clearText(row[i])
            return
          }

          if (!obj[val]) obj[val] = {}

          obj = obj[val]
        })
      } else {
        if (row[0]) data[row[0]] = clearText(row[i])
      }
    })

    if (!translateObject[lang]) translateObject[lang] = {}

    Object.assign(translateObject[lang], data)
  }
}

;(async () => {
  const doc = new GoogleSpreadsheet(spreadsheet, {
    apiKey: key,
  })

  await doc.loadInfo()

  for await (const tabName of Object.keys(tabs)) {
    console.log(`\x1b[33m ${tabName}: start \x1b[0m`)

    try {
      const xlsxBuffer = await doc._downloadAs(fileExt, tabs[tabName])
      await fsp.writeFile(getUrlToFile(tabName), Buffer.from(xlsxBuffer))
      await transformFile(tabName)
      await fsp.unlink(getUrlToFile(tabName))
    } catch (error) {
      console.log(`\x1b[31m ${tabName}: error \x1b[0m`)
      console.log(error)
      return
    }

    console.log(`\x1b[32m ${tabName}: done \x1b[0m`)
  }

  Object.keys(translateObject).forEach((key) => {
    fs.writeFileSync(
      getUrlToFile(key, 'json', 'locales'),
      JSON.stringify(translateObject[key], null, 2),
    )
    // fs.copyFile('./en.json', './en.json', (err) => {
    //   if (err) throw err
    // })
  })
})()
