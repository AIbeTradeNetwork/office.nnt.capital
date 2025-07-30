import { obfuscator } from 'rollup-obfuscator'

/**
 * Obfuscator Config
 *
 * @see {@link https://github.com/javascript-obfuscator/javascript-obfuscator#preset-options}
 */
const obfuscatorConfig = {
  optionsPreset: 'high-obfuscation',
  log: true,
  debugProtection: false,
  debugProtectionInterval: 0,
  disableConsoleOutput: true,
  deadCodeInjection: false,
  domainLock: [],
  stringArray: false,
  identifierNamesGenerator: 'hexadecimal',
  renamePropertiesMode: 'safe',
  stringArrayIndexesType: ['hexadecimal-number'],
}

export { obfuscator, obfuscatorConfig }
