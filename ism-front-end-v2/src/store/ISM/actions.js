import {getDisplayModelLayerData,setDisplayModelLayerData,getLayerDataStructByToken} from "@/services/displayModel";
import {GetDisplayLoginPage} from "@/services/system";
import { uuid } from 'vue-uuid';
import {getDisplayModelPagerLayerData} from "../../services/displayModel";
import {normalizeISMScene} from "@/pages/ISMDisPlay/utils/ismSceneNormalizer";

function pascalToKebab(str) {
    if (!str) return ''
    // DvBorderBox1 → dv-border-box1, DeviceTree → device-tree
    return str
        .replace(/([A-Z])([A-Z][a-z])/g, '$1-$2')
        .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
        .replace(/([A-Z]+)([A-Z][a-z])/g, '$1-$2')
        .toLowerCase()
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

export const getLayerDataStruct = (ctx,data) => {
    let params={
        muid:data.uuid
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    let isPopUp = data.isPopUp?data.isPopUp:false
    getDisplayModelLayerData(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer
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
                    pageInfo.children=[]

                    pageLayer[i].name = pageLayer[i].PageName
                    let tempConfigData = pageLayer[i]
                    try{
                        tempConfigData.layer = JSON.parse(tempConfigData.layer)
                        if (tempConfigData.components=="")
                        {
                            tempConfigData.components=[]
                        }
                        else{
                            tempConfigData.components = JSON.parse(tempConfigData.components)
                        }
                        tempConfigData = normalizePageConfigData(tempConfigData)

                    }catch (e) {
                        console.error('[getLayerDataStruct] parse/normalize error:', e.message || e, 'stack:', (e.stack || '').slice(0,400))
                        tempConfigData.components.cells=[]
                        continue
                    }
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
                        let components = {
                            isComponents:true,
                            title:tempConfigData.components.cells[k].data.detail.name,
                            key:tempConfigData.components.cells[k].id,
                            cellid:tempConfigData.components.cells[k].data.detail.identifier
                        }
                        pageInfo.children.push(components)
                    }

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
                                    if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                                        for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                                bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                            }
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                            }
                                        }
                                    }

                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                                    }

                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                                        if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                                    }
                                    else{
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
                                    if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                                        for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                                bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                            }
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                            }
                                        }
                                    }

                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                                    }
                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!=="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                                        if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                                    }
                                    else{
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
                    let tempConfigData = pageLayer[0]
                    ctx.state.selectPageUuid = pageLayer[0].PageId
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
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                            for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                }
                            }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
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
            data.cb(0,res.data.Display.project_uuid,bangDingData,newbangDingDeviceSN)
        }
        else
        {
            data.cb(-1,"",bangDingData,bangDingDeviceSN)
        }

    })
}
export const updateAllLayerDataStruct = (ctx,data) => {
    let params={
        muid:data.uuid
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    let isPopUp = data.isPopUp?data.isPopUp:false
    getDisplayModelLayerData(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer
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
                    pageInfo.children=[]

                    pageLayer[i].name = pageLayer[i].PageName
                    let tempConfigData = pageLayer[i]
                    try{
                        tempConfigData.layer = JSON.parse(tempConfigData.layer)
                        if (tempConfigData.components=="")
                        {
                            tempConfigData.components=[]
                        }
                        else{
                            tempConfigData.components = JSON.parse(tempConfigData.components)
                        }
                        tempConfigData = normalizePageConfigData(tempConfigData)

                    }catch (e) {
                        console.error('[getLayerDataStruct] parse/normalize error:', e.message || e, 'stack:', (e.stack || '').slice(0,400))
                        tempConfigData.components.cells=[]
                        continue
                    }
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
                        let components = {
                            isComponents:true,
                            title:tempConfigData.components.cells[k].data.detail.name,
                            key:tempConfigData.components.cells[k].id,
                            cellid:tempConfigData.components.cells[k].data.detail.identifier
                        }
                        pageInfo.children.push(components)
                    }
                    pageInfo.pageLayerData = tempConfigData
                    if(pageLayer[i].PageType==1)
                    {
                        pcPageData.push(pageInfo)
                    }else{
                        phonePageData.push(pageInfo)
                    }
                    if(pageInfo.id==ctx.state.selectPageUuid)
                    {
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
                tempConfigData.layer = JSON.parse(tempConfigData.layer)
                if (tempConfigData.components=="")
                {
                    tempConfigData.components=[]
                }
                else{
                    tempConfigData.components = JSON.parse(tempConfigData.components)
                }
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
                    if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                        for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                            }
                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                            }
                        }
                    }

                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                    {
                        if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                        }
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                    }
                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                    {
                        if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                        }
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                        if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                        }
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                    }
                    else{
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

    let pageid = page.page.pageUuid
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
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsAlen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsAlen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsAlen= tempConfigData.components.cells[k].data.detail.active.length; kv <componentsAlen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    getLayerDataStruct(ctx,{uuid:page.page.displayUUID,isPopUp:true,cb:function () {
            let PCPageInfo = ctx.state.PCPageList
            let PhonePageInfo = ctx.state.PhonePageList

            let pageid = page.page.pageUuid
            let bangDingData = []
            let bangDingDeviceSN = []
            for (let i = 0, PCPageInfoLen = PCPageInfo.length; i < PCPageInfoLen; i++) {
                if (PCPageInfo[i].pageUuid == pageid) {
                    let tempConfigData = PCPageInfo[i].pageLayerData
                    tempConfigData.name = page.page.title
                    for (let k = 0, componentsLen = tempConfigData.components.cells.length; k < componentsLen; k++) {
                        if (typeof tempConfigData.components.cells[k].data.detail.active != "undefined") {
                            for (let kv = 0, componentsAlen = tempConfigData.components.cells[k].data.detail.active.length; kv < componentsAlen; kv++) {
                                if (tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID != "") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                                if (tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN != "") {
                                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                }
                            }
                        }

                        if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && tempConfigData.components.cells[k].data.detail.animate.condition.dataID != "") {
                            if (tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN != "") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && typeof tempConfigData.components.cells[k].data.detail.animate.move != "undefined") {
                            if (tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN != "") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if (tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN != "") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        } else {
                            if (typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") {
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
                        if (typeof tempConfigData.components.cells[k].data.detail.active != "undefined") {
                            for (let kv = 0, componentsAlen = tempConfigData.components.cells[k].data.detail.active.length; kv < componentsAlen; kv++) {
                                if (tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID != "") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                            }
                        }

                        if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && tempConfigData.components.cells[k].data.detail.animate.condition.dataID != "") {
                            if (tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN != "") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && typeof tempConfigData.components.cells[k].data.detail.animate.move != "undefined") {
                            if (tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN != "") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if (tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN != "") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        } else {
                            if (typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") {
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

    let pageid = page.page.pageUuid
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
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    getLayerDataStruct(ctx,{uuid:page.page.displayUUID,cb:function (){
            let PCPageInfo = ctx.state.PCPageList
            let PhonePageInfo = ctx.state.PhonePageList

            let pageid = page.page.pageUuid
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
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                            for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                }
                            }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
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
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                            for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                }
                            }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
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
        token:data.token
    }
    let bangDingData=[]
    let bangDingDeviceSN=[]
    let isPopUp = data.isPopUp?data.isPopUp:false
    getLayerDataStructByToken(params).then(function (res){
        if(res.data.code==0)
        {
            let pageLayer = res.data.layer
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
                    pageInfo.children=[]

                    pageLayer[i].name = pageLayer[i].PageName
                    let tempConfigData = pageLayer[i]
                    try{
                        tempConfigData.layer = JSON.parse(tempConfigData.layer)
                        if (tempConfigData.components=="")
                        {
                            tempConfigData.components=[]
                        }
                        else{
                            tempConfigData.components = JSON.parse(tempConfigData.components)
                        }
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
                                    if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                                        for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                                bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                            }
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                            }
                                        }
                                    }

                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                                    }
                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                                        if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                                    }
                                    else{
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
                                    if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                                        for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                                bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                            }
                                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                            }
                                        }
                                    }

                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                                    }
                                    if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                                    {
                                        if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                                        if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                                        }
                                        bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                                    }
                                    else{
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
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                            for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                }
                            }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
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
                    tempConfigData.layer = JSON.parse(tempConfigData.layer)
                    if (tempConfigData.components=="")
                    {
                        tempConfigData.components=[]
                    }
                    else{
                        tempConfigData.components = JSON.parse(tempConfigData.components)
                    }
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
    let params={
        muid:page.uuid,
        pageid:page.pageid,
        saveData:page.LayerData
    }
    return setDisplayModelLayerData(params).then(function (res){
        return res
    })
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
                let tempConfigData = PCPageInfo[i].pageLayerData
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
        }
    }else if(pagetype==0)
    {
        for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
        {
            if(PhonePageInfo[i].pageUuid==pageid) {
                let tempConfigData = PhonePageInfo[i].pageLayerData
                tempConfigData.name = page.title
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
            }
        }
    }
}

export const selectDisplayPageDataStructFromDb = (ctx,page) => {
    getLayerDataStruct(ctx,{uuid:page.page.displayUUID,cb:function (){
            let PCPageInfo = ctx.state.PCPageList
            let PhonePageInfo = ctx.state.PhonePageList

            let pageid = page.page.pageUuid
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
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                            for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                }
                            }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
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
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                            for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                                }
                                if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                                }
                            }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                            if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                            if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                                bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                            }
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
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

    let pageid = page.page.pageUuid
    let bangDingData=[]
    let bangDingDeviceSN=[]
    if(!pageid)
    {
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,cb:function (errno, project_uuid, uuids, devices){
                page.callback(uuids,devices,errno == 0)
        }});
        return
    }
    if(!isDisplayPagesLoaded(ctx, page.page.displayUUID))
    {
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,cb:function (errno){
                if(errno == 0 && isDisplayPagesLoaded(ctx, page.page.displayUUID))
                {
                    selectDisplayPageDataStruct(ctx,page)
                    return
                }
                page.callback(bangDingData,bangDingDeviceSN,false)
        }});
        return
    }
    for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
    {
        if(PCPageInfo[i].pageUuid==pageid)
        {
            let tempConfigData = PCPageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            ctx.state.selectPageUuid = tempConfigData.PageId
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
            }
            ctx.state.LayerData = tempConfigData
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(bangDingData,newbangDingDeviceSN)
            return
        }
    }

    ctx.state.selectPageUuid=""
    page.callback(bangDingData,bangDingDeviceSN,false)
}

export const selectPopUpDisplayPageDataStruct = (ctx,page) => {
    let PCPageInfo = ctx.state.PCPageList
    let PhonePageInfo = ctx.state.PhonePageList

    let pageid = page.page.pageUuid
    let bangDingData=[]
    let bangDingDeviceSN=[]
    if(!pageid)
    {
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,isPopUp:true,cb:function (errno, project_uuid, uuids, devices){
                page.callback(errno == 0 ? 0 : -1,uuids,devices)
        }});
        return
    }
    if(!isDisplayPagesLoaded(ctx, page.page.displayUUID))
    {
        getLayerDataStruct(ctx,{uuid:page.page.displayUUID,isPopUp:true,cb:function (errno){
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
            let tempConfigData = PCPageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            ctx.state.selectPageUuid = tempConfigData.PageId
            ctx.state.PopUpConfigData = tempConfigData
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsAlen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsAlen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
    {
        if(PhonePageInfo[i].pageUuid==pageid) {
            let tempConfigData = PhonePageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            ctx.state.selectPageUuid = tempConfigData.PageId
            ctx.state.PopUpConfigData = tempConfigData
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsAlen= tempConfigData.components.cells[k].data.detail.active.length; kv <componentsAlen; kv++) {
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                            bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                        }
                        if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                        }
                    }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
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
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,bangDingData,newbangDingDeviceSN)
            return
        }
    }

    ctx.state.selectPageUuid=""
    ctx.state.PopUpConfigData={ "name": "--", "layer": { "backColor": "", "backgroundImage": "", "widthHeightRatio": "", "width": 300, "height": 600 }, "components": [] }
    page.callback(-1,bangDingData,bangDingDeviceSN)
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
    if (data && data.cells && Array.isArray(data.cells)) {
      data = { ...data, cells: data.cells.filter(cell => cell && cell.shape) }
    }
    ctx.state.ISMCavasContainer.fromJSON(data)
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
