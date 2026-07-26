/**
 * 虚拟列头柜前缀分组纯函数自测（不依赖 webpack / navContext）
 * 用法: node scripts/test_virtual_cabinet.mjs
 */

function extractPointPrefix(name) {
  const n = String(name || '').trim()
  const i = n.indexOf('_')
  if (i <= 0) return ''
  return n.slice(0, i)
}

function compareCabinetPrefix(a, b) {
  const rank = (p) => {
    const m = String(p).match(/^([A-Za-z])列(头|尾)$/)
    if (!m) return { series: 2, letter: 99, raw: p }
    const letter = m[1].toUpperCase().charCodeAt(0) - 65
    const series = m[2] === '头' ? 0 : 1
    return { series, letter, raw: p }
  }
  const ra = rank(a)
  const rb = rank(b)
  if (ra.series !== rb.series) return ra.series - rb.series
  if (ra.letter !== rb.letter) return ra.letter - rb.letter
  return String(ra.raw).localeCompare(String(rb.raw), 'zh')
}

function groupDatapointsByPrefix(points) {
  const map = Object.create(null)
  ;(points || []).forEach((p) => {
    const prefix = extractPointPrefix(p && (p.name || p.Name || p.label))
    if (!prefix) return
    if (!map[prefix]) map[prefix] = { prefix, count: 0 }
    map[prefix].count += 1
  })
  return Object.keys(map).sort(compareCabinetPrefix).map(k => map[k])
}

const pts = []
for (const s of ['头', '尾']) {
  for (const L of 'ABCDEFGHIJ') {
    pts.push({ name: `${L}列${s}_主路A相电压` })
    pts.push({ name: `${L}列${s}_支路电流` })
  }
}
const groups = groupDatapointsByPrefix(pts)
const names = groups.map(g => g.prefix)
const expected = [
  ...'ABCDEFGHIJ'.split('').map(L => `${L}列头`),
  ...'ABCDEFGHIJ'.split('').map(L => `${L}列尾`),
]
let ok = names.length === 20 && names.join(',') === expected.join(',')
if (!ok) {
  console.error('FAIL expected', expected.join(','))
  console.error('GOT     ', names.join(','))
  process.exit(1)
}
console.log('PASS virtual cabinet prefixes x20:', names.join(', '))
