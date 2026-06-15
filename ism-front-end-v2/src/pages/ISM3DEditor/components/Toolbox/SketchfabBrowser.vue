<template>
  <div class="sketchfab-browser">
    <!-- API Token 配置栏 -->
    <div v-if="showTokenInput" class="sfb-token-bar">
      <a-input
        v-model="apiToken"
        placeholder="输入 Sketchfab API Token（可选，用于直接下载模型）"
        size="small"
        type="password"
        class="sfb-token-input"
      >
        <a-icon slot="prefix" type="lock" style="color:#555588" />
      </a-input>
      <a-button size="small" type="primary" @click="saveToken">保存</a-button>
      <a-button size="small" @click="showTokenInput = false">取消</a-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="!loaded && !searchMode" class="sfb-loading">
      <div class="sfb-loading-spinner"></div>
      <div class="sfb-loading-text">加载资产清单...</div>
    </div>

    <!-- 搜索 + 分类 -->
    <div v-show="loaded || searchMode" class="sfb-toolbar">
      <div class="sfb-search-row">
        <a-input
          v-model="search"
          placeholder="搜索本地资产或 Sketchfab 模型..."
          size="small"
          :allow-clear="true"
          class="sfb-search"
          @pressEnter="doSearch"
        >
          <a-icon slot="prefix" type="search" style="color:#555588" />
        </a-input>
        <a-button
          size="small"
          type="primary"
          class="sfb-search-btn"
          :loading="searching"
          @click="doSearch"
        >
          Sketchfab
        </a-button>
        <a-button
          v-if="!showTokenInput"
          size="small"
          class="sfb-token-btn"
          @click="showTokenInput = true"
          title="配置 Sketchfab API Token"
        >
          <i class="fas fa-key"></i>
        </a-button>
      </div>

      <!-- 搜索结果统计 -->
      <div v-if="searchMode" class="sfb-search-info">
        <span v-if="searching">搜索中...</span>
        <span v-else-if="searchResults.length > 0">找到 {{ searchResults.length }} 个 Sketchfab 模型</span>
        <span v-else-if="searched">未找到匹配模型</span>
        <a-button
          v-if="searchMode"
          size="small"
          type="link"
          @click="clearSearch"
        >返回本地资产</a-button>
      </div>

      <!-- 本地分类（仅在非搜索模式显示） -->
      <div v-if="!searchMode" class="sfb-categories">
        <button
          v-for="cat in categories"
          :key="cat.key"
          class="sfb-cat-btn"
          :class="{ active: currentCat === cat.key }"
          @click="currentCat = cat.key"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- 资产网格（本地） -->
    <div v-show="!searchMode && loaded" class="sfb-grid scrollbar">
      <div
        v-for="asset in filtered"
        :key="'local_' + asset.id"
        class="sfb-card"
        :class="{ disabled: !isLoadableAsset(asset) }"
        draggable="true"
        @dragstart="onDragStart($event, asset, 'local')"
        @click="onCardClick(asset, 'local')"
        :title="getAssetDisplayName(asset) + (asset.description ? ' - ' + asset.description : '')"
      >
        <div class="sfb-card-thumb">
          <img
            v-if="asset.thumbnail && !asset._thumbFailed"
            :src="asset.thumbnail"
            alt="thumb"
            @error="onThumbError(asset)"
          />
          <div v-else class="sfb-card-placeholder">
            <i class="fas fa-cube"></i>
          </div>
          <!-- 面数标签 -->
          <div v-if="asset.polygons" class="sfb-card-poly">
            {{ formatPoly(asset.polygons) }}面
          </div>
          <div v-if="!isLoadableAsset(asset)" class="sfb-card-status">
            未下载
          </div>
          <!-- 下载按钮：打开 Sketchfab 页面 -->
          <a
            v-if="isSketchfabUrl(asset.modelUrl)"
            class="sfb-card-dl"
            :href="asset.modelUrl"
            target="_blank"
            title="在 Sketchfab 中打开并下载"
            @click.stop
          >
            <i class="fas fa-external-link-alt"></i>
          </a>
        </div>
        <div class="sfb-card-label">{{ getAssetDisplayName(asset) }}</div>
      </div>

      <div v-if="!searchMode && loaded && filtered.length === 0" class="sfb-empty">
        无匹配资产
      </div>
    </div>

    <!-- 搜索结果网格（Sketchfab） -->
    <div v-show="searchMode" class="sfb-grid scrollbar">
      <div
        v-for="asset in searchResults"
        :key="'sf_' + asset.uid"
        class="sfb-card sfb-card-sketchfab"
        draggable="true"
        @dragstart="onDragStart($event, asset, 'sketchfab')"
        @click="onCardClick(asset, 'sketchfab')"
        @dblclick="onCardDblClick(asset, 'sketchfab')"
        :title="asset.name + (asset.author ? ' - by ' + asset.author.username : '')"
      >
        <div class="sfb-card-thumb">
          <img
            v-if="asset.thumbnail"
            :src="asset.thumbnail"
            alt="thumb"
          />
          <div v-else class="sfb-card-placeholder">
            <i class="fas fa-cube"></i>
          </div>
          <!-- 多边形数量 -->
          <div v-if="asset.polygonCount" class="sfb-card-poly">
            {{ formatPoly(asset.polygonCount) }}面
          </div>
          <!-- 动画标识 -->
          <div v-if="asset.animated" class="sfb-card-animated">
            <i class="fas fa-film"></i>
          </div>
        </div>
        <div class="sfb-card-label">{{ asset.name }}</div>
        <div v-if="asset.author" class="sfb-card-author">by {{ asset.author.username }}</div>
      </div>

      <div v-if="searchMode && !searching && searchResults.length === 0 && searched" class="sfb-empty">
        未找到 Sketchfab 模型，请尝试其他关键词
      </div>
    </div>
  </div>
</template>

<script>
import { getAssetList, getCategories, searchSketchfab, getSketchfabDownloadUrl } from '@/pages/ISM3DEditor/utils/SketchfabAssets'

const CATEGORY_LABELS = {
  available: '可用',
  all: '全部',
  sample: '示例',
  furniture: '家具',
  equipment: '设备',
  pipe: '管道',
  vehicle: '车辆',
  tank: '罐体'
}

const ASSET_NAME_LABELS = {
  Barrel_01: '红色油桶',
  'Barrel 02': '蓝色油桶',
  Barrel_02: '蓝色油桶',
  CashRegister_01: '收银机',
  CheeseBox_01: '木质货箱',
  'Drill 01': '电钻',
  Drill_01: '电钻',
  'Lantern 01': '煤油灯',
  Lantern_01: '煤油灯',
  'Megaphone 01': '扩音器',
  Megaphone_01: '扩音器',
  'Wet Floor Sign 01': '小心地滑警示牌',
  WetFloorSign_01: '小心地滑警示牌',
  'All Purpose Cleaner': '多用途清洁剂',
  all_purpose_cleaner: '多用途清洁剂',
  'Barrel 03': '蓝色金属桶',
  barrel_03: '蓝色金属桶',
  'Barrel Stove': '油桶炉',
  barrel_stove: '油桶炉',
  'Bench Vice 01': '台虎钳',
  bench_vice_01: '台虎钳',
  bleach_bottle: '漂白剂瓶',
  bolt_cutters_01: '断线钳',
  brass_blowtorch: '黄铜喷灯',
  caged_hanging_light: '防护罩吊灯',
  can_rusted: '生锈铁罐',
  cannon_01: '火炮',
  cardboard_box_01: '纸箱',
  ceiling_fan: '吊扇',
  cleaner_tin_01: '清洁剂铁罐',
  compost_bag_02: '编织袋',
  concrete_road_barrier: '混凝土路障',
  concrete_road_barrier_02: '混凝土路障 02',
  covered_car: '覆盖车辆',
  cross_pein_hammer: '横头锤',
  crowbar_01: '撬棍',
  drain_cleaner: '管道疏通器',
  drawer_cabinet: '抽屉柜',
  exterior_aircon_unit: '室外空调机',
  fire_hydrant: '消防栓',
  flathead_screwdriver: '一字螺丝刀',
  garden_sprinkler_01: '喷淋器',
  hand_plane_no4: '木工刨',
  handsaw_wood: '木工手锯',
  hanging_industrial_lamp: '工业吊灯',
  hatchet: '短柄斧',
  industrial_pipe_lamp: '工业管道灯',
  industrial_wall_lamp: '工业壁灯',
  korean_fire_extinguisher_01: '灭火器',
  ladder_sectioned_01: '分段梯子',
  large_castle_door: '大型金属门',
  leather_cleaner_can: '清洁剂罐',
  lightbulb_01: '灯泡',
  lightbulb_led: 'LED 灯泡',
  lubricant_spray: '润滑喷剂',
  measuring_tape_01: '卷尺',
  metal_detector: '金属探测器',
  metal_jerrycan: '金属油桶',
  metal_tool_chest: '金属工具柜',
  metal_toolbox: '金属工具箱',
  metal_trash_can: '金属垃圾桶',
  modular_airduct_circular_01: '圆形风管模块',
  modular_airduct_rectangular_01: '矩形风管模块',
  modular_chainlink_fence: '链式围栏模块',
  modular_electric_cables: '电缆模块',
  modular_electricity_poles: '电线杆模块',
  modular_factory_facade: '工厂外墙模块',
  modular_fire_escape: '消防梯模块',
  modular_industrial_pipes_01: '工业管道模块',
  modular_metal_gutter: '金属排水槽模块',
  modular_pipes: '管道模块',
  modular_pipes_plastic_01: '塑料管道模块',
  mounted_fluorescent_lights: '荧光灯组',
  multi_cleaner_5_litre: '5L 清洁剂桶',
  multi_cleaner_bottle: '清洁剂瓶',
  old_tyre: '旧轮胎',
  painted_wooden_cabinet: '木柜',
  painted_wooden_shelves: '木质货架',
  plastic_bottle_gallon: '加仑塑料瓶'
}

export default {
  name: 'SketchfabBrowser',
  props: {},
  data() {
    return {
      search: '',
      currentCat: 'available',
      categories: [],
      assets: [],
      loaded: false,
      // Sketchfab 搜索相关
      searchMode: false,
      searchResults: [],
      searching: false,
      searched: false,
      apiToken: localStorage.getItem('sketchfab_api_token') || '',
      showTokenInput: false,
      // 点击后因无 Token 中断的模型，保存 Token 后自动重试
      _pendingAsset: null,
      _skipNextSketchfabClick: false,
    }
  },
  computed: {
    filtered() {
      let list = this.assets.slice()
      if (this.currentCat === 'available') {
        list = list.filter(a => this.isLoadableAsset(a))
      } else if (this.currentCat && this.currentCat !== 'all') {
        list = list.filter(a => a.category === this.currentCat)
      }
      if (this.search && !this.searchMode) {
        const q = this.search.toLowerCase()
        list = list.filter(a => {
          const hay = [
            a.name || '',
            this.getAssetDisplayName(a),
            (a.tags || []).join(' '),
            a.description || ''
          ].join(' ')
          return hay.toLowerCase().includes(q)
        })
      }
      return list
    }
  },
  created() {
    this.load()
  },
  methods: {
    async load() {
      try {
        const [assets, cats] = await Promise.all([
          getAssetList(),
          getCategories(),
        ])
        this.assets = this.sortAssets(assets || [])
        this.categories = [{ key: 'available', label: '可用' }, { key: 'all', label: '全部' }]
          .concat((cats || []).map(cat => this.translateCategory(cat)))
      } catch (e) {
        console.warn('[SketchfabBrowser] 加载资产清单失败', e)
        this.assets = []
        this.categories = [{ key: 'available', label: '可用' }, { key: 'all', label: '全部' }]
      }
      this.loaded = true
    },
    sortAssets(list) {
      return list.slice().sort((a, b) => {
        const aLoadable = this.isLoadableAsset(a) ? 0 : 1
        const bLoadable = this.isLoadableAsset(b) ? 0 : 1
        return aLoadable - bLoadable
      })
    },
    getAssetModelPath(asset) {
      return asset ? (asset.localModelUrl || asset.fileUrl || asset.downloadedModelUrl || asset.modelUrl || '') : ''
    },
    isLoadableAsset(asset) {
      return this.isLoadableModelUrl(this.getAssetModelPath(asset))
    },
    toPayload(asset, source) {
      if (source === 'sketchfab') {
        return {
          type: 'sketchfab-model',
          uid: asset.uid || '',
          name: asset.name || 'Sketchfab 模型',
          thumbnail: asset.thumbnail || '',
          polygonCount: asset.polygonCount || 0,
          author: (asset.author && asset.author.username) || '',
          viewerUrl: asset.viewerUrl || '',
          isSketchfabModel: true,
          source: 'sketchfab',
        }
      }
      // 本地资产
      const modelPath = this.getAssetModelPath(asset)
      return {
        type: 'gltf',
        modelPath,
        thumbnail: asset.thumbnail || '',
        fitSize: asset.fitSize || 2,
        name: this.getAssetDisplayName(asset) || '外部模型',
        assetId: asset.id || '',
        isExternalModel: true,
        sourceAssetUrl: asset.sourceUrl || asset.sketchfabUrl || asset.modelUrl || '',
      }
    },
    translateCategory(cat) {
      const key = typeof cat === 'string' ? cat : (cat && cat.key)
      return {
        key,
        label: CATEGORY_LABELS[key] || (cat && cat.label) || key
      }
    },
    getAssetDisplayName(asset) {
      if (!asset) return ''
      const rawName = asset.name || ''
      const sourceKey = this.getAssetSourceKey(asset)
      const mapped = ASSET_NAME_LABELS[sourceKey] || ASSET_NAME_LABELS[rawName]
      if (mapped) return mapped
      return this.humanizeAssetName(rawName)
    },
    getAssetSourceKey(asset) {
      if (!asset) return ''
      const url = asset.sourceUrl || asset.modelUrl || ''
      const match = String(url).match(/\/a\/([^/?#]+)/)
      if (match && match[1]) return decodeURIComponent(match[1])
      const id = asset.id || ''
      if (id.indexOf('polyhaven-') === 0) return id.replace(/^polyhaven-/, '')
      return String(asset.name || '').replace(/（Poly Haven）|\(Poly Haven\)/g, '').trim()
    },
    humanizeAssetName(name) {
      return String(name || '')
        .replace(/（Poly Haven）|\(Poly Haven\)/g, '')
        .replace(/_/g, ' ')
        .trim()
    },
    isLoadableModelUrl(url) {
      return /\.(glb|gltf)(\?|$)/i.test(url || '')
    },
    async doSearch() {
      const q = this.search.trim()
      if (!q) {
        this.$message && this.$message.warning('请输入搜索关键词')
        return
      }
      this.searching = true
      this.searched = false
      this.searchMode = true
      try {
        const results = await searchSketchfab(q, {
          limit: 24,
          token: this.apiToken || undefined
        })
        this.searchResults = results || []
      } catch (e) {
        console.warn('[SketchfabBrowser] 搜索失败', e)
        this.searchResults = []
      }
      this.searching = false
      this.searched = true
    },
    clearSearch() {
      this.searchMode = false
      this.searchResults = []
      this.searched = false
      this.search = ''
    },
    saveToken() {
      const hadToken = !!(this.apiToken && this.apiToken.trim())
      if (hadToken) {
        localStorage.setItem('sketchfab_api_token', this.apiToken.trim())
        this.$message && this.$message.success('API Token 已保存（本次会话有效）')
      } else {
        localStorage.removeItem('sketchfab_api_token')
        this.$message && this.$message.info('已清除 API Token')
      }
      this.showTokenInput = false
      // 如果有因无 Token 中断的待下载模型，自动重试
      if (hadToken && this._pendingAsset) {
        const asset = this._pendingAsset
        this._pendingAsset = null
        this.$nextTick(() => this.loadSketchfabModel(asset))
      }
    },
    onDragStart(e, asset, source) {
      const payload = this.toPayload(asset, source)
      if (source === 'sketchfab') {
        if (e && e.preventDefault) e.preventDefault()
        this.$message && this.$message.info('Sketchfab 模型：请先点击卡片加载到场景')
        return
      }
      // 本地资产
      const modelPath = payload.modelPath
      if (!this.isLoadableModelUrl(modelPath)) {
        if (e && e.preventDefault) e.preventDefault()
        this.$message && this.$message.warning('该资产还未下载为 .glb/.gltf 文件，请先点击外部链接到 Sketchfab 下载。', 3)
        return
      }
      if (e && e.dataTransfer) {
        e.dataTransfer.effectAllowed = 'copy'
        e.dataTransfer.setData('text/plain', 'gltf')
        try {
          e.dataTransfer.setData('application/x-ism3d-object', JSON.stringify(payload))
        } catch (ex) { /* ignore */ }
      }
    },
    onCardClick(asset, source) {
      if (source === 'sketchfab') {
        if (this._skipNextSketchfabClick) {
          this._skipNextSketchfabClick = false
          return
        }
        this.loadSketchfabModel(asset)
        return
      }
      // 本地资产
      const payload = this.toPayload(asset, source)
      if (!this.isLoadableModelUrl(payload.modelPath)) {
        this.$message && this.$message.warning('该资产还未下载为 .glb/.gltf 文件，请先点击外部链接到 Sketchfab 下载。', 3)
        return
      }
      this.$emit('add-object', payload)
    },
    onCardDblClick(asset, source) {
      if (source !== 'sketchfab') {
        this.onCardClick(asset, source)
        return
      }
      this._skipNextSketchfabClick = true
      const targetUrl = asset && (asset.viewerUrl || asset.modelUrl || asset.sourceUrl)
      if (!targetUrl) {
        this.$message && this.$message.warning('未找到 Sketchfab 页面链接')
        return
      }
      const link = document.createElement('a')
      link.href = targetUrl
      link.target = '_blank'
      link.rel = 'noopener noreferrer'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    },
    async loadSketchfabModel(asset) {
      if (!asset || !asset.uid) {
        this.$message && this.$message.warning('无效的 Sketchfab 模型')
        return
      }
      const token = this.apiToken || localStorage.getItem('sketchfab_api_token') || ''
      if (!token) {
        // 无 Token：记录待下载模型，展开 Token 输入栏并提示，保存 Token 后自动重试
        this._pendingAsset = asset
        this.showTokenInput = true
        this.$message && this.$message.warning('请先配置 Sketchfab API Token，点击下方"密钥"图标粘贴 Token 后即可自动下载', 4)
        return
      }
      // 有 token：清空待下载记录，尝试获取下载链接并加载
      this._pendingAsset = null
      this.$message && this.$message.loading('正在获取模型下载链接...', 0)
      try {
        const dlUrl = await getSketchfabDownloadUrl(asset.uid, token)
        this.$message.destroy()
        if (!dlUrl) {
          this.$message && this.$message.error('无法获取下载链接')
          return
        }
        // 直接以 gltf 类型 emit，复用现有加载流程（loadGLTFModelUrl 支持全 URL）
        const payload = {
          type: 'gltf',
          modelPath: dlUrl,
          name: asset.name || 'Sketchfab 模型',
          thumbnail: asset.thumbnail || '',
          isExternalModel: true,
          isSketchfabModel: true,
          fitSize: 2,
        }
        this.$emit('add-object', payload)
        this.$message && this.$message.success('模型已添加到场景！')
      } catch (err) {
        this.$message.destroy()
        console.warn('[SketchfabBrowser] 加载模型失败', err)
        this.$message && this.$message.error('模型加载失败：' + (err.message || err))
      }
    },
    onThumbError(asset) {
      asset._thumbFailed = true
    },
    formatPoly(n) {
      if (!n || isNaN(Number(n))) return ''
      n = Number(n)
      if (n >= 10000) return (n / 10000).toFixed(1) + '万'
      return String(n)
    },
    isSketchfabUrl(url) {
      return url && /^https?:\/\/(www\.)?sketchfab\.com\//i.test(url)
    },
  },
}
</script>

<style scoped>
.sketchfab-browser {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  background: #fafafa;
}
/* ---- Token 栏 ---- */
.sfb-token-bar {
  display: flex;
  gap: 4px;
  padding: 6px 8px;
  border-bottom: 1px solid #eee;
  background: #fffbe6;
  align-items: center;
}
.sfb-token-input {
  flex: 1;
}
/* ---- 工具栏 ---- */
.sfb-toolbar {
  padding: 6px 8px 5px;
  border-bottom: 1px solid #eee;
  flex-shrink: 0;
  background: #fff;
}
.sfb-search-row {
  display: flex;
  gap: 4px;
  margin-bottom: 5px;
}
.sfb-search {
  flex: 1;
}
.sfb-search >>> .ant-input {
  height: 24px;
  font-size: 12px;
}
.sfb-search >>> .ant-input:focus,
.sfb-search >>> .ant-input:hover {
  border-color: #13c2c2;
  box-shadow: 0 0 0 2px rgba(19,194,194,0.15);
}
.sfb-search-btn {
  flex-shrink: 0;
}
.sfb-token-btn {
  flex-shrink: 0;
  width: 32px;
  padding: 0;
}
/* ---- 搜索信息 ---- */
.sfb-search-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  color: #666;
  margin-bottom: 4px;
  flex-wrap: wrap;
}
/* ---- 分类胶囊按钮 ---- */
.sfb-categories {
  display: flex;
  flex-wrap: wrap;
  gap: 3px;
}
.sfb-cat-btn {
  border: 1px solid #d9d9d9;
  border-radius: 10px;
  background: #fff;
  font-size: 11px;
  line-height: 16px;
  min-width: 34px;
  padding: 1px 8px;
  cursor: pointer;
  color: #666;
  outline: none;
  transition: all 0.2s;
  white-space: nowrap;
}
.sfb-cat-btn.active {
  border-color: #13c2c2;
  background: #13c2c2;
  color: #fff;
}
.sfb-cat-btn:hover:not(.active) {
  border-color: #13c2c2;
  color: #13c2c2;
  background: #e6fffb;
}
/* ---- 资产网格 ---- */
.sfb-grid {
  flex: 1;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(78px, 1fr));
  align-content: flex-start;
  padding: 6px;
  gap: 6px;
}
/* ---- 卡片 ---- */
.sfb-card {
  min-width: 0;
  cursor: pointer;
  border: 1px solid #eee;
  border-radius: 6px;
  padding: 3px;
  background: #fff;
  transition: all 0.2s;
  position: relative;
}
.sfb-card:hover {
  border-color: #13c2c2;
  box-shadow: 0 2px 8px rgba(19,194,194,0.18);
  transform: translateY(-2px);
}
.sfb-card.disabled {
  cursor: not-allowed;
  opacity: 0.62;
}
.sfb-card.disabled:hover {
  border-color: #eee;
  box-shadow: none;
  transform: none;
}
.sfb-card-thumb {
  width: 100%;
  height: 74px;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}
.sfb-card-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.sfb-card-placeholder {
  color: #ccc;
  font-size: 28px;
}
/* 面数标签 */
.sfb-card-poly {
  position: absolute;
  bottom: 2px;
  right: 2px;
  background: rgba(0,0,0,0.55);
  color: #fff;
  font-size: 9px;
  padding: 1px 4px;
  border-radius: 3px;
  pointer-events: none;
}
.sfb-card-status {
  position: absolute;
  left: 2px;
  bottom: 2px;
  background: rgba(120,120,120,0.85);
  color: #fff;
  font-size: 9px;
  padding: 1px 4px;
  border-radius: 3px;
  pointer-events: none;
}
.sfb-card-animated {
  position: absolute;
  top: 2px;
  left: 2px;
  background: rgba(19,194,194,0.85);
  color: #fff;
  font-size: 9px;
  padding: 1px 4px;
  border-radius: 3px;
  pointer-events: none;
}
.sfb-card-label {
  margin-top: 3px;
  font-size: 11px;
  color: #555;
  text-align: center;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  line-height: 14px;
  height: 14px;
}
.sfb-card-author {
  font-size: 10px;
  color: #999;
  text-align: center;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
/* 下载按钮：右上角悬浮 */
.sfb-card-dl {
  position: absolute;
  top: 3px;
  right: 3px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: rgba(0,0,0,0.45);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  text-decoration: none;
  transition: background 0.15s;
  z-index: 2;
}
.sfb-card-dl:hover {
  background: rgba(19,194,194,0.85);
}
/* ---- 空状态 ---- */
.sfb-empty {
  width: 100%;
  text-align: center;
  color: #bbb;
  padding: 32px 0;
  font-size: 12px;
}
/* ---- 加载状态 ---- */
.sfb-loading {
  width: 100%;
  text-align: center;
  padding: 24px 0;
}
.sfb-loading-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2px solid #e8e8e8;
  border-top-color: #13c2c2;
  border-radius: 50%;
  animation: sfb-spin 0.6s linear infinite;
  margin-bottom: 8px;
}
@keyframes sfb-spin {
  to { transform: rotate(360deg); }
}
.sfb-loading-text {
  font-size: 12px;
  color: #999;
}
/* 滚动条 */
.scrollbar::-webkit-scrollbar {
  width: 5px;
}
.scrollbar::-webkit-scrollbar-thumb {
  border-radius: 10px;
  background-color: #91d5ff;
}
.scrollbar::-webkit-scrollbar-track {
  background: #ededed;
  border-radius: 10px;
}
</style>
