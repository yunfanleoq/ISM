<template>
  <div class="hx-dcim-frame">
    <iframe title="新大屏" src="/hx-dcim/index.html#/power"></iframe>
  </div>
</template>

<script>
export default {
  name: 'HxDcimScreen'
}
</script>

<style scoped>
.hx-dcim-frame {
  width: 100%;
  height: calc(100vh - 64px);
  background: #040c1c;
}
.hx-dcim-frame iframe {
  width: 100%;
  height: 100%;
  border: 0;
  background: #040c1c;
}
</style>
