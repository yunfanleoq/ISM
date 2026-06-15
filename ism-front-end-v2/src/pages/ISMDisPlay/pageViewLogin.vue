<template>
  <div class="ism-pageview">
    <ISMRenderLogin :showUuid="$route.params.uid" showDeviceUuid="" />
  </div>
</template>

<script>

import ISMRenderLogin from './ISMRenderLogin';

export default {
  name: 'AppRun',
  components: {
    ISMRenderLogin
  },
  methods: {

  },
  created() {

  },
  mounted() {

  }
}
</script>

<style lang="less">
.ism-pageview {
  height: 100vh;
  width: 100%;
}
</style>


