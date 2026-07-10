<template>
  <AppRunShell :uid="$route.params.uid" />
</template>

<script>
/**
 * 路由入口：禁止 sync import 同目录任何 SFC（含 Loader），
 * 否则与 ISMRender/ISMRunTreeNav 链形成循环依赖 → export default 为 undefined。
 */
export default {
  name: 'AppRun',
  components: {
    AppRunShell: () => import(/* webpackChunkName: "app-run-shell" */ './AppRunShell.vue'),
  },
}
</script>
