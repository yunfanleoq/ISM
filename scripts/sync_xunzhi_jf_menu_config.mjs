#!/usr/bin/env node
/**
 * 将 §8 菜单树 merge 进 ism_server_user/conf/*MenuConfig.json
 *
 * 行为：仅在 root.children 内 upsert 一条 `XunzhiJfSection`（插入在 dashboard 之后，或替换同名节点）。
 * **禁止**整文件替换、禁止清空原 ISM 侧栏（设备管理/数据模型等须保留）。
 */
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const repoRoot = path.resolve(__dirname, '..')

const SECTION = 'XunzhiJfSection'

/** Admin 菜单锚点（User/Opt 菜单树不同，不做同一套锚点校验） */
const ANCHOR_ROUTERS_ADMIN = [
  'dashboard',
  'DataWarehouse',
  'DeviceLibraryConfig',
  'DeviceModel',
  'SCADAMonitor',
  'RealTimeAlarm',
]

function menuCodeToRouterKey(menuCode) {
  return 'XunzhiJf_' + menuCode.replace(/\./g, '_')
}

const menuFlatPath = path.join(repoRoot, 'ism-front-end-v2/src/xunzhiJf/menuFlat.js')
const src = fs.readFileSync(menuFlatPath, 'utf8')
const rows = []
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

function topLevelRouters(children) {
  return new Set((children || []).map(c => c.router).filter(Boolean))
}

function assertMergeSafe(beforeChildren, afterChildren, rel) {
  const before = topLevelRouters(beforeChildren)
  const after = topLevelRouters(afterChildren)
  for (const r of before) {
    if (r === SECTION) continue
    if (!after.has(r)) {
      throw new Error(`${rel}: merge 后丢失原菜单 router="${r}"，已中止写入`)
    }
  }
  if (afterChildren.length < beforeChildren.length) {
    throw new Error(
      `${rel}: root.children 数量 ${beforeChildren.length} → ${afterChildren.length}，已中止写入`
    )
  }
  if (rel.endsWith('conf/MenuConfig.json') || rel.endsWith('OptMenuConfig.json')) {
    for (const anchor of ANCHOR_ROUTERS_ADMIN) {
      if (!after.has(anchor)) {
        throw new Error(`${rel}: 缺少锚点 router="${anchor}"，已中止写入`)
      }
    }
  }
}

const configFiles = [
  'ism_server_user/conf/MenuConfig.json',
  'ism_server_user/conf/router/OptMenuConfig.json',
  'ism_server_user/conf/router/UserMenuConfig.json',
]

for (const rel of configFiles) {
  const filePath = path.join(repoRoot, rel)
  const raw = fs.readFileSync(filePath, 'utf8')
  const data = JSON.parse(raw)
  const root = data.find(item => item.router === 'root')
  if (!root || !root.children) {
    console.error('skip', rel, 'no root.children')
    continue
  }
  const beforeChildren = JSON.parse(JSON.stringify(root.children))
  const idx = root.children.findIndex(c => c.router === SECTION)
  if (idx >= 0) {
    root.children[idx] = xunzhiBlock
  } else {
    const dashIdx = root.children.findIndex(c => c.router === 'dashboard')
    const insertAt = dashIdx >= 0 ? dashIdx + 1 : 0
    root.children.splice(insertAt, 0, xunzhiBlock)
  }
  assertMergeSafe(beforeChildren, root.children, rel)
  fs.writeFileSync(filePath, JSON.stringify(data, null, 4) + '\n')
  console.log(
    `updated ${rel} (upsert ${SECTION}, children ${beforeChildren.length} → ${root.children.length})`
  )
}
