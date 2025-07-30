export enum EAlerts {
  'BASE' = 'base',
  'INFO' = 'info',
  'SUCCESS' = 'success',
  'WARNING' = 'warning',
  'ERORR' = 'error',
}

export enum ETeamType {
  'auto' = 'auto',
  'left' = 'left',
  'right' = 'right',
}

export enum ESizes {
  'BASE' = 'base',
  'XS' = 'xs',
  'SM' = 'sm',
  'MD' = 'md',
  'LG' = 'lg',
  'XL' = 'xl',
  '2XL' = '2xl',
}

export enum ELocales {
  en_US = 'en-US',
  ru_RU = 'ru-RU',
}

export enum ERoles {
  'client' = 'client',
  'distributor' = 'distributor',
}

export enum EPlans {
  'start',
  'advanced',
  'professional',
  'silver',
  'gold',
  'platinum',
  'brilliant',
}

export enum EPlansVIP {
  'silver',
  'gold',
  'platinum',
  'brilliant',
}

export enum ERangs {
  'silver',
  'gold',
  'platinum',
  'brilliant',
  'fly',
  'keen',
  'smart',
  'savvy',
  'hero',
  'ace',
  'brain',
  'whiz',
  'sage',
  'elite',
  'legend',
}

export enum ECurrencies {
  '' = '',
  udex = 'udex',
  usd = 'usd',
  usdHidden = 'usd',
  usdt = 'usdt',
  cv = 'cv',
  ton = 'ton',
}

export enum EPayoutCodes {
  usdt_bep_20 = 'USDT_BEP20',
  usdt_ton = 'USDT_TON',
}

export enum EStrategyType {
  classic = 'classic',
  ico = 'ico',
}

export enum EStrategyCategory {
  spot = 'spot',
  futures = 'futures',
}

export enum EBotTradeType {
  limit = 'limit',
  percent = 'percent',
}

export enum ETradeSide {
  buy = 'buy',
  sell = 'sell',
}

export enum ECoinFarmingActions {
  'farming' = 'farming',
  'invite' = 'invite',
  'team' = 'team',
  'subscription' = 'subscription',
  'tasks' = 'tasks',
}

export enum EABTRanks {
  'Novice' = 'Novice',
  'Farmer' = 'Farmer',
  'Supervisor' = 'Supervisor',
  'Expert' = 'Expert',
  'Master' = 'Master',
  'Champion' = 'Champion',
  'Veteran' = 'Veteran',
  'Legionary' = 'Legionary',
  'Grandmaster' = 'Grandmaster',
  'Archon' = 'Archon',
  "archon-I" = 'archon-I',
  'archon-II' = 'archon-II',
  'archon-III' = 'archon-III',
  'archon-IV' = 'archon-IV',
  'archon-V' = 'archon-V',
  'archon-VI' = 'archon-VI',
  "archon-VII" = "archon-VII",
  "archon-VIII" = "archon-VIII",
  "archon-IX" = "archon-IX",
  "archon-X" = "archon-X",
}

export const CABTCounts = {
  [EABTRanks.Novice]: {
    id: 0,
    code: EABTRanks.Novice,
    min: 0,
    max: 50,
    invitingMax: 5,
    farming: 0.9,
    img: new URL('coin/assets/abt_ranks/tier-1/Novice.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Novice.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Farmer]: {
    id: 1,
    code: EABTRanks.Farmer,
    min: 50,
    max: 150,
    invitingMax: 10,
    farming: 1.2,
    img: new URL('coin/assets/abt_ranks/tier-1/Farmer.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Farmer.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Supervisor]: {
    id: 2,
    code: EABTRanks.Supervisor,
    min: 150,
    max: 500,
    invitingMax: 20,
    farming: 1.5,
    img: new URL('coin/assets/abt_ranks/tier-1/Supervisor.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Supervisor.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Expert]: {
    id: 3,
    code: EABTRanks.Expert,
    min: 500,
    max: 1000,
    invitingMax: 30,
    farming: 1.8,
    img: new URL('coin/assets/abt_ranks/tier-1/Expert.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Expert.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Master]: {
    id: 4,
    code: EABTRanks.Master,
    min: 1000,
    max: 2500,
    invitingMax: 40,
    farming: 2.1,
    img: new URL('coin/assets/abt_ranks/tier-1/Master.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Master.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Champion]: {
    id: 5,
    code: EABTRanks.Champion,
    min: 2500,
    max: 5000,
    invitingMax: 50,
    farming: 2.4,
    img: new URL('coin/assets/abt_ranks/tier-1/Champion.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Champion.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Veteran]: {
    id: 6,
    code: EABTRanks.Veteran,
    min: 5000,
    max: 8000,
    invitingMax: 60,
    farming: 2.7,
    img: new URL('coin/assets/abt_ranks/tier-1/Veteran.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Veteran.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Legionary]: {
    id: 7,
    code: EABTRanks.Legionary,
    min: 8000,
    max: 15000,
    invitingMax: 80,
    farming: 3.0,
    img: new URL('coin/assets/abt_ranks/tier-1/Legionary.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Legionary.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Grandmaster]: {
    id: 8,
    code: EABTRanks.Grandmaster,
    min: 15000,
    max: 30000,
    invitingMax: 150,
    farming: 4.2,
    img: new URL('coin/assets/abt_ranks/tier-1/Grandmaster.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Grandmaster.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },
  [EABTRanks.Archon]: {
    id: 9,
    code: EABTRanks.Archon,
    min: 30000,
    max: 60000,
    invitingMax: 300,
    farming: 5.1,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: null,
    img_bottom: null,
  },

  [EABTRanks["archon-I"]]: {
    id: 10,
    code: EABTRanks["archon-I"],
    min: 60000,
    max: 150000,
    invitingMax:  -1,
    farming: 12,
    img: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t1.png', import.meta.url).href,
    img_bottom: null
  },
  [EABTRanks["archon-II"]]: {
    id: 11,
    code: EABTRanks["archon-II"],
    min: 150000,
    max: 300000,
    invitingMax:  -1,
    farming: 15,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t2.png', import.meta.url).href,
    img_bottom: null,
  },
  [EABTRanks["archon-III"]]: {
    id: 12,
    code: EABTRanks["archon-III"],
    min: 300000,
    max: 600000,
    invitingMax:  -1,
    farming: 18,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t3.png', import.meta.url).href,
    img_bottom: null,
  },
  [EABTRanks["archon-IV"]]: {
    id: 13,
    code: EABTRanks["archon-IV"],
    min: 600000,
    max: 1500000,
    invitingMax:  -1,
    farming: 21,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t4.png', import.meta.url).href,
    img_bottom: null,
  },
  [EABTRanks["archon-V"]]: {
    id: 14,
    code: EABTRanks["archon-V"],
    min: 1500000,
    max: 3000000,
    invitingMax:  -1,
    farming: 24,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t5.png', import.meta.url).href,
    img_bottom: null,
  },
  [EABTRanks["archon-VI"]]: {
    id: 15,
    code: EABTRanks["archon-VI"],
    min: 3000000,
    max: 6000000,
    invitingMax:  -1,
    farming: 27,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t6.png', import.meta.url).href,
    img_bottom: null
  },
  [EABTRanks["archon-VII"]]: {
    id: 16,
    code: EABTRanks["archon-VII"],
    min: 6000000,
    max: 12000000,
    invitingMax:  -1,
    farming: 30,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t7.png', import.meta.url).href,
    img_bottom: new URL('coin/assets/abt_ranks/tier-2/b1.png', import.meta.url).href,
  },
  [EABTRanks["archon-VIII"]]: {
    id: 17,
    code: EABTRanks["archon-VIII"],
    min: 12000000,
    max: 20000000,
    invitingMax:  -1,
    farming: 42,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t8.png', import.meta.url).href,
    img_bottom: new URL('coin/assets/abt_ranks/tier-2/b2.png', import.meta.url).href,
  },
  [EABTRanks["archon-IX"]]: {
    id: 18,
    code: EABTRanks["archon-IX"],
    min: 20000000,
    max: 40000000,
    invitingMax:  -1,
    farming: 51,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t9.png', import.meta.url).href,
    img_bottom: new URL('coin/assets/abt_ranks/tier-2/b3.png', import.meta.url).href,
  },
  [EABTRanks["archon-X"]]: {
    id: 19,
    code: EABTRanks["archon-X"],
    min: 40000000,
    max: null,
    invitingMax: -1,
    farming: 200,
    img: new URL('coin/assets/abt_ranks/tier-1/Archon.png', import.meta.url).href,
    img_premium: new URL('coin/assets/abt_ranks/tier-1/premium_Archon.png', import.meta.url).href,
    img_top: new URL('coin/assets/abt_ranks/tier-2/t10.png', import.meta.url).href,
    img_bottom: new URL('coin/assets/abt_ranks/tier-2/b4.png', import.meta.url).href,
  },
}

export enum EHideClasses {
  autotrade = 'autotrade',
}

export enum BrokerTransactionType {
  'all' = '',
  'deposit' = 'deposit',
  'withdraw' = 'withdraw',
  'trade' = 'trade',
}

export enum SearchParamsModals {
  'payin' = 'payin',
}

export enum Precision {
  'udex' = 9,
  'usd' = 2,
}