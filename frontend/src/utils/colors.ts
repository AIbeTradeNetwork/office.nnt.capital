const $planColors = {
  transparent: 'rgba(255,255,255,0)',
  vip: '#f08c78',
  sstart: '#dedede',
  advanced: '#dedede',
  professional: '#dedede',
  silver: '#D3D6D8',
  gold: '#D9A400',
  platinum: '#B0B5B9',
  brilliant: '#4F98CB',
}

const $riskColors = {
  low: '#00ac46',
  medium: '',
  high: '#dc0000',
} as const

const $riskLevelBadge = {
  low: 'badge-success',
  medium: '',
  high: 'badge-error',
}

// const $textColors = {
//   silver: 'text-silver',
//   vip: '#333333',
//   gold: 'text-gold',
//   platinum: 'text-platinum',
//   brilliant: 'text-brilliant',
// }

export { $planColors, $riskColors, $riskLevelBadge }
