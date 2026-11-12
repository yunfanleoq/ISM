import { AUTH_TYPE, getAuthorization, setAuthorization } from '@/utils/request'
import { ProjectList } from '@/services/project'

let pending = null

function uuidOf(row) {
  const info = row && (row.ProjectInfo || row)
  return (info && (info.uuid || info.Uuid)) || ''
}

/** 数据仓库请求发出前补上当前项目。只有一个项目时才写入，不写死某一套 uuid。 */
export function ensureProjectUuid() {
  const existing = getAuthorization(AUTH_TYPE.AUTH1)
  if (existing) return Promise.resolve(existing)
  if (pending) return pending
  pending = ProjectList().then(function (res) {
    const body = res && res.data
    const list = (body && body.code == 0 && body.list) ? body.list : []
    if (list.length !== 1) return ''
    const uuid = uuidOf(list[0])
    if (!uuid) return ''
    setAuthorization({ token: uuid }, AUTH_TYPE.AUTH1)
    return uuid
  }).catch(function () {
    return ''
  }).finally(function () {
    pending = null
  })
  return pending
}
