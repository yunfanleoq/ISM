/**
 * 循安电力监控大屏 — 变电所 key → 组态 page_id（遗留映射，仅作兜底）
 * 主路径已改为 build_ncc_dashboard.py 的 uuid5 page_id_room；勿再优先使用本表。
 * 与 scripts/rebuild_xunan_dashboard.py 中 ROOM_PAGES 保持同步（该脚本已弃用）
 */
export const XUNAN_MODEL_ID = 'b8b4c094-faa9-a22a-1d0d-037539b27a6c'
export const XUNAN_PROJECT_UUID = '3ec5821f-b512-2adb-3e1c-473720d0a93e'

/** @type {Record<string, string>} 变电所编码 → display_model_layer.page_id */
export const XUNAN_SUBSTATION_PAGE_MAP = {
  '1A': 'a212682e-70a1-bcd4-c91f-49a5c6786f5f',
  '1B': '406a4f48-5c75-49b1-0019-0e6f80e584cd',
  '2A1': '07c7fa5b-f71a-093b-25d2-7e9f2e2eb658',
  '2A2': '8f1193e9-dbbd-f183-aaab-0ecb2373412c',
  '2A3': '4fb22468-e607-93da-5909-5d729b1ff731',
  '2A4': '5badaa89-bbb1-9159-4d9e-807e3c305d4d',
  '2B1': '5c2be637-be7d-b3ec-bb42-7a36f290cdf3',
  '2B2': '9387c9e9-c49b-4bb0-be5a-6c4c305c2b91',
  '2B3': '914a33f2-66ac-533e-4fc9-cff7e5447fad',
  '2B4': '8dc72a27-4d54-2f24-d257-447088b36d43',
  '3A1': 'b82cbaa7-8942-3b6f-424a-439d24702bc9',
  '3A2': 'ffa295da-576d-19e1-4b36-0845ea14faf1',
  '3A3': '86e5c739-59cb-f418-c337-59e4191439d5',
  '3A4': 'e39454ed-e4de-8dac-6263-e0bfd04ca8cb',
  '4A1': '2f468fd5-def4-8109-53e4-ca460b74a520',
}

const XUNAN_PAGE_KEYS_DESC = Object.keys(XUNAN_SUBSTATION_PAGE_MAP)
  .sort((a, b) => b.length - a.length)

export function resolveXunanSubstationPageId(modelId, zoneKey) {
  if (modelId !== XUNAN_MODEL_ID) return null
  if (!zoneKey || zoneKey === '_other') return null
  if (XUNAN_SUBSTATION_PAGE_MAP[zoneKey]) return XUNAN_SUBSTATION_PAGE_MAP[zoneKey]
  // 树节点可能是 1A1/2A11 等，按 ROOM_PAGES 前缀归并到 1A/2A1 页
  for (let i = 0; i < XUNAN_PAGE_KEYS_DESC.length; i++) {
    const prefix = XUNAN_PAGE_KEYS_DESC[i]
    if (zoneKey === prefix || zoneKey.startsWith(prefix)) {
      return XUNAN_SUBSTATION_PAGE_MAP[prefix]
    }
  }
  return null
}
