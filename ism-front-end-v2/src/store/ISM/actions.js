import {getDisplayModelLayerData,setDisplayModelLayerData,getLayerDataStructByToken} from "@/services/displayModel";
import {GetDisplayLoginPage} from "@/services/system";
import { uuid } from 'vue-uuid';
import {getDisplayModelPagerLayerData} from "../../services/displayModel";
import {normalizeISMScene} from "@/pages/ISMDisPlay/utils/ismSceneNormalizer";
import {ismDebug} from "@/utils/ismDebug";
import {resolvePageComponentsAsync, ensureTemplatePageLayer} from "@/pages/ISMDisPlay/utils/navContextBinding";
import {applyDeviceListPagination, isDeviceListNav, resolveDeviceListTemplateId} from "@/pages/ISMDisPlay/utils/deviceListPager";
import {resolveDeviceSignalTemplateId} from "@/pages/ISMDisPlay/utils/deviceSignalTemplate";
import {resolveOldPageTarget} from "@/pages/ISMDisPlay/utils/navTreeIndex";
import {applyDeviceDetailPagination} from "@/pages/ISMDisPlay/utils/deviceDetailPager";
import {sanitizeGraphComponents} from "@/pages/ISMDisPlay/utils/graphCellSanitizer";
import {applyDatapointPagination} from "@/pages/ISMDisPlay/utils/navContext";

/**
 * 若存在 navContext，对页 components 做相对绑点解析（不写回缓存原页，避免污染模板）。
 */
async function applyNavContextToPageConfig(ctx, tempConfigData) {
    if (!tempConfigData) {
        return tempConfigData
    }
    ensureTemplatePageLayer(tempConfigData)
    if (!tempConfigData.components) {
        tempConfigData.components = { cells: [] }
    }
    const finish = (cfg) => ({
        ...cfg,
        components: sanitizeGraphComponents(cfg.components, { tag: 'applyNavContextToPageConfig' }),
    })
    let nav = ctx.state.navContext
    if (!nav) {
        return finish(tempConfigData)
    }
    try {
        if (isDeviceListNav(nav)) {
            nav = applyDeviceListPagination(nav)
        } else if (nav.signalMode || nav.routeMode === 'signal') {
            // 信号层测点表：不要走设备详情行分页，否则会冲掉 datapoint 分页元数据
            nav = applyDatapointPagination(nav)
        } else if (nav.kind === 'device' || nav.kind === 'registerGroup') {
            const peek = applyDeviceDetailPagination(
                JSON.parse(JSON.stringify(tempConfigData.components.cells)),
                { ...nav, detailPageIndex: nav.detailPageIndex || 0 },
            )
            nav = peek.nav
        }
        try {
            ctx.commit('setNavContext', nav)
        } catch (e) { /* ignore */ }
        const resolved = await resolvePageComponentsAsync(tempConfigData.components, nav, {
            treeIndex: ctx.state.navTreeIndex,
            templateMap: ctx.state.navTemplateMap,
        })
        // resolve 内可能再次分页，把最新 total/pageIndex 写回 store
        try {
            const latest = ctx.state.navContext
            if (latest && (latest.signalMode || latest.routeMode === 'signal')) {
                ctx.commit('setNavContext', applyDatapointPagination({
                    ...latest,
                    datapointPageIndex: (nav && nav.datapointPageIndex) || latest.datapointPageIndex || 0,
                }))
            }
        } catch (e) { /* ignore */ }
        return finish({ ...tempConfigData, components: resolved })
    } catch (e) {
        console.warn('[applyNavContextToPageConfig]', e && e.message)
        return finish(tempConfigData)
    }
}

/**
 * 安全收集 active 条件绑定，避免 condition 缺失时抛错打断大屏加载。
 */
function collectActiveBindings(activeList, bangDingData, bangDingDeviceSN) {
    if (!activeList || !activeList.length) {
        return
    }
    for (let kv = 0; kv < activeList.length; kv++) {
        const condition = activeList[kv] && activeList[kv].condition
        if (!condition) {
            continue
        }
        if (condition.dataID) {
            bangDingData.push(condition.dataID)
        }
        if (condition.deviceSN) {
            bangDingDeviceSN.push(condition.deviceSN)
        }
    }
}

function pascalToKebab(str) {
    if (!str) return ''
    // DvBorderBox1 → dv-border-box1, DeviceTree → device-tree
    return str
        .replace(/([A-Z])([A-Z][a-z])/g, '$1-$2')
        .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
        .replace(/([A-Z]+)([A-Z][a-z])/g, '$1-$2')
        .toLowerCase()
}

/**
 * 安全收集 animate 绑定。旧组态数据常有 animate 但缺 condition，
 * 直接读 animate.condition.dataID 会抛错，导致大屏 callback 中断、pageLoading 永不关闭。
 */
function collectAnimateBindings(animate, bangDingData, bangDingDeviceSN) {
    if (!animate || typeof animate !== 'object') {
        return
    }
    const condition = animate.condition
    if (condition && condition.dataID) {
        if (condition.deviceSN) {
            bangDingDeviceSN.push(condition.deviceSN)
        }
        bangDingData.push(condition.dataID)
    }
    const move = animate.move
    if (!move || typeof move !== 'object') {
        return
    }
    if (move.x) {
        if (move.x.deviceSN) {
            bangDingDeviceSN.push(move.x.deviceSN)
        }
        if (move.x.dataID) {
            bangDingData.push(move.x.dataID)
        }
    }
    if (move.y) {
        if (move.y.deviceSN) {
            bangDingDeviceSN.push(move.y.deviceSN)
        }
        if (move.y.dataID) {
            bangDingData.push(move.y.dataID)
        }
    }
}

/** 解析 API 返回的 layer/components 原始字段（metaOnly 模板页 layer 可能为空字符串） */
function parseRawPageLayerFields(pageData) {
    if (typeof pageData.layer === 'string') {
        pageData.layer = pageData.layer !== '' ? JSON.parse(pageData.layer) : {}
    } else if (!pageData.layer || typeof pageData.layer !== 'object') {
        pageData.layer = {}
    }
    if (pageData.components === '' || pageData.components == null) {
        pageData.components = { cells: [] }
    } else if (typeof pageData.components === 'string') {
        pageData.components = JSON.parse(pageData.components)
    }
    ensureTemplatePageLayer(pageData)
    return pageData
}

function normalizePageConfigData(pageData) {
    // 转换扁平组件格式 → ISM cells 格式 (带 shape + data.detail.style)
    let componentsInput = { cells: [] }
    if (Array.isArray(pageData.components)) {
        const cells = pageData.components.map(item => {
            const pos = (item.style && item.style.position) || {}
            const kebabType = pascalToKebab(item.type)
            return {
                shape: kebabType,
                id: item.identifier || uuid.v1(),
                x: pos.x || 0,
                y: pos.y || 0,
                width: pos.w || 100,
                height: pos.h || 40,
                zIndex: item.style && item.style.zIndex,
                visible: item.style && item.style.visible !== 0,
                position: { x: pos.x || 0, y: pos.y || 0 },
                size: { width: pos.w || 100, height: pos.h || 40 },
                data: {
                    detail: {
                        type: kebabType,
                        identifier: item.identifier,
                        name: item.name,
                        style: item.style || {},
                        animate: item.animate,
                        action: item.action,
                        active: item.active,
                        dataBind: item.dataBind
                    }
                }
            }
        })
        componentsInput = { cells }
    } else if (pageData.components && pageData.components.cells) {
        componentsInput = pageData.components
    }

    // metaOnly 占位页：无 cells 时跳过重型 normalize，减少首屏 CPU
    if (!componentsInput.cells || componentsInput.cells.length === 0) {
        if (typeof pageData.layer === 'string') {
            try {
                pageData.layer = pageData.layer ? JSON.parse(pageData.layer) : {}
            } catch (e) {
                pageData.layer = {}
            }
        }
        pageData.layer = pageData.layer || {}
        pageData.components = { cells: [] }
        return pageData
    }

    const normalized = normalizeISMScene({
        layer: pageData.layer,
        components: componentsInput
    })
    pageData.layer = {
        ...pageData.layer,
        ...normalized.layer
    }
    pageData.components = normalized.components
    return pageData
}

function isDisplayPagesLoaded(ctx, displayUUID) {
    const pageList = [...ctx.state.PCPageList, ...ctx.state.PhonePageList]
    return pageList.some(page => page.pageModelUuid == displayUUID)
}

function buildPageTreeChildren(tempConfigData) {
    const children = []
    const cells = tempConfigData && tempConfigData.components && tempConfigData.components.cells
    if (!cells || !cells.length) {
        return children
    }
    for (let k = 0; k < cells.length; k++) {
        if (typeof cells[k].data === 'undefined' || typeof cells[k].data.detail === 'undefined') {
            continue
        }
        children.push({
            isComponents: true,
            title: cells[k].data.detail.name,
            key: cells[k].id,
            cellid: cells[k].data.detail.identifier
        })
    }
    return children
}

function markMetaOnlyLazyState(pageInfo, tempConfigData, metaOnly, isHome) {
    const hasCells = tempConfigData.components && tempConfigData.components.cells && tempConfigData.components.cells.length > 0
    const isLazyPlaceholder = metaOnly && !hasCells && isHome != 1
    pageInfo._lazyLoaded = !isLazyPlaceholder
    if (isLazyPlaceholder) {
        tempConfigData.components = {cells: []}
    } else if (!metaOnly) {
        pageInfo._lazyLoaded = true
    }
    return tempConfigData
}

function mergePageLayerFromExisting(pageInfo, oldPc, oldPhone) {
    const lists = [oldPc || [], oldPhone || []]
    for (let li = 0; li < lists.length; li++) {
        const oldList = lists[li]
        for (let j = 0; j < oldList.length; j++) {
            if (oldList[j].pageUuid === pageInfo.pageUuid && oldList[j]._lazyLoaded && oldList[j].pageLayerData) {
                pageInfo.pageLayerData = oldList[j].pageLayerData
                pageInfo._lazyLoaded = true
                pageInfo.children = buildPageTreeChildren(oldList[j].pageLayerData)
                return
            }
        }
    }
}

function findPageInLists(ctx, pageid) {
    const lists = [ctx.state.PCPageList || [], ctx.state.PhonePageList || []]
    for (let i = 0; i < lists.length; i++) {
        const hit = lists[i].find(p => p && p.pageUuid === pageid)
        if (hit) return hit
    }
    return null
}

function upsertPCPageInfo(ctx, pageInfo) {
    const list = ctx.state.PCPageList ? ctx.state.PCPageList.slice() : []
    const idx = list.findIndex(p => p && p.pageUuid === pageInfo.pageUuid)
    if (idx >= 0) {
        list[idx] = { ...list[idx], ...pageInfo }
    } else {
        list.push(pageInfo)
    }
    ctx.state.PCPageList = list
    return pageInfo
}

function buildPageInfoFromConfig(pageid, displayUUID, cfg) {
    return {
        id: cfg.ID || 0,
        key: 0,
        isEdit: false,
        pageUuid: pageid,
        pageModelUuid: displayUUID || cfg.modelId || '',
        isNewItem: false,
        title: cfg.PageName || cfg.name || pageid,
        depth: 1,
        pageType: cfg.PageType != null ? cfg.PageType : 1,
        scopedSlots: { title: 'custom' },
        isComponents: false,
        IsHome: cfg.IsHome || 0,
        IsLogin: cfg.IsLogin || 0,
        templateKind: cfg.templateKind || cfg.TemplateKind || '',
        templateModelUuid: cfg.templateModelUuid || cfg.TemplateModelUuid || '',
        children: buildPageTreeChildren(cfg),
        _lazyLoaded: true,
        pageLayerData: cfg,
    }
}

/** PCPageList 缺页或整表为空时，按需拉取并注入条目（导航模板页 fallback 依赖此路径） */
function ensurePageRegistered(ctx, pageid, displayUUID) {
    if (!pageid) return Promise.resolve(null)
    const existing = findPageInLists(ctx, pageid)
    if (existing) return Promise.resolve(existing)
    return loadSinglePageLayer(pageid).then(function (cfg) {
        if (!cfg) return null
        const pageInfo = buildPageInfoFromConfig(pageid, displayUUID, cfg)
        if (pageInfo.pageType === 1) {
            upsertPCPageInfo(ctx, pageInfo)
        }
        return pageInfo
    })
}

// 按需加载: 运行态(AppRun)与编辑器首屏只拉取页面元数据 + 首页 components，
// 其余页面在下钻/点选时按 page_id 单独拉取(见 loadSinglePageLayer / selectDisplayPageDataStruct)。
// 禁止 metaOnly 失败后回退全量加载（会导致 30s 超时与内存打满）。
const pendingPageLayerLoads = new Map()

function loadSinglePageLayer(pageid) {
    if (pendingPageLayerLoads.has(pageid)) {
        return pendingPageLayerLoads.get(pageid)
    }
    const promise = getDisplayModelPagerLayerData({pageid: pageid}).then(function (res) {
        if (!res || !res.data || res.data.code != 0) {
            return null
        }
        let cfg = res.data.layer
        try {
            cfg = parseRawPageLayerFields(cfg)
            cfg = normalizePageConfigData(cfg)
        } catch (e) {
            console.error('[loadSinglePageLayer] parse error:', e && e.message)
            return null
        }
        return cfg
    }).catch(function (e) {
        console.error('[loadSinglePageLayer] request error:', e && e.message)
        return null
    }).finally(function () {
        pendingPageLayerLoads.delete(pageid)
    })
    pendingPageLayerLoads.set(pageid, promise)
    return promise
}

// ---------- 大屏按需加载 + 空闲 LRU 预取 ----------
// 原则：用户展开/下钻时才同步拉取；浏览器空闲时按链接关系预测下一页并静默预取。
// 预取失败不重试，避免打爆后端；受限网络下完全禁用，正常网络也最多预取 1 页。
const PREFETCH_MAX_CANDIDATES = 1
const PREFETCH_TRIGGER_DELAY_MS = 800
let prefetchQueue = []
let prefetchTimer = null
let prefetchRunning = false

function cancelIdlePrefetch() {
    if (prefetchTimer) {
        clearTimeout(prefetchTimer)
        prefetchTimer = null
    }
    prefetchQueue = []
    prefetchRunning = false
}

function scheduleIdleTask(fn) {
    if (typeof requestIdleCallback === 'function') {
        requestIdleCallback(fn, {timeout: 2500})
    } else {
        setTimeout(fn, 150)
    }
}

function collectPageUUIDs(obj, out, depth) {
    if (!obj || depth > 10) {
        return
    }
    if (typeof obj !== 'object') {
        return
    }
    if (typeof obj.pageUUID === 'string' && obj.pageUUID) {
        out.push(obj.pageUUID)
    }
    if (typeof obj.pageUuid === 'string' && obj.pageUuid) {
        out.push(obj.pageUuid)
    }
    if (typeof obj.PageID === 'string' && obj.PageID) {
        out.push(obj.PageID)
    }
    const keys = Object.keys(obj)
    for (let i = 0; i < keys.length; i++) {
        collectPageUUIDs(obj[keys[i]], out, depth + 1)
    }
}

function extractPageLinkTargets(pageLayerData) {
    const targets = []
    const cells = pageLayerData && pageLayerData.components && pageLayerData.components.cells
    if (!cells || !cells.length) {
        return targets
    }
    for (let i = 0; i < cells.length; i++) {
        const detail = cells[i] && cells[i].data && cells[i].data.detail
        if (detail) {
            collectPageUUIDs(detail, targets, 0)
        }
    }
    return targets
}

function applyLoadedPageLayer(ctx, pageid, cfg) {
    if (!cfg) {
        return
    }
    const lists = [ctx.state.PCPageList || [], ctx.state.PhonePageList || []]
    for (let li = 0; li < lists.length; li++) {
        const pageList = lists[li]
        for (let i = 0; i < pageList.length; i++) {
            if (pageList[i].pageUuid === pageid) {
                pageList[i]._lazyLoaded = true
                pageList[i].pageLayerData = cfg
                pageList[i].children = buildPageTreeChildren(cfg)
                return
            }
        }
    }
}

function isPagePrefetchable(ctx, pageid) {
    const pageList = [...(ctx.state.PCPageList || []), ...(ctx.state.PhonePageList || [])]
    for (let i = 0; i < pageList.length; i++) {
        const page = pageList[i]
        if (page.pageUuid !== pageid) {
            continue
        }
        if (page.IsHome == 1 || page._lazyLoaded) {
            return false
        }
        if (pendingPageLayerLoads.has(pageid)) {
            return false
        }
        return true
    }
    return false
}

function scorePrefetchCandidates(ctx, anchorPageId) {
    const pageList = [...(ctx.state.PCPageList || []), ...(ctx.state.PhonePageList || [])]
    const scores = new Map()
    const addScore = (pageid, score) => {
        if (!pageid) {
            return
        }
        scores.set(pageid, (scores.get(pageid) || 0) + score)
    }

    const anchor = pageList.find(p => p.pageUuid === anchorPageId)
    if (anchor && anchor.pageLayerData) {
        extractPageLinkTargets(anchor.pageLayerData).forEach((pid, idx) => addScore(pid, 100 - idx))
    }
    const home = pageList.find(p => p.IsHome == 1)
    if (home && home.pageLayerData) {
        extractPageLinkTargets(home.pageLayerData).forEach((pid, idx) => addScore(pid, 60 - idx))
    }
    const anchorIdx = pageList.findIndex(p => p.pageUuid === anchorPageId)
    if (anchorIdx >= 0) {
        for (let d = 1; d <= 3; d++) {
            if (pageList[anchorIdx + d]) {
                addScore(pageList[anchorIdx + d].pageUuid, 24 - d)
            }
            if (pageList[anchorIdx - d]) {
                addScore(pageList[anchorIdx - d].pageUuid, 24 - d)
            }
        }
    }
    for (let i = 0; i < lazyPageLRU.length; i++) {
        const entry = lazyPageLRU[i]
        const recent = pageList.find(p => p.pageUuid === entry.pageid)
        if (recent && recent.pageLayerData) {
            extractPageLinkTargets(recent.pageLayerData).forEach((pid, idx) => addScore(pid, 18 - idx - i))
        }
    }
    return scores
}

function pumpIdlePrefetch(ctx, anchorPageId) {
    if (prefetchRunning) {
        return
    }
    const scores = scorePrefetchCandidates(ctx, anchorPageId)
    prefetchQueue = [...scores.entries()]
        .filter(([pageid]) => isPagePrefetchable(ctx, pageid))
        .sort((a, b) => b[1] - a[1])
        .slice(0, PREFETCH_MAX_CANDIDATES)
        .map(([pageid]) => pageid)
    if (!prefetchQueue.length) {
        return
    }
    prefetchRunning = true

    const runSlice = (deadline) => {
        while (prefetchQueue.length > 0) {
            const timeLeft = deadline && typeof deadline.timeRemaining === 'function'
                ? deadline.timeRemaining()
                : 50
            if (timeLeft < 4 && !(deadline && deadline.didTimeout)) {
                scheduleIdleTask(runSlice)
                return
            }
            const pageid = prefetchQueue.shift()
            if (!pageid || !isPagePrefetchable(ctx, pageid)) {
                continue
            }
            loadSinglePageLayer(pageid).then(function (cfg) {
                applyLoadedPageLayer(ctx, pageid, cfg)
            }).finally(function () {
                if (prefetchQueue.length > 0) {
                    scheduleIdleTask(runSlice)
                } else {
                    prefetchRunning = false
                }
            })
            return
        }
        prefetchRunning = false
    }
    scheduleIdleTask(runSlice)
}

function triggerIdlePrefetch(ctx, anchorPageId) {
    if (!anchorPageId) {
        return
    }
    if (typeof navigator !== 'undefined') {
        const connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection
        if (!navigator.onLine || (connection && (connection.saveData
            || connection.effectiveType === 'slow-2g' || connection.effectiveType === '2g'))) {
            cancelIdlePrefetch()
            return
        }
    }
    if (typeof document !== 'undefined' && document.visibilityState === 'hidden') {
        cancelIdlePrefetch()
        return
    }
    if (prefetchTimer) {
        clearTimeout(prefetchTimer)
    }
    prefetchTimer = setTimeout(function () {
        prefetchTimer = null
        pumpIdlePrefetch(ctx, anchorPageId)
    }, PREFETCH_TRIGGER_DELAY_MS)
}

// 按需加载页面的 LRU 缓存：下钻访问过的非主页页面会把 components 常驻内存，
// 页面极多时(数千页)会持续累积。这里按"最久未访问"顺序逐步释放，只保留主页 +
// 当前页 + 最近 LAZY_PAGE_CACHE_LIMIT 个页面，被释放的页面下次访问会重新按需拉取。
const LAZY_PAGE_CACHE_LIMIT = 12
let lazyPageLRU = [] // [{pageid, at}] 访问记录。自封顶: 超限即按最久未访问回收, 切换大屏后旧 id 也会被优先淘汰

function touchAndReleaseLazyPages(ctx, currentPageId, isHome) {
    if (isHome) {
        return // 主页常驻, 不纳入回收
    }
    const now = Date.now()
    let hit = null
    for (let i = 0; i < lazyPageLRU.length; i++) {
        if (lazyPageLRU[i].pageid === currentPageId) { hit = lazyPageLRU[i]; break }
    }
    if (hit) {
        hit.at = now
    } else {
        lazyPageLRU.push({pageid: currentPageId, at: now})
    }
    if (lazyPageLRU.length <= LAZY_PAGE_CACHE_LIMIT) {
        return
    }
    lazyPageLRU.sort((a, b) => a.at - b.at) // 旧 → 新
    const pcList = ctx.state.PCPageList || []
    while (lazyPageLRU.length > LAZY_PAGE_CACHE_LIMIT) {
        const victim = lazyPageLRU[0]
        if (victim.pageid === currentPageId) {
            break // 当前页不释放
        }
        lazyPageLRU.shift()
        for (let i = 0; i < pcList.length; i++) {
            if (pcList[i].pageUuid === victim.pageid && pcList[i].IsHome != 1) {
                if (pcList[i].pageLayerData && pcList[i].pageLayerData.components) {
                    pcList[i].pageLayerData.components = {cells: []} // 释放组件, GC 回收
                }
                pcList[i]._lazyLoaded = false // 下次访问重新按需拉取
                break
            }
        }
    }
}

export const getLayerDataStruct = (ctx,data) => {
    // 兼容旧调用：直接传 displayUUID 字符串
    if (typeof data === 'string') {
        data = { uuid: data, metaOnly: true, cb: function () {} }
    } else if (!data || typeof data !== 'object') {
        data = { uuid: '', metaOnly: true, cb: function () {} }
    }
    // 默认走 metaOnly 按需加载；仅显式 metaOnly:false 才全量（一般禁止）
    if (data.metaOnly !== false) {
        data.metaOnly = true
    }
    let params={
        muid:data.uuid,
        metaOnly: !!data.metaOnly
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    let isPopUp = data.isPopUp?data.isPopUp:false
    if (!isPopUp) {
        lazyPageLRU = []
        pendingPageLayerLoads.clear()
        cancelIdlePrefetch()
    }
    getDisplayModelLayerData(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer || []
            let is_find_home = 0
            if(pageLayer.length>0)
            {
                let pcPageData = []
                let phonePageData = []
                for(let i=0;i<pageLayer.length;i++)
                {
                    let pageInfo = {
                        id: 9,
                        key: 0,
                        isEdit: false,
                        pageUuid: "",
                        pageModelUuid: "",
                        isNewItem: false,
                        title: '',
                        depth: 1,
                        pageType:1,
                        scopedSlots: { title: 'custom' },
                    }
                    pageInfo.id = pageLayer[i].ID
                    pageInfo.key = i
                    pageInfo.isComponents=false
                    pageInfo.title = pageLayer[i].PageName
                    pageLayer[i].AppName = res.data.Display.name
                    pageInfo.IsHome = pageLayer[i].IsHome
                    pageInfo.IsLogin = pageLayer[i].IsLogin
                    pageInfo.AppName = pageLayer[i].AppName
                    pageInfo.pageUuid = pageLayer[i].PageId
                    pageInfo.pageType = pageLayer[i].PageType
                    pageInfo.pageModelUuid = pageLayer[i].modelId
                    pageInfo.templateKind = pageLayer[i].templateKind || pageLayer[i].TemplateKind || ''
                    pageInfo.templateModelUuid = pageLayer[i].templateModelUuid || pageLayer[i].TemplateModelUuid || ''
                    pageInfo.children=[]

                    pageLayer[i].name = pageLayer[i].PageName
                    let tempConfigData = pageLayer[i]
                    try{
                        tempConfigData = parseRawPageLayerFields(tempConfigData)
                        tempConfigData = normalizePageConfigData(tempConfigData)
                        tempConfigData = markMetaOnlyLazyState(pageInfo, tempConfigData, !!data.metaOnly, pageLayer[i].IsHome)

                    }catch (e) {
                        console.error('[getLayerDataStruct] parse/normalize error:', e.message || e, 'stack:', (e.stack || '').slice(0,400))
                        tempConfigData.layer = {}
                        tempConfigData.components = { cells: [] }
                        tempConfigData = markMetaOnlyLazyState(pageInfo, tempConfigData, !!data.metaOnly, pageLayer[i].IsHome)
                    }
                    pageInfo.children = buildPageTreeChildren(tempConfigData)

                    if(pageLayer[i].IsHome==1&!isPopUp)
                    {
                        if(data.pageType)
                        {
                            if(pageLayer[i].PageType==0)
                            {
                                ctx.state.selectPageUuid = pageLayer[i].PageId
                                is_find_home=1
                                for(let k=0;k<tempConfigData.components.cells.length;k++)
                                {
                                    if(typeof(tempConfigData.components.cells[k].data)=="undefined")
                                    {
                                        continue
                                    }
                                    if(typeof(tempConfigData.components.cells[k].data.detail)=="undefined")
                                    {
                                        continue
                                    }
                                    collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                                    collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)

                                }
                                ctx.state.LayerData = tempConfigData
                            }
                        }
                        else
                        {
                            if(pageLayer[i].PageType==1) {
                                ctx.state.selectPageUuid = pageLayer[i].PageId
                                is_find_home=1
                                for(let k=0;k<tempConfigData.components.cells.length;k++)
                                {
                                    if(typeof(tempConfigData.components.cells[k].data)=="undefined")
                                    {
                                        continue
                                    }
                                    if(typeof(tempConfigData.components.cells[k].data.detail)=="undefined")
                                    {
                                        continue
                                    }
                                    collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                                    collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                                }
                                ctx.state.LayerData = tempConfigData
                            }
                        }


                    }
                    pageInfo.pageLayerData = tempConfigData
                    if(pageLayer[i].PageType==1)
                    {
                        pcPageData.push(pageInfo)
                    }else{
                        phonePageData.push(pageInfo)
                    }
                }
                if(is_find_home==0&!isPopUp)
                {
                    const fallbackPage = pcPageData[0] || phonePageData[0]
                    const tempConfigData = fallbackPage && fallbackPage.pageLayerData
                    if (tempConfigData) {
                        ctx.state.selectPageUuid = fallbackPage.pageUuid || pageLayer[0].PageId
                        const cells = (tempConfigData.components && tempConfigData.components.cells) ? tempConfigData.components.cells : []
                        for(let k=0;k<cells.length;k++)
                        {
                            if(typeof(cells[k].data)=="undefined")
                            {
                                continue
                            }
                            if(typeof(cells[k].data.detail)=="undefined")
                            {
                                continue
                            }
                            collectActiveBindings((cells[k].data && cells[k].data.detail) ? cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                            collectAnimateBindings((cells[k].data && cells[k].data.detail) ? cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                        }
                        ctx.state.LayerData = tempConfigData
                    }
                }
                ctx.state.PCPageList = pcPageData
                ctx.state.PhonePageList = phonePageData
            }
            else{
                ctx.state.PCPageList = []
                ctx.state.PhonePageList = []
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            data.cb(0,res.data.Display.project_uuid,bangDingData,newbangDingDeviceSN)
            if (data.metaOnly && !isPopUp) {
                triggerIdlePrefetch(ctx, ctx.state.selectPageUuid)
            }
        }
        else
        {
            // 禁止 metaOnly 失败后回退全量加载：大屏页多时会 30s+ 超时并打满内存
            console.error('[getLayerDataStruct] metaOnly/full response code!=0, code=', res && res.data && res.data.code)
            data.cb(-1,"",bangDingData,bangDingDeviceSN)
        }

    }).catch(function (e) {
        console.error('[getLayerDataStruct] request error:', e && e.message)
        data.cb(-1,"",bangDingData,bangDingDeviceSN)
    })
}
export const updateAllLayerDataStruct = (ctx,data) => {
    const useMetaOnly = data.metaOnly !== false
    let params={
        muid:data.uuid,
        metaOnly: useMetaOnly
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    const oldPc = ctx.state.PCPageList || []
    const oldPhone = ctx.state.PhonePageList || []
    getDisplayModelLayerData(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer || []
            let is_find_home = 0
            if(pageLayer.length>0)
            {
                let pcPageData = []
                let phonePageData = []
                for(let i=0;i<pageLayer.length;i++)
                {
                    let pageInfo = {
                        id: 9,
                        key: 0,
                        isEdit: false,
                        pageUuid: "",
                        pageModelUuid: "",
                        isNewItem: false,
                        title: '',
                        depth: 1,
                        pageType:1,
                        scopedSlots: { title: 'custom' },
                    }
                    pageInfo.id = pageLayer[i].ID
                    pageInfo.key = i
                    pageInfo.isComponents=false
                    pageInfo.title = pageLayer[i].PageName
                    pageLayer[i].AppName = res.data.Display.name
                    pageInfo.IsHome = pageLayer[i].IsHome
                    pageInfo.IsLogin = pageLayer[i].IsLogin
                    pageInfo.AppName = pageLayer[i].AppName
                    pageInfo.pageUuid = pageLayer[i].PageId
                    pageInfo.pageType = pageLayer[i].PageType
                    pageInfo.pageModelUuid = pageLayer[i].modelId
                    pageInfo.templateKind = pageLayer[i].templateKind || pageLayer[i].TemplateKind || ''
                    pageInfo.templateModelUuid = pageLayer[i].templateModelUuid || pageLayer[i].TemplateModelUuid || ''
                    pageInfo.children=[]

                    pageLayer[i].name = pageLayer[i].PageName
                    let tempConfigData = pageLayer[i]
                    try{
                        tempConfigData = parseRawPageLayerFields(tempConfigData)
                        tempConfigData = normalizePageConfigData(tempConfigData)
                        tempConfigData = markMetaOnlyLazyState(pageInfo, tempConfigData, useMetaOnly, pageLayer[i].IsHome)

                    }catch (e) {
                        console.error('[getLayerDataStruct] parse/normalize error:', e.message || e, 'stack:', (e.stack || '').slice(0,400))
                        tempConfigData.layer = {}
                        tempConfigData.components = { cells: [] }
                        tempConfigData = markMetaOnlyLazyState(pageInfo, tempConfigData, useMetaOnly, pageLayer[i].IsHome)
                    }
                    pageInfo.children = buildPageTreeChildren(tempConfigData)
                    mergePageLayerFromExisting(pageInfo, oldPc, oldPhone)
                    if (!pageInfo.pageLayerData) {
                        pageInfo.pageLayerData = tempConfigData
                    }
                    if(pageLayer[i].PageType==1)
                    {
                        pcPageData.push(pageInfo)
                    }else{
                        phonePageData.push(pageInfo)
                    }
                    if(pageInfo.pageUuid==ctx.state.selectPageUuid)
                    {
                        ctx.state.LayerData = pageInfo.pageLayerData
                    }
                }
                ctx.state.PCPageList = pcPageData
                ctx.state.PhonePageList = phonePageData
            }
            else{
                ctx.state.PCPageList = []
                ctx.state.PhonePageList = []
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            data.cb(0,res.data.Display.project_uuid,bangDingData,newbangDingDeviceSN)
        }
        else
        {
            data.cb(-1,"",bangDingData,bangDingDeviceSN)
        }

    })
}
export const getLayerPagerContainerDataStruct = (ctx,data) => {
    let params={
        pageid:data.pageid
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    getDisplayModelPagerLayerData(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer
            let tempConfigData = pageLayer
            try{
                tempConfigData = parseRawPageLayerFields(tempConfigData)
                tempConfigData = normalizePageConfigData(tempConfigData)
                for(let k=0;k<tempConfigData.components.cells.length;k++)
                {
                    if(typeof(tempConfigData.components.cells[k].data)=="undefined")
                    {
                        continue
                    }
                    if(typeof(tempConfigData.components.cells[k].data.detail)=="undefined")
                    {
                        continue
                    }
                    collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                    collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                }
                ctx.state.LayerContainerData = tempConfigData
                let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                data.cb(0,tempConfigData,bangDingData,newbangDingDeviceSN)

            }catch (e) {
                console.log(e)
                data.cb(-3,null,bangDingData,null)
            }
        }
        else
        {
            data.cb(-1,null,bangDingData,bangDingDeviceSN)
        }

    })
}
export const selectPopUpPagerContainerDisplayPageDataStruct = (ctx,page) => {
    let PCPageInfo = ctx.state.PCPageList
    let PhonePageInfo = ctx.state.PhonePageList

    let pageid = page.page ? (page.page.pageUuid || page.page.pageUUID || '') : ''
    let bangDingData=[]
    let bangDingDeviceSN=[]

    for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
    {
        if(PCPageInfo[i].pageUuid==pageid)
        {
            const cells = PCPageInfo[i].pageLayerData && PCPageInfo[i].pageLayerData.components
                && PCPageInfo[i].pageLayerData.components.cells
            if (!PCPageInfo[i]._lazyLoaded && PCPageInfo[i].IsHome != 1
                && Array.isArray(cells) && cells.length === 0) {
                const pageIndex = i
                ctx.state.pageLayerLoading = true
                loadSinglePageLayer(pageid).then(function (cfg) {
                    ctx.state.pageLayerLoading = false
                    if (!cfg) {
                        page.callback(-1, bangDingData, bangDingDeviceSN)
                        return
                    }
                    PCPageInfo[pageIndex]._lazyLoaded = true
                    PCPageInfo[pageIndex].pageLayerData = cfg
                    selectPopUpPagerContainerDisplayPageDataStruct(ctx, page)
                })
                return
            }
            let tempConfigData = PCPageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
    {
        if(PhonePageInfo[i].pageUuid==pageid) {
            const cells = PhonePageInfo[i].pageLayerData && PhonePageInfo[i].pageLayerData.components
                && PhonePageInfo[i].pageLayerData.components.cells
            if (!PhonePageInfo[i]._lazyLoaded && PhonePageInfo[i].IsHome != 1
                && Array.isArray(cells) && cells.length === 0) {
                const pageIndex = i
                ctx.state.pageLayerLoading = true
                loadSinglePageLayer(pageid).then(function (cfg) {
                    ctx.state.pageLayerLoading = false
                    if (!cfg) {
                        page.callback(-1, bangDingData, bangDingDeviceSN)
                        return
                    }
                    PhonePageInfo[pageIndex]._lazyLoaded = true
                    PhonePageInfo[pageIndex].pageLayerData = cfg
                    selectPopUpPagerContainerDisplayPageDataStruct(ctx, page)
                })
                return
            }
            let tempConfigData = PhonePageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    getLayerDataStruct(ctx,{uuid:page.page.displayUUID,isPopUp:true,metaOnly:true,cb:function () {
            let PCPageInfo = ctx.state.PCPageList
            let PhonePageInfo = ctx.state.PhonePageList

            let pageid = page.page ? (page.page.pageUuid || page.page.pageUUID || '') : ''
            let bangDingData = []
            let bangDingDeviceSN = []
            for (let i = 0, PCPageInfoLen = PCPageInfo.length; i < PCPageInfoLen; i++) {
                if (PCPageInfo[i].pageUuid == pageid) {
                    let tempConfigData = PCPageInfo[i].pageLayerData
                    tempConfigData.name = page.page.title
                    for (let k = 0, componentsLen = tempConfigData.components.cells.length; k < componentsLen; k++) {
                        collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                        collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                    }
                    let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                    page.callback(0,tempConfigData, bangDingData, newbangDingDeviceSN)
                    return
                }
            }

            for (let i = 0, PhonePageInfoLen = PhonePageInfo.length; i < PhonePageInfoLen; i++) {
                if (PhonePageInfo[i].pageUuid == pageid) {
                    let tempConfigData = PhonePageInfo[i].pageLayerData
                    tempConfigData.name = page.page.title
                    for (let k = 0, componentsLen = tempConfigData.components.cells.length; k < componentsLen; k++) {
                        collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                        collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                    }
                    let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                    page.callback(0, tempConfigData,bangDingData, newbangDingDeviceSN)
                    return
                }
            }
            const tempConfigData ={ "name": "--", "layer": { "backColor": "", "backgroundImage": "", "widthHeightRatio": "", "width": 300, "height": 600 }, "components": [] }
            page.callback(-1,tempConfigData,bangDingData,bangDingDeviceSN)
        }});
}
export const selectDisplayPageContainerDataStruct = (ctx,page) => {
    let PCPageInfo = ctx.state.PCPageList
    let PhonePageInfo = ctx.state.PhonePageList

    let pageid = page.page ? (page.page.pageUuid || page.page.pageUUID || '') : ''
    let bangDingData=[]
    let bangDingDeviceSN=[]
    for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
    {
        if(PCPageInfo[i].pageUuid==pageid)
        {
            let tempConfigData = PCPageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
    {
        if(PhonePageInfo[i].pageUuid==pageid) {
            let tempConfigData = PhonePageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            ctx.state.selectPageUuid = tempConfigData.PageId
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    getLayerDataStruct(ctx,{uuid:page.page.displayUUID,metaOnly:true,cb:function (){
            let PCPageInfo = ctx.state.PCPageList
            let PhonePageInfo = ctx.state.PhonePageList

            let pageid = page.page ? (page.page.pageUuid || page.page.pageUUID || '') : ''
            let bangDingData=[]
            let bangDingDeviceSN=[]
            for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
            {
                if(PCPageInfo[i].pageUuid==pageid)
                {
                    let tempConfigData = PCPageInfo[i].pageLayerData
                    tempConfigData.name = page.page.title
                    for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
                    {
                        collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                        collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                    }
                    let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                    page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
                    return
                }
            }

            for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
            {
                if(PhonePageInfo[i].pageUuid==pageid) {
                    let tempConfigData = PhonePageInfo[i].pageLayerData
                    tempConfigData.name = page.page.title
                    for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
                    {
                        collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                        collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                    }
                    let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                    page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
                    return
                }
            }
            page.callback(-1,null,bangDingData,bangDingDeviceSN)
        }});
}

export const getLayerDataStructByTokenData = (ctx,data) => {
    let params={
        muid:data.uuid,
        token:data.token,
        metaOnly:data.metaOnly !== false
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    let isPopUp = data.isPopUp?data.isPopUp:false
    getLayerDataStructByToken(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer || []
            let is_find_home = 0
            let PC_home_index = 0
            let Phone_home_index = 0
            if(pageLayer.length>0)
            {
                let pcPageData = []
                let phonePageData = []
                for(let i=0;i<pageLayer.length;i++)
                {
                    let pageInfo = {
                        id: 9,
                        key: 0,
                        isEdit: false,
                        pageUuid: "",
                        pageModelUuid: "",
                        isNewItem: false,
                        title: '',
                        depth: 1,
                        pageType:1,
                        scopedSlots: { title: 'custom' },
                    }
                    pageInfo.id = pageLayer[i].ID
                    pageInfo.key = i
                    pageInfo.isComponents=false
                    pageInfo.title = pageLayer[i].PageName
                    pageLayer[i].AppName = res.data.Display.name
                    pageInfo.IsHome = pageLayer[i].IsHome
                    pageInfo.IsLogin = pageLayer[i].IsLogin
                    pageInfo.AppName = pageLayer[i].AppName
                    pageInfo.pageUuid = pageLayer[i].PageId
                    pageInfo.pageType = pageLayer[i].PageType
                    pageInfo.pageModelUuid = pageLayer[i].modelId
                    pageInfo.templateKind = pageLayer[i].templateKind || pageLayer[i].TemplateKind || ''
                    pageInfo.templateModelUuid = pageLayer[i].templateModelUuid || pageLayer[i].TemplateModelUuid || ''
                    pageInfo.children=[]

                    pageLayer[i].name = pageLayer[i].PageName
                    let tempConfigData = pageLayer[i]
                    try{
                        tempConfigData = parseRawPageLayerFields(tempConfigData)
                        tempConfigData = normalizePageConfigData(tempConfigData)

                    }catch (e) {
                        continue
                    }

                    // for(let k=0;k<tempConfigData.components.cells.length;k++)
                    // {
                    //     let components = {
                    //         isComponents:true,
                    //         title:tempConfigData.components.cells[k].data.detail.name,
                    //         key:tempConfigData.components.cells[k].data.detail.identifier
                    //     }
                    //     pageInfo.children.push(components)
                    // }

                    if(pageLayer[i].IsHome==1&!isPopUp)
                    {
                        if(data.pageType)
                        {
                            if(pageLayer[i].PageType==0)
                            {
                                ctx.state.selectPageUuid = pageLayer[i].PageId
                                is_find_home=1
                                Phone_home_index = i
                                for(let k=0;k<tempConfigData.components.cells.length;k++)
                                {
                                    if(typeof(tempConfigData.components.cells[k].data)=="undefined")
                                    {
                                        continue
                                    }
                                    if(typeof(tempConfigData.components.cells[k].data.detail)=="undefined")
                                    {
                                        continue
                                    }
                                    collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                                    collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                                }
                                ctx.state.LayerData = tempConfigData
                            }
                        }
                        else
                        {
                            if(pageLayer[i].PageType==1) {
                                ctx.state.selectPageUuid = pageLayer[i].PageId
                                is_find_home=1
                                PC_home_index = i
                                for(let k=0;k<tempConfigData.components.cells.length;k++)
                                {
                                    if(typeof(tempConfigData.components.cells[k].data)=="undefined")
                                    {
                                        continue
                                    }
                                    if(typeof(tempConfigData.components.cells[k].data.detail)=="undefined")
                                    {
                                        continue
                                    }
                                    collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                                    collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                                }
                                ctx.state.LayerData = tempConfigData
                            }
                        }
                    }
                    pageInfo.pageLayerData = tempConfigData
                    if(pageLayer[i].PageType==1)
                    {
                        pcPageData.push(pageInfo)
                    }else{
                        phonePageData.push(pageInfo)
                    }
                }
                if(is_find_home==0&!isPopUp)
                {
                    let tempConfigData={}
                    if(data.pageType)
                    {
                        tempConfigData= pageLayer[PC_home_index]
                        ctx.state.selectPageUuid = pageLayer[PC_home_index].PageId
                    }
                    else
                    {
                        tempConfigData= pageLayer[0]
                        ctx.state.selectPageUuid = pageLayer[0].PageId
                    }
                    for(let k=0,componentsLen =tempConfigData.components.cells.length ;k<componentsLen;k++)
                    {
                        collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                        collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                    }
                    ctx.state.LayerData = tempConfigData
                }
                ctx.state.PCPageList = pcPageData
                ctx.state.PhonePageList = phonePageData
            }
            else{
                ctx.state.PCPageList = []
                ctx.state.PhonePageList = []
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            data.cb(0,res.data.Display.project_uuid,res.data.expireAt,res.data.token,bangDingData,newbangDingDeviceSN)
        }
        else
        {
            data.cb(res.data.code,"","","",bangDingData,bangDingDeviceSN)
        }

    })
}

export const getLoginLayerDataStruct = (ctx,data) => {
    let params={
        muid:data.uuid,
        pageType:data.pageType
    }
    let isPopUp = data.isPopUp?data.isPopUp:false
    GetDisplayLoginPage(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer
            let is_find_home = 0

                let pageInfo = {
                    id: 9,
                    key: 0,
                    isEdit: false,
                    pageUuid: "",
                    pageModelUuid: "",
                    isNewItem: false,
                    title: '',
                    depth: 1,
                    pageType:1,
                    scopedSlots: { title: 'custom' },
                }
                pageInfo.id = pageLayer.ID
                pageInfo.key = 0
                pageInfo.isComponents=false
                pageInfo.title = pageLayer.PageName
                pageLayer.AppName = res.data.Display.name
                pageInfo.IsHome = pageLayer.IsHome
                pageInfo.IsLogin = pageLayer.IsLogin
                pageInfo.AppName = pageLayer.AppName
                pageInfo.pageUuid = pageLayer.PageId
                pageInfo.pageType = pageLayer.PageType
                pageInfo.pageModelUuid = pageLayer.modelId
                pageInfo.children=[]

                pageLayer.name = pageLayer.PageName
                let tempConfigData = pageLayer
                try{
                    tempConfigData = parseRawPageLayerFields(tempConfigData)
                    tempConfigData = normalizePageConfigData(tempConfigData)

                }catch (e) {
                    data.cb(-2,"")
                    return
                }

                for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
                {
                    // let components = {
                    //     isComponents:true,
                    //     title:tempConfigData.components.cells[k].data.detail.name,
                    //     key:tempConfigData.components.cells[k].data.detail.identifier
                    // }
                    if(typeof tempConfigData.components.cells[k].data.detail.animate!=="undefined"&&typeof tempConfigData.components.cells[k].data.detail.animate.move=="undefined")
                    {
                        if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                            tempConfigData.components.cells[k].data.detail.animate.move = {
                                x: {
                                    deviceSN: "",
                                    selectVideoType: 0,
                                    isBandDevice: false,
                                    bandType: 1,
                                    dataID: "",
                                    dataName: "",
                                },
                                y: {
                                    deviceSN: "",
                                    selectVideoType: 0,
                                    isBandDevice: false,
                                    bandType: 1,
                                    dataID: "",
                                    dataName: "",
                                },
                            }
                        }
                    }
                    // pageInfo.children.push(components)
                }
                ctx.state.LayerData = tempConfigData
                ctx.state.selectPageUuid = pageLayer.PageId

                data.cb(0,"")
                return
        }
        else{
            ctx.state.PCPageList = []
            ctx.state.PhonePageList = []
        }
        data.cb(-1,"")
    })
}

export const saveLayerDataStruct = (ctx,page) => {
    if (ctx.state.editorRuntimePreview && ctx.state.editorRuntimePreview.active) {
        return Promise.resolve({
            data: {
                code: 4090,
                message: '运行态预览不可直接保存，请编辑对应模板后再保存'
            }
        })
    }
    let params={
        muid:page.uuid,
        pageid:page.pageid,
        saveData:page.LayerData
    }
    return setDisplayModelLayerData(params).then(function (res){
        return res
    })
}

function applyEditorPageLayer(ctx, tempConfigData, page) {
    tempConfigData.name = page.title
    for(let k=0,componentsLen =tempConfigData.components.cells.length ;k<componentsLen;k++)
    {
        if(typeof(tempConfigData.components.cells[k].data)=="undefined")
        {
            continue
        }
        if(typeof(tempConfigData.components.cells[k].data.detail)=="undefined")
        {
            continue
        }
        if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined"&& typeof tempConfigData.components.cells[k].data.detail.animate.move=="undefined")
        {
            tempConfigData.components.cells[k].data.detail.animate.move = {
                x:{
                    deviceSN:"",
                    selectVideoType:0,
                    isBandDevice:false,
                    bandType:1,
                    dataID: "",
                    dataName: "",
                },
                y:{
                    deviceSN:"",
                    selectVideoType:0,
                    isBandDevice:false,
                    bandType:1,
                    dataID: "",
                    dataName: "",
                },
            }
        }
    }
    ctx.state.selectPageUuid = tempConfigData.PageId
    ctx.state.LayerData = tempConfigData
}

function ensureEditorPageLoaded(ctx, list, idx, page, pageid) {
    if (list[idx]._lazyLoaded) {
        return false
    }
    ctx.state.pageLayerLoading = true
    loadSinglePageLayer(pageid).then(function(cfg){
        ctx.state.pageLayerLoading = false
        list[idx]._lazyLoaded = true
        if(cfg){
            list[idx].pageLayerData = cfg
            list[idx].children = buildPageTreeChildren(cfg)
        }
        selectLayerDataStruct(ctx, page)
    }).catch(function(){
        ctx.state.pageLayerLoading = false
    })
    return true
}

export const selectLayerDataStruct = (ctx,page) => {
    let PCPageInfo = ctx.state.PCPageList
    let PhonePageInfo = ctx.state.PhonePageList

    let pagetype = page.pageType
    let pageid = page.pageUuid
    ctx.state.curPageUuid = pageid
    if(ctx.state.curPageUuid!=ctx.state.prePageUuid)
    {
        ctx.state.prePageUuid = ctx.state.curPageUuid
    }
    if(pagetype==1)
    {
        for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
        {
            if(PCPageInfo[i].pageUuid==pageid)
            {
                if(ensureEditorPageLoaded(ctx, PCPageInfo, i, page, pageid)) {
                    return
                }
                let tempConfigData = PCPageInfo[i].pageLayerData
                applyEditorPageLayer(ctx, tempConfigData, page)
                touchAndReleaseLazyPages(ctx, pageid, PCPageInfo[i].IsHome==1)
            }
        }
    }else if(pagetype==0)
    {
        for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
        {
            if(PhonePageInfo[i].pageUuid==pageid) {
                if(ensureEditorPageLoaded(ctx, PhonePageInfo, i, page, pageid)) {
                    return
                }
                let tempConfigData = PhonePageInfo[i].pageLayerData
                applyEditorPageLayer(ctx, tempConfigData, page)
                touchAndReleaseLazyPages(ctx, pageid, PhonePageInfo[i].IsHome==1)
            }
        }
    }
}

/**
 * 编辑器运行态虚拟页预览。
 * 始终对缓存中的原始模板做深拷贝后解析，resolved 结果只进入 LayerData，
 * 不回写 PCPageList.pageLayerData，避免把设备上下文和分页数据保存进模板。
 */
export const selectEditorRuntimePreview = async (ctx, payload) => {
    const pageid = payload && payload.pageUuid
    if (!pageid) {
        return null
    }
    const list = ctx.state.PCPageList || []
    const idx = list.findIndex(item => item.pageUuid === pageid)
    if (idx < 0) {
        return null
    }

    const pageInfo = list[idx]
    if (!pageInfo._lazyLoaded) {
        ctx.state.pageLayerLoading = true
        try {
            const cfg = await loadSinglePageLayer(pageid)
            if (!cfg) {
                return null
            }
            pageInfo.pageLayerData = cfg
            pageInfo.children = buildPageTreeChildren(cfg)
            pageInfo._lazyLoaded = true
        } finally {
            ctx.state.pageLayerLoading = false
        }
    }
    if (!pageInfo.pageLayerData) {
        return null
    }

    ctx.commit('setNavContext', payload.navContext || null)
    const rawTemplate = JSON.parse(JSON.stringify(pageInfo.pageLayerData))
    const resolved = await applyNavContextToPageConfig(ctx, rawTemplate)
    if (!resolved) {
        return null
    }
    resolved.name = payload.virtualTitle || pageInfo.title
    ctx.state.curPageUuid = pageid
    ctx.state.prePageUuid = pageid
    ctx.state.selectPageUuid = pageid
    ctx.state.LayerData = resolved
    ctx.commit('setEditorRuntimePreview', {
        active: true,
        virtualKey: payload.virtualKey,
        virtualTitle: payload.virtualTitle,
        templatePageUuid: pageid,
        templateTitle: pageInfo.title,
    })
    touchAndReleaseLazyPages(ctx, pageid, pageInfo.IsHome == 1)
    return resolved
}

export const selectDisplayPageDataStructFromDb = (ctx,page) => {
    getLayerDataStruct(ctx,{uuid:page.page.displayUUID,metaOnly:true,cb:function (){
            let PCPageInfo = ctx.state.PCPageList
            let PhonePageInfo = ctx.state.PhonePageList

            let pageid = page.page ? (page.page.pageUuid || page.page.pageUUID || '') : ''
            let bangDingData=[]
            let bangDingDeviceSN=[]
            for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
            {
                if(PCPageInfo[i].pageUuid==pageid)
                {
                    let tempConfigData = PCPageInfo[i].pageLayerData
                    tempConfigData.name = page.page.title
                    ctx.state.selectPageUuid = tempConfigData.PageId
                    for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
                    {
                        collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                        collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                    }
                    ctx.state.LayerData = tempConfigData
                    let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                    page.callback(bangDingData,newbangDingDeviceSN)
                    return
                }
            }

            for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
            {
                if(PhonePageInfo[i].pageUuid==pageid) {
                    let tempConfigData = PhonePageInfo[i].pageLayerData
                    tempConfigData.name = page.page.title
                    ctx.state.selectPageUuid = tempConfigData.PageId
                    for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
                    {
                        collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                        collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
                    }
                    ctx.state.LayerData = tempConfigData
                    let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                    page.callback(bangDingData,newbangDingDeviceSN)
                    return
                }
            }
            page.callback(bangDingData,bangDingDeviceSN)
        }});
}
export const selectDisplayPageDataStruct = (ctx,page) => {
    let PCPageInfo = ctx.state.PCPageList
    let PhonePageInfo = ctx.state.PhonePageList

    let pageid = page.page ? (page.page.pageUuid || page.page.pageUUID || '') : ''
    ismDebug('SCADA.selectPage.enter', {
        pageid,
        displayUUID: page.page && page.page.displayUUID,
        pcPages: PCPageInfo.length,
        phonePages: PhonePageInfo.length,
    })
    console.log('[selectDisplayPageDataStruct] pageid=', pageid, 'displayUUID=', page.page?.displayUUID, 'PCPageInfo.length=', PCPageInfo.length)
    if (PCPageInfo.length > 0) {
        console.log('[selectDisplayPageDataStruct] PCPageInfo pageUuids=', PCPageInfo.map(p => ({ name: p.title, pageUuid: p.pageUuid, pageModelUuid: p.pageModelUuid })))
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    if(!pageid)
    {
        console.log('[selectDisplayPageDataStruct] pageid is empty, calling getLayerDataStruct with metaOnly')
        ismDebug('SCADA.selectPage.emptyPageId', {displayUUID: page.page && page.page.displayUUID})
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,metaOnly:true,cb:function (errno, project_uuid, uuids, devices){
                ismDebug('SCADA.selectPage.emptyPageId.cb', {errno, uuids: (uuids||[]).length, devices: (devices||[]).length})
                page.callback(uuids,devices,errno == 0)
        }});
        return
    }
    if(!isDisplayPagesLoaded(ctx, page.page.displayUUID))
    {
        ismDebug('SCADA.selectPage.notLoaded', {displayUUID: page.page.displayUUID, pageid})
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,metaOnly:true,cb:function (errno){
                if(errno == 0 && isDisplayPagesLoaded(ctx, page.page.displayUUID))
                {
                    selectDisplayPageDataStruct(ctx,page)
                    return
                }
                if (errno == 0 && pageid) {
                    ensurePageRegistered(ctx, pageid, page.page.displayUUID).then(function (pageInfo) {
                        if (pageInfo) {
                            selectDisplayPageDataStruct(ctx, page)
                            return
                        }
                        ismDebug('SCADA.selectPage.notLoaded.fail', {errno, reason: 'lazyRegisterFailed'})
                        page.callback(bangDingData, bangDingDeviceSN, false)
                    })
                    return
                }
                ismDebug('SCADA.selectPage.notLoaded.fail', {errno})
                page.callback(bangDingData,bangDingDeviceSN,false)
        }});
        return
    }
    for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
    {
        if(PCPageInfo[i].pageUuid==pageid)
        {
            console.log('[selectDisplayPageDataStruct] FOUND pageUuid=', pageid, 'at index=', i, 'title=', PCPageInfo[i].title)
            // 按需加载: 元数据模式下非首页 components 为空, 首次下钻时单独拉取目标页
            if(!PCPageInfo[i]._lazyLoaded && PCPageInfo[i].IsHome!=1
                && PCPageInfo[i].pageLayerData && PCPageInfo[i].pageLayerData.components
                && PCPageInfo[i].pageLayerData.components.cells
                && PCPageInfo[i].pageLayerData.components.cells.length===0)
            {
                let _idx = i
                ctx.state.pageLayerLoading = true
                ismDebug('SCADA.selectPage.lazyLoad', {pageid, title: PCPageInfo[i].title})
                loadSinglePageLayer(pageid).then(function(cfg){
                    ctx.state.pageLayerLoading = false
                    PCPageInfo[_idx]._lazyLoaded = true
                    if(cfg){
                        PCPageInfo[_idx].pageLayerData = cfg
                    }
                    selectDisplayPageDataStruct(ctx,page)
                    triggerIdlePrefetch(ctx, pageid)
                }).catch(function(err){
                    ctx.state.pageLayerLoading = false
                    ismDebug('SCADA.selectPage.lazyLoad.fail', {message: err && err.message})
                    page.callback(bangDingData,bangDingDeviceSN,false)
                })
                return
            }
            let tempConfigData = PCPageInfo[i].pageLayerData
            applyNavContextToPageConfig(ctx, tempConfigData).then(function (resolvedCfg) {
            try {
                tempConfigData = resolvedCfg
                tempConfigData.name = page.page.title
                ctx.state.selectPageUuid = tempConfigData.PageId
                const cells = (tempConfigData.components && tempConfigData.components.cells) ? tempConfigData.components.cells : []
                for(let k=0,componentsLen=cells.length;k<componentsLen;k++)
                {
                    const detail = cells[k] && cells[k].data && cells[k].data.detail
                    if (!detail) {
                        continue
                    }
                    collectActiveBindings(detail.active, bangDingData, bangDingDeviceSN)
                    collectAnimateBindings(detail.animate, bangDingData, bangDingDeviceSN)
                }
                ctx.state.LayerData = tempConfigData
                let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                touchAndReleaseLazyPages(ctx, pageid, PCPageInfo[i].IsHome==1)
                triggerIdlePrefetch(ctx, pageid)
                ismDebug('SCADA.selectPage.ok', {
                    pageid,
                    cells: cells.length,
                    bindData: bangDingData.length,
                    bindDevices: newbangDingDeviceSN.length,
                })
                page.callback(bangDingData,newbangDingDeviceSN)
            } catch (e) {
                ismDebug('SCADA.selectPage.crash', {
                    pageid,
                    message: e && e.message,
                    stack: e && e.stack,
                })
                page.callback(bangDingData,bangDingDeviceSN,false)
            }
            })
            return
        }
    }
    console.log('[selectDisplayPageDataStruct] NOT FOUND pageUuid=', pageid, 'in PCPageInfo, total pages=', PCPageInfo.length)

    const navCtx = ctx.state.navContext
    const navIndex = ctx.state.navTreeIndex
    const templateMap = ctx.state.navTemplateMap

    // 合成 uuid5 page_id（导航树 legacy）→ 模板页 + navContext
    const oldTarget = resolveOldPageTarget(pageid, navIndex, templateMap)
    if (oldTarget && oldTarget.pageUuid && oldTarget.pageUuid !== pageid) {
      console.log('[selectDisplayPageDataStruct] oldPageId remap →', oldTarget.pageUuid)
      if (oldTarget.navContext) {
        try {
          ctx.commit('setNavContext', oldTarget.navContext)
        } catch (e) { /* ignore */ }
      }
      page.page.pageUuid = oldTarget.pageUuid
      selectDisplayPageDataStruct(ctx, page)
      return
    }

    const displayUUID = page.page && page.page.displayUUID
    let targetPageId = pageid
    if (navCtx && (isDeviceListNav(navCtx) || navCtx.routeMode === 'childrenList' || navCtx.routeMode === 'org')) {
      const fallbackId = resolveDeviceListTemplateId(templateMap, PCPageInfo)
      if (fallbackId) {
        targetPageId = fallbackId
        if (fallbackId !== pageid) {
          console.log('[selectDisplayPageDataStruct] nav template fallback →', fallbackId)
          page.page.pageUuid = fallbackId
        }
      }
    } else if (navCtx && (navCtx.signalMode || navCtx.routeMode === 'signal')) {
      const fallbackId = resolveDeviceSignalTemplateId(templateMap, PCPageInfo, navCtx.modelUuid || navCtx.muid || '')
      if (fallbackId) {
        targetPageId = fallbackId
        if (fallbackId !== pageid) {
          console.log('[selectDisplayPageDataStruct] device signal template fallback →', fallbackId)
          page.page.pageUuid = fallbackId
        }
      }
    }

    for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
    {
        if(PhonePageInfo[i].pageUuid==pageid) {
            let tempConfigData = PhonePageInfo[i].pageLayerData
            applyNavContextToPageConfig(ctx, tempConfigData).then(function (resolvedCfg) {
            try {
                tempConfigData = resolvedCfg
                tempConfigData.name = page.page.title
                ctx.state.selectPageUuid = tempConfigData.PageId
                const cells = (tempConfigData.components && tempConfigData.components.cells) ? tempConfigData.components.cells : []
                for(let k=0,componentsLen=cells.length;k<componentsLen;k++)
                {
                    const detail = cells[k] && cells[k].data && cells[k].data.detail
                    if (!detail) {
                        continue
                    }
                    collectActiveBindings(detail.active, bangDingData, bangDingDeviceSN)
                    collectAnimateBindings(detail.animate, bangDingData, bangDingDeviceSN)
                }
                ctx.state.LayerData = tempConfigData
                let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                ismDebug('SCADA.selectPage.phone.ok', {
                    pageid,
                    cells: cells.length,
                    bindData: bangDingData.length,
                })
                page.callback(bangDingData,newbangDingDeviceSN)
            } catch (e) {
                ismDebug('SCADA.selectPage.phone.crash', {pageid, message: e && e.message})
                page.callback(bangDingData,bangDingDeviceSN,false)
            }
            })
            return
        }
    }

    ensurePageRegistered(ctx, targetPageId, displayUUID).then(function (pageInfo) {
        if (pageInfo) {
            console.log('[selectDisplayPageDataStruct] lazy register →', targetPageId)
            selectDisplayPageDataStruct(ctx, page)
            return
        }
        ctx.state.selectPageUuid = ""
        ismDebug('SCADA.selectPage.notFound', {pageid: targetPageId, pcPages: PCPageInfo.length, phonePages: PhonePageInfo.length})
        page.callback(bangDingData, bangDingDeviceSN, false)
    })
}

export const selectPopUpDisplayPageDataStruct = (ctx,page) => {
    let PCPageInfo = ctx.state.PCPageList
    let PhonePageInfo = ctx.state.PhonePageList

    let pageid = page.page ? (page.page.pageUuid || page.page.pageUUID || '') : ''
    let bangDingData=[]
    let bangDingDeviceSN=[]
    if(!pageid)
    {
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,isPopUp:true,metaOnly:true,cb:function (errno, project_uuid, uuids, devices){
                page.callback(errno == 0 ? 0 : -1,uuids,devices)
        }});
        return
    }
    if(!isDisplayPagesLoaded(ctx, page.page.displayUUID))
    {
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,isPopUp:true,metaOnly:true,cb:function (errno){
                if(errno == 0 && isDisplayPagesLoaded(ctx, page.page.displayUUID))
                {
                    selectPopUpDisplayPageDataStruct(ctx,page)
                    return
                }
                page.callback(-1,bangDingData,bangDingDeviceSN)
        }});
        return
    }

    for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
    {
        if(PCPageInfo[i].pageUuid==pageid)
        {
            const cells = PCPageInfo[i].pageLayerData && PCPageInfo[i].pageLayerData.components
                && PCPageInfo[i].pageLayerData.components.cells
            if (!PCPageInfo[i]._lazyLoaded && PCPageInfo[i].IsHome != 1
                && Array.isArray(cells) && cells.length === 0) {
                const pageIndex = i
                ctx.state.pageLayerLoading = true
                loadSinglePageLayer(pageid).then(function (cfg) {
                    ctx.state.pageLayerLoading = false
                    if (!cfg) {
                        page.callback(-1, bangDingData, bangDingDeviceSN)
                        return
                    }
                    PCPageInfo[pageIndex]._lazyLoaded = true
                    PCPageInfo[pageIndex].pageLayerData = cfg
                    selectPopUpDisplayPageDataStruct(ctx, page)
                })
                return
            }
            let tempConfigData = PCPageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            ctx.state.selectPageUuid = tempConfigData.PageId
            ctx.state.PopUpConfigData = tempConfigData
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
    {
        if(PhonePageInfo[i].pageUuid==pageid) {
            const cells = PhonePageInfo[i].pageLayerData && PhonePageInfo[i].pageLayerData.components
                && PhonePageInfo[i].pageLayerData.components.cells
            if (!PhonePageInfo[i]._lazyLoaded && PhonePageInfo[i].IsHome != 1
                && Array.isArray(cells) && cells.length === 0) {
                const pageIndex = i
                ctx.state.pageLayerLoading = true
                loadSinglePageLayer(pageid).then(function (cfg) {
                    ctx.state.pageLayerLoading = false
                    if (!cfg) {
                        page.callback(-1, bangDingData, bangDingDeviceSN)
                        return
                    }
                    PhonePageInfo[pageIndex]._lazyLoaded = true
                    PhonePageInfo[pageIndex].pageLayerData = cfg
                    selectPopUpDisplayPageDataStruct(ctx, page)
                })
                return
            }
            let tempConfigData = PhonePageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            ctx.state.selectPageUuid = tempConfigData.PageId
            ctx.state.PopUpConfigData = tempConfigData
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                collectActiveBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.active : null, bangDingData, bangDingDeviceSN)

                collectAnimateBindings((tempConfigData.components.cells[k].data && tempConfigData.components.cells[k].data.detail) ? tempConfigData.components.cells[k].data.detail.animate : null, bangDingData, bangDingDeviceSN)
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    ensurePageRegistered(ctx, pageid, page.page && page.page.displayUUID).then(function (pageInfo) {
        if (pageInfo) {
            selectPopUpDisplayPageDataStruct(ctx, page)
            return
        }
        ctx.state.selectPageUuid=""
        ctx.state.PopUpConfigData={ "name": "--", "layer": { "backColor": "", "backgroundImage": "", "widthHeightRatio": "", "width": 300, "height": 600 }, "components": [] }
        page.callback(-1,bangDingData,bangDingDeviceSN)
    })
}

export const updateLayerDataStruct = (ctx,layerData) => {

    // let PCPageInfo = ctx.state.PCPageList
    // let PhonePageInfo = ctx.state.PhonePageList
    //
    // if(layerData.PageType==1)
    // {
    //     for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
    //     {
    //         if(PCPageInfo[i].pageUuid == layerData.PageId)
    //         {
    //             ctx.state.PCPageList[i].children=[]
    //             for(let j=0,componentsLen=layerData.components.length;j<componentsLen;j++)
    //             {
    //                 let components = {
    //                     isComponents:true,
    //                     title:layerData.components[j].name,
    //                     key:layerData.components[j].identifier
    //                 }
    //                 ctx.state.PCPageList[i].children.push(components)
    //             }
    //         }
    //     }
    // }
    // else if(layerData.PageType==0)
    // {
    //     for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
    //     {
    //         if(PhonePageInfo[i].pageUuid == layerData.PageId) {
    //             ctx.state.PhonePageList[i].children = []
    //             for (let j = 0,componentsLen=layerData.components.length; j < componentsLen; j++) {
    //                 let components = {
    //                     isComponents: true,
    //                     title: layerData.components[j].name,
    //                     key: layerData.components[j].identifier
    //                 }
    //                 ctx.state.PhonePageList[i].children.push(components)
    //             }
    //         }
    //     }
    // }




}

export const selectParentLayerDataStruct = (ctx,childKey) => {
    let PCPageInfo = ctx.state.PCPageList
    let PhonePageInfo = ctx.state.PhonePageList

    for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
    {
        for(let j=0,childrenLen=PCPageInfo[i].children.length;j<childrenLen;j++)
        {
            if(PCPageInfo[i].children[j].key == childKey)
            {
                let tempConfigData = PCPageInfo[i].pageLayerData
                for(let k=0,componentsLen =tempConfigData.components.cells.length ;k<componentsLen;k++)
                {
                    if(typeof tempConfigData.components.cells[k].data.detail.animate!=="undefined"&&typeof tempConfigData.components.cells[k].data.detail.animate.move=="undefined")
                    {
                        tempConfigData.components.cells[k].data.detail.animate.move = {
                            x:{
                                deviceSN:"",
                                selectVideoType:0,
                                isBandDevice:false,
                                bandType:1,
                                dataID: "",
                                dataName: "",
                            },
                            y:{
                                deviceSN:"",
                                selectVideoType:0,
                                isBandDevice:false,
                                bandType:1,
                                dataID: "",
                                dataName: "",
                            },
                        }
                    }
                }
                ctx.state.selectPageUuid = tempConfigData.PageId
                ctx.state.LayerData = tempConfigData
                return
            }
        }
    }
    for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
    {
        for(let j=0,childrenLen=PhonePageInfo[i].children.length;j<childrenLen;j++)
        {
            if(PhonePageInfo[i].children[j].key == childKey)
            {
                let tempConfigData = PhonePageInfo[i].pageLayerData
                for(let k=0,componentsLen =tempConfigData.components.cells.length ;k<componentsLen;k++)
                {
                    if(typeof tempConfigData.components.cells[k].data.detail.animate!=="undefined"&&typeof tempConfigData.components.cells[k].data.detail.animate.move=="undefined")
                    {
                        tempConfigData.components.cells[k].data.detail.animate.move = {
                            x:{
                                deviceSN:"",
                                selectVideoType:0,
                                isBandDevice:false,
                                bandType:1,
                                dataID: "",
                                dataName: "",
                            },
                            y:{
                                deviceSN:"",
                                selectVideoType:0,
                                isBandDevice:false,
                                bandType:1,
                                dataID: "",
                                dataName: "",
                            },
                        }
                    }
                }
                ctx.state.selectPageUuid = tempConfigData.PageId
                ctx.state.LayerData = tempConfigData
                return
            }
        }
    }
}
export const SyncLayerData = (ctx,data) => {
    ctx.state.LayerData = data
}
export const setLayerData = (ctx,data) => {
    const safe = sanitizeGraphComponents(data, { tag: 'setLayerData' })
    ctx.state.ISMCavasContainer.fromJSON(safe)
}
export const SetEquidistantStateValue = (ctx,data) => {


}

export const setGroupList = (ctx) => {
    let tempGroupList=[]
    if(ctx.state.ISMCavasContainer==null){
        return
    }
    const cells = ctx.state.ISMCavasContainer.getCells()
    for (let k = 0,componentsLen=cells.length; k < componentsLen; k++) {
        if(typeof cells[k].data !="undefined" &&typeof cells[k].data.detail !="undefined")
        {
            let item=cells[k]
            if(typeof item!=="undefined")
            {
                let GroupObj={}
                GroupObj.Name=item.data.detail.name
                GroupObj.ID=item.id
                tempGroupList.push(GroupObj)
            }
        }
    }
    let GroupList=[]
    for(let key in tempGroupList) {
        GroupList.push(tempGroupList[key])
    }
    ctx.state.GroupList=GroupList
}
export const SyncLayerComponents = (ctx,data) => {
    // ctx.state.LayerData.components = ctx.state.ISMCavasContainer.toJSON()
    // console.log(ctx.state.LayerData.components)
}
export const  lockScreen = (ctx) => {
    ctx.state.isLocked = true
    localStorage.setItem("LockState",true)
}
export const  unLockScreen = (ctx) => {
    ctx.state.isLocked = false
    localStorage.setItem("LockState",false)
}
