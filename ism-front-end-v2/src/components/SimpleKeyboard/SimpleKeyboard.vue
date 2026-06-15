<template>
  <div :class="keyboardClass"></div>
</template>

<script>
import Keyboard from "simple-keyboard";
import "simple-keyboard/build/css/index.css";
import layout from 'simple-keyboard-layouts/build/layouts/chinese'; // 中文输入法

export default {
  name: "SimpleKeyboard",
  props: {
    keyboardClass: {
      default: "simple-keyboard",
      type: String,
    },
    input: {
      type: String,
    },
  },
  data: () => ({
    keyboard: null,
    displayDefault: {
      '{bksp}': 'backspace',
      '{lock}': 'caps',
      '{enter}': '> enter',
      '{tab}': 'tab',
      '{shift}': 'shift',
      '{change}': '英文',
      '{space}': ' ',
      '{clear}': '清空',
      '{close}': '关闭',
    },
  }),
  mounted() {
    this.keyboard = new Keyboard(this.keyboardClass, {
      onChange: this.onChange,
      onKeyPress: this.onKeyPress,
      layoutCandidates: layout.layoutCandidates,
      layout: {
        // 默认布局
        default: [
          '` 1 2 3 4 5 6 7 8 9 0 - = {bksp}',
          '{tab} q w e r t y u i o p [ ] \\',
          "{lock} a s d f g h j k l ; ' {enter}",
          '{shift} z x c v b n m , . / {clear}',
          '{change} {space} {close}',
        ],
        // shift布局
        shift: [
          '~ ! @ # $ % ^ & * ( ) _ + {bksp}',
          '{tab} Q W E R T Y U I O P { } |',
          '{lock} A S D F G H J K L : " {enter}',
          '{shift} Z X C V B N M &lt; &gt; ? {clear}',
          '{change} {space} {close}',
        ],
      },
      // 按钮展示文字
      display: this.displayDefault,
      // 按钮样式
      buttonTheme: [
        {
          class: 'hg-red close',
          buttons: '{close}',
        },
        {
          class: 'change',
          buttons: '{change}',
        },
      ],
    });
  },
  methods: {
    onChange(input) {
      this.$emit("onChange", input);
    },
    onKeyPress(button,$event) {
      this.$emit("onKeyPress", button);
     if (button === '{change}') {
       // 切换中英文输入法
       if (this.keyboard.options.layoutCandidates !== null) {
         this.$set(this.displayDefault, '{change}', '英文')
         // 切换至英文
         this.keyboard.setOptions({
           layoutCandidates: null,
           display: this.displayDefault,
         });
       } else {
         // 切换至中文
         this.$set(this.displayDefault, '{change}', '中文')
         this.keyboard.setOptions({
           layoutCandidates: layout.layoutCandidates,
           display: this.displayDefault,
         });
       }
      }
      else if (button === '{close}') {
        return false;
      }
      else  if (button === '{clear}') {
        this.keyboard.clearInput()
      }
      /**
       * If you want to handle the shift and caps lock buttons
       */
      else if (button === "{shift}" || button === "{lock}") this.handleShift();
    },
    handleShift() {
      let currentLayout = this.keyboard.options.layoutName;
      let shiftToggle = currentLayout === "default" ? "shift" : "default";

      this.keyboard.setOptions({
        layoutName: shiftToggle,
      });
    },
  },
  watch: {
    inputName(inputName) {
      this.keyboard.setOptions({ inputName });
    },
    input(input) {
      this.keyboard.setInput(input);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="less">

.simple-keyboard {
  position: fixed;
  width: 100%;
  //   font-size: 42px;
  left: 18%;
  bottom: 0;
  //   transform: translateX(-50%);
  background: rgba(256, 256, 256);
  z-index: 9999;
  max-width: 1000px;
  //   padding: 15px;
  //   margin: 0 auto;
  //   margin-top: 380px;
  box-shadow: 0 4px 0 #b2b2b2, 0 7px 16px rgba(0, 0, 0, 0.6);
  .hg-candidate-box {
    left: 0;
    top: 0;
    position: relative;
    font-size: 30px;
    margin-top: 0;
    transform: translateY(0);
    max-width: 100%;
  }
}

</style>
