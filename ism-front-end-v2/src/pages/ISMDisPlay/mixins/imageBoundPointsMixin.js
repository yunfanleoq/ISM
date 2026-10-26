import BoundPointsPopup from '../BoundPointsPopup.vue'

export default {
  components: {
    BoundPointsPopup
  },
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
      if (!this.canShowBoundPoints || !this.$refs.boundPointsPopup) return
      this.$refs.boundPointsPopup.cancelSuppressClear()
      this.$refs.boundPointsPopup.show({
        binds: this.boundPointList,
        clientX: evt.clientX,
        clientY: evt.clientY,
        ownerKey: this.boundPointsOwnerKey()
      })
    },
    onBoundPointsMove(evt) {
      if (!this.canShowBoundPoints || !this.$refs.boundPointsPopup) return
      this.$refs.boundPointsPopup.move({
        clientX: evt.clientX,
        clientY: evt.clientY
      })
    },
    onBoundPointsLeave(evt) {
      const popup = this.$refs.boundPointsPopup
      if (!popup) return
      const to = evt && evt.relatedTarget
      if (to && popup.$el && popup.$el.contains(to)) return
      popup.hide(false)
      popup.armSuppressClear()
    },
    onBoundPointsClick(evt) {
      if (!this.canShowBoundPoints || !this.$refs.boundPointsPopup) return
      this.$refs.boundPointsPopup.togglePin({
        binds: this.boundPointList,
        clientX: evt.clientX,
        clientY: evt.clientY,
        ownerKey: this.boundPointsOwnerKey()
      })
    }
  }
}
