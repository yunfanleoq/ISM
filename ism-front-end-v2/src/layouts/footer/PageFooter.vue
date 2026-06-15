<template>
  <div class="footer">
    <div class="links">
      <a target="_blank" :key="index" :href="item.link ? item.link : 'javascript: void(0)'" v-for="(item, index) in linkList">
        <a-icon v-if="item.icon" :type="item.icon"/>{{item.name}}
      </a>
    </div>
    <div class="copyright">
      {{ $t('copyright.Version') }}:{{version}} Operated by <a :href="'http://'+systemUrl" target="_blank">{{systemCompany}}</a>
    </div>
    <div class="copyright">
      Copyright<a-icon type="copyright" />{{copyright}}
    </div>
  </div>
</template>

<script>
export default {
  name: 'PageFooter',
  i18n: require('@/i18n/language'),
  props: ['copyright', 'linkList','version'],
  computed: {
    systemUrl () {
      return this.$store.state.setting.systemUrl
    },
    systemCompany () {
      return this.$store.state.setting.systemCompany
    },
  },
}
</script>
<style lang="less" scoped>
.footer {
  padding: 22px 40px 18px; /* 左右加宽内边距，横向更舒展 */
  text-align: center;
  background-color: #fff;
  border-top: 1px solid #eee;
  box-shadow: 0 -1px 3px rgba(0, 0, 0, 0.02); /* 极浅阴影，增加层次感不压抑 */

  .links {
    margin-bottom: 15px;

    a {
      color: #2d3748; /* 深灰增强对比度，更显沉稳 */
      font-size: 15px;
      font-weight: 500; /* 加粗提升文字分量 */
      margin: 0 22px; /* 加宽链接间距，横向更大气 */
      padding: 4px 0;
      text-decoration: none;
      transition: color 0.3s ease;

      &:hover {
        color: #165dff; /* 高质感主色，hover反馈更清晰 */
      }

      .anticon {
        margin-right: 6px;
        vertical-align: middle;
        font-size: 16px; /* 图标稍大，增强视觉平衡 */
      }
    }
  }

  .copyright {
    color: #718096;
    font-size: 14px;
    line-height: 1.6;
    margin: 6px 0;
    letter-spacing: 0.4px; /* 轻微加宽字距，提升精致感 */

    a {
      color: #4a5568;
      font-weight: 500;
      text-decoration: none;
      transition: color 0.3s ease;

      &:hover {
        color: #13c2c2;
      }
    }

    .anticon {
      margin: 0 6px;
      vertical-align: middle;
      color: #718096;
    }
  }

  .copyright:last-child {
    margin-top: 10px;
    font-size: 12px;
    color: #a0aec0;
  }
}
</style>