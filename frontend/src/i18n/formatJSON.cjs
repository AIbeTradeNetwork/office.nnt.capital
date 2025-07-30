// https://theoephraim.github.io/node-google-spreadsheet/#/guides/exports

'use strict'

const fs = require('fs')
const fsp = require('fs/promises')
const { resolve } = require('path')
const fileExt = 'tsv'
const files = process.argv.splice(2)[0].split(',')
let fileString = ''

function getUrlToFile(name, ext, path = '') {
  return resolve(__dirname, `${path ? path + '/' : ''}${name}.${ext || fileExt}`)
}

async function transformFile(name) {
  const rawdata = fs.readFileSync(getUrlToFile(name, 'json', 'locales'))
  const content = JSON.parse(rawdata)
  let fileString = ''

  function parse(obj, path) {
    Object.keys(obj).forEach((key) => {
      if (typeof obj[key] === 'object') {
        parse(obj[key], !path ? key : path + '.' + key)
        return
      }

      fileString += `${path ? path + '.' + key : key}\t${obj[key]}\n`
    })
  }

  parse(content)

  fs.writeFileSync(getUrlToFile(name), fileString)
}

async function saveFile(name) {
  fs.writeFileSync(getUrlToFile(name + '_new', 'json'), fileString)
  fileString = ''
}

;(async () => {
  for await (const fileName of files) {
    await transformFile(fileName)
    // await saveFile(fileName)
  }
})()
