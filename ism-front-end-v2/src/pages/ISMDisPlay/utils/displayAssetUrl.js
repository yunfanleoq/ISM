/**
 * 组态图片 URL：编辑态预览与运行态同一解析。
 * 自定义图走后端 static/customPel；系统图仍走前端 public/static/ISM。
 */
export function resolveDisplayAssetUrl(raw) {
  if (!raw || typeof raw !== 'string') {
    return ''
  }
  const url = raw.trim()
  if (!url) {
    return ''
  }
  if (/^(https?:)?\/\//i.test(url) || url.indexOf('data:') === 0 || url.indexOf('blob:') === 0) {
    return url
  }
  if (url.indexOf('/api/') === 0 || url.indexOf('api/') === 0) {
    return url.indexOf('/') === 0 ? url : '/' + url
  }
  let path = url.indexOf('/') === 0 ? url : '/' + url
  const isBackendStatic = path.indexOf('/static/customPel/') === 0 ||
    path.indexOf('/static/upload/') === 0 ||
    path.indexOf('/static/DiyUpload/') === 0
  if (isBackendStatic) {
    return '/api' + path
  }
  return path
}
