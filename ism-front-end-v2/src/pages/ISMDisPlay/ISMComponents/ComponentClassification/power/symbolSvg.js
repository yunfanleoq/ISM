/**
 * 把图元描述画成 SVG 字符串，或交给 Vue render 使用的节点列表。
 * 元素: [tag, attrs, extra?]
 * extra.when: 仅在 symbolState 等于该值时绘制
 * attrs.fill: 'body' 用填充色, 'stroke' 用线条色, 'none' 或不写为无填充
 */

function num(v, fallback) {
  const n = Number(v)
  return Number.isFinite(n) ? n : fallback
}

function resolveFill(fill, style) {
  if (fill === 'body') return style.fill
  if (fill === 'stroke') return style.stroke
  if (fill == null || fill === '') return 'none'
  return fill
}

function describe(symbol, style) {
  const state = num(style.state, symbol.state ? symbol.state.default : 1)
  const nodes = []
  const elements = symbol.elements || []
  for (let i = 0; i < elements.length; i++) {
    const el = elements[i]
    const tag = el[0]
    const raw = el[1] || {}
    if (raw.when != null && Number(raw.when) !== state) continue
    const attrs = {}
    const keys = Object.keys(raw)
    for (let k = 0; k < keys.length; k++) {
      if (keys[k] === 'when' || keys[k] === 'text' || keys[k] === 'size') continue
      attrs[keys[k]] = raw[keys[k]]
    }
    if (tag === 'text') {
      attrs.fill = style.stroke
      attrs.stroke = 'none'
      attrs['font-size'] = raw.size || 8
      attrs['font-family'] = 'sans-serif'
      attrs['text-anchor'] = raw.anchor || 'middle'
      attrs['dominant-baseline'] = 'middle'
      nodes.push({ tag: 'text', attrs: attrs, text: raw.text || '' })
      continue
    }
    attrs.fill = resolveFill(raw.fill, style)
    attrs.stroke = style.stroke
    attrs['stroke-width'] = raw.sw != null ? raw.sw : style.sw
    attrs['stroke-linecap'] = 'round'
    attrs['stroke-linejoin'] = 'round'
    nodes.push({ tag: tag, attrs: attrs })
  }
  return nodes
}

function escapeXml(s) {
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function toSvg(symbol, style) {
  const nodes = describe(symbol, style)
  let body = ''
  for (let i = 0; i < nodes.length; i++) {
    const n = nodes[i]
    const attrs = n.attrs
    let attrText = ''
    const keys = Object.keys(attrs)
    for (let k = 0; k < keys.length; k++) {
      attrText += ' ' + keys[k] + '="' + escapeXml(attrs[keys[k]]) + '"'
    }
    if (n.tag === 'text') {
      body += '<text' + attrText + '>' + escapeXml(n.text) + '</text>'
    } else {
      body += '<' + n.tag + attrText + '/>'
    }
  }
  return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="' + symbol.viewBox + '">' + body + '</svg>'
}

function toDataUri(symbol, style) {
  return 'data:image/svg+xml;charset=UTF-8,' + encodeURIComponent(toSvg(symbol, style))
}

const DEFAULT_STYLE = {
  stroke: '#5EEAD4',
  fill: '#0C2E36',
  sw: 1.6,
  state: 1
}

module.exports = {
  describe: describe,
  toSvg: toSvg,
  toDataUri: toDataUri,
  DEFAULT_STYLE: DEFAULT_STYLE
}
