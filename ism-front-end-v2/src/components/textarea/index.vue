<template>
  <div class="textarea-wrapper">
    <a-textarea class="m-textarea"
        v-bind="$attrs"
        v-model="$attrs.value"
        @change="onChange"
    />
    <span class="m-count" v-if="showWordLimit"
    >{{ textLength }}/<template v-if="$attrs.maxLength"
    >{{ $attrs.maxLength }}</template
    ></span
    >
  </div>
</template>
<script>

export default {
  props: {
    // 是否展示字数统
    showWordLimit: {
      type: Boolean,
      default: false,
    },
  },
  // v-model处理
  model: {
    prop: "value",
    event: "change",
  },
  computed: {
    // 长度控制
    textLength() {
      return (this.$attrs.value || "").length;
    },
  },
  methods: {
    onChange(e) {
      // v-model 回调函数
      this.$emit("change", e.target.value);
    },
  },
};

</script>


<style scoped lang="less">

.textarea-wrapper {
  position: relative;
  display: block;
  .m-textarea {
    padding: 8px 12px;
    height: 100%;
  }
  .m-count {
    color: #808080;
    background: #fff;
    position: absolute;
    font-size: 12px;
    bottom: 8px;
    right: 12px;
  }
}
</style>