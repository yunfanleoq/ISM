<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px'}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2">
        <div class="view-chart-gauge" :ref="detail && detail.identifier ? detail.identifier : 'chart_default'" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px',}">
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
import BaseView from '../../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-category-3',
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
        "text": "configComponent.category3.title",
        "icon": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/4QBaRXhpZgAATU0AKgAAAAgABQMBAAUAAAABAAAASgMDAAEAAAABAAAAAFEQAAEAAAABAQAAAFERAAQAAAABAAAOw1ESAAQAAAABAAAOwwAAAAAAAYagAACxj//bAEMACAYGBwYFCAcHBwkJCAoMFA0MCwsMGRITDxQdGh8eHRocHCAkLicgIiwjHBwoNyksMDE0NDQfJzk9ODI8LjM0Mv/bAEMBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIAbMBtAMBIgACEQEDEQH/xAAfAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgv/xAC1EAACAQMDAgQDBQUEBAAAAX0BAgMABBEFEiExQQYTUWEHInEUMoGRoQgjQrHBFVLR8CQzYnKCCQoWFxgZGiUmJygpKjQ1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4eLj5OXm5+jp6vHy8/T19vf4+fr/xAAfAQADAQEBAQEBAQEBAAAAAAAAAQIDBAUGBwgJCgv/xAC1EQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/APf6KKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACioy5fiPp3b/AAqC181YWJLOPMfhuuN5xj8KLq9gLdFIrBhkGloAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKYz4O1Rlv5fWgBzMFGSaZhpPvcL/d9frSqmDuY5b19KUkKCScAdaG0ldgMlbC7QcZHX0Hc0ltjyRgYGWwPTk1ExySWGehI/ktS2+fK+bGdzZx9TXkUKzq43m/uu33r+vXR7I0atEeyZO5Ttb19aFfJ2sMN/P6U6kZQwwRmvXMx1FR7mj+9yv8Ae7j61J1GRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFGcDJprOF69T0A6mm7S5y/4L2FABuaT7vyr/AHu5+lOVQowBilooAKhlfJ2jkA8j1PYVI7bVyBkngD1NVvo3r8382rzcfXsvZR+f9f0uj3uXBdQ5zwecnBPr3b8Kksypg+XO3c2M/WoWxg5GFwMj/Z7L9TUlkWMUm4AMJGyB9a4sG7YqPo/0/wAvXSz1i76SXulmiiivfMApm0ocp+K9qfRQAiuH6dR1B6inUxkDc9GHQigOQQr8HsexoAfRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUhIUZJwKAFphcklU5Pc9hSfNJ6qn6n/CnAADAGAKAEVApz1Y9SadRRQAUUVFM2BtzjIyT6Csq1VUoObGldkbvvbOcDBwfQdz/SmE+q+mV/9BX+tKc+nOR8vv2X8OtRu+3o2TztP/oTf0FfPtynJye7/r/ga/4X0ZskKSRyGBbccHsW7n6CjTHRopQr7h5rYJ6ketV2IZcv8qhckD+FOw+pp2nZaOcldjCZuB/D0rXDL/a6fz/L+vX1TKkv3bNSiollwcPx6Hsalr6A5gooooAKQgEYIyDS0UAM+aP1ZP1H+NPBDDIORRTChB3IcMevoaAJKKarhjgjDDqDTqACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooqMuX4j6d2/wAKAHM4U4Ayx6AU0ISdznJ7egpVUL079T606gAooooAKKKKAEZgqlj0FViSWyRliRx79h+HWnTSDPsp/Nu35darsdw6kJj8SO5+prxMbX9pPlWy/Pr/AJfnozWCsDSDkK3AB+f0Hc/U1EcAHKk9Ny9/9laeSc/dBOR8vYnsv4dTUZPICsM5O1j/AOPP+HQVzRXT+u3/AALf9uv7LNEhCW3dnffx/tyf4KP5U3SHzJeJ5nmbZM7/AF/zimkrt67E2dT1ji7n6t/npUOjSj7fexiIxZIYIeqgf/rpxly4yg77t/8ApL/q/wCtzVxvSn8vz/r+rG31pFLR9PmX07ilor6U4CVWDDKnIpagwQdynB/nUiSBjg8N6UrDH0UUUgCiiigBrKGHPbofSkDlOJOnZv8AGn0UALRUeGj+5yv93/CnqwYZBoAWiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKRmCjJNNZ+dqjLfoPrQqYO5jub1oATDSff4X+7/jT6KKACiiigAooooAKilmCDaOXPQUjSF+E4X+9/hVYkHJ7Y/HH/ANeuTG1/Y07Ld/03/Xz0Kirsb2yxJ46+3+JoJOeg3Z6dt3+ApSTntuz+Gf8AACo2K7ec7cfjt/xY14aX9f12/Dzi9NkhpIx97C7T83ovdvqe1RuRht6nGBvQdcfwxj696cxbd90M+4fL2Z+w+ijmosnKiNgTk+Wzev8AHKfYdBW8Y9P6/rpb5bOLNYr+v6/rr3EdiCdw8xg43L/z0l7KP9lf6exqnpkjx+IpklkWV5FYb16Ejn+hH4VPnhfLYxAxny2b/llF/FIf9pu3/wCuseCeKHXreS3RoojIFCuMFQeP5EH8a5MZU9nWoT/vf1/Xz3bOulT5oTj5P+v6/Kx2tFFFfWnjBSEAjmlooAFkKcPyv97/ABqaoaQbk+7yP7v+FKwE9FNVw4yPxHpTqQwooooAKayZO5Ttb1p1FADVfnaww38/pT6ayhhgjIpu5o/vfMvr3H1oAkooByMiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKazhevJPQDqaAHEgDJ4qPc0n3flX17n6UbSxzJ+C9hT6AEVQowBgUtFFABRRRQAUUVG8uDtTlu/oKAHM4QZP4D1qI7n+9wOy/40Ac5Jyx7mhiFUk9BQ2opyfQW4yVsDb+f0qE5z2Bz+AP+AFKSc575/X/61MJGOhI449fQfj1NfOV6zrVHN/1/wfxXmjeMbIRiMdCVx07kdh9SaYS27hgWycE9N3c/RRQzd9+Op3+n95vw6ConIwQUOMAFB1x/Cg9z1NKK/r+v6W6uro2ihrMu3qVTZknukfc/7zGopCMOJV+XC+bGOuP4IR7nv9fenMzbuMSP5nHo8v8A8Sg/l61AN7mPyHBzkwu/f+/O3t2H19+OmEf6/r+umzTW0V/X9f113TElc5k8xDKQ6mZUPEsv8EK+w7//AK6w9ZFzHdMbop54O4sgwGGe30JI/Fa3E6w/ZR5ZKsLUuf8AVRfx3D57ntn19zXP3Udq0Mn2Fna23/8ALUkFGI4Y552uO/rg+teVnUW4U/V/1/Xruz0MGrTv/X/A/pdFfvLWcXVpDOuMSIG47ZFTVz3hC5NxpLwBv3ls5Uqe6nkfQ9R+FdADnjoR1Br6rCVvbUI1O6/Hr+J4eJpOjWlT7MWiiiugwCiiigBCOcg4YdxT0lydr8N29DTaCARgjNAE1FQq7Jwcsvr3FSghhkHINSMWiiigAooooAZtKHMf4r2NOVw3TgjqD1FLTWQNz0YdCKAH0UwOQdr8Hsexp9ABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFISFGScCmfNJ6qnp3NAClyTtTk9z2FCoF56sepNKAAMAYApaACiiigAooooAKQkAZJwBSPIE46k9AKi5Y5f8AAdhTsArOz8DKr69zQAAMAUUUxBUMjZbAPQ/r/wDWqR22rx1PSq/4Z9vX0H49a8nMcR/y6j8/6/H7tGrmkI9RD7Anpx/Ifj1NMZv9r1O7+Z/oKViOpJxgkkenc/j0FRsTnGBuzjHbcOn4KP1ry4r+v6/D74u+hukNYnONg3ZACdif4U+g6moHfABV/wC8Vk9P78p/kKVmVlbBYIF5buFJ6/7zmmkFTtKB5NwBTPEkn8Mf+6o5NdEF/X9f11VndPeMf6/r+vlqRuqbSZQ0cYjG8L1jiz8qD/bc/wCelMmbIlFwgZflFwkfc/8ALO2T8+fr78O3MXXynDSF2MTOPlkk/jmb/ZXoP/1VCpJMP2bAYhjavJ/An/LS6fPc5O3Pr7nHTCP9f1/WvZq3RGP9f1/WndapOVIuPtmHQlftxj581/4LaP1HIz65/wBo4zdR+1tLi9gggu5CUiZDmKQHrC59fQ9+o7irse5jbLZRjzDu/s2OXnap+/dSfXPHrn1Y4zlSGbzl028/tWHn7Ra3EuTJ6upPQnr/AHT2x1ry87j7tP1/rX+vxOuirO/b+vRdtdNLfZuQeHNQGmeI0jfesNz+4cSdVbsD7g4Ge4Oe9ekOgfrwR0I7V5FfxCWBnjkeRUIXfJ8skbHokue45w/pweua9D8K62ut6OkjsPtUP7ude+4d/wAf559K7cmr8qdF+q/VGWcYZyisRH0f6f16GqcocP8Ag3Y0tTEAjBGQahZGTlcsvp3Fe/c+fCikBBGQaWmAUUUUAFJypyn4g9DS0UASI4fjoR1Bp1QEA+xHQinLKQcSfg3Y0rAS0UUUhhRRRQAhAIwRkGm/NH6snp3H+NPooAAQwyDkUtRlCDuQ4Y9fQ05XDHBGGHY0AOooooAKKKKACiiigAooooAKKKKACiiigAprOFOAMsewpu8vwnT+9/hSqoUcd+p9aAECEnc5y3b0FPoooAKKKKACiikZgoyxwKAFqJpSeI/xbsKaxaTrwvp3P1pelOwhAMe5PUmloopgFFFRStxt/P8Aw/Gsa9ZUabmxpXdhjtuOcZB6D29PxqNjxyT3yR+p/HoKUn1P1I/U/wBBUbMQTjAYc9eBjv8AQD9a+ZcnOTlLd/1/w3/kr6HQkI7beSQpz+AI/oo/WoSM9VbZgDZ/EwPRfqx5PtT9vQFdzcAIT36qv9TUTuMbvMIGC3mY5x0Z/qfurWsV/X9fj+Nnqbwj/X9f12AsQw2um/cSrH7pYfef/dQcD3qs7oIicukIjyTzviiJ/wDRkh/H8ae/9wxZyVTyQerdUhHsPvMagLs0i+XIjTF2aOR/uySDh52/2EHC++PY11Qj/X9f187N9MI/1/X9fK6GzEbZRMmUwqTxR9cf8s7VPc5Bb6+h4hnIIuPtSeau5RdiI58+X+C1j/2R3/HPVsG8ZjMLGHEbvBLJ/wAsYv8AlpdPn+JudufX/ewyDJktvswa3byWNoJv+XO3P37iTP8Ay0bnGfx/jrrhH+v6/r73bojG2v8AX9afK3VpcyzdLp7yUAfL/as8WT/uWkWOT1wccnPq3Fe8LSrv1LR5LVUBMM9o+9rcdgdoDLjvgFeOeKs25CG0W2t9wJP9l2spOXP8V1N375555/vNwy3DyTSJp/iNLq4WQiVJtkis27nhcMv0Bx7V5Gcxuqb83+RcXyu/b19Ol7dtbrTlWzby7uFpoxc+cLuHBQajaKGcL0IljHDrgNwPyHWsXTdWn8M6xHqCBWtJ/llSJ9ySJ/sE91/unkYI5HNdLdW7eaZrvT57WZ+GvdNbzAQeMsuNzffHBVqxru1+3pLNEYbxn5lkswA5Yjd+9hJwSPnGQQ2cYxwKypXVpR0Z6WHnCUHCorxej7fetPy8orc9TtbqC9tYrq2lWWGVdyOvQipq8b8OeJp/CV40MpNxo7vh9mSYT/eAOGGeDtYA8jqc169aXdvf2kd1aTJNBKu5JEOQRX0mHxCrR8z5zMctqYOfeD2f6PzHPHk7lOG/Q0wHnBGG9KnprIrjBH09q6bnmkdFIQ0f3uV/vf40tMQUUUUAFFFFACKWj6cr6dxUysGGVORUVJgg7lOG/nSAnopiSBjtIw3pT6QwooooAKayhhz26H0p1FADA5Th+nZv8akpKZgx/c5X+7/hQBJRSKwYZFLQAUUUUAFFFFABRRTGfnaoy36D60AOZgoyTTMGT7/C/wB3/GlVMHcx3N606gAooooAKKKKACiioWkL8Jwv97/CgB7yBTtHLelR4JO5jlv5UABRgUtUIKKKKACiiigBGYKpJqsx6kn3J/n/AICnSv8AN1wB0+tR7eu7AA7Htj1+g/WvnsfiPa1OVbL8/wCttu6eptCNkNJJxghR1J7DH+A/Wkxt4Vecj5T69QD9PvGnE9gPQAH8wD9PvGoHYEdGYEYx3YHt9WPJ9BXJFf1/X+Xqup0RiNdht5LFSOo+8QT/AOhOeB7VCztu4ZA+4kH+EOo5P+5GPzNK8nVjJjqxkA6Y4aQf+gqKrOSTs8nJJVBDnqw5WHPoPvOa6oR/r+v67XWh0wj/AF/X9d7bjZHjEZJMiQ7OT/GkbH/0bKfxAqCZvvpPDvBZIpIoj99v+Wdqn+yBy56de2cDzYxKs687pRO4+UY4e5YegHyoP6dK7N5S8NJbxRxZLHl7SBu/qZ5T+Iz69eynH+v6/r5bdcIf1/X/AA9/71uZ5LOx3MLzfPh1XgX1yOkY9IY8c/Q+hykjIIGe4ka7iM4E5jHzajcjpEg/55Jj6cc8Bsxs6ReYHQ2higCTeX/y425xtgTH/LV+Mkc9P9nMirKuAPLsriODBPGzSrbHT081gP09Bz0Jf1/X9fiU9P6/y+W3lb7A+UPMbw3N0I8gf2peoSAijpbxHr35xzzn7zcRyCKaFEn8LyG2Q4iULExjUdBtLAqcdhnFLG4RLMxWbFQT/ZenuSGc97iYnkdc5PIzn7xAEvlXYZw3iIG4Vm8xfLi2BhnIAxuAz2JJrx85dvZvz8yFePX8/T7P3dlsr2bKBmsLZD5WoanpjqpwlwjmPIUH/lopXGV7EUTQPqL74m0rVmUsIpoJDDKCP3g2spbtkcFf0rWEOshhtvrCdQwHzWzKfvMvUOfX0rPubO7kTzLzRNKuyqrJ5nmkMNhw3WM9R71NKNld9fy/4JrCtHm0d36/5qL/ABMDUbAPjzUuolwAjXkPmBI2yUxNESQp5Q7s5yPTjK0nVdW8LTm40oi4sZCZJbTzBJHt6E7l+6w6E7V/hJBB4617B4Cyf2Zq0EauYpTa3yuoR+UwrPjgkDgfrWbeWrTtvvDciYuFmlutKMhjlGNuGhwQHBA79R1610xjKLutH/X9ff2PUo4mE4OnNKUXut/wV/W977dWdp4c8aaR4lRUt5vJvMZa1lOHHrj+8OOo/HFdFXg+o6XbSkP5kEMxbaN87RmFgOFbzFyO21t/A4PatHSvHPiPRkxK66tZK+wecwMvXGNyFsZ7bs57Zrup4zpUXzPPxPD6qe/hJf8Abst/k/8AO3qz2eoWjK8p0/u/4VyWnfE3QbpvJvjPplyDtaO5Q4B7jcOn44rrLW9tb6HzrS5huIj/ABwyBx+YrshVhL4WfP4jBYjDv99Br8vv2AEHp+Ipae8YbkcN6io8kHawwe3oa1OUWiiigAooooAQgMMGlWQpw/K/3v8AGiigCaioBlPu8juv+FSq4cZH4j0pWGOooopAFFFFADWTJ3KdretCvztYYb9D9KdSMoYYIyKAHUVHlo/vfMvr3H1qQEEZHIoAKKKKAIbkyiMeUVBLAHd6VEtwYRiW3eMf3k+df05/MVZddyFfUUiNuQH1FACRzRTLuikVx/snNPqGS1hmbc8Y3dmHDD8RzTPJnj/1VwWH92UZ/Uc/nmgCzRVb7TJH/r7d1H96P5x+nP6VLFPFMP3UitjqAeRQBJTWcIMn8B6015Odqct39BTAOck5Y9zTsAHc/wB/gdl/xpaKKYgooooAKKKCQBk0AFRySY+VT8x79hTwC4z91PXuaiJBPygAdh/L/E1wZhifY0+WPxP+v6/4BcI3YwDaPcdz2/8A1dTTSew/AH8xn/0I0pPpz6Z79+f5moWbPH3s9A3G7PIH49T7Cvnor+v6/T5xOmKGsc8YLZ4APVs8gH3Y8n2FV5JAQSXJBGSy9SDwSPdj8q+1Pdsj+JwfTq2T292PHsoqrLLjLGQjGWLqucdi4H/jiD6muiEf6/r+vNHTCP8AX9f16iPIQc7o1YHdknKKydT/ALkY6erVUlZQrLsl2hQPLH+s2uciMf8ATSU8sey+nWnuxDbAiBg2zy2PyAqMhSf7kf3mPduKoy3CRqsjSSqoBJkx+8COcb8f89pT8qj+FfTpXbTh/X9f137nZTh/X9f12/lJHkcyjCxTNJJjbnEcsqdvaCHv6t79YEaQtE1rKHBLTwy3A4Y/x3kg7KOiD6Y7bY/LLu8EsAZWK25toiNrEcraof7q/ekbp1+gkAkduAt6Z5fujhb2ZegH923j/X3/AIuyKt/X9f1+O9klb+v6/C3l8EkO/EItkG47pbRbjt/evJ/c87Qf0/hljERtoNkclzbvIXtoXP7zUJupmk9EHBHbgH+6KhwJ0/eM13FPL8+zhtTmH8K/3YE/LA9M7re5t11JJdqjKNl/qCcLEB/ywg9+2euT/eOBa/r+v6+7fGb/AK/rX9b/AN5+6qhka7Zr1VkA/wCJlqmdoiA/5YxZ6Y/TqcseKqLYvDiDwrPJBg+WzQwgkdjh3DfmM1sWOl/axBJc2/2exgwbSwI+7jo8nq3cDt1OT0uHmZx/tN/OvGzpNKl5s5vrKjJqOr/D8P8Ahlsu75mWysDub/hC2BAZsiK1B4IP/PT0NH2C0EoH/CJ3yASshCSwLwwyBxN9K6jbv+X+/lfzQH+lRyNmGSQf88o5x9Ryf5CqpW+L+tFd/jYaxtTb9Zdf+3jlUslxEv8AYWtKJFMEmzUFXLrnB4n9m68UeQ7qrSabr4SQfZ7kf2ihzLkBT80x75HPqMe/Vzo3+lpHy6lbmL3Pp+JX/wAepskUU9w8ZJ8i/i3KwOMMB/MjB/4DXZGl0/rt+a/8mKWYT3a/GXr37fkcm9rqTLMslnrrSRrsvkM9u/mw4O04LHJxz0P8QqlNpV1MYpGttRlaVdljLOljKkikZ2yZ+8fQk/QZBz2v2mXy471lJvLMmG5VR/rE74HrjDgfh3p0sNvF/o0uH0q/IMZB4ikPIwewY8g9m+oq1TT6/wBf1o+z1N4ZjOL0ivu++2vzXdaHmj+F55pJLWDS7xWhGbmHEAkiB5ymHHmrnoCMenTFZw8MXUdu2p2cFzDCp2xzwBSQw6lwJt0ZB65yB7V6x9mkluI7K9maPUIQWsr9QMyL3B7E/wB5eh6j2YEefUW2Muna6i5YDmG7UdyP4l9/vL9Or+rxf9f1/wAHodsM8rx00t8/v9H3tps1fU88juPHekQW7DU7spMAfNuGhlh/CVmYfmVrTbxZ45sWWC902wZTjbLcgRK+emH8wIT9DXVW8QN28NvjSdVdS0tk/wA9tcjuyjgEHuy4b+8O1RpDFbzC2y+h3chwIuJbO4J67QcDn0GxqpU5R2k/v/r9CJ42lVf7yjBv/D+OmtvNc3nYxf8AhO/EFmuL/wAMBOMiR7sQofozLg/ganX4gahtDy+FbpYz0lWcGM/8DIC/rWk1mumljPaz6XnrdaWxe3Pu8RBC/ipH+1UY0dmUXtpbWWoRvyLnTJTazP7/ACttY/8AAgPatF7Xbm/L/K5zP6g1d0V98v8A5Ll/G/kVYviJ5xxFoF9Oe4tXSbH/AHyakf4hRRDM2jXsH/Xd44/5tSSRWd1KIJr5BOeBa67ZKxPsrfKW+oZqsf2Zc2fXTLmMD+PStQbH18pyo/AbqadX+b+vkTKGBW9K3za/FtX+VyKLx/FP/qdJu5f+uTI//oJNTnxnMPveHdWX3a2cD89tVZFsZXEc+qokh4EetacnJ9ASqbvwJqwugSRqGj0fSJUPIezme1J+gAI/8epqVV7P+vlciVLBR+KFvm1+MuVC/wDCaj+LS7hP+ujBf50h8bW+cmKCNh3e/hX+ZpfJ+z/6yz8RWn+1FdtcL+A3uf8Ax2lW/tkYKPFlzbsekd/DHGfyeNW/WnzTW8vy/WxPscM9YU7/ADk//SeZE0XjOOXAWGzY/wCzqUBB/wDHqtp4gupf9Xp0T/7t9Ef61ELPUriPel9pN5Ge72R5/wCBCTH6VE2nXiD994f0O7X13lW/IxnP51Xv9/6/EwthXoor73+riaQ1LVWGV0Qkeou0pf7Q1f8A6AZ/8C0rIbTrBF8ybwRbkd3t1tz+pKmmWi+Hru8+yQ6NcxXGwyBApQYBxncrbevvSvLv/X3A6dFK6je3bX8qhtf2hrH/AEBP/JtP8KQ6hq4GToigD1u1/wAKo3ei2draS3U0lzZwRKWYxX05b8MOBn25qzoGkvb6cGv3uJ5pXMojupmmMIOMJlieQByfXNUua9m/6+4xk6Hs+eMV801/7eyKXxDfxRSSHSFaONSzut2pVQOpJxgfnV7Sru6vra3upLT7GsylzC7bmA7E8DHbj3qvcD+2dT+wqP8AiX2bK1yR0lkHKxfQcM3/AAEetbCcyM3p8oqoXu9dDKvyKCXKlJ66X26bt77/AHElFFFWcoVGvyu6++4fj/8AXqSo34dG/wCAn8aAH0UUUAFRS20M/MkaluzY5H0NS0UAZ6QTRr+6nPHBWQbh+fX9ad58sf8Arrdsf3ozuH5df0qww2ysPXmiqERxzxTcRyKxHUZ5H4VJUckEU3+sjVsdCRyKj+zyJ/qZ2A/uyfOP8f1oAsUVX86aP/WwEj+9Ed36df506O4jnfZE43DqDwR+B5oAlJxxjJPQCnLFzufk9h2FPRAnTknqT1NDtsUnr6D1qJzUIuT2Q0iOd/4Bz6+/oKrE598+vfP+PX6ClJyfUnrjv/8Ar6fQVGxz/tZ9O+f8f5Cvla9Z16jm+v5f16r0OmMbIaxz/tZ7H+LP+PX6CoHbPq27/wAez/8AFfoopztn/a3en8Wf5Z/RRVaR93q+7+7xvzxx6FsYHooJogv6/r/g/I6IRGSPuGcl88kpwWzxke7fdX0AJqrLLty28JjLb15VccFwO4X7iDuxJpZZQfm3bx13R8E5+XK+7fcT0AJrOmn43bvKVfn3oMqu35dwHdU+6g/ick120o/1/X+fzO2nT/r+v8vkLPIc+UsGTxH5Tv8ALx8wiJ/ur9+Ru549qqNJzJP57Aj9+LgrlwG+Uzlf77/diTsOfamuN2YfI4x5Rhd/l4+YxFv7q/flbucLz0pfMxvmMxXpN55TLAt8vnlf77fdij7Dmu6EbL+v6/rodkY2X9f8H9b/AN7ZSOgVMGJ4woFu0UJy6A8i2jPeRusj9vXuJPNDK4kHmByLeRbb/lsw6WkHoi/xtx0PTnDFAjiCbJItn7nZE254t3Pkxn+Kd+rv/D696dGoUliwiCD7Mz2oz5Y/59bb1Y4+Z/b2+WyXa39f1/XXRucyFRPNLOI9oENzcwDiIdra2A6t2JHOffAXZ0rSnk8i5vLdbeOAf6JYrytuP7zf3pPft27kyaXo/lmG5u4o0kiXbbWqcx2q+g9XPdvwHHXZrphT6v8Ar+v68vIxWLXwU/v/AMv8/ktLuRWeP+Phv9//ANmFXy2DgcsegFZ4OJnzx8//ALNXi57vR9X+Rx0eo6Pjyj6FD+akURKN0Cn7pWSE/gf/AKxo6J9BGfyc0rfKc/8APO5H/j3/AO1Sw3uxi36/jf8AKJq9f6/ruMhYr9hlbrgwOff/APWv61GY3FnLFGMy2cu+Ieo+8B+KkrUkqkQXiKMtFJ5yD8m/nmpSwTUIZlP7u4TZn3HzL+m6u6EdLP0/9t/NJj5uq/rr+TZC8iR3cF6hBt7tVjc9sn7jfrt/EelRxxx28kuk3Cb7OdWaAHsP4k/AnI9j7VJFAjx3emS5CDlMdQjcjH0OR+ApoWTUdO2MwS+t3xvxwsq9Dj0IOcejVsrvp/w/2l891947pdf6+y/lt6aBEgvIn0m+kZriICSG4HDOoPyyKezA8H39jTQq6lnTdTBj1C3/AHkU8Xylh0EsZ7HsR26HIIy5kOqWUdzbt9nvIWyu7+CQcMrex6H259KkBh1yzV13215A5wcfPBIOo9x+hB96pLtr+q/zX9dS1O2r01+59/R9V/wCB/Luiuk62i+eTm3uEyglI/iQjlJB3APuMimXEsmnwSWuuIl7pTLt+1ugOB6TLjGP9sceoHWrSNBrNvNp2owKtzFjzYwSMf3ZEPXHGQRyCPUUkN1NpsqWWpyeZG52wXbDAf0R+wf36N7HiqXf8f0ZSb+G2q1t+sX0/Lrqto1tL7TUWTSZheWZGRaTyZIH/TOTnjHQNke4qCCDTNSupXsXn0zVF+aZEHlyfV0OVcf7WCPQ1YfTrrSXafRgHhY5ksHbanuYz/Afb7p9utPA0vxHBkq3nW7YOcxz2z/zU/ofcU7dP6+Q/aac97rut/SS2f693sQXL38URh1XT49RtTwZ7WPLY/2oTkn/AICW+gqtZ2VlLGZND1OS2VTgpC++IH0aJs7D7Daau+dqmkZFyr6lZjpNEn79B/tIOH+q4P8AsnrUj2Wk69Gt7AytJ91bq2fZIvsWHP1U/iKe7/q/3i5uWN9k+q1i/WL6/d6Fd59WtkKXunw6hAeC9odrEe8bnH5MfpVK3t/Dt3cGOzZtPvDyYome1kz6lON34girjnVtHTdI6anaL1bKxTqP0R//AB36U9brTvEFgGWz+2QkkFZocBWBwQd3cEEcU7X0/P8Ar/MXO4LmWi7xenzT/L3Q+xaxbf8AHtqqXCj+C9gBP/fSbcfiDTZNTvreNhqOjs0f8UlrKsqfiG2t+QNUIbS6t/EEFlZXk0VusZluYPM81UU5CAF8lSTkjHGFPHStHVWWxtV8lBLfTuIbbzTuy5789gAWOOymhbN7W+YpazjGylfXbla9bWXn101KVvbeHtXRrmy0xXfcUaSKAwOGHUFvlORVe7sr21vLO10/UruCed8+W0xuAkS43sfMzx0HHdhW/awW+j6UsZk2w28ZaSVz16lnPuTkn61U0iKWUTatcoVuLsDy426xRD7i+xOdx92x2o5U7K2oKvJOU+ZuC2T1v20/F6eXUtSrZ2NnPfXeXjiBYs/zHj0HqT0xSaNZzRpLfXigXt2Q8i/881H3I/8AgIPPqSx71XK/2trCwjmw05wX44luOoH0Tqf9oj+6asaxeTIIrCyYC+u8rGf+eSD78h/3QePUkDvQ3d83RGcabUVSW8tX5Lf/AIL+XYhP/E61jb10/T5OeOJbgfzCf+hf7tWNXvpYEis7Ir/aF0SsO4ZCAfekPso/MkDvUqrZ6FpGM+Xa2sfJJycDv7kn8STVfSLWZml1O9QpeXQGIm/5YRj7sf17n3J7AUWe3Vlc0X7/ANmOiXd/1q/u7Fuys4dLsEt4QdkYJJY5Z2PJYnuSck+5qzGu2MA9e/1pr8lU9Tk/QVJWiVlZHHKTlJyluwooooEFNddyEdz0p1FADVbcgb1FLTE4Z19Dkfj/AJNPoAKKKKAI5hwreh/nTalYBlKnuMVCpJUZ69DTQhaKKKYBTDDFK5WWNXBGRuGcGn01jtw3905/CgBPsrx/6i4dP9l/nX9ef1pkn2njzIwwHeM/0PerlFcuJoe3pune3oVF2dzM81GO0Hnpt6H6D+X50x2z/tZ/u988fr09gCa0ZraGdcSID796yrrTLmIF7WUuuOUb0/8A1cV87XweJw+rjzx7rf5r/h36HVTlCTtexDI+7/a3enG7P8t2PwUH1qlNKM45cv02cb938t2Mf7KAnvVW51GWIus8BDc7gO/qMds8D2AxVZr2GbcEkV2bOVkO3fnOc+gO0kn+FFx1NGGqU6vwO/8AX9dz0oUZLW39f16D5pWk+ZsyljkGLgtu4G30LfdX+6gLVQknUnck+w8MHjXIX+BXRe+PuRL3OXoml8wnaSzNnIY7C+4ZOf7pZRkn+CIAdWquJizjcHZ2IIZBsZiy4BUfwsy8KP8AlnGCx5NerTjbU76dOy/r+v08mSYRVKsYggBUoz5jCpyVLd4kPMjfxvxUis27KtMH3+YZGTMoZxw+3vO44ROka8nHSoUeEKGXywuEIKR5TAOIyF7qDxGnV2y54qSLzJ5lhhWQszMgWOTLkn76q/8AfP8Ay0m6L91ea2Rcl3/r+vP56aEwblY41ZQpNuq27bmB/ighPdz1kl7c8+nU6Nov2QR3N2sX2lU2RRRD93bJ/cT+rdT7DAqTSNGj05FlkEbXWwICi4SJO0cY7KPzJ5PtqEgDJNdVOlbVnhYzG894U9u/f+vxFpBlzhOndu1KsZfluF9O5+tTAADAGBW1zzBqIEHHXuT1NZg/102f75/ma1ayv+Wlwf8Aab+tfPZ9vR9X/wCks3o9QkRljk28/u5OPTDUTEMt3s5/drIPqM//ABNWAuZNvr5q/mQaZGgkkgYcGS3598Y6/nWlGN/dfkv/ACW35yLvbX+v60AEfb1I+7ND/wCgn/7Kq2xv7KZFGZbV/lHrsOQPxXH505W2wWTtw0bhGz36r/PFTx/u9RmTtKgkH1Hyn9Ntdsff30v+Ule/3oL8vy/R/wCRFcuqSWt+h+T7jn1R8YP4Hafpmib/AETUo7jpFcYil9m/gP8ANfxFLbRI9pcWEgysbNHj/YPK/ocfhSQr9v0xre4JMgBilI67h3H6EfUVqry1W71+a0a/T7w0W+y0+T1X+f3DZv8AiX6gLoDFvckJP6K/RX/HhT/wGmX8bWdx/acKsVAxcxoOWUdGH+0v6jI9KntmF/YPBdKGcZhnXsSOv4EYI+opNOmkHmWdwxae3wNx6yIfuv8Aj0PuDVqz22eq8n/X6ju1q9WtH5r+v0Eu7cXqQ3UEqpdRjdb3A5Ug/wALeqnjP4EcgVNBNBq9nNb3MADD93c20nO0+nuD1B7iqKSJo12LaRgLK4YmEk8RP1KfQ8kfiPSq+rSvbtFqNikgnjZYhuG1ZlZgBGc8nJPBxwfbILvo5fehxV2qd/8AC/0+/wC56+tuOWbQiIbuVptNJxFcuctB6LIe6+j/AJ+tT6np9vOyXYuPsd3GMR3SEA4/usDwy+x/DB5o1KSS3025ubqQeVHGxMUS/f4+7k5znp0HWmaXYQaLokH2jYHt4AZpTztwMtg9gOcD0qrJPl6fkHPJpVY6Svb176fn0dyPTNZubuW4tJrPN3bbd7xHEThs7SC3Izjpg49T1qpfWch8Q2Jt3W2vJmL3DW+RmFRzv7NklVGRkZJB4rQ0CGQaebydStzfObmQHqu7G1f+AoFX8KZov+mz3esNyty3l2/tAmQp/wCBEs30I9KW6SZo7U6k5x0SVvJt6fde7+QawV0+wLW6B76dhBbNIdx8xuh57AZY+ympI0tdC0UKWK21pDlmPJIA5J9Sf1JqtGf7T8RSXHW308GGL0aU/fb8BhfrupL7/iZavBpw5t7fbc3XoTn92n4kFj/uj1qr/a+4hQVlTey96X+X3bebsTaLbSxWr3V0u28vH86YH+DIwqf8BUAfXJ71Dp//ABM9Um1Q828O63tPQ8/PJ+JGAfRfepNankZIdNtnKXV6SgdescY++/4A4Huy1LdzwaJpI8mH5YlWK3gT+Jvuog+pwKNFp0QXlJcy+Kei8lt/wF5JlW//AOJrqcelrk20G2e8PY85SP8AEjJHoAP4qn1e8mgjhtrTBvrlvLgDDIU45cj+6o5P0A707TbT+zNPY3MqtO5M91N0DOfvH2AxgegAqHQ0a+nl1udSDOPLtUYYMcGc8+7H5j7bR2o1+b/r+vMfub7xh+L/AOD+StuXoo7TQdHwX2W9uhZ3Y5LHqzH1JOT7k1Bo9tMxl1O9Rku7vGI26wRj7sf17n3J9BUMn/E61fyeun2EgMmRxNOOQv0Tgn/ax/dNTavdTO0Wl2TlLu6BzKv/ACwjH3pPr2X3PoDSut+iDlk/cb96WrfZb/8ABfy80Q/8hzVc8NpthJx6TXCn/wBBT/0L/drbqG1tYbK1itrdAkUShVUelSsQqlj0AzVRVt9znq1FJ2j8K2/ruxq8yM3p8oqSmRgrGM9TyfrT6oyCiiigAooooAjfiRW9flNPpJF3RkDr2+tCtuUMO4zQAtFFFABUBG2Vh68ip6ilH3W9Dg/jTQCUUUUxBR1GKKKAHxHMYz1HBp9QxnEjD+8M/wCf0qakxhRRRSApajp1vfQMJEG8Dhx1Fefa9oc+n5cjzIGPDgZB+o7V6dVV4Y5YnglQOh+Uqw4IrzMZldLEP2kPdn3XX1/q53YPH1MNLvHseLi9ZJcPKcdWL/MMFgzEj0JGSONxCjoKvx3UtwqpJEsrv8jZP3izRAhvXLyfPz823aMKKt+NPC0umqb6yDNa5yR3iPv7eh/yeSsb9kUMmQ6Y2jn7wztUZ5PzNuJ9q8uOIr0H7OstUfa0YUsVRVaj/T/zOntDJfTRJAJGlmO5AzbWJZQSS38JKsu5x91SqJ1zXoei6JFpMGSVkuWUK8gXaAo6Ig/hQdh+Jyea4nwjOn9uRx4WSMRTRgEdQBDGP/RIrvhFHt/cSPC3ZQfl/I8flXp4TGYabs5Wl2f9WPms5qVYyVFaK2v9f15lzPO1RlvSnpHg7mOW/lUEZuYkH7uOZTzlDtb8jwfzp4vIchZCYmPaQbf16H8K9Vs8EsUUUUgCsr+K4/3n/k1atZY/1k3vI/8AJq+fz1XdH/E/yNqPUsjidf8Arsw/Nc0yL5RaH0Zo/wBD/wDE07+PP/TZT+agU0/LDn+5c/zb/wCyrWm7Nvs2/ucf/kWV/X5/5jJI99lexAfNG7Mv1++P1NJcHY9pc5zHu2lv9lh3/HbUpnit76ZZHA3orAdyeQePwFUmllbRJUSBtsKMA7nbjYTjjrngdhXV8Kduif8A5K9PzGtWv63WpZf9zqcb/wAM6FD/ALy8j9N35VFJKllqZLsFjuU5/wB9f8R/6DUd3DKzQCSY8zJhYxtAPU89envUjwRJqlokcYBCySE9zgBeT/wKuh6N26NP79/1+8mOq17P8Nv68iJp3TV0NvCxFzGQ/mfIMrjDevQ46elNnhlk1q0Dy7X8qRn8obfk+UYJ69SD+FWR+81tj2gtwPxduf8A0AUlt++1e9m7RBIF+oG4/wDoY/Kq30XWX5av8UyttX0j+ei/Bogmt4m1mygRB+7D3DseScDaoJPPVif+A0+6H2vXLC16pDuupPw+VQfxbP8AwGnWH7/UNQuz03i3Q/7KDn/x5n/KjRh9oub/AFA8iWXyYz/sR5H/AKEXp7/N/gv+G/Er4df5V+L/AMrv7hdT/wBL1Gw04cqX+1TD/YjI2j8XKf8AfJpuu/6W9ppA/wCXyTdMP+mKYL/gTtT/AIHTtJ/0u8v9SPKySeRCf+mcZI/Vy5+mKbpX+m6nf6oeU3fZbc/7CE7j+L7vwVaN/n+X9fmaL927/wAi/wDJn+q/9tH69LIbOOwgYrcX8n2dWHVFIy7fgobHvj1p+pXA0nSAlnGvm4W3tYu28/Ko+g6n2BqGx/4mGt3eoHmG2zZ2/wBQQZW/76AX/gBqPcdR1+Sbrb6aPKj9GncfMf8AgKkD/gbDtTu3quuglFRtGW0dX6vZfkvLUfH9n0DQ2ZmZordCWbGWkPc+7MxJ+rUukWr2Vg894VF1OxuLps8KxHTPooAUey1BOBqOtw2Y5trDbPN6NL/yzX8Pv/8AfNLq/wDxMLqHRkPySjzbsjtCD93/AIGfl+gaqvbVdNEJJyXLJ6y96T8t/wDg+egujI13JNrMqkNdALArDlIB93jsWyWP1A7U2D/ib6y10fms7Fmjg9Hm6O//AAHlB7lvaptYuZY4YrGzIW8uyY4iP+Waj78n/AR09yo70txLBoOjolvEW8sLDbwg8yOeFH1J6n6mjRaPpuF5S96K1lpFdlt/wL+pX1QnVL4aSnNvGBLenHVf4YvqxGT/ALI/2hVzVLqa2ghsbIj7fdZSEkZCAfekPso59yQO9RWcEej6a0l7Ohf5rm8uDwGbqx+nQAegAp+kW80ry6reRlLm6ACRN1hiH3U+p6t7nHYUnfbqwvFa7xjt5v8Ar8El1JWaz8O6KB83kwKFUZ3PIxPA92Zj+JNJpFjLAkl3ebTqF0Q8+05CY+6i+yjj3OT3qtbf8TrUxfHJsLRitsp6SydGl+g5Vf8AgR9K26cVd36dCKsnBOD+J6y/y/V+foFMk52p/ePP0p9MX5pWP90bR/X+lWcpJRRRQAUUUUAFFFFABUacFk9DkfQ1JUbcSK3r8poAfRRRQAUjLuQr6ilooAgU5UE9e9LSEbZGHryKWqEFFFFACE7Sreh5+lT1CRkEHvT4jmMZ6jg0mA+iiikMKhkGJQf7w/UVNUcw/d5/u800BFLFHPE8UqB43BVlI4INeGeMdAk8K68rQnNpcZeCQrnb6qfcf4V7tWL4r0FPEfh+4sDgTY8yByfuyDofp1B9ia5MZhlXh5r+rHsZLmTwWIXP8EtH/n8jzPwRdq2vQhGyip5Sk9ztdifzNeqJJkda8Q8IXUtv4mtLeVSjq7qwYYIIRhgj1r2GCcEDn3zn/P8An0r47EU2qh7HEOH5cSmuq/VmtDK6AgMRj16f5/KrAmVwUkTOeoxn9Kzo3+dffjp/n/ParAIIHoTx06/y/ka7MNiq9HSD07HzMoLqTi1hxut3aL/rm3H5dP0pc3cXURzr7fI3+B/Sowec55HU5PH49R+NSrMy4yQQem7jP49DXtUczi9Kit5r+vu3MXDsAvYgcS7oW9JRgfn0/WqIkRZZMt1kbAHJPWrlxcHy/LjU+a/ADLnA7n3qq0H2XaUVlj4GW6g+tcubzjUjCpTd+R3duit17fMqno7MDLKyMY4sAbH3OcdO+O/SiaCQxXfmzMSo8zag2qePz7etSrh8KeFYFCPQNyPyIIpy5lYZ6ywYP1H/AO1Soy9pC3f9b/rJL5F7O/8AX9aCpDFBfoIkVQ0TE4HXBX/Gq8v/ACBr4erTD82ap4m3SWMn96Fh+YU/0qtOf+JXcL/fnZfzeu1y5uZrqn+MYscVZr1X5skuPmubIHqZWc/98N/iKE+fWJT2jgUfizNn/wBBFLL82qW4/uxSN+qj+pptqw+0X8zHCiULn2VB/XNdf/Lx/wCL/wBtIXwfL9RLD55r6c9GnKg+ygL/ADDVDYzrbaJJqEv3XEl03+6SWH/juBUW54vC24fLNPHx7PKeP1eptSjX7NZ6egws0qR4/wBhfmb9Fx+NTFuyfW1/nI1cU5cr2vb5RWv4ERaTSfDO4jNyIs4/vSv2/Fmqe53aL4bWC3O6dI1ghJ/ilbCg/wDfRyfxpt6PtWsafZ9VQm6k+icL/wCPFT/wGn3H+m+IbW36x2SG5k/32yiD8vMP4CtGraLpov6/rYcXe0pdbyf6L77r5jb0nRfDyWtl/rgiWttnvI2FUn6feP0NLdkaF4fjtrIZmVUtrUN/FIcKpP48n2BNI3/Ew8TKnWDTY959DM4IH4qmf+/go/5CXiTPW30xcD0M7r/7Kh/8iU/T0X9f1sUunPr9qXn2Xz/9uHTFPD3h5IrdfNkiRYoVY8yyscLn6sck/U1G3leHdCUOWmaFS7kfenlY/wA2dvzNO/5CfiL1ttN/Jp2H/sqH839qrv8A8TbxFt62mnMGf0ecj5R/wFSW+rr6U13Xov6/rYLXVp/4pfovnf8A8m8iazRdE0WS4vpAZcNcXUg/ic8nH8gPQAUaTA9taTX99iO6uj58+48Rrj5U+irgfXJ71Fd/8TTWY7Ec2tmVnufR5Osafh98/RfWjU/+JnfR6OnMOBNent5eflj+rEf98qfUU9tum39f11FZy0k7OWsvJdF+tv8ACO0dGvJZdZmVg1yAtujDBjgHK8di33j9QO1Nsv8Aibao2pNzaWxaKzGPvN0eT/2Ue2T/ABU/V5pJ5ItItXKTXIJlkXrFCOGPsT90e5z2pdTuTptjBY6eirdTDybVNuVjAHLEf3VHP5DvRovl+Ye9LVaOWi8o/wDDael77kbA61rDW/JsLJ1Mx7TSjlU+i8MffaOxqfVppL25XRbVyrSrvu5FODDCeOD2ZsED/gR7U4mDw3occcUbSspCRxg/PPKx9fVmJJP1Pap9KsHsbdmnkEt5O3mXEoH3n9B/sgYAHoKVm3b7w54xXOtlpHzfVv8AP7lqkXIoo4IY4YkCRxqFRVHCgcACn0UVocLbbuxCQASeg5pIwRGM9TyaSTnan94/pUlABRRRQAUUUUAFFFFABTZAWjOOvUfWnUUANUhlDDoRmlpkfG5P7p4+lPoAKKKKAIpRgq3ocH8aSpHXcjL6iolO5QaaELRRRTAKWM4kZfXkf5/KkpCdrK3oefpQBPRRRUjCiiigCuvA2/3TinUOMS5/vD9RRVCPGPGGnronxIjuo12w3o88ADgMcq456knk/wC9XW6dd7sck/5/+v8Ar7/NU+LlsP7K0vUAx8yC6MYHs65J/wDHB+dZWgXvmQIMZJAyP06HHv6flk18zmGHtXbR9sm8Xl1Kq90uV/J6fgd3HKCvBHQf/W9v6emDxV9H3AH+9x9fbnr9Dz6GsW1mzgdT29Tn6+uO/XGDhhWhbyblIznjAwM5X6fxAdweRXMqen9f1/SPn6tOzNAN3J6cZyRj8eq/Q8UrSCIEsSPUcAn8OjVEjZK4ySR8u05JH+yf4h7Hmo4mFxJlcNGp+VVxgn12k8Cr5bf1/X/AXmzlcSxEhX55F2yN1GGAUenHapP3bgqRA2exkOf5VHuCHkqn1LR/z4qTc+OfMI91Dj9Oa0hFbWXySf3t6ENFYB4nMTctjKHPXuP1H61OjASxkdBKcfRhn+dMlRZVwjQLIDlTjYfyqJJcxs2CpXkg9tjZ/kf0rOl+4ldfD9+2tr7b66Fb6E0XH2Mf3ZHj/IN/hVeXm3kHYXQH5yirB+WRc/wXR/VSf/ZqrZ3W/wDvXgP/AI/mvUp6RcfJr8Ir9B9U/T82WB82rv8A7EC/qx/+JqkzH/hH7l1OGuGkCn3dyq/zFWfN8q51Cc9I0Ufkpb+tQeWVtdJtD1LIW/4ApbP5gfnXS9b/APb34uyCOjXy/BXZNeKGutPtFHyiQykf7KDj/wAeKUg/0jX2P8NpBgf7znJ/IKP++qdF++1u4f8AhghWIfVjub9NlUobv7Nol7qm3c87vKg/vD7sY/EBfzra6vfpf8F/wRqL5bLeyXzlr+Wha0nFxeajqLEbWk8iM542R5B/8eL/AJCmaZcRw6Vea3ctsS4LXTMf4YVGE/8AHFB+pNMu7drLw7a6RG5M91ttd46ktzI/12h2+tSasi3E+n6NGoEUjebMo6CGPBx+LbF+hNCuv66s0tGTt0f/AKTH/P8ANDbKRtH8OzaheoftEu66mQdS79EHuPlQfQUsRbw/4baWYedeNl3Cn/W3EjfdH1ZgB6DFPvf+JhrtpYjmG0xd3HoW5ES/mGb/AIAPWkk/4mfiNIuttpoEj+jTsPlH/AVJP/A1Pantt6Dvzaz6+8/Tovn+qEJPh3w4S37+7746z3Dn+rH8B9KjQDw7oGWzPcdTjrPO56D6scD0H0pXb+1PEQ72umn8HnI5/wC+VP5sfSmx/wDE111pTzZ6cxRPR5yPmP8AwEHA92b0ql5eiJev8T/FL9F+P4+Q6Lb4f0N5rpvOuCTJMyjmaZz0X6sQoHpgUtqq6LpE97qD5uHzPdOozlj/AAr6gDCqPYdzTE/4m2tl+tlpzlV9JJ8cn6IDj/eJ/u0N/wATjWtnWx09wW44kuOoHuEBz/vEf3aPT0QNXup9fel+i/ruuw/TozY2dzqmpssVxOPOuGJ4iQD5U+ij8ySe9N0qGS4nfV7xTHLcL+7Rj/qYByoPoTwzfgP4abc/8TjVfsS82Nmwe5I6SSdVj+g4Y/8AAR60t4W1bUW0mIkW6ANfOO6HpED6t39FH+0KNtun5hZyvfRy1flHovy/BdWS6cp1a/8A7YkB+zoClgvqp+9KR6t0H+z/ALxrapFUKoVQAoGAAOlLVxVkclWpzy00S2Xl/X4hRRQTgZPamZjF+aUn+6Mf5/SpKZED5YJ6tyafQAUUUUAFFFFABRRRQAUUUUARtxKp/vDaf6f1p9NkBMZx1HIpQQQCOhoAWiiigAqDG2Rl/EfjU9RSjDK34GmgEooopiCkIyCD3paKAJI23Rgnr0NOqKI4dl9eRUtSxhRRRQBHMPk3f3Tmm1N1GDVdeBtPVTimhHJfE2BJfAt47DLQvE6/Xeq/yY1514elKBUIIB5AwP5Dr6f7RwOhNen+P4/N8DaovpGrfk6n+leUaOrIFGzB4ISM8nt1/vHO1fQtntx5eMjeqvQ+2yJ82Wyi/wCZ/kj0SykBXLFdoXkk7hjAzn1GMZ9QVcc5rXhLebgZyxCkM2CWHQE9nx0b+IVzmnyjYrb+g3b0XgAc7wPQBg4HeNyP4a1jKsCLuVEI/d7Wb5F77GP9zB3K3YcVyOml/X9f13R5eJp+9/X9f10NRnaYmJSzLnMrCMkf8CUDKvxyR9asIwkTav7xV4wpWYL+BwwrNiLRoJHUtnkzPliT/wBdY+g+tT+c0qh9glQdHIE6j6MmGH41Dg9v6+fZf15nC4f1/W5eSTnajYb+6km1v++HpxKo2WCIfV1MR/76HBqms7yoQg81B1CyCdR9VYbhTo5znZFgP/cilKn/AL9uKOTb8P8AgLr6szcf6/zZf3Pt583b7qJV/TmqsqqH81DEVAIkWM4OCMZ2+2TTBKFf51iRj67oGP49DUjyy4/eIdv/AE2QSL+aj+dNx5vi1t5/rsvlqRa39f0wZ98YJPLvCx/Hg/ypP+ea+t4f6mq+fLfOYzESmwo+4AhwcfzqwP8AXQj/AKfHP/jjVtRvon1/WTB9f66EVx81lqX/AE2mEQ/EKlWT+81tB2gtyce7sAP/AEA1WHzw26/89r5m/BWZh/6CKJZ2ij1e8T76nyo/chRj/wAeYiuyEl8T20f5y/VD5W/dW/8Aw0f8yLznTQry6jOJruVhEfdm8uM/ltqW5gRrvS9KjH7qMiZx/sR42/8AjxT8qV4FS50rTk+5Aplb6IoVf1YH/gNFnPGLnVdWnbEEI8lW9EjBLn/vokf8BrVKys/Jfq/vKv8Aaj5tfP3Y/duTJ/pviSR+sVhF5Y/66vhm/JQv/fZqLTZo5JNS12dwsDExxOeiwxZyfxbe30xVYtc2PhlQuU1PUZMD1WWU5/8AHFz+CVNf20W3TfDtsuLcgNKvpBHjg/7x2r7gt6VSfX+tdi+VfDfTb5R1k/m9V80JZTnTdCudYvI2+0XTG4aL+LLYEcQ98bF+ufWnAzaD4eyQs2pTv07S3Eh/9BBP4KvtUtx/xMvEENqObewxPN6NKf8AVr+Ay3/fNMjb+0dbmu+ttp26GH0aYjDt/wABHyfUvTt0Xp/mxNp+9NafE15fZj/XR+RBIG0TR4rOzbzb6ZvJhdhy8rZLyMPQHc5+lOuQdH0m203Tzm7mPkwFufmPLSN645Y+p470aaRf302rOR9miDQWhPTaD88n/AiMA+ig96NNZbuafXrj5ImQpa7+NkA5L/8AAyN30C+lV6f0v6/QTur8+tnd+cnsvl1+fkLdH+x9MtdM03/j6m/c2+/5sHq0jeuOWPqeO9LcEaLpdvp+nruupT5VuH5yxyWkb1A5Ynv9TTdNYTGfX7z90skZ8gOMeTbjnJ9C33j+A7UunfvTNr19+6DxnyQ/Hk24559C33j+A7Ub7f0gatfm1s7vzl2+X+fdCzMnh/SIra1Xz7uRtkSueZpTkszH82Y+mfatDStPXTbBYd5klYmSaVuskh5Zj/h2AA7VnaXE+oXf9r3KFfNG21iYYMUPXJHZmIBPoMDsa36qK6/cYV5uN4Xu95Pz7fL87+QUUUVRzBTJOQE/vHH4d6fTB80xP90Y/E0ASUUUUAFFFFABRRRQAUUUUAFFFFABUcfAZP7p/SpKjPyyg/3hj/P60APooooAKa67kK+o4p1FAECncoPrS0mNrsvvkfjS1QgooooATO1lb0OD+NT1ARuUj1qWNt0YJ69/rSYDqKKKQwqFxtl/3hn8RU1RzD5N3905poDE8WoJPCWqqen2Zz+QzXkWkxxJHtDMI1J3uOudpLEf7qBiPcivY/Ei7/C+rIOpspgP++DXk+lw/Zo1VU3mFR8v95seawP12xr/AMDrgxS/eJn1mQzthKkfP9DobR2hGZHSBkOSx5WJtxGf91ZC6H/YlXsK04m+zOssoNo+3ZEksjQGNMkhFkwUcc8A+1ZFvGHnSFZd5jOEVZI1lkI+XcqyDbIjIsRIzncPWr8chtpPIjIhkb/ljETZu3sIJcxP9QRUcl9Sa0eZ/wBfj5fn0030mfyG8yUCFm58yVTbk/8AbWPMbfiKlZyMTSjGeksyYz/22i4A+orMSZbSbywRazMcBBmxkY+yNmGQ+4qXzBay5fbbSsesgNk7H/eGYpDWbp/1/W7OV0v6/wAuy8zSLmRFlcGRP4ZHQTr+EkfzD6kVIshmi4LSxD+6RdRj+T1nO/kP5k4ELt/y0uFNsx/7bRZRvoRU0j42zXAI4+WW5j7e08PQfUVLhv8A197/AERi6e1v69O/qy/DKWysLFsdVglD4+sb8j6ClR1WTanlrIf4UJgc/wDAG4NVC5kiV5A0kX8LyILmP8JE+YfU1JHK0sJ2F5Ie/lsLuL8QfnH0FRyf1/ktl6sycP6/4PX0RZm2sCLjaCe9zHsI/wCBrxTrV3NxAkpXzBK7nacggqcEfrVaCTJK27ZI6ray8j6xSdPwprNiTdALYXQBwGUwOc9RtPDfWiPuyUvNfn97fpoZypvb+v8AJFuz+dtMH/TB5z9Tt/8AijUUf76y06P/AJ+7k3Df7uTJ/RR+NNaXZbXLx/ehsUiT2c7hj8wtWCY7a/ZmOIdPs+T6buT+QjH511U1ok/T8k/wixPuv63a/GSI3uhA+r6mwJWBRBGB32jJx9WYj8KbLaNFp2l6IxzJcOHuSO6r88h+jMQv/A6jjhd4tI06QfvJnN5cj6He3/j7KPxqyLuNLvVtZmJ+z2aG3jx3CfM5HuWIX6pW613/AKvq/wBC0mnePy87e7H8bslX/T/EzP1g02PYPQzOAT+Kpj/v4arWV5HHZ6h4kuMmOYfuAOSYUyEA92JLD/fFRPBPb6FbaazFdR1WQ+ey9UL5aUj/AHVyo/4DVq4jS81m00yJQLSwVbiZR03DiJPwILe21fWqu/67v/JDaja32f8A22O//gUtvPQjzc6PoIHytq9/L06gzv8AzVAP++Upt/D9j0+x8PWLsJbgFWkz8yxjmSQ+5z1/vMKntnW/1e51SVgLOyDW9uxPBYf62T8xtHptb1qrZ3aJbXviO9DqsyjyU2/MsI+4oH95ic49WA7U1b+uwnzXu1d3vb+8/hXyXT1Q/UI0uHt/D9qPLhMYa5Cf8s4BwE/4Hjb9A3pTr8DU79NHjA+ywhZLzA4K/wAEX44yR/dGP4qZG8mj6VLe3UfmaneSAmJTndK3CRg/3VGBn2J9aU+ZoelJDERcapeSHB6CSZuWc+iqBn2VQPSn6/0uwkmrcurWi831l6Lo/R9x15/xOdT/ALOXmztWV7wg/ffgpF9OjN7bR3NF1/xOtSNgvNjasGumHSSTqsX0HDN+A7mmyhtG02302xfzdRuSwSSQZLOeXmf2Gcn8B3FaVjZQ6bYpbxZKoCWZvvOx5Zie5JyT9apJt2fz/wAv6/UzlNU4qUenw/rL/L/7UtxDLs3pwKlpka7YwD16n60+tGcIUUUUgDoM02L7mT1Y5pJOV2/3jj/GpKACiiigAooooAKKKKACiiigAooooAKZKD5ZI6ryKfRQAgORkd6KZHwCn904/DtT6ACiiigCKUYdW/A0lSSLujIHXt9aiB3KCO9NCFooopgFLEcMy/iKSkztdW98H8aAJ6KKKkYUEAgg9DRRQBl6mhk0S+jPU28in/vkivMNMgmnl2wCTzMtKWiTeyL53DBcgtjyEGBzhq9VvyFhlUjJlQqqjuTxj+VcMmhXWmWiRXdqXC4LOYBdw7u7LtCyxn6ZA/nz143kme3lVdQpThezbX9W6hAqtCbK3HmRIPmtoj5yqP8Aatp8Og9kNLEd4a1hJYAZe2gbeMf7VrcfMo9kNH/H3B/z8wRn/sIRRn9J0P54oH+lwbf+PmCM8gf8TCJD7qcTofpnFZHe/P8Az18+78ugQkgta25OSPmt7ZtpP1tLjgD/AHWpYD5Un2a2PlSHj7Pat9nc/wDbrPlCPdT9O1C5uoGjXNzAn3kU/b4kPvG+JkPsM4/KhP8ASYXgiJnhT78UbC9jX/egkxKn0U0f1/X+Yn1v87/r3faKHRObabyoSsMzf8soCbKVv+2MuYn+oNPEq20+07Ladj0O6wlY/Q5ilP6VHETLG9vATLGPvwQOLhAP9u2mw6D2Q0W7Z3W9qxPHz29o/b/atJ+VHspzSsv6/T/Mlpa3+f8Awey7RRZeVYJt84WCZv47lTZyH/ttHmNz7VNNJtIluhtOOJLyPYce1xDwB9eaoW58uQwWh8uTHMNm/wBnfHvaz/Jj/dP9KdC4guPLhKQ3DH/V27GxmY/9cZMxSH3zUuC/r+tWQ6ev9LT/ANtX4mm7s8SvMHeLqrTxi5j/AOAyR/Mv1apI5WkgyhkeDuY2F7D+X3/wrL81ILj5/LguGP8Ay0DWEzH/AHhmOU/pU00gjlEl2BHKej3qG2k/C4i+Q/SpcOv9ff8AojF0tl/XyX6stxFS4SAWpt2ljafyJD8u1gQSh5XoBjmrEv8ApNvKvU6hfeV9Y04b8CsbfnVGZyY1a7UmPHytewiVMf7M8f3fq3NS2tz9mVbmVIls7K3laBophKrk46N1JABHIyd1OC+y/wCv8tG9+5hKLj7y/p/rqlsraF03nkzavq+3eLdBbQJ/fYckD/eZlX/gNIbPb/ZOhbt+z/Srtv7wQ55/3pSD9A1Jb2rhtI0uTl4wb67x3fOQPxkYkf7lRpeSfYb7WIMG61GUW1jn+6CVjP0yXk+hroXd/wBdX+iEo20j8v8A0mPp9qRZjuoWvdR124P+i2aNbQEDOQpzKw9ywC/8AqIfa7DRgowus6rNk9/Ldhz+EaD/AMcHrTzZxm50/QYMm1skS4uSed2D+7U+pLAsf9z3pYruGS4vtfuWP2O1RoLbjOVB/eOB3LMNo9Qox1qv6/z/AMhadFdaaeW0V837zXVajb23hKWvhy2H+jRxCS6Gc/uRwFPu7DB9QHp23+1NZWLrZ6cwZ/SS4xkD6IDn/eI/u1Cz3OnaZ5rIG1nVJsiMnIViPlX/AHUUZP8Ause9E9uLe2tvDtk7CSdS1xN/Esef3jk9nckge5J7U1/X6Im3RPvr/wClS/Rd/VEtpImpX8mrysFsbUPHaluAe0kv042g+gJ/iptnKjifxFfnyovLItlcYMUHXOP7zkA/TaO1LcxR6hdx6LAgWwtVU3YXgYxlIh9cAn2wP4qc3/E71Xbw2nWMnzZHE06nj6qn/oX+7T1/rv8A8ATtbXRW+6PRest/x2bJdJtpneTVL1Cl1cgBYm/5YRD7qfXu3ufYVqEbmVfU8/SlpYxmRm9Bgf5/KtkrI4KlR1JczJaKKKRAUUUUAMHzTeyj9T/n9akpkX3N3945p9ABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFAEZ+WYH+8MfiKfTZfuZHVTmndRmgAooooAKgA2sy+hyPxqeopRh1b14NNAJRRRTEFIw3KR60tFAEiNuRW9RzTqiiOCy++RUtSxhQSACScAdSaKqy/6VMYBzCn+tPqey/40AJEgupPtMi/IOIlPp/e+p/lVjyY/7i/lT6KAM+80PTL9hJPaIZQMLMmUkX6OuGH4Gsi88Ku7B4poror90XiYkUeizR4dfqd1dPRUyhF7nRSxVWn8L/ry7fI4G8sZoSGv7WQbOj3UXnqv+7cRYkT6uKiaNLiBJ5EE8K/cllUXsS/7s8WJU+rCvQ6zrrQtPupzcGDybk9bi3YxSH6suCfocispUX0Z30syjopq3p+i/XVnGlFuLdZHX7Rbr92R1GoQqf8AZlTEyH3YcfnS+WtzbZ2/aLZD1wNSgU/UYmQ/XOPwrau/Dd0JvtEMkN3IP+Wkube4/wC/0WM/Qqfesm8ikt5fOv42jkXpNfJ5bD6XcHAH++M1lKMlv/X9djup1qdS3I/u/RPr3kyNUW4t2VF+0Wyn5gmNRgX6o2JkPsOlLEizwvHAPOhX78duwvYl/wB+CTEifRTSzjcEuboZGPkuLtcjH+zdwcqPdhmiYGWNJ7kGWIf6ua6QXEYH+xcw4dB7uKX9f12L0/rXX/25+eyC3VWDw2nzAD95FZSCQAf7drP8yj2U5pLZUWUw2YVJQPmisH8iTHvaT/Lj6HNOkzcWySy5nthykk6i/gH+7NHiVf8AeYUpzc2gJ3XFop43AanbA+xGJlPuen4Uf1/XYGlrf/P7/wCZ+S0QkQjjufLiWJLlj9yAmwuD/wBspP3ch984pZlt/OAuY7ZbgsCou4jZTMR0+cZjl+nSiMtPassRee1HDrCw1G3H+9E+JV/3V6dKW2Zmjkjs2do14kjsJhcIP9+2m+ZB7JzSsn/X9fexNbt/16tbf4Y6mi8t1LaXjyIbfUtUmSzjQMC0SBeoIyOB5rj61bLWyam8zERaboduVH90SFef++Y8D/gZrBsnSzvhc6fYWl1dQow8i1ka3cA9c20h+Q+4OeoxWtBamZdO0YyLKSPt+oyL0cltwH0aTPH91CK0i7/1/XU46kFB2ei+52tZ6Xf2VZX1bYuLuLS1jBaLVtalLMf4oFI5/wC/cYA9N2PWrRtobrUbbSoEC6dpao8qDoZAP3Sf8BGHP/AKgS+G288QFPN3/wCi6fEOrruwMf778/7qqabNavBZQ6EkrPdXYae/nUkHYT+8IPYsTtUdh0+7Vf1/XqZu99dHr8nbX/wBaJd9hYr2GWS78R3JP2OFDFZgDOUB5cDuXbAHqAvrQj3GmWL3csayazqMgCQ543Y+VP8AdRckn2Y96cojv9RCfKmm6UeeytOB+W1B/wCPH/ZpkN1GVm8R3odYQmyziI+YRkjBA/vSHGB6bR61f9f5v/Ii1+nbT/0mPz3l/mOkik0+0g0axlZr+7LPLc4G5AfvzH3ycAepA6Ctm0tYbK0itoE2xRKFUf4+p96p6TZSwrLeXoX7fdENLtOQgH3YwfRQfxJJ71pVpBdTjxFS75E79W+7/wAui+/qBOBk9qfEMRjPU8mo2G4qv94/pU9UzmCiiikMKZJ9zaOrHFPpn3pv90fqaAJKKKKACiiigAooooAKKKKACiiigAooooAKKKKACo4+F2/3Tj/CpKjPyzezD9R/n9KAH0UUUAFMkXdGQOvUU+igCAHIBHelpANrMvoePpS1QgooooATO2RW/A/jU9QMMqRTmnWO385+AB0HXPpSYDbiVl2xRY86T7uew7k+wp8MSwxCNOg6k9Se5NMt4mG6aX/XP1H90dl/z3qekMKKKKACiiigAooooAKKKKAMuXw/YNK01sj2U7HJls3MRY/7QHyt/wACBrIn8O3tvM09v5NxIeTLCxs7g/Up8kh9mUCuroqHTizqp4ytDrf1/wA9/wATz+cfZrnzboCC4Y48y6BspmP/AF8RZikPsR/OluB5cqz3YEcp+7Ner9mkP0uofkI9iOe9d8yq6lWUMpGCCMg1kv4ds0LNYPNp7tyfsrbUP1jIKH/vmsnRfQ7qeYwfxq39fgvRN+Zy9yMbZ70YIHyTX6bSB/s3kH3R/vDP5UtwC8cct2C8Q5jlvYhcRj/cuYfmQe7jNaMmg39i7SW0SMScmTTpPssh+sbZic+5xWYNtvdfwQXTntnTrhz+sUx/T8qyaa3/AK/zO6nUhNXg727fp/KvN6j5N09orTBprQco1wg1G3HusqYlX/ebpTrESfZprSwt0WDUHWOW+gvftCRjo2GY7wSMgAjhj1qOTbDdb5gsN0xxunBsLhj/ANdUzFKfbpSXaIJN98qLKwwHv0+yyn2F1D8jfQihf1/X+QON1ZfLf8rq/rLQ3jLbG8lvHxHpejIY4lA4MoGGIH+yPkHuzDtUDteWtoGwF1rWJPlUjIgUDgH2jXr6sT/eqLT0aZbWyurVrPTdPX7SxlkVxcEHKncCcgH5yT1O33p63sgWXW3iLXl7i3062f5SEPKg+hP329AP9mtl3/r+uiOFrotenT5J+rvKXTzHvZxTND4dtc/YrdQ96zclweRGT3Zzlm9v96rFuP7Z1MXZ/wCQfZuVt1/hmkHBk+i8qvvk+lVjbPDHHoVtMzXVxme/uQcMqsfmb2Zj8qjsB/s1vwwx28EcEKKkUahERRwoHAFXGN3/AF93yOatV5I6O7d7fPeXrLp5fJklFFHQZrY88IxmRj/dGP8AP6VNTIhiME9TyafSYwooopAFNi5Ut/eOf8KST7mB1bgVIBgYHagAooooAKKKKACiiigAooooAKKKKACiiigAooooAKZL9zd/dOafRQAlFMj+5tPVTin0AFFFFAEUgxIrevB/z+dJT5RmM46jkUwHIyO9NCCiiimAVViLTzeZjdDGx8sZ+8e7fzxSzkzSfZkJAIzIw7D0+pqwAFUKBgAYAFADvNb+4P8AvqjzW/uD/vqkooAXzW/uD/vqjzW/uD/vqkooAXzW/uD/AL6o81v7g/76pKKAF81v7g/76o81v7g/76pKKAF81v7g/wC+qPNb+4P++qSigBfNb+4P++qPNb+4P++qSigBfNb+4P8AvqjzW/uD/vqkooAXzW/uD/vqmTBLiJop7eOWNuGR8MD+BFOopWQ02ndGNJoMSRsthJNZKRgwo4eE+3lsCoH+7isttK1HTlbyLUqh6tpcoRT9beTKY+hzXW0VDpRe2h108dWjpJ39f16v56HnvkwSyiEWkDzbtxt4HaykkwckNbuQjg9yCM/jW5FePiXXL22lWVSbezsmxv3E88f3nP5KAfWugurO2vYTDd28U8Z6pKgYfrWUuhWenX0eoQidliBAheVnSPPBZQScEDj6E4qFScdjoljadRXne/3/ACvdWvtoti5pVg9lbu9w4lvJ28y4kHdvQf7IGAPYfWr9AIIBByD0IorZKysjzpzc5OUgprDdhf7xxTqIxmUn+6MUyCaiiipGFFFFADPvTAdlGfx/zmpKZFyC/wDeOfwp9ABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAR/dm9mH6in02XhQ3905/wAadQAUUUUAFV1G0sv90/pVioZBiRT/AHhj/P600AVFPKYkAUbpGOEX1P8AhT3dY0Z3ICqMkntUMCM7m4kGGYYRT/Cv+J70xEkEIhjxncxOWY9z61JRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAVo/9GlEJ/wBU5/dn0P8Ad/wqzTJYlmjKP0PcdQfWmW8rNuilx5yfex3HY/jQBNToRiPPduajbkBf7xxVikwCiiikMKbISEIHU8CnUw/NKo/ujJ/p/WgB4AAAHQUtFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUABGRg96jj+5g9V4NSVH92Yj+8M/iP8AIoAfRRRQAUyUZjJHUcin1XuJGZhbxEiRxksP4F9fr6UAV8/a5sDmCM8/7bf4D+f0q1SpbQxoEVAFAwBS+TH/AHRTuIbRTvJj/uijyY/7oouA2ineTH/dFHkx/wB0UXAbRTvJj/uijyY/7oouA2ineTH/AHRR5Mf90UXAbRTvJj/uijyY/wC6KLgNop3kx/3RR5Mf90UXAbRTvJj/ALoo8mP+6KLgNop3kx/3RR5Mf90UXAbRTvJj/uijyY/7oouA2ineTH/dFHkx/wB0UXAbRTvJj/uijyY/7oouA2obiNjtlix5qdP9odxVjyY/7oo8mP8Auii4EVvIs5EifdA79j/iKs1UcCzl81RiBziQf3T/AHv8at0mMKKKKACmx87n/vHj6USEhDjqeB9acoCqFHQDFAC0UUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUyXgB/wC6c/hT6QgEEHoaACimxklAD1HBpxIAJJwB1NAEc8wgiLkFj0VR1Y9hTbeExqzOd0rnLn+g9hUcINxL9pYEKMiIH0/vfj/L61aoAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKAEIDKVIBBGCDVaAmCT7K5JGMxMe49PqP5VaqKeETx7clWB3Kw/hI6GgCWiobeYyoQ42yocOvof8KmoAYfmlUf3fmP8ASpKZHyWf1PH0FPoAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKAIx8srD+8Mj+v9Kgl/0qY245hX/Wn1PZf8f/AK9S3CymPMG3zB03dKSGMwxBFj6dSW5J7k0wJqKZuk/55/8Aj1G6T/nn/wCPUgH0UzdJ/wA8/wDx6jdJ/wA8/wDx6gB9FM3Sf88//HqN0n/PP/x6gB9FM3Sf88//AB6jdJ/zz/8AHqAH0UzdJ/zz/wDHqN0n/PP/AMeoAfRTN0n/ADz/APHqN0n/ADz/APHqAH0UzdJ/zz/8eo3Sf88//HqAH0UzdJ/zz/8AHqN0n/PP/wAeoAfRTN0n/PP/AMeo3Sf88/8Ax6gB9FM3Sf8APP8A8eo3Sf8APP8A8eoAfRTN0n/PP/x6jdJ/zz/8eoAfRTN0n/PP/wAeo3Sf88//AB6gB9FM3Sf88/8Ax6jdJ/zz/wDHqAH0UzdJ/wA8/wDx6jdJ/wA8/wDx6gB9FM3Sf88//HqN0n/PP/x6gCG4RkcXMS5dRh1H8a+n1Hb/AOvUvmq0HmRkMGHykd89KXdJ/wA8/wDx6oIoJEuDgAQE79ueje3setMC0o2qFHYYpaKKQBRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAVDdzNb2skqruKjOKmpkiu0ZCPsbs2M/pQBDHcFYPNmljkUn5GhUnd+HP6VDLeiSa1SGSQJIx3MsZJ47cjjnrTW0x23nzkDO4ZlEeEbA6Fc8/nT7fTvIEQ83Pls5+VMfe/ligCU39sN+ZCNiliShAIHccc/hRFcLNc/JKdpjyIzGR3+9k1T/sdjndcKWMbIWEfzNnuTnk1eW3xd+fu/wCWXl7ce+c0AT0UUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFAEF7M1vZyzJt3IuRu6VBHc72iJvYXyxGIk4bjOOpx61YuoPtNtJDu27xjOM4pJrbznhbfjyiT068YoAauoWzxq6s7KxwuI2Jb6DHNQterJeQJHKyxFC5YJwfYkjgdc0jaYrQW6bkZoBgGSPcrA+oz/WpBp6YVSwC+U0ZCqF69SMdKAHDULYgne2Au/mNhkeo45H0qaOeOZmEbbtuMkA45561V+xTnYWuVLRoUjIixjPBJ55OPpUtnaCyjaJHzHnKgjkevPegCzRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQB//2Q==",
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
              "w": 342,
              "h": 300
            },
            "backColor": "transparent",
            "zIndex": 1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.ChartPublic.ChartTitle",
                "type":4,
                "value":"Title",
                "key":"ChartTitle",
              },
              {
                "name":"configComponent.ChartPublic.ChartUnit",
                "type":4,
                "value":"%",
                "key":"ChartUnit",
              },
              {
                "name":"configComponent.ChartPublic.ChartTitleFontSize",
                "type":1,
                "value":20,
                "key":"ChartTitleFontSize",
              },
              {
                "name":"configComponent.ChartPublic.ChartTitleFontColor",
                "type":2,
                "value":"#000000",
                "key":"ChartTitleFontColor",
              },
              {
                "name":"configComponent.ChartPublic.ChartAxisTickColor",
                "type":2,
                "value":"#77EADF",
                "key":"ChartAxisTickColor",
              },
              {
                "name":"configComponent.category3.AreaColor",
                "type":2,
                "value":"#5B7AD8",
                "key":"AreaColor",
              },
              {
                "name":"configComponent.ChartPublic.EchartsWidth",
                "type":1,
                "value":2,
                "key":"EchartsWidth",
              },
            ]
          }
        }
      },
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      fillOpacity:1,
      strokeOpacity:1,
      animateType:"blink",
      EchartsWidth:6,
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      eventValue: '0.00',
      eventUnit: '',
      option: {
        title: {
          text: 'Waterfall Chart',
          textStyle:{

          }
        },
        color: ['#67F9D8'],
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          },
          formatter: function (params) {
            let tar = params[0];
            return tar.name +' : ' + tar.value;
          }
        },
        grid: {
          left: '15%',
          right: '15%',
          top:'15%',
          bottom: '15%',
          containLabel: true
        },
        radar: {
          shape: 'circle',
          indicator: [
            { text: 'Indicator1', max: 150 },
            { text: 'Indicator2', max: 150 },
            { text: 'Indicator3', max: 150 },
            { text: 'Indicator4', max: 120 },
            { text: 'Indicator5', max: 108 },
            { text: 'Indicator6', max: 72 }
          ],
          axisName: {
            color: '#fff',
            borderRadius: 3,
          }
        },
        series: [
          {
            name: 'Budget vs spending',
            type: 'radar',
            itemStyle: {
              normal: {
                lineStyle: {
                  width:6
                }
              }
            },
            data: [
              {
                value: [5, 75, 50, 36, 20, 69],
                areaStyle: {
                  color: '#5B7AD8'
                },
                label: {
                  show: true,
                  color:'#ffffff',
                  formatter: function (params) {
                    return params.value;
                  }
                }
              },
            ]
          }
        ]
      }
    }
  },
  beforeDestroy () {
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.dispose()
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
      this.option.color[0]=[]
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.title.text=option.style.diy[i].value
          if(this.option.title.text!="")
          {
            this.option.grid.top='15%'
          }
          else
          {
            this.option.grid.top='4%'
          }
        }
        else if(option.style.diy[i].key=="ChartUnit")
        {
          // this.option.series[0].detail.formatter='{value}'+option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartTitleFontSize")
        {
          this.option.title.textStyle.fontSize = option.style.diy[i].value
          this.option.radar.axisName.fontSize = option.style.diy[i].value
          this.option.series[0].data[0].label.fontSize = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartTitleFontColor")
        {
          this.option.title.textStyle.color = option.style.diy[i].value
          this.option.radar.axisName.color = option.style.diy[i].value
          this.option.series[0].data[0].label.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartAxisTickColor")
        {
          this.option.color[0] = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="AreaColor")
        {
          this.option.series[0].data[0].areaStyle.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="EchartsWidth")
        {
          this.option.series[0].itemStyle.normal.lineStyle.width = parseInt(option.style.diy[i].value)
          this.option.series[0].symbolSize = parseInt(option.style.diy[i].value)
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
      if(!this.editMode)
      {
        this.option.radar.indicator=[]
        this.option.series[0].data[0].value=[]
        for(let i =0;i<this.detail.active.length;i++)
        {
          if(this.detail.active[i].condition.dataName=="")
          {
            continue
          }
          // this.option.yAxis.data.push(this.detail.active[i].condition.dataName)
          let series= {
            name: this.detail.active[i].condition.DeviceName+"-"+this.detail.active[i].condition.dataName,
            dataID:this.detail.active[i].condition.deviceSN+this.detail.active[i].condition.dataID,
          }
          this.option.radar.indicator.push(series)
        }
      }
      this.echartsView.setOption(this.option,true);
    },
    UpdateChartData : function(data) {
      for(let i=0;i<this.option.radar.indicator.length;i++)
      {
        if(this.option.radar.indicator[i].dataID == (data.DeviceSN+data.dataID))
        {
          this.option.series[0].data[0].value[i] = data.result
          this.option.radar.indicator[i].max= parseFloat(data.result)+100
        }
      }
      // 重新将数组赋值给echarts选项
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
    this.detail = node.getData().detail
    this.editMode = node.getData().editMode
    this.showDeviceUuid = node.getData().showDeviceUuid
    this.IsToolBox = node.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
      // _t.initComponents(_t.detail)
    })
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据

      _t.$EventBus.$on(activeEvent, (data) => {
        _t.UpdateChartData(data)
        // if(data.ID == "ShowData")
        // {
        //   this.option.series[0].data[0].value = data.result
        //   this.echartsView.setOption(this.option, true);
        // }
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
