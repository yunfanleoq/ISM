<script>

export default {
    name: 'View',
    inject: ['getNode'],
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        selected: {
            type: Boolean,
            default: false
        },
        detail: {
            type: Object,
            default: function (){

            }
        },
        showDeviceUuid: {
          type: String,
          required: true
        },
        IsToolBox: {
          type: Boolean,
          default: false,
          required: true
        }
    },
    computed: {
        textAlign: function(){
            if(this.detail.style.textAlign == undefined) {
                return "center";
            } else {
                return this.detail.style.textAlign;
            }
        },
        lineHeight: function() {
            if(this.detail.style.lineHeight == undefined) {
                return this.detail.style.position.h;
            }
            return this.detail.style.lineHeight;
        }
    },
    data() {
        return {

        }
    },
    methods: {
        refreshData(val,sceneName) {

        },
        onResize() {

        },
    },
    created() {
      const node = this.getNode()
      this.detail = node.getData().detail
      this.editMode = node.getData().editMode
      this.showDeviceUuid = node.getData().showDeviceUuid
      this.IsToolBox = node.getData().IsToolBox
    }
}
</script>
