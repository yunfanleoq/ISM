import { getBoundPointsPopup } from './boundPointsHost'

export default {
  computed: {
    boundPointList() {
      const binds = this.detail && Array.isArray(this.detail.dataBind) ? this.detail.dataBind : []
      return binds.filter((item) => item && (item.dataID || item.dataName))
    },
    canShowBoundPoints() {
      return !this.editMode && !this.IsToolBox && this.boundPointList.length > 0
    }
  },
  methods: {
    boundPointsOwnerKey() {
      const detail = this.detail || {}
      return String(detail.identifier || detail.id || detail.uuid || this._uid)
    },
    onBoundPointsEnter(evt) {
      if (!this.canShowBoundPoints) return
      const popup = getBoundPointsPopup(this)
      if (!popup) return
      popup.show({
        binds: this.boundPointList,
        clientX: evt.clientX,
        clientY: evt.clientY,
        ownerKey: this.boundPointsOwnerKey()
      })
    },
    onBoundPointsMove() {
      // 卡片出现后保持原位，方便点击钉住或关闭。
    },
    onBoundPointsLeave() {
      const popup = getBoundPointsPopup(this, false)
      if (!popup) return
      popup.clearSuppress(this.boundPointsOwnerKey())
      popup.hide(false)
    },
    onBoundPointsClick(evt) {
      if (!this.canShowBoundPoints) return
      const popup = getBoundPointsPopup(this)
      if (!popup) return
      popup.togglePin({
        binds: this.boundPointList,
        clientX: evt && evt.clientX,
        clientY: evt && evt.clientY,
        ownerKey: this.boundPointsOwnerKey()
      })
    }
  }
}
