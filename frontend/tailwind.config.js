/** @type {import('tailwindcss').Config} */

// import defaultTheme from 'tailwindcss/defaultTheme'
import plugin from 'tailwindcss/plugin'
import { $defaultsScreenToTilwind } from './src/utils/screen'
import { $planColors, $riskColors } from './src/utils/colors'

const defaultsCustomScreen = $defaultsScreenToTilwind()
// sm: '640px',
// md: '768px',
// lg: '1024px',
// xl: '1280px',
// '2xl': '1536px'

export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],

  plugins: [
    plugin(function ({ matchUtilities, theme }) {
      matchUtilities(
        {
          'text-shadow': (value) => ({
            textShadow: value,
          }),
        },
        { values: theme('textShadow') },
      )
    }),
    require('@tailwindcss/typography'),
    require('daisyui'),
  ],

  theme: {
    screens: {
      mobile: { min: defaultsCustomScreen.zr, max: defaultsCustomScreen.lg },
      pc: defaultsCustomScreen.lg,
      ...defaultsCustomScreen,
    },

    safelist: ['overflow-hidden', 'fixed'],

    // 'padding-x': ({ theme }) => ({
    //   gray: theme('colors.gray')
    // })
    //   @apply mobile:px-2 pc:px-6;
    // },
    // 'padding-x': {
    //   @apply mx-auto mobile:px-2 pc:px-6;
    // }
    // spacing: {
    //   1: '8px',
    //   2: '12px',
    //   3: '16px',
    //   4: '24px',
    //   5: '32px',
    //   6: '48px',
    // },
    // boxShadow: {
    // DEFAULT: '0 1px 5px 1px rgb(0 0 0 / .8)',
    // },
    extend: {
      colors: {
        ...$planColors,
        ...$riskColors,
        white: '#ffffff',
        black: {
          DEFAULT: '#222222',
          300: '#101117',
        },
      },

      dropShadow: {
        custom: `0 3px 5px rgba(0 0 0 / .4)`,
      },

      textShadow: {
        sm: '0 1px 2px var(--tw-shadow-color)',
        DEFAULT: '0 2px 4px var(--tw-shadow-color)',
        lg: '0 8px 16px var(--tw-shadow-color)',
      },
    },
  },

  daisyui: {
    // https://daisyui.com/docs/config/?lang=ru
    // styled: true,
    // base: true,
    // utils: true,
    // logs: true,
    themes: ['dim'],
    // [
    //   {
    //     dim: {
    //       ...require('daisyui/src/theming/themes')['dim'],
    //     },
    //   },
    // ],
    darkTheme: 'dim',
    // themeRoot: 'rem',
    // [
    // 'dim',
    // {
    //   light: {
    //     ...require('daisyui/src/theming/themes')['emerald'],
    //   },
    // },
    // {
    //   dark: {
    //     ...require('daisyui/src/theming/themes')['dim'],
    //   },
    // },
    // {
    //   storm: {
    //     primary: '#9aa5ce',
    //     secondary: '#565f89',
    //     accent: '#bb9af7',
    //     neutral: '#111827',
    //     'base-100': '#24283b',
    //     info: '#2ac3de',
    //     success: '#9ece6a',
    //     warning: '#e0af68',
    //     error: '#f7768e',
    //   },
    // },
    // ], // false: only light + dark | true: all themes | array: specific themes like this ["light", "dark", "cupcake"]
    // darkTheme: 'luxury', // name of one of the included themes for dark mode
    // base: true, // applies background color and foreground color for root element by default
    // styled: true, // include daisyUI colors and design decisions for all components
    // utils: true, // adds responsive and modifier utility classes
    // prefix: '', // prefix for daisyUI classnames (components, modifiers and responsive class names. Not colors)
    // logs: true, // Shows info about daisyUI version and used config in the console when building your CSS
    // themeRoot: ':root', // The element that receives theme color CSS variables
  },
}
