import {
    PROJECTADD,PROJECTEDIT,PROJECTDEL,PROJECTLIST,EXPORTPROJECT,IMPORTPROJECT
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 项目添加
 */
export async function ProjectAdd(params) {
    return request(PROJECTADD, METHOD.POST,params)
}

/**
 * 项目编辑
 */
export async function ProjectEdit(params) {
    return request(PROJECTEDIT, METHOD.POST,params)
}

/**
 * 项目删除
 */
export async function ProjectDel(params) {
    return request(PROJECTDEL, METHOD.POST,params)
}
/**
 * 项目获取
 */
export async function ProjectList() {
    return request(PROJECTLIST, METHOD.POST)
}

/**
 * 项目导出
 */
export async function ExportProject(params) {
    return request(EXPORTPROJECT, METHOD.POST,params)
}
/**
 * 项目导入
 */
export async function ImportProject(params) {
    return request(IMPORTPROJECT, METHOD.POST,params)
}
/**
 * 项目更新
 */
export async function UpdateProject(params) {
    return request(IMPORTPROJECT, METHOD.POST,params)
}
export default {
    ProjectList,
    ProjectDel,
    ProjectEdit,
    ProjectAdd,
    ExportProject,
    UpdateProject,
    ImportProject
}
