#!/usr/bin/env node
/**
 * 将 §8 菜单树写入 ism_server_user/conf/*MenuConfig.json（root.children 内插入循智段）
 */
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const repoRoot = path.resolve(__dirname, '..')

// 与 ism-front-end-v2/src/xunzhiJf/menuFlat.js 保持同步（Node 侧轻量副本键名）
const SECTION = 'XunzhiJfSection'

function menuCodeToRouterKey(menuCode) {
  return 'XunzhiJf_' + menuCode.replace(/\./g, '_')
}

// 动态 import ESM from frontend - use read+eval of flat table via regex on menuFlat.js
const menuFlatPath = path.join(repoRoot, 'ism-front-end-v2/src/xunzhiJf/menuFlat.js')
const src = fs.readFileSync(menuFlatPath, 'utf8')
const rows = []
const re = /menu_code:\s*'([^']+)'[^}]*title:\s*'([^']+)'[^}]*parent_code:\s*'([^']+)'/g
let m
const blocks = src.match(/\{ sort:[^}]+\}/g) || []
for (const block of blocks) {
  const menu_code = /menu_code:\s*'([^']+)'/.exec(block)?.[1]
  const parent_code = /parent_code:\s*'([^']+)'/.exec(block)?.[1]
  if (menu_code && parent_code) {
    rows.push({ menu_code, parent_code })
  }
}

function buildMenuChildren(parentCode) {
  return rows
    .filter(r => r.parent_code === parentCode)
    .map(r => {
      const node = { router: menuCodeToRouterKey(r.menu_code) }
      const kids = buildMenuChildren(r.menu_code)
      if (kids.length) node.children = kids
      return node
    })
}

const xunzhiBlock = {
  router: SECTION,
  name: '循智·机房环境',
  icon: 'cloud-server',
  children: buildMenuChildren('(root)'),
}

const configFiles = [
  'ism_server_user/conf/MenuConfig.json',
  'ism_server_user/conf/router/OptMenuConfig.json',
  'ism_server_user/conf/router/UserMenuConfig.json',
]

for (const rel of configFiles) {
  const filePath = path.join(repoRoot, rel)
  const data = JSON.parse(fs.readFileSync(filePath, 'utf8'))
  const root = data.find(item => item.router === 'root')
  if (!root || !root.children) {
    console.error('skip', rel, 'no root.children')
    continue
  }
  const idx = root.children.findIndex(c => c.router === SECTION)
  if (idx >= 0) {
    root.children[idx] = xunzhiBlock
  } else {
    const dashIdx = root.children.findIndex(c => c.router === 'dashboard')
    const insertAt = dashIdx >= 0 ? dashIdx + 1 : 0
    root.children.splice(insertAt, 0, xunzhiBlock)
  }
  fs.writeFileSync(filePath, JSON.stringify(data, null, 4) + '\n')
  console.log('updated', rel)
}
