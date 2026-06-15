<template>
  <div :style="{'height':height+'px','width':width+'px'}">
    <iframe src="https://www.yuque.com/jiahangseng-qsywc/zzwr7r"
            width="100%"
            height="100%"
            frameborder="0"
            scrolling="yes"> </iframe>
  </div>
</template>

<script>
export default {
  name: "help",
  i18n: require('../../i18n/language'),
  data() {
    return {
     width:1024,
      height:768
    }
  },
  mounted() {
    this.width = window.innerWidth-260
    this.height = window.innerHeight
  },
  methods: {
    handleAnchorClick(e,link){
      // 阻止点击的默认事件修改路由
      e.preventDefault();
      var srcolls = document.getElementById(link.href);
      srcolls.scrollIntoView({block: 'start', behavior: 'smooth'});
    },
  }
}
</script>

<style lang="less">
.project-list {

.card-title {
span{
  vertical-align: middle;
&:last-child{
   margin-left: 12px;
 }
}
}
.project-item {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  overflow: hidden;
  font-size: 12px;
  color: inherit;
.group{
  color: @text-color;
  flex: 1 1 0;
&:hover {
   color: @primary-color;
 }
}
.datetime {
  color: @text-color-second;
  flex: 0 0 auto;
}
}
.ant-card-meta-description {
  height: 44px;
  line-height: 22px;
  overflow: hidden;
}
}
.item-group{
  padding: 20px 0 8px 24px;
  font-size: 0;
a{
  color: inherit;
  display: inline-block;
  font-size: 14px;
  margin-bottom: 13px;
  width: 25%;
}
}
.members {
a {
  display: block;
  margin: 12px 0;
  color: @text-color;
&:hover {
   color: @primary-color;
 }
.member {
  vertical-align: middle;
  margin-left: 12px;
}
}
}
.markdown h1 {
  color: #0d1a26;
  font-weight: 500;
  margin-bottom: 20px;
  margin-top: 8px;
  font-family: Avenir,-apple-system,BlinkMacSystemFont,Segoe UI,PingFang SC,Hiragino Sans GB,Microsoft YaHei,Helvetica Neue,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol;
  font-size: 30px;
  font-variant: tabular-nums;
  line-height: 38px;
}
.markdown p, .markdown pre {
  margin: 1em 0;
  text-align: left;
}
.main-container {
  padding: 0 200px 144px 64px;
  min-height: 500px;
  overflow: hidden;
}
.main-container, .main-wrapper {
  background: #fff;
  position: relative;
}
.toc-affix {
  width: 150px;
  position: fixed;
  top: 100px;
  right: 10px;
}
</style>
