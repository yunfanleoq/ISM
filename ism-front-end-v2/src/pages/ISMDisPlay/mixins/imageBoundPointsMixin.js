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
    onBoundPointsEnter(evt) {
      if (!this.canShowBoundPoints || !this.$refs.boundPointsPopup) return
      this.$refs.boundPointsPopup.show({
        binds: this.boundPointList,
        clientX: evt.clientX,
        clientY: evt.clientY
      })
    },
    onBoundPointsMove(evt) {
      if (!this.canShowBoundPoints || !this.$refs.boundPointsPopup) return
      this.$refs.boundPointsPopup.move({
        clientX: evt.clientX,
        clientY: evt.clientY
      })
    },
    onBoundPointsLeave() {
      if (!this.$refs.boundPointsPopup) return
      this.$refs.boundPointsPopup.hide(false)
    },
    onBoundPointsClick(evt) {
      if (!this.canShowBoundPoints || !this.$refs.boundPointsPopup) return
      this.$refs.boundPointsPopup.togglePin({
        binds: this.boundPointList,
        clientX: evt.clientX,
        clientY: evt.clientY
      })
    }
  }
}
