import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
const { describe, toDataUri, DEFAULT_STYLE } = require('./symbolSvg')
const { DC_SYMBOLS } = require('./dcEnvSymbols')
const { SLD_SYMBOLS } = require('./sldSymbols')

const SYMBOLS = {}
DC_SYMBOLS.concat(SLD_SYMBOLS).forEach(function (item) {
  SYMBOLS[item.id] = item
})

function animateConfig() {
  return {
    selected: [],
    condition: {
      deviceSN: '',
      selectVideoType: 0,
      isBandDevice: false,
      bandType: 1,
      dataID: '',
      dataName: '',
      operator: '',
      OperatorValue: '',
      OperatorMaxValue: ''
    },
    isExpression: false,
    animateList: [
      { id: 'Forbidden', name: 'component.public.Forbidden' },
      { id: 'blink', name: 'component.public.animateBlink' },
      { id: 'millcolorGrad', name: 'component.public.millcolorGrad' },
      { id: 'Zoom', name: 'component.public.Zoom' },
      { id: 'animateSpin', name: 'component.public.animateSpin' }
    ],
    animateElement: [
      {
        id: 'blink',
        elementList: [
          { name: 'component.public.animateSpeed', type: 7, value: 1, min: 0.1, key: 'blinkSpeed' }
        ]
      },
      {
        id: 'millcolorGrad',
        elementList: [
          { name: 'component.public.startColor', type: 2, value: '#74f808', key: 'startColor' },
          { name: 'component.public.stopColor', type: 2, value: '#f30b0b', key: 'stopColor' },
          { name: 'component.public.animateSpeed', type: 7, value: 1, min: 0.1, key: 'animateSpeed' }
        ]
      },
      {
        id: 'animateSpin',
        elementList: [
          { name: 'component.public.animateSpinSpeed', type: 7, value: 1, min: 0.1, key: 'spinSpeed' },
          {
            name: 'configComponent.bigScreen.border.border89Direction',
            type: 6,
            value: 0,
            enumList: [
              { value: 0, option: 'configComponent.bigScreen.border.border89DirectionForward' },
              { value: 1, option: 'configComponent.bigScreen.border.border89DirectionNegative' }
            ],
            min: 1,
            key: 'spinDirection'
          }
        ]
      }
    ]
  }
}

function diyOf(meta) {
  const diy = [
    { name: 'component.Electric.ElectronicDeviceWidth', type: 7, value: DEFAULT_STYLE.sw, min: 0.4, key: 'strokeWidth' },
    { name: 'component.Electric.ElectronicDeviceColor', type: 2, value: DEFAULT_STYLE.stroke, key: 'strokeColor' },
    { name: 'component.Electric.ElectronicDeviceFill', type: 2, value: DEFAULT_STYLE.fill, key: 'fillColor' },
    { name: 'component.public.fillOpacity', type: 7, value: 1, min: 0.1, max: 1, key: 'fillOpacity' }
  ]
  if (meta.state) {
    diy.push({
      name: meta.state.name,
      type: 6,
      value: meta.state.default,
      enumList: meta.state.options,
      key: 'symbolState'
    })
  }
  return diy
}

function cssSize(v, fallback) {
  if (v == null || v === '') return fallback + 'px'
  if (typeof v === 'number') return v + 'px'
  const s = String(v)
  if (/px|%|em|rem|vw|vh$/.test(s)) return s
  return s + 'px'
}

function selectedList(option) {
  const selected = option && option.animate && option.animate.selected
  return Array.isArray(selected) ? selected : []
}

export function createPowerSymbol(id) {
  const meta = SYMBOLS[id]
  if (!meta) {
    throw new Error('unknown power symbol: ' + id)
  }
  const thumb = toDataUri(meta, {
    stroke: DEFAULT_STYLE.stroke,
    fill: DEFAULT_STYLE.fill,
    sw: DEFAULT_STYLE.sw,
    state: meta.state ? meta.state.default : 1
  })

  return {
    name: meta.name,
    mixins: [ISMChildAutoMixin],
    inject: {
      getNode: { default: null }
    },
    data: function () {
      return {
        detail: {
          style: {
            position: { x: 0, y: 0, w: meta.w, h: meta.h },
            visible: 1,
            backColor: 'transparent',
            zIndex: -1,
            transform: 0,
            diy: diyOf(meta)
          },
          animate: animateConfig()
        },
        IsToolBox: false,
        editMode: true,
        strokeColor: DEFAULT_STYLE.stroke,
        fillColor: DEFAULT_STYLE.fill,
        strokeWidth: DEFAULT_STYLE.sw,
        fillOpacity: 1,
        symbolState: meta.state ? meta.state.default : 1,
        animateType: [],
        startColor: '#74f808',
        stopColor: '#f30b0b',
        animateSpeed: 0.5,
        animateSpinSpeed: 0.5,
        spinDirection: 0,
        blinkSpeed: 0.5,
        isStart: false,
        base: {
          text: meta.text,
          icon: thumb,
          isFontIcon: true,
          info: {
            type: meta.name,
            action: [],
            dataBind: [],
            animate: animateConfig(),
            style: {
              position: { x: 0, y: 0, w: meta.w, h: meta.h },
              visible: 1,
              backColor: 'transparent',
              zIndex: -1,
              transform: 0,
              diy: diyOf(meta)
            }
          }
        }
      }
    },
    watch: {
      detail: {
        handler: function (val) {
          if (this.editMode) this.initComponents(val)
        },
        deep: true
      }
    },
    methods: {
      initComponents: function (option) {
        if (!option || this.IsToolBox) return
        const diy = (option.style && option.style.diy) || []
        for (let i = 0; i < diy.length; i++) {
          const item = diy[i]
          if (item.key === 'strokeWidth') this.strokeWidth = item.value
          else if (item.key === 'strokeColor') this.strokeColor = item.value
          else if (item.key === 'fillColor') this.fillColor = item.value
          else if (item.key === 'fillOpacity') this.fillOpacity = item.value
          else if (item.key === 'symbolState') this.symbolState = Number(item.value)
        }
        this.animateType = selectedList(option)
        this.isStart = !(option.animate && option.animate.isExpression)
        const blocks = (option.animate && option.animate.animateElement) || []
        for (let i = 0; i < blocks.length; i++) {
          const block = blocks[i]
          const list = block.elementList || []
          for (let k = 0; k < list.length; k++) {
            const el = list[k]
            if (block.id === 'millcolorGrad' && el.key === 'startColor') this.startColor = el.value
            else if (block.id === 'millcolorGrad' && el.key === 'stopColor') this.stopColor = el.value
            else if (block.id === 'millcolorGrad' && el.key === 'animateSpeed') this.animateSpeed = el.value
            else if (block.id === 'blink' && el.key === 'blinkSpeed') this.blinkSpeed = el.value
            else if (block.id === 'animateSpin' && el.key === 'spinSpeed') this.animateSpinSpeed = el.value
            else if (block.id === 'animateSpin' && el.key === 'spinDirection') this.spinDirection = el.value
          }
        }
      }
    },
    render: function (h) {
      const vm = this
      const pos = (vm.detail && vm.detail.style && vm.detail.style.position) || {}
      const nodes = describe(meta, {
        stroke: vm.strokeColor,
        fill: vm.fillColor,
        sw: vm.strokeWidth,
        state: vm.symbolState
      })
      const children = nodes.map(function (n) {
        if (n.tag === 'text') return h('text', { attrs: n.attrs }, n.text)
        return h(n.tag, { attrs: n.attrs })
      })
      const types = vm.animateType || []
      if (vm.isStart && !vm.IsToolBox && types.indexOf('blink') !== -1) {
        children.push(h('animate', {
          attrs: {
            attributeName: 'opacity',
            values: '0.1;1;0.1',
            dur: vm.blinkSpeed + 's',
            repeatCount: 'indefinite'
          }
        }))
      }
      if (vm.isStart && !vm.IsToolBox && types.indexOf('millcolorGrad') !== -1) {
        children.push(h('animate', {
          attrs: {
            attributeName: 'stroke',
            values: vm.startColor + ';' + vm.stopColor + ';' + vm.startColor,
            dur: vm.animateSpeed + 's',
            repeatCount: 'indefinite'
          }
        }))
      }
      const vb = (meta.viewBox || '0 0 80 96').split(/\s+/).map(Number)
      const cx = (vb[2] || 80) / 2
      const cy = (vb[3] || 96) / 2
      if (vm.isStart && !vm.IsToolBox && types.indexOf('Zoom') !== -1) {
        children.push(h('animateTransform', {
          attrs: {
            attributeName: 'transform',
            type: 'scale',
            values: '0.9;1;0.9',
            dur: '0.6s',
            repeatCount: 'indefinite',
            additive: 'sum'
          }
        }))
      }
      if (vm.isStart && !vm.IsToolBox && types.indexOf('animateSpin') !== -1) {
        const from = vm.spinDirection == 1 ? '360 ' + cx + ' ' + cy : '0 ' + cx + ' ' + cy
        const to = vm.spinDirection == 1 ? '0 ' + cx + ' ' + cy : '360 ' + cx + ' ' + cy
        children.push(h('animateTransform', {
          attrs: {
            attributeName: 'transform',
            type: 'rotate',
            from: from,
            to: to,
            dur: vm.animateSpinSpeed + 's',
            repeatCount: 'indefinite'
          }
        }))
      }
      return h('svg', {
        attrs: {
          xmlns: 'http://www.w3.org/2000/svg',
          viewBox: meta.viewBox,
          preserveAspectRatio: 'xMidYMid meet'
        },
        style: {
          overflow: 'visible',
          width: cssSize(pos.w, meta.w),
          height: cssSize(pos.h, meta.h)
        }
      }, [
        h('g', {
          style: { opacity: vm.fillOpacity }
        }, children)
      ])
    },
    mounted: function () {
      const vm = this
      this.$nextTick(function () {
        vm.initComponents(vm.detail)
        if (!vm.detail || !vm.detail.identifier) return
        vm.$EventBus.$on(vm.detail.identifier + 'activeEvent', function () {})
        vm.$EventBus.$on(vm.detail.identifier + 'animateEvent', function (data) {
          vm.isStart = data
        })
      })
    },
    created: function () {
      const vm = this
      const node = typeof this.getNode === 'function' ? this.getNode() : null
      if (!node || typeof node.getData !== 'function') return
      node.on('change:data', function ({ current }) {
        if (current) vm.detail = current.detail
      })
      node.on('change:size', function ({ current }) {
        if (!vm.detail.style) vm.detail.style = { position: {} }
        if (!vm.detail.style.position) vm.detail.style.position = {}
        vm.detail.style.position.w = current.width
        vm.detail.style.position.h = current.height
      })
      const data = node.getData() || {}
      if (data.detail) this.detail = data.detail
      this.editMode = data.editMode
      this.IsToolBox = data.IsToolBox
      vm.$EventBus.$on('cell-editMode', function (payload) {
        vm.editMode = payload.edit
        vm.IsToolBox = payload.toolbox
        vm.initComponents(vm.detail)
      })
    }
  }
}
