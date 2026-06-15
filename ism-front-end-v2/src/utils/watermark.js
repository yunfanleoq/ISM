let watermark = {}

// 注入动画 keyframes（只注入一次）
var _animStylesInjected = false
function injectAnimationStyles() {
  if (_animStylesInjected) return
  var style = document.createElement('style')
  style.id = 'wm-anim-styles'
    style.textContent =
      '@keyframes wm-breathe {' +
      '  0%, 100% { opacity: 0.85; border-color: rgba(24, 144, 255, 0.10); box-shadow: 0 4px 12px rgba(20, 35, 52, 0.06); }' +
      '  50% { opacity: 1; border-color: rgba(24, 144, 255, 0.55); box-shadow: 0 6px 24px rgba(24, 144, 255, 0.18); }' +
      '}' +
      '@keyframes wm-text-glow {' +
      '  0%, 100% { text-shadow: 0 0 0 transparent, 0 1px 0 rgba(255,255,255,0.50); color: rgba(46,55,70,0.72); }' +
      '  50% { text-shadow: 0 0 14px rgba(24,144,255,0.35), 0 0 4px rgba(24,144,255,0.20), 0 1px 0 rgba(255,255,255,0.55); color: rgba(24,100,220,0.85); }' +
      '}' +
      '.wm-container {' +
      '  animation: wm-breathe 3s ease-in-out infinite;' +
      '}' +
      '.wm-text {' +
      '  animation: wm-text-glow 2.5s ease-in-out infinite;' +
      '}'
  document.head.appendChild(style)
  _animStylesInjected = true
}

let setWatermark = (str,date) => {
    let id = '1.23452384164.123412416'

    if (document.getElementById(id) !== null) {
        document.body.removeChild(document.getElementById(id))
    }

    injectAnimationStyles()

    let div = document.createElement('div')
    div.id = id
    div.className = 'wm-container'
    div.innerHTML = '<div class="wm-text">' + str + '</div><div class="wm-text">' + date + '</div>'
    div.style.pointerEvents = 'none'
    div.style.right = '18px'
    div.style.bottom = '18px'
    div.style.position = 'fixed'
    div.style.zIndex = '100000'
    div.style.padding = '8px 12px'
    div.style.border = '1px solid rgba(24,144,255,0.28)'
    div.style.borderRadius = '6px'
    div.style.background = 'rgba(255,255,255,0.75)'
    div.style.boxShadow = '0 6px 18px rgba(20,35,52,0.12)'
    div.style.backdropFilter = 'blur(6px)'
    div.style.color = 'rgba(46,55,70,0.72)'
    div.style.fontSize = '12px'
    div.style.fontWeight = '600'
    div.style.letterSpacing = '0.3px'
    div.style.lineHeight = '1.35'
    div.style.textAlign = 'right'
    div.style.userSelect = 'none'
    //   document.body.appendChild(div)
    div.style.opacity = '1' // 水印的透明度
    let show = document.getElementById("WaterMarkShow") // 控制水印显示的区域，设置id = show,表示在该范围内有水印
    if(show!=null)
    {
        show.appendChild(div)
    }

    return id
}

// 该方法只允许调用一次
watermark.set = (str,date) => {
    let id = setWatermark(str,date) // str,date代表的是两行水印。如果需好几个的话再追加。
    setInterval(() => {
        if (document.getElementById(id) === null) {
            id = setWatermark(str,date)
        }
    }, 500);
    window.onresize = () => {
        setWatermark(str,date)
    };
}

export default watermark
