<template>
  <div v-show="detail.style.visible==1 ||isStart? true:false" :style="animatedStyle">
    <div :class="{
          'animated':true,[`${detail.style.animate}`]: true
        }"
         :style="{
                                width: detail.style.position.w + 'px',
                                height: detail.style.position.h + 'px',
                                'background-color': detail.style.backColor,
                                'border-radius':detail.style.BorderEdges+'px',
                                opacity:detail.style.opacity,
                                borderWidth: detail.style.borderWidth + 'px',
                                borderStyle: detail.style.borderStyle,
                                borderColor: detail.style.borderColor,
                                transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
                            }">
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"   xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">

      <image  :class="{'spin-element':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,'spin-element-reverse':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1,}" preserveAspectRatio="none meet"   :width="detail.style.position.w" :height="detail.style.position.h"   :href="imageURL"></image>
      <!--      闪烁-->
      <animate v-if="isStart&&animateType.includes('blink')&&!IsToolBox" attributeName="opacity"
               values="0.1;1;0.1" :dur="blinkSpeed+'s'"
               repeatCount="indefinite"/>
      <!--渐变-->
      <animate v-if="isStart&&animateType.includes('millcolorGrad')&&!IsToolBox" attributeName="fill"
               :values="startColor+';'+stopColor+';'+startColor" :dur="animateSpeed+'s'"
               repeatCount="indefinite"/>
      <!--缩放      -->
      <animateTransform v-if="isStart&&animateType.includes('Zoom')&&!IsToolBox" attributeName="transform"   begin="0s" dur="0.6s" type="scale" values="0.9;1;0.9" repeatCount="indefinite"/>
  </g>
</svg>
    </div>
  </div>

</template>

<script>

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'ism-view-png-image',
    inject: ['getNode'],
    props: {

    },
    computed: {
        imageURL: function () {
            if (this.detail.style.imageURL == undefined || this.detail.style.imageURL == '') {
                return '';
            } else {
                return this.detail.style.imageURL;
            }
        },
      animatedStyle(){
        return {
          "--blinkSpeed":this.blinkSpeed+'s',
          "--stopColor":this.stopColor,
          "--startColor":this.startColor,
          "--animateSpeed":this.animateSpeed+'s',
          "--animateSpinSpeed":this.animateSpinSpeed+'s'
        }
      },
    },
    components: {},
    data() {
        return {
          detail:null,
          IsToolBox:false,
          editMode:true,
          strokeColor:"#000000",
          fill:"#A1BFE2",
          strokeWidth:0.3,
          fillOpacity:1,
          strokeOpacity:1,
          animateType:"blink",
          startColor:"#74f808",
          stopColor:"#74f808",
          animateSpeed:0.5,
          animateSpinSpeed:0.5,
          spinDirection:0,
          blinkSpeed:0.5,
          isStart:false,
        }
    },
    watch: {
    detail: {
      handler(newVal, oldVal) {
        if(this.editMode) {
          this.initComponents(newVal);
        }
      },
      deep: true
    }
  },
  methods: {
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      let i=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="strokeWidth")
        {
          this.strokeWidth=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="strokeFill")
        {
          this.fill=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="strokeColor")
        {
          this.strokeColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="fillOpacity")
        {
          this.fillOpacity=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="strokeOpacity")
        {
          this.strokeOpacity=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="imageURL")
        {
          this.imageURL=option.style.diy[i].value
        }
      }
      i=0
      this.animateType = option.animate.selected
      if(option.animate.isExpression)
      {
        this.isStart = false
      }
      else
      {
        this.isStart = true
      }
      for( i=0;i<option.animate.animateElement.length;i++)
      {
        if(option.animate.animateElement[i].id=="millcolorGrad")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++)
          {
            if(option.animate.animateElement[i].elementList[k].key=="startColor")
            {
              this.startColor=option.animate.animateElement[i].elementList[k].value
            }
            else if(option.animate.animateElement[i].elementList[k].key=="stopColor")
            {
              this.stopColor=option.animate.animateElement[i].elementList[k].value
            }
            else if(option.animate.animateElement[i].elementList[k].key=="animateSpeed")
            {
              this.animateSpeed=option.animate.animateElement[i].elementList[k].value
            }
          }
        }
        else if(option.animate.animateElement[i].id=="blink")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
            if (option.animate.animateElement[i].elementList[k].key == "blinkSpeed") {
              this.blinkSpeed = option.animate.animateElement[i].elementList[k].value
            }
          }
        }
        else if(option.animate.animateElement[i].id=="animateSpin")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
            if (option.animate.animateElement[i].elementList[k].key == "spinSpeed") {
              this.animateSpinSpeed = option.animate.animateElement[i].elementList[k].value
            }
            else if (option.animate.animateElement[i].elementList[k].key == "spinDirection") {
              this.spinDirection = option.animate.animateElement[i].elementList[k].value
            }
          }
        }
      }
    }
  },
  mounted() {
    let _t = this

  },
  created(){
    let _t = this
    this.GetNodeObj = this.getNode()
    this.GetNodeObj.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    this.GetNodeObj.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
    });
    this.detail = this.GetNodeObj.getData().detail
    this.editMode = this.GetNodeObj.getData().editMode
    this.showDeviceUuid = this.GetNodeObj.getData().showDeviceUuid
    this.IsToolBox = this.GetNodeObj.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
	  // _t.initComponents(_t.detail)
    })
    this.$nextTick(function(){
      this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {


        })
        _t.$EventBus.$on(animateEvent, (data) => {
          _t.isStart = data
        })
    });
  }
}
</script>

<style lang="less">
.spin-element {
  backface-visibility: hidden;
  perspective: 1000px; /* 配合 backface-visibility 使用，增强效果 */
  transform: translateZ(0);
  will-change: transform;
  animation: spin linear infinite;
  animation-duration: var(--animateSpinSpeed);
  transform-origin: 50% 50%;
  animation-delay:0.5s;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.spin-element-reverse {
  backface-visibility: hidden;
  perspective: 1000px; /* 配合 backface-visibility 使用，增强效果 */
  transform: translateZ(0);
  will-change: transform;
  animation: spin-reverse linear infinite;
  animation-duration: var(--animateSpinSpeed);
  transform-origin: 50% 50%;
  animation-delay:0.5s;
}

@keyframes spin-reverse {
  from { transform: rotate(0deg); }
  to { transform: rotate(-360deg); /* 负角度=逆时针 */ }
}
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
.view-image {
    height: 100%;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
}
.animated-svg {
  will-change: transform; /* 告诉浏览器该元素即将动画，提前准备优化 */
  transform: translateZ(0); /* 强制创建合成层（老浏览器兼容） */
}
</style>
