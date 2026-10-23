import { createBuilding3Fixture } from './fixture-building3'

/** 页面只读这一处。换夹具时改返回值，不改画面。 */
export function loadDemoSource() {
  return createBuilding3Fixture()
}
