<template>
    <div class="mark-line" :style="styleVar">
        <div
            v-for="line in lines"
            :key="line"
            class="show-maker-line"
            :class="line.includes('x') ? 'show-xline' : 'show-yline'"
            :ref="line"
            v-show="lineStatus[line] || false"
        ></div>
    </div>
</template>

<script>
import { mapActions, mapState, mapMutations } from 'vuex'
import { getComponentRotatedStyle } from '@/utils/style';
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: "ISMMarkLine",
  i18n: require('@/i18n/language'),
  data() {
      return {
          lines: ['xt', 'xc', 'xb', 'yl', 'yc', 'yr'], // 分别对应三条横线和三条竖线
          diff: 10, // 相距 dff 像素将自动吸附
          lineStatus: {
              xt: false,
              xc: false,
              xb: false,
              yl: false,
              yc: false,
              yr: false,
          },
      };
  },
  props: {
    IsShowMakerLine: {
      type: Boolean,
      default: false
  }
},
  computed: {
    ...mapState({
          curComponent: state => store.state.ISMDisPlayEditorTool.selectedComponent,
          componentData: state => store.state.ISMDisPlayEditorTool.LayerData,
          selectedValue:state => store.state.ISMDisPlayEditorTool.selectedValue,
          selectedComponentMap: state => store.state.ISMDisPlayEditorTool.selectedComponentMap,
        }
    ),
    styleVar() {
      return {
        "--lineHeight": this.componentData.layer.height+"px" ,
        "--lineWidth": this.componentData.layer.width+'px' ,
      };
    }
  },
  mounted() {
      // 监听元素移动和不移动的事件
    this.$EventBus.$on('DragMoveLine', (isDownward, isRightward) => {
          if(!this.IsShowMakerLine)
          {
            this.hideLine();
          }
          else {
            this.showLine(isDownward, isRightward, "drag", false);
          }
      });
    this.$EventBus.$on('KeyMoveLine', (isDownward, isRightward) => {
      if(!this.IsShowMakerLine)
      {
        this.hideLine();
      }
      else {
        this.showLine(isDownward, isRightward, "key", false);
      }
    });

    this.$EventBus.$on('unmove', () => {
          if(!this.IsShowMakerLine)
          {
            this.hideLine();
          }
          else {
            this.showLine(false, false, "key", true);
            this.hideLine();
          }
      });
  },
  methods: {
      ...mapMutations('ISMDisPlayEditorTool',[
        'execute',
      ]),
        hideLine() {
            Object.keys(this.lineStatus).forEach((line) => {
                this.lineStatus[line] = false;
            });
        },

        showLine(isDownward, isRightward,type,isover) {
            const lines = this.$refs;
            const zoom = this.selectedValue/100
            if(zoom!=1)
            {
              return
            }
            // if(this.curComponent.style.transform!=0)
            // {
            //   return
            // }
            const components = this.componentData.components;
            if((this.curComponent==null)||(this.curComponent==undefined)||(this.curComponent.style==undefined))
            {
              return
            }
            const curComponentStyle = getComponentRotatedStyle(this.curComponent.style,this.selectedValue);
            const curComponentHalfwidth = curComponentStyle.position.w / 2;
            const curComponentHalfHeight = curComponentStyle.position.h / 2;
            this.hideLine();
            components.forEach((component) => {
                if (component.identifier == this.curComponent.identifier) {
                    return;
                }
                const componentStyle = getComponentRotatedStyle(component.style,this.selectedValue);
                const { y, x, b, r } = componentStyle.position;
                const componentHalfwidth = parseInt(componentStyle.position.w / 2);
                const componentHalfHeight = parseInt(componentStyle.position.h / 2);

                const conditions = {
                    top: [
                        {
                            isNearly: this.isNearly(curComponentStyle.position.y, y),
                            lineNode: lines.xt[0], // xt
                            line: 'xt',
                            dragShift: parseInt(y),
                            lineShift: parseInt(y)+27,
                        },
                        {
                            isNearly: this.isNearly(curComponentStyle.position.b, y),
                            lineNode: lines.xt[0], // xt
                            line: 'xt',
                            dragShift: parseInt(y) - parseInt(curComponentStyle.position.h),
                            lineShift: parseInt(y)+27,
                        },
                        {
                            // 组件与拖拽节点的中间是否对齐
                            isNearly: this.isNearly(
                                parseInt(curComponentStyle.position.y) + curComponentHalfHeight,
                                y + componentHalfHeight
                            ),
                            lineNode: lines.xc[0], // xc
                            line: 'xc',
                            dragShift: parseInt(y) + componentHalfHeight - curComponentHalfHeight,
                            lineShift: y + componentHalfHeight+27,
                        },
                        {
                            isNearly: this.isNearly(parseInt(curComponentStyle.position.y), b),
                            lineNode: lines.xb[0], // xb
                            line: 'xb',
                            dragShift: b,
                            lineShift: b+27,
                        },
                        {
                            isNearly: this.isNearly(parseInt(curComponentStyle.position.b), b),
                            lineNode: lines.xb[0], // xb
                            line: 'xb',
                            dragShift: b - curComponentStyle.position.h,
                            lineShift: b+27,
                        },
                    ],
                    left: [
                        {
                            isNearly: this.isNearly(curComponentStyle.position.x, x),
                            lineNode: lines.yl[0], // yl
                            line: 'yl',
                            dragShift: x,
                            lineShift: x+27,
                        },
                        {
                            isNearly: this.isNearly(curComponentStyle.position.r, x),
                            lineNode: lines.yl[0], // yl
                            line: 'yl',
                            dragShift: x - curComponentStyle.position.w,
                            lineShift: x+27,
                        },
                        {
                            // 组件与拖拽节点的中间是否对齐
                            isNearly: this.isNearly(
                                curComponentStyle.position.x + curComponentHalfwidth,
                                x + componentHalfwidth
                            ),
                            lineNode: lines.yc[0], // yc
                            line: 'yc',
                            dragShift: x + componentHalfwidth - curComponentHalfwidth,
                            lineShift: x + componentHalfwidth+27,
                        },
                        {
                            isNearly: this.isNearly(curComponentStyle.position.x, r),
                            lineNode: lines.yr[0], // yr
                            line: 'yr',
                            dragShift: r,
                            lineShift: r+27,
                        },
                        {
                            isNearly: this.isNearly(curComponentStyle.position.r, r),
                            lineNode: lines.yr[0], // yr
                            line: 'yr',
                            dragShift: r - curComponentStyle.position.w,
                            lineShift: r+27,
                        },
                    ],
                };

                const needToShow = [];
                const  rotate = this.curComponent.style.transform;
                Object.keys(conditions).forEach((key) => {
                    // 遍历符合的条件并处理
                    conditions[key].forEach((condition) => {
                        if (!condition.isNearly) {
                            return;
                        }
                        if(type=="drag"||isover)
                        {
                          this.execute({
                            op: 'adsorption',
                            key: key,
                            value:  rotate != 0 ? this.translatecurComponentShift(key, condition, curComponentStyle) : parseInt(condition.dragShift),
                            items: this.selectedComponentMap
                          });
                        }
                        condition.lineNode.style[key] = `${condition.lineShift}px`;
                        needToShow.push(condition.line);
                    });
                });

                if (needToShow.length) {
                    this.chooseTheTureLine(needToShow, isDownward, isRightward);
                }
            });
        },

        translatecurComponentShift(key, condition, curComponentStyle) {
            const { w, h } = this.curComponent.style.position;
            if (key == 'top') {
                return Math.round(parseInt(condition.dragShift) - (parseInt(h) - parseInt(curComponentStyle.style.position.h)) / 2);
            }

            return Math.round(parseInt(condition.dragShift) - (parseInt(w) - parseInt(curComponentStyle.style.position.w)) / 2);
        },

        chooseTheTureLine(needToShow, isDownward, isRightward) {
            if (isRightward) {
                if (needToShow.includes('yr')) {
                    this.lineStatus.yr = true;
                } else if (needToShow.includes('yc')) {
                    this.lineStatus.yc = true;
                } else if (needToShow.includes('yl')) {
                    this.lineStatus.yl = true;
                }
            } else {
                if (needToShow.includes('yl')) {
                    this.lineStatus.yl = true;
                } else if (needToShow.includes('yc')) {
                    this.lineStatus.yc = true;
                } else if (needToShow.includes('yr')) {
                    this.lineStatus.yr = true;
                }
            }

            if (isDownward) {
                if (needToShow.includes('xb')) {
                    this.lineStatus.xb = true;
                } else if (needToShow.includes('xc')) {
                    this.lineStatus.xc = true;
                } else if (needToShow.includes('xt')) {
                    this.lineStatus.xt = true;
                }
            } else {
                if (needToShow.includes('xt')) {
                    this.lineStatus.xt = true;
                } else if (needToShow.includes('xc')) {
                    this.lineStatus.xc = true;
                } else if (needToShow.includes('xb')) {
                    this.lineStatus.xb = true;
                }
            }
        },

        isNearly(dragValue, targetValue) {
            return Math.abs(parseInt(dragValue) - parseInt(targetValue)) <= this.diff;
        },
    },
};
</script>

<style lang="less" scoped>
.mark-line {
  height: 100%;
}
.show-maker-line {
  background: #59c7f9;
  position: absolute;
  z-index: 1000;
}
.show-xline {
  width: var(--lineWidth);
  left:30px;
  height: 1px;
}
.show-yline {
  width: 1px;
  top:20px;
  height: var(--lineHeight);
}
</style>
