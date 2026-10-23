/** 机房 4-1 机柜在舞台上的程序化位置。跟随镜头与柜体共用这一份。 */
export function rackWorld(col, row) {
  return {
    x: (col - 1.5) * 0.98,
    y: 1.05,
    z: row === 0 ? -0.9 : 0.9,
  }
}

export const ALARM_COL = 2
export const ALARM_ROW = 0
