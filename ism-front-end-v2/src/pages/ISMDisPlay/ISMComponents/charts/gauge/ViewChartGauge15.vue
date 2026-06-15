<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px'}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2">
        <div class="view-chart-gauge" :ref="detail.identifier" :style="{'overflow': 'visible'}">
            Click to bind data.
        </div>
      </foreignObject>
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
      <!--      顺时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="0 0 0" to="360 0 0" repeatCount="indefinite" />
      <!--      逆时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="360 0 0" to="0 0 0" repeatCount="indefinite" />
  </g>
</svg>
</template>

<script>
import * as echarts from 'echarts';
import 'echarts-liquidfill'
import BaseView from '../../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-gauge-15',
  inject: ['getNode'],
  props: {

  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        if(this.editMode){
          this.initComponents(newVal)
        }
        this.onResize()
      },
      deep: true
    }
  },
  data() {
    return {
      detail:{},
      IsToolBox:false,
      editMode:true,
      base:{
        "text": "configComponent.chartGauge.CylindricalFlow",
        "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAANAAAADQCAYAAAB2pO90AAAACXBIWXMAAAsTAAALEwEAmpwYAAAFyGlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMxNDggNzkuMTY0MDM2LCAyMDE5LzA4LzEzLTAxOjA2OjU3ICAgICAgICAiPiA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtbG5zOnhtcE1NPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvbW0vIiB4bWxuczpzdEV2dD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL3NUeXBlL1Jlc291cmNlRXZlbnQjIiB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1lbnRzLzEuMS8iIHhtbG5zOnBob3Rvc2hvcD0iaHR0cDovL25zLmFkb2JlLmNvbS9waG90b3Nob3AvMS4wLyIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgMjEuMCAoV2luZG93cykiIHhtcDpDcmVhdGVEYXRlPSIyMDIzLTAyLTEzVDE0OjUwOjIxKzA4OjAwIiB4bXA6TWV0YWRhdGFEYXRlPSIyMDIzLTAyLTEzVDE0OjUwOjIxKzA4OjAwIiB4bXA6TW9kaWZ5RGF0ZT0iMjAyMy0wMi0xM1QxNDo1MDoyMSswODowMCIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDoyNjFiZTNkYy0yNGJkLWQzNDAtODA5NS02OGJlNDZiZDI1NjYiIHhtcE1NOkRvY3VtZW50SUQ9ImFkb2JlOmRvY2lkOnBob3Rvc2hvcDo1OTcxYmFiMC0wNzJhLTE5NDgtYjQ2MC02YWE0OTI1NzBiNTAiIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDowODgxZWIwMy0xMTMwLTA4NDAtYWZiMy0yNTg3YjdlODQyODEiIGRjOmZvcm1hdD0iaW1hZ2UvcG5nIiBwaG90b3Nob3A6Q29sb3JNb2RlPSIzIj4gPHhtcE1NOkhpc3Rvcnk+IDxyZGY6U2VxPiA8cmRmOmxpIHN0RXZ0OmFjdGlvbj0iY3JlYXRlZCIgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDowODgxZWIwMy0xMTMwLTA4NDAtYWZiMy0yNTg3YjdlODQyODEiIHN0RXZ0OndoZW49IjIwMjMtMDItMTNUMTQ6NTA6MjErMDg6MDAiIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCAyMS4wIChXaW5kb3dzKSIvPiA8cmRmOmxpIHN0RXZ0OmFjdGlvbj0ic2F2ZWQiIHN0RXZ0Omluc3RhbmNlSUQ9InhtcC5paWQ6MjYxYmUzZGMtMjRiZC1kMzQwLTgwOTUtNjhiZTQ2YmQyNTY2IiBzdEV2dDp3aGVuPSIyMDIzLTAyLTEzVDE0OjUwOjIxKzA4OjAwIiBzdEV2dDpzb2Z0d2FyZUFnZW50PSJBZG9iZSBQaG90b3Nob3AgMjEuMCAoV2luZG93cykiIHN0RXZ0OmNoYW5nZWQ9Ii8iLz4gPC9yZGY6U2VxPiA8L3htcE1NOkhpc3Rvcnk+IDwvcmRmOkRlc2NyaXB0aW9uPiA8L3JkZjpSREY+IDwveDp4bXBtZXRhPiA8P3hwYWNrZXQgZW5kPSJyIj8+GE4oXAAAIeRJREFUeJztnXd8FNXax3+zO9mUTXYTkpAAISEQSOhFQZAAgqEjBLk2EAsgRZCqSFcsl4igoHDBelEvevW9VwlI9yIgUgSkS6gBBOkm1BBS5v1j67TdaZsEfL588jk885w5c2bm+e0p05hpWZ9wIAhCEaMH9wYAREVFMQDA2Gv2IAERhEJyd3wOwCMgtlxrQxB3GC7huDCVV0UI4m6ABEQQOiABEYQOSEAEoQMSEEHogAREEDogARGEDkhABKEDEhBB6IAERBA6IAERhA5IQAShAxIQQeiABEQQOiABEYQOSEAEoQMSEEHogAREEDogARGEDkhABKEDEhBB6IAERBA6IAERhA5IQAShAxIQQeiABEQQOiABEYQOSEAEoQMSEEHogAREEDogARGEDkhABKEDEhBB6IAERBA6IAERhA5IQAShAxIQQeiABEQQOiABEYQOSEAEoQMSEEHogAREEDogARGEDkhABKEDEhBB6IAERBA6IAERhA5IQAShAxIQQeiABEQQOiABEYQOSEAEoQMSEEHogAREEDogARGEDkhABKEDEhBB6IAERBA6IAERhA5IQAShAxIQQeiABEQQOiABEYQOSEAEoQMSEEHogAREEDpgy7sC2uEMKIMxoAzir8wdJCAjBOOrTBIToZ4yFlAgRGAUSupGIrv70BeTZSSgiiwcNVCLdXdgXDwGSEB3i2B8QS3WnUPg4tEgAf0VBKMFarHKh7KLR50CIuEoh8QUeMo+HjUKiISjDxKTsZRfPKoQUAUXzR07JCExaaNixKMfAZVjJQN92UcJZR7PrgqSkKSpePEoIaAyqGTF+PHwj1Q9yyS273wh1amVgK4P3ofO7e/F6bOXMGbKfNy4eUtDSRU7HlndJSjlThGNP8q0x8WVxUYMITbajgdaN0GPTq3Qolka4mIrwWRicDnvKhYsWqZSPHdOPLIBr+zdIhwpyqSFqtitUfOmqZj5yhDUT0uGJYjfoeE4DmvX78Sy1ZsVlFRGgWLwZoy/kFoOgunVtTXenj4MlWMiDS+b4zis3bATL05bgFOnz8vma940DXPeHIH6aTUMrwMAHDl+Gs+Pn4Ptuw5pWt8aFoI3Jw9Ery6tYTI5bsLfuecwRk+Zh1OnL2iu18B+3dGkQQoYRizwU6cvYOGipT7WvnNaGjn0P87ACf4MQ1iw/F/2yk0YPv5dnwGuqQYchxVrt2H4S3PQ8p66eLZvV8l8Q5/tCQAYOGomduzWFuC+2HPgGAaOehvbd+VAy8G2hoXglZeexpOPdERUZATsNivsNivapzfBvKxRSEyorKlej/Z6AN073icpnuKSEnyd/SN27z8q8AQkWKSLD0g88tEmIMMq6EsY6li7fgeGvvgOck+d1VspAEBpKYelq37G4LGzkNnN0cIlJ8ZL5k1v0RALZ49FpSgbBoyciW2/HgTHGXP2duw+hCFjZ2Pvb8ckvMqO19Rx/TGgX1ewZjNvOcMwaNOyoSYRWcNC8PRjnRERHibp33vgOP71f2tV1VMz5RiP6gSku5IB+3kAAGz+ZT8Gj5mFI8dP6yqnuKQE//1+A55/aQ4GPtkNr778DOw2q9MrXf9aNapi4awxSE2pjudGz8JPW/fqEhHHcdi4ZQ8GjJqJnKOnlKwhWa/HMtuj398yROJxwTAMWjWvh0d6PqCqfkOf6YnmzdIkfTdu3sJnX69y9ggC2IeqAPHoX0C64z3Avz4Ctu/KwcBRM7Fnv9Qvtn+KS0rw6eIVGD1pHkYO7oOJo/rBGhYqk5u/X4kJcZg/czSaNkzBwJFvY8XabZpE5Bp3jZjwnsZuKQdrWDB6d0/H92u2IOu9r1BQUCiZs7ikBP/8chUWLlqKhzrfr6glSkupjkd7PSCaNHCx/ufd+M/SDRrqrYAKFo/SAqpglVTL3gPHMGTcLNXjkdtFxfjos+8x/tWFmDK2P8YMewShocGqyqgcE4k5b45A5w7NMXjsLCxdtRmlpcqPg/e4S+uYzjHmeQYfzB6Hf8wcjUVfrcQb734hEpHjx2Ilpr/9Gca/8Dg+enccZk0f5ldEA/p1Q0pyNUnfhUv5+GTxco3XfHygOZwC2+sRzDvqKap856s5wfYPHjmJBx8ew1vGKJgKtoaF4v2sUXiiz4Oy3R5/REVG4M3JgxAaYsHTw/8O3VPQKlZ3icc15uneqSXmFI3A6MnzUFBQiEcz28PsnIVbv3k35iz8D0YO7oPB/XsgJNiCju3uwYwpz2H05Hm4eDlfVH6HNk2R2S0dZrP4t5fjOCxfuwXrftqldU8FBZbZSpph9W8zcBUWiqIsyku/rwEe6ny/ZvG4sNuseH5gJrbsOID9B3MFXpWCUnHxduq4p3gTBqzZjN7d26C0lMOQcbPx6ZcrefmHD8zEyEEPu1tahmHQLeM+lJaWYuiL74hakgF9u8leLjh87DQ+WLRM1a5JUkGFIxU/JvXbNbZJ5Hz8Kw+M3KqJMckIUcfx87PafffUFW2TNZtRp1aCZP5W99YXdVMZhkGjerVQuyZ/nacf74z26U0ly7ldVIxvsn9UOOEhgeZDYmzXTG08lund2OUlCqMQHkglXUIlpXpQUZ5rtTK6QaFalRj0f6QTrGEhkv5d+47g39+tU15gBej1GBGPAb8bu1xE45r5krjAZ+hmnPtm3D5qEJNG/anl6ce7oEnDFEnftes38cm/luPM2UvyBRhyjUZvCYGeRPDalB4CIhot11TUrhNgwalDQxMToN+qJg1S8Fiv9rLjwg2b9+Cb7PUBqFPFFI03XgKqIKIx6Aq+9m0buX0jxnJl3FeT4OnHO8tObf9x7pL4fre7qbXxE48mPYMwXQN+jpP+q+houTBqmJD8l7Hg02xcuJTPW5aXfw2fLF4hmf+zf68WXW8qKCjEp1+uwO79R9HlwRbo2aW15P1uJSWlyF75MzZt3Vchrh2WRzxquhdOUyXvNKH4QuO+GDPL6DvQvslej9GT5+HsucsAgMt5V/Hy9A+w+n/bsXjBFOQdXYb8Y98j/9j32LB0Ls5fzMPA0W+7b3+6cfMWXpv9OeZ++F8AQP6V67h9u0hyW0dzz+Czf6/WuB/+90VZCRqPp0HxqFhAmioaUMEIp9OV/FUsAiWmFWu3YuSk9/Hztv14YcJ7+HHTbsx5cwS6Ce6cbly/Ft6fMRKFhUV4/qU52LB5D16e/gEWfJrtLnbr9t+Qe1J8g+6twtv451crkXNE7bR1OYkmQD/gPmfhNDeFhmF00Psvr6SkxJA7qQsKCvHh58skbueXRv/0OH86bu36HVi7fgesYSFYOHusSDwuGjdwiOiFie+h15OTJUveuecwWjWvz1t/5+7D+OLrNRrqpo3yj0VHLYRItkBl39JUnBbjhw078OHnS2VvvpQi99RZFJeUuO2CgkJkzV2MeR9/q+m46B8zeY5hu9aNRcEvpE6tBHRoI32BFAB27z+K6zcK3HZe/jV8/C8l97uVc0ujY8tK45HXAqmuqGYCP0MyfcJADH02EyHBFlVFOm6wXI7Pvl6FQf17uKduCwoKceaPi5CaFcuasxiRtnD0/VsGim4XI2vuYny7fCN++PZdNG8qfcu/Lw7k5GLU5PedD9Dpu2C7Yu1WBLEs3nljOKKjbCJ/QUEh3vvoW8xZ+B/ZMvb9dhxnzl5CpSgb9h/MxbLVm/Hd8p8kchpzXsu2tdFXZ8ae3I0rG+FoXE/l9qzWUGRNG4q+fTqCZbXdz1ZcUoKv/vsDoiIj0L1jKxQVl+DjL5Zh4usfyKzBwBoWihlTB+PQ0VNYtnoz5r81Bm1aNfL5y++LYyf+wAsT5uDnX/Z7bUW7kB7t1R4zpg3micglnhlzFqssLTC9A00TARq3pJUruat4J4GxJXf1X1pZVVRnnzUxIQ6zXx+Bjg801xy4LoqLS7BkxUaYzWacPX/Zh3j4WMNC8cHsF9Gj8/2663Dq9HmMnTYfa9fvEPm0iKlX13S0al7fbR86egr/FNxcWtaUXWtjTCxeObGad+DlJxECLRrdAzyp9Y39ZeQ4YPj4d3DDq/8PoELcsaBl0iF75SZkr9ykKG/zpmno81A7tG7RAEtWbMLsf3ytqZ5SVDjR6IhFsYACVdEy+NU4dfo8nhnxJrKmDdPXhSsuwZf/XYsJry0UiwcQ74uXoG7cLMCkNz+E3WZFm1aNdXThzuCFCXN5XTg5jLrBtVfXdDzbtyuaNaoNW4TVvfz0Hxdl10lMiEPfPhmItIcDAI4eP42vvl2HGzf5x+3uEI14HXNwVO1Xy180xrUcRUXFWPnDFljDQtGkYW3VIiouLsGni7/HlL9/hOgoGxbNn4zWLRpg49a9KCoq9l8Aw+DK1Rv4+Zd9aN60LqpViVG9DwdycjF8/LvYuuM3AMCAJ7sjs1sbbN+Vo6wOrqqoEFNa7US88/pw3NskFcGCiReOAzZs3o3LeVd5yxMT4jDvrVHo2ycDzZum4d4mqejQthkirKH4YeNO7TOJquPR6B9w+bwTR/ef7m2zRhXsyVL2opHilayP8UrWx5rXd4ynXkD79GYoKSkFYzJhwmsLvFokmeB07v+p388h4+HRCrbkO8ife+ohTB33NKzWEISFWvDqzM9Ev+5yqGmZfD2mbQsPQ0x0JA4f87ysxSWeNi35EyWs2YxBT/UAB2CCwnGjo7LlKRrtsajwTgQF8/l+596Nucbj64EnDhw4TvCn4J8Ql3hckxEsa0bfPh0xfODDMvsjV1kl1yTky3GJx26zgjWbMaBfd0wY1VdT98PXfndo0wyZ3drIPqa98n/bsNmrKyknHhes2YznnuqBrKlD/FRKyzUbI2LRuxx9P+Q+BGREgGirpC9RyF7jktuMv/wceOXHRNsxY9pQ0Uwey5oRbpV+B5qqY+Vnz11leIvHXQezGcmJVfjl6Xl1lvPfs327+nxMe+Fn2bzzcenPK8g5cgolpaWyZV+7dhMHcnJlNqyl3kbGovZ4FMIKs/ovTUkQKMddKc1xoGVF+e5Ms0apSL9P++CfXx+JMryPn8w2GtRNxvABvXni8b1J/2XK8fTjXWTvQrhdVIyvl6xz3++WmBCHf7w1Bjv3HsarMxcBcHT9hM8J5eVfw9QZn2DVul/w+fxJOHv+Ml57W3nX04OClkbP+pJrcN6GX5R/ncFvN0QZ6gQTqHGS1n1xtoLOGPU/SNcmJpY1w2TSKGAfM4RCqlWJxVOPdpZ97533Y9qJCXH4dO7LaN40zf1CxZdeXYDi4hIM6t/D/Y44t3j+tw1z/j4S3Tu2RElpKUwmBq++tUiBiMr2R1xvPPr+OoNBFVVWSS1i0SowNYEv7+YQADEZfY3Jh6Ae790BjerXklzt2vWb+PiL73Hm7EWeeAAgJNiC5wdkwhLEYuIbH2LiGx/ytpOYEId/vD0WHR+41zGGdI7fAPgQUdm1NkbGo9So0ZBxjff4RTq7kj6p2gGPUvSUyYlN1zhKpp8sv30pl7qxkmqc5TdpkIJ+f8vw+XbRb7J/BAAMerK76J4+SxCLJx5+EBlt7xHVd0Dfbm7xuGDNZvR/tAue6POgxH5oHTsqP2/GxKMYz9EzQOG+lW3AL0ggenSafvCFFWF4i5W3TK5yAjTe8sEzj3dBzaSqkj7HY9rZfmOCYRjJmTuT2SQ5hmQYIIg1o6xam4DHI/xeBwpkRZV1k5RsXT1eJ1e4ut9mnYM4WCWEoEpM2oXgsw4ydGjTDN07tZJ9THvJik3YtHWvxnroIOCx6MuhLR41v5XHUNGo3lm96C1XLuBllgvE5L9V8hNIHOdjrCS1rievNSwUwwf2RuWYKMm1HY9pl+ENpuXa81Ebj2JUv5WHA6dfNIYLxrGeOSgcQcF2BFkiUXQ7H0WFV1BSdN0rn5Jfd1/9cS83471Qqnx5MambfJCri1f5ficePHn/1rMd745sb24V3sY/v1yh8jFtDeMxI7to5fwDrugbqfIK1yMapQedn4+12GC1JSPUlowwew0Eh1VBUIhDNCaGBQcODBh3WsoVO8R06woKb57FzSsnUHA1Fzeu5qL49lXBtpQKjJHRhwox+WiViotLUOrjIqW4SsrElJgQh2ee6Co7bb1j9yF8/vUqCI/5th0H8EfP9qgaH+21SQ6//HoQBw+fFJWzbNUmdMtoyXs1MMdx2LP/KJav8fW91EC0NkpFo+3HW/ZxBkMqqrqSYn9wWDxssY1gi2kMa2QKLMGVPE6GvwrHcJ5feKfNwAxLcDQswdGwRtZEpSqt3flvF/6JG3mHkX9hJ65fPoDCgotYvW4LvvzPagx6KlN2hkokCE1ikm+V9h88jtnzv8LrkwYjKjLCx/b9XFsCeILq/2hn2WnrvPxr+PiLZZJTzMvXbkH+leuYP3MskpOqoLSUw7LVmzDsxXfQvVMrLP9qpvui75lzlzBywhwMeGEG5r01Bo0bpMDxseHtGDd1nsQnWwLR2ijt9aiLR3NQuCiHKEJ0N4uqRCPts8c2RURMA9hiGiLMluwRilQKabtVQwumDIhAaLCvViUGQB0APRy14YC1m45j4Oj3cfXaTYwa8pjgxev+ThZfTNbwUCyYNR49u7RRfGGU4zgsX7MZoybNQeHtImRNG4boSna53PxtSxcIAGjSsDYey+wg+XZRjuOw+sdf8N3yjbL1+vmXfXhuzFt4P2sMtu/KwYTXFuLhHm3x+sRBPJHbbeFYMOtFDBs3C4PHvIV5M8fi/MU8jJ441+tzKUZMCATiR1zss9prwRbbGBExDWGLaSDys4ZU1ADRMAyDSgntEF2tHSJi6jsFwQjekCNocmTs7q1DMfJRKyxB6sYXDAN0alMTG5a+jaHTd6Dysj144qHGPkQkNxsHxEZHYW7WWHTvpO7JVIZh0L3T/bBYgjB2ynsYM+U9zHpthOygX1m9gJ5dWqN6tTjJNc+eu4wvvl7lt27bfz2Ilp0GA3C0ZkLxuKhdMwELZjtElNF7tET9fGN8a6M+Hm2VmyI6oa2jx+IMM6mcJndleV4OEgvFy2TXk6qktC8o2I64Wj1Rt+0sJDcZjoiYBo6AYxhH74NhVNnd00Mx+rFw1eLxpmY1Fp/PaIEVu1KwdX8hiorlxiPS+5WYEI9P509RLR4XDMOg4wPN8cl7E/HrnkMYNm6Wyq/Viet1IOeEZPesuLgE//q/1Y5pa+G704R/TmokVsGwAb1lupcOUpKrYeTgPqJ6yNdYcKFTcn98LBPtsr8Ylq5XpWptUaflVNRuMQnR1dLF8SZA8H0gfxvUWlHvxY6TEWSJRJU6jyCt7Uwk1HsSobbqznEApznNbBeC0Y+Hg1Xx0RY54qNNWPByJOZ9W4Sf9xaB44Bm97ZFaESiKKC89zcxIQ6fvD8ZbVs10fVOBIZhcN899bFo3hT8fuY8hox9C7kn/3BuSekdD5567fvtKM6evyzyHjp6Cl98o+7toidOncUrWZ/IiprjOPy0ZQ8mv/mhn5r5ujtAKn4klumJRUdlAY5DTOKDSGubheRmIxAR21A+zgSYg+3Jr8rtnrwpd/JklnsFGxsShbiUnkhqMgSR8ffCHBTiaCKdrYnn/+rsQb2sGNjDCp0fluNhDWGQ0SIEMxdfh4kBTEFWHMxPBxtsQ+HNiygp9JrFYxik1U7CR3Mno0WzerpfKOKiapUYtGhWF0uWb8TytVtgsQRhi/dj3gqnwy/nXUW71k2RmpLoXlZQUIg5H3yDdRvFLy3hb0B8Xo+fOIO9+4/i/hYNeS2RSzzDx8+WFJgh1w8Vd9H8x2N0UgaSmg5BbI2OsIRU8htvz3QLm+5dFBOR+KBAzr7qoLAPKqgkByAoJAoxNTIQm5SBoJBK3hkhPzugzD/s4Qg80j4UWm9g9kf+9VJM/egK9h0rdm+/tKQQF3PX4OKJNSi8cQ4tmtXD3KxxaFBXepZLL0eO/46hY2di+66DziUSO+tHTJPGPIUXR/R1337zw4bteOr5N6Tf+6CQ5k3rYvzIfm4RnT5zAdOyPuaJR9vdAf5iUcV6gDseXd6YpAzEJndEmD0ZauJv/bxY3sF1CkiLwpVVNCgkCjHJGYhN6ghLSJSsHO7UtKjoJi6dWINLuWtQeNPx8g33EZZthaSWM/6zyLYyMssZ7/8y6JrREgtnj0ekPRyX/7yCUZPmYtkqZW/pUYO+u52Nj0XvHDE1XMKpqel8//g+X0D8C6lau2kC0QAe4cQkZcASEuX08/PdDSkbFIr42r0Qk9QeF3N/wMXc1Si6lQcAYGQvcLqWSy1zLpfKIr1Qfjnv1HLYe+Aozpy9CLvNiuxVPxkiHtG4wAjRSC5S2kqJRQO4hJOBMHtNt1frefeGVVZR35UVemJrdkJ8aiYsobHOPE4NM4yX7Vp8d/jZYBuq1O2D6KS2OHd4CS4eX8OXhOTzPgLRiJYzMln8rSflc7yeau+Bowi3hmLRlyv4lwi0dH99isVfBq2ikVtXOhbt8c0Qn5qJ8Eppxpx/AX6eSFVeUQCwRqcivk4m7PH3OLMKtCuycdf5g0JjUL3xINjimuLckSW4cemQICsn08VT0SrxssmJSehzsG3nbziQk4vd+474y6oBf4UYLBpANh4tYdGIq5OJ2OTOHr9h598DE1G9vZK98iEcDmbWiri0TFRJzZQX9F/UPntoCc7nLEFJ8U3xwXf/R+XYRtF4yadDAf7F6B+5oFeRV+WPOMAhpmZnxKdmIjgsxvDzuX5OZeEYSHtFASAqMR3xdXoj1F4dntwMHPPmzF/ejq/TC/b4Zjh3+DvknXKNORj+cVXVKgnPiFzLJMwnmcEHBgpG1mVMa+NaYo1Jc/SAqtwDwHvEbeT55cNEVH9Aoj7+hRNqS0J83UxEVU+X3llCRN7vm3Du4HcouOp6XEDihLj/o7ZVkvEpboSUZFQhKtmsRonGs9TMWhFXNxPxqZmKq6eVDe/KtUAysxdSS+JSMxGXlglzUBikrs4S0kRWb43w+Ca4kJON84eWQKrpUN4qedbx6VPcCOk8j6oF48enQDiAsweU2hsh9urlEous0ooCgL1KM1ROzUR4rOMFEwzEp4xs3zZrsaJqo76wVb0HFw8tRf4f251ecXfNvURyBs87B389sU+QJ6BxpqRw7a2Ni1B7IuJSeyEqqY17WVmcPyGsv9YGACxhMaic1guxtTpL5HSNssRrkl/eHx6TivCYl3A5dx3OH1qKwmt/8NZyoLRVkl5PjFxwa5lsUKNCrS2NtCcuLROV03qBDbKiPM8voGASISalE+Lq9oYlTO4rA/4OPvl9EV2zA+wJzXHhYDbO52T7nURQLyYl9TC6SfJTns/hgvxSe3xTVK7b290DclC+50/6kW6OgzU2DXF1M2Gr0syxiMY6AcNsCUeVxv1gq9YcF3KW4MqZHYrGPZJiAmQmIIxsffyVKZVVm2jAcY4eUN1MxKR0cuasOLEomkQwWcIQVzcTcXUzUTFGDX8d2xqbiuTY8fjz+I84n5ONwqvObp0aMQH8W4hE6wsJUDD6FYyfbTvXj07phLh6vWEJi0Z5nx+pHxvee+GiktJRuV4vhEYmObt+nKcLSHaZ2ZVqtYct4V5cOJiNCweXeoJR4bhHFJbCFkqyLI1IXJ33L0n51saFNTYNcfWcPaCKdH4EsAAQEd8I0bU7I7J6c88KoLQ8UzY4AlWbPgl7YktcOrwKebkb+cGqchJBSlTSOZWjvO3y39K4CLEnILpOZ8SmdPYKXKezoqResEn3v4DIGukAV7H6loSDsEq1kNhqOKJqtMHlw6sd4yNAhZjcmSTLN/6MKyhRotWyhFdGTO1OiK7dGSbWIngXRsWFjUxOB8A4mydKK2oaUbURIqo2Rv6pLbh0aBVuXHA9WAdxQKqaRAC0tUMqA1xGEGyI3dHi1O4Ec4gN7ianghx3USqsv2Oxa+fEKfkrlj8ysSUiE1vi8tF1uHRoJW7lS7xFVNUkgvf2DEJB62EKCkFM7c6IrtMJFqvrEkn5H19lfg+sv2Pr79CTv3z80bU7IKpWW1zOWYlLh1bj9g35T9HLBrQREwkqu1oMY0J0nU6IrtMJIXbpjxrz8ldwP+vSGCi941LGxCKmXg9E1myDP4+uR/6JTbiV/7vUeZamDMcZpqBQRNVqh0o12yK0Ui14RtzlfxzVpXxYz68QpXdqyoZGonLDTFRumIn83E3Iz92Eq2d+hQfxiQ88DnlYIuIRVasdomq2hcUaK1GnOy3lwwqeeqD0Dk/tya0Rmdwa1y/k4Oqpbbh+7gBu5Z2ENEYIS9yKMWwwIuIbIiKhCSKT28LMBrtbm/I+Psa2P64xEMM4rl5Tetek4XFpCI+vC3AcCvJO4vrZfbjy+w7cPO81eycR/FqxRMQhvGpjRFRtiPCqjWE2B8P9SGcFOB5GpSIB8a/9cIKDSvbdYIdUSkJIVBJi6vdASfEt3PrzJG7lncStvFO49edJFOSdBFdcCKVYIuIRWqkGQqJrILSS448NiXL+VDsCzTG8qhj7b6zNxzEG4nXzyL6bbTMbAmtcKqyVU+HdNykuvI7iW1dQXJCPkgJnWlIINjgSQWF2sKGRYENsYEMjwZhYH32cirW/xtt8WMc3dACOARgwIPuvaZuDrWCDw8HZq8KjB9eHyvg2wJV7fcvLFgmIcaqMLzKyySZbyhbifCJV2Lcjm2yy5W0PrKy0CILwi8lzPZhSSin1lwphXd+xYUAppZT6S4WwcsoiCMI/jjEQA8c4iVJKKfWdCjC52h/O6SSbbLLlbSEs4/jan1tcrtt9yCabbLEthOU4ju8VNlVkk022LCxMcOZw5nQ/+0022WSLbAHOMVB5z65TSumdkQpx34ng0pYwhcxy8pP/r+gXwnrU5combLLIT37ye/x8vN6JIMxU3lonP/krop+P404Ep+C8nsDlpeQnP/khbmvgHAO5HYwnJ+NegyM/+cnv9vNxjoGcKwIAhDMO5Cc/+T02H887ETwyI5tssmVtPizHCJQlzEc22WTL4hgDwfXyCEoppdRXKsQEePf6KKWUUn+pN/ROBILQAX2dgVJKVaV8WDAMGAbgOEoppdRfKhIQBw6urh3n/CObbLKlbSHOrzM4vYyzkSKbbLIlbSHO60CMx0k22WTL2wIcYyB4BOaa7SabbLLFthDe80CUUkqp/9QbwXUgRpCNbLLJ5tt8WL6qhBojm2yy5W33vXAOF6WUUuo7FcJ7NzYHvsbIJptsvi3E/YU6l8YYgebIJpts+TbIMQZiOKefUkop9ZkKBQQGjnt9GC+dkU022ZK2EBMH/tMOZJNNtrwthGVMLl05UtcX68gmm2yxLYS+UEcQOqB3IlBKqYpUJCBhD8+Tuizyk5/8rlSIj690Q2Y5+cn/1/UL4Y+BGIAnNLLJJptvC3B+pZsBOA4M4+jrOWxXwpGf/OR3+kUC4gA4V4Orp+exAfKTn/wevxD3OxFc4mJANtlky9lCPO9EcGoNZJNNtrwtwOsLdYwzG9lkky1nC2E5ZxZ+n49sssmWsoWwvL4epZRS6jMV4vVubDnIT37yy8H/SjcDflvFAODIT37yu20BbKOT//Nak9I7MWUYU4Woh6+U40orRD30H9/a8IbhOA4EQWjDVN4VIIg7GRIQQejg/wGXlRI1l4lyIAAAAABJRU5ErkJggg==",
        "isFontIcon": true,
        "info": {
          "type": "chart-gauge",
          "action": [],
          "active": [
            {
              id:"ShowData",
              name:"configComponent.ChartPublic.ShowData",
              result:"",
              isExpression:false,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            }
          ],
          "dataBind": [],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 150,
              "h": 150
            },
            borderWidth:0,
            "backColor": "transparent",
            foreColor:"#ffffff",
            fontSize:40,
            fontFamily:"数字字体-2",
            "zIndex": 1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.ChartPublic.ChartUnit",
                "type":4,
                "value":"%",
                "key":"ChartUnit",
              },
              {
                "name":"configComponent.CircularFlow.progressMax",
                "type":1,
                "value":100,
                "min":1,
                "key":"progressMax",
              },
              {
                "name":"configComponent.CircularFlow.isShowAnimation",
                "type":6,
                "enumList":[
                  {option:'AlarmTips.WEBSpeech.Enable',value:1},{option:'AlarmTips.WEBSpeech.Disable',value:0}
                ],
                "value":0,
                "key":"isShowAnimation",
              },
              {
                "name":"configComponent.CircularFlow.HighColor",
                "type":2,
                "value":"#446bf5",
                "key":"HighColor",
              },
              {
                "name":"configComponent.CircularFlow.LowColor",
                "type":2,
                "value":"#2ca3e2",
                "key":"progressColor",
              },
              {
                "name":"configComponent.CircularFlow.surplusColor",
                "type":2,
                "value":"#0F224C",
                "key":"surplusColor",
              }
            ]
          }
        }
      },
      progressMax:100,
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      fillOpacity:1,
      tempValue:"default",
      strokeOpacity:1,
      animateType:"blink",
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      MaxValue:100,
      echart: null,
      eventValue: '0.00',
      eventUnit: '',
      option:{
        backgroundColor: '#0F224C',
        series: [
          {
            type: 'liquidFill',
            shape: 'container',
            color: [
              {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  {
                    offset: 0,
                    color: '#446bf5',
                  },
                  {
                    offset: 1,
                    color: '#2ca3e2',
                  }
                ],
                globalCoord: false,
              },
            ],
            waveAnimation: false,
            amplitude:0,
            data: [0.45,0.45,0.45,0.45], // data个数代表波浪数
            backgroundStyle: {
              borderWidth: 1,
              color: 'RGBA(51, 66, 127, 0.7)',
            },
            label: {

              normal: {
                rotate:-45,
                textStyle: {
                  fontSize: 50,
                  color: '#fff',
                },
              },
            },
            outline: {
              show: false,
              borderDistance: 0,
              itemStyle: {
                borderWidth: 2,
                borderColor: '#112165',
              },
            },
          },
        ],
      }
    }
  },
  methods: {
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      // 确保 option 存在
      if (!option || !option.style) {
        console.warn('Chart initComponents: option or option.style is undefined')
        return
      }
      let refObj = this.detail && this.detail.identifier ? this.detail.identifier : 'chart_' + Date.now()
      let view = this.$refs[refObj]
      // 确保 DOM 元素存在
      if (!view) {
        console.warn('Chart initComponents: cannot find DOM element with ref:', refObj)
        return
      }
      let i=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartUnit")
        {
          this.chartUnit = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="progressColor")
        {
          this.progressColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="HighColor")
        {
          this.option.series[0].color[0].colorStops[0].color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="surplusColor")
        {
          this.surplusColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartWidth")
        {
          this.option.series[0].barWidth = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="progressMax")
        {
          this.progressMax = parseFloat(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="isShowAnimation")
        {
          if(parseInt(option.style.diy[i].value))
          {
            this.option.series[0].waveAnimation = true
            this.option.series[0].amplitude = '4%'
          }
          else
          {
            this.option.series[0].waveAnimation = false
            this.option.series[0].amplitude = 0
          }
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

      if (!this.echartsView) {
        this.echartsView = echarts.init(view, null);
      }
      else
      {
        this.echartsView.resize()
      }

      this.option.backgroundColor = this.detail.style.backColor

      this.option.series[0].label.normal.formatter = '45'+this.chartUnit
      this.option.series[0].label.normal.textStyle.color = this.detail.style.foreColor
      this.option.series[0].label.normal.textStyle.fontSize = this.detail.style.fontSize
      this.option.series[0].label.normal.textStyle.fontFamily = this.detail.style.fontFamily
      this.option.series[0].color[0].colorStops[1].color = this.progressColor

      this.option.series[0].backgroundStyle.color = this.surplusColor
      let _t = this
      // this.option.series[0].itemStyle.normal.color = this.progressColor
      setTimeout(function (){
        _t.echartsView.setOption(_t.option,true)
        _t.echartsView.resize()
      }, 100)
    },

    onResize() {
      if (this.echartsView) {
        this.echartsView.resize();
      }
    },
    updateView() {
      this.setOption(this.option);
    },
  },
  beforeDestroy () {
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.dispose()
    }
  },
  created(){
    let _t = this
    const node = this.getNode()
    node.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    node.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
    });
    node.on('change:visible', ({ current }) => {
      if(current)
      {
        _t.option.series[0].data[0].value =_t.ShowDataResult
        _t.initComponents(_t.detail);
      }
    });
    this.detail = node.getData().detail
    this.editMode = node.getData().editMode
    this.showDeviceUuid = node.getData().showDeviceUuid
    this.IsToolBox = node.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
      // _t.initComponents(_t.detail);
    })
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据

      _t.$EventBus.$on(activeEvent, (data) => {
        if(data.ID == "ShowData")
        {
          if( _t.tempValue!=data.result)
          {
            let value = parseFloat(data.result)/_t.progressMax
            _t.option.series[0].label.normal.formatter = data.result+this.chartUnit
            _t.option.series[0].data[0] = value
            _t.option.series[0].data[1] = value
            _t.option.series[0].data[2] = value
            _t.option.series[0].data[3] = value
            setTimeout(this.echartsView.setOption(_t.option,true), 500)
            _t.tempValue = data.result
          }
        }
      })
      _t.$EventBus.$on(animateEvent, (data) => {
        _t.isStart = data
      })
    });
  }
}
</script>

<style lang="less">
.view-chart-gauge {
  height: 100%;
  width: 100%;
}
</style>
