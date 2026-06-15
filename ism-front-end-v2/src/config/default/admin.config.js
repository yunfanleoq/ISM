// admin 配置
const ADMIN = {
  palettes: [],
  animates: require('./animate.config').preset,
  Canimates: require('./animate.config').ComponentPreset,
  theme: {
    mode: {
      DARK: 'dark',
      LIGHT: 'light',
      NIGHT: 'night'
    }
  },
  layout: {
    SIDE: 'side',
    HEAD: 'head'
  }
}

module.exports = ADMIN
