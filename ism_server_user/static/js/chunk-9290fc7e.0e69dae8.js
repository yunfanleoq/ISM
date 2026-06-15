(window.webpackJsonp=window.webpackJsonp||[]).push([["chunk-9290fc7e"],{
/***/"58be":
/***/function(t,e,a){"use strict";
/* harmony import */a("74d64");
/* harmony import */},
/***/5908:
/***/function(t,e,a){
// extracted by mini-css-extract-plugin
/***/},
/***/5985:
/***/function(t,e,a){"use strict";
/* harmony import */a("5908");
/* harmony import */},
/***/"6d14":
/***/function(t,e,a){"use strict";
/* harmony import */a("cfe6");
/* harmony import */},
/***/"74d64":
/***/function(t,e,a){
// extracted by mini-css-extract-plugin
/***/},
/***/c8c3:
/***/function(t,e,a){"use strict";
// ESM COMPAT FLAG
a.r(e);
// EXTERNAL MODULE: ./node_modules/core-js/modules/es.function.name.js
a("b0c0");
// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"50a92b3c-vue-loader-template"}!./node_modules/cache-loader/dist/cjs.js??ref--13-0!./node_modules/thread-loader/dist/cjs.js!./node_modules/babel-loader/lib!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/cache-loader/dist/cjs.js??ref--1-0!./node_modules/vue-loader/lib??vue-loader-options!./src/layouts/PageView.vue?vue&type=template&id=38fc72ea&scoped=true
var i=a("5530"),s=(a("4de4"),a("caad"),a("14d9"),a("d3b7"),a("2532"),a("0643"),a("2382"),a("4e3e"),a("159b"),a("2f62")),n={name:"PageHeader",props:{title:{type:[String,Boolean],required:!1},breadcrumb:{type:Array,required:!1},logo:{type:String,required:!1},avatar:{type:String,required:!1}},computed:Object(i.a)({},Object(s.e)("setting",["layout","showPageTitle","pageWidth"])),mounted:function(){}},r=(a("6d14"),a("2877")),c=Object(r.a)(n,(function(){var t=this,e=t._self._c;return e("div",{class:["page-header",t.layout,t.pageWidth]},[e("div",{staticClass:"page-header-wide"},[e("div",{staticClass:"breadcrumb"},[e("a-breadcrumb",t._l(t.breadcrumb,(function(a,i){return e("a-breadcrumb-item",{key:i},[e("span",[t._v(t._s(a))])])})),1)],1),e("div",{staticClass:"detail"},[e("div",{staticClass:"main"},[e("div",{staticClass:"row"},[t.showPageTitle&&t.title?e("h1",{staticClass:"title"},[t._v(t._s(t.title))]):t._e(),e("div",{staticClass:"action"},[t._t("action")],2)]),e("div",{staticClass:"row"},[this.$slots.content?e("div",{staticClass:"content"},[t.avatar?e("div",{staticClass:"avatar"},[e("a-avatar",{attrs:{src:t.avatar,size:72}})],1):t._e(),t._t("content")],2):t._e(),this.$slots.extra?e("div",{staticClass:"extra"},[t._t("extra")],2):t._e()])])])])])}),[],!1,null,"1f76a225",null)
/* harmony default export */.exports,o=a("89a5"),u={name:"PageLayout",components:{PageHeader:c},props:["desc","logo","title","avatar","linkList","extraImage"],data:function(){return{page:{},pageHeaderHeight:0}},watch:{$route:function(){this.page=this.$route.meta.page}},updated:function(){this._inactive||this.updatePageHeight()},activated:function(){this.updatePageHeight()},deactivated:function(){this.updatePageHeight(0)},mounted:function(){this.updatePageHeight()},created:function(){this.page=this.$route.meta.page},beforeDestroy:function(){this.updatePageHeight(0)},computed:Object(i.a)(Object(i.a)({},Object(s.e)("setting",["layout","multiPage","pageMinHeight","pageWidth","customTitles"])),{},{pageTitle:function(){var t=this.page&&this.page.title;return this.customTitle||t&&this.$t(t)||this.title||this.routeName},routeName:function(){var t=this.$route;return this.$t(Object(o.b)(t.matched[t.matched.length-1].path))},breadcrumb:function(){var t=this,e=this.page,a=e&&e.breadcrumb;if(a){var i=[];return a.forEach((function(e){i.push(t.$t(e))})),i}return this.getRouteBreadcrumb()},marginCorrect:function(){return this.multiPage?24:0}}),methods:Object(i.a)(Object(i.a)({},Object(s.d)("setting",["correctPageMinHeight"])),{},{getRouteBreadcrumb:function(){var t=this,e=this.$route.matched,a=this.$route.path,i=[];e.filter((function(t){return a.includes(t.path)})).forEach((function(e){var a=0===e.path.length?"/home":e.path;i.push(t.$t(Object(o.b)(a)))}));var s=this.page&&this.page.title;return(this.customTitle||s)&&(i[i.length-1]=this.customTitle||s),i},
/**
     * 用于计算页面内容最小高度
     * @param newHeight
     */
updatePageHeight:function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:this.$refs.pageHeader.$el.offsetHeight+this.marginCorrect;this.correctPageMinHeight(this.pageHeaderHeight-t),this.pageHeaderHeight=t}})},l=u,d=(a("58be"),Object(r.a)(l,(function(){var t=this,e=t._self._c;return e("div",{staticClass:"page-layout"},[e("page-header",{ref:"pageHeader",style:"margin-top: ".concat(t.multiPage?0:-24,"px"),attrs:{breadcrumb:t.breadcrumb,title:t.pageTitle,logo:t.logo,avatar:t.avatar}},[t._t("action",null,{slot:"action"}),t._t("headerContent",null,{slot:"content"}),!this.$slots.headerContent&&t.desc?e("div",{attrs:{slot:"content"},slot:"content"},[e("p",[t._v(t._s(t.desc))]),this.linkList?e("div",{staticClass:"link"},[t._l(t.linkList,(function(a,i){return[e("a",{key:i,attrs:{href:a.href}},[e("a-icon",{attrs:{type:a.icon}}),t._v(t._s(a.title))],1)]}))],2):t._e()]):t._e(),this.$slots.extra?t._t("extra",null,{slot:"extra"}):t._e()],2),e("div",{ref:"page",class:["page-content",t.layout,t.pageWidth]},[t._t("default")],2)],1)}),[],!1,null,null,null)
/* harmony default export */.exports),g={name:"PageView",components:{PageToggleTransition:a("7664").a,PageLayout:d},data:function(){return{page:{}}},computed:Object(i.a)(Object(i.a)({},Object(s.e)("setting",["isMobile","multiPage","animate"])),{},{desc:function(){return this.page.desc},linkList:function(){return this.page.linkList},extraImage:function(){return this.page.extraImage}}),mounted:function(){this.page=this.$refs.page},updated:function(){this.page=this.$refs.page}},h=(a("5985"),Object(r.a)(g,(function(){var t=this,e=t._self._c;return e("page-layout",{attrs:{desc:t.desc,linkList:t.linkList}},[this.extraImage&&!t.isMobile?e("div",{staticClass:"extraImg",attrs:{slot:"extra"},slot:"extra"},[e("img",{attrs:{src:t.extraImage}})]):t._e(),e("page-toggle-transition",{attrs:{disabled:t.animate.disabled,animate:t.animate.name,direction:t.animate.direction}},[e("router-view",{ref:"page"})],1)],1)}),[],!1,null,"38fc72ea",null)
/* harmony default export */);e.default=h.exports},
/***/cfe6:
/***/function(t,e,a){
// extracted by mini-css-extract-plugin
/***/}}]);