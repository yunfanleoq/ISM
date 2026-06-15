<template>
  <div class="three-model-editor">
    <a-collapse :defaultActiveKey="['basic', 'transform', 'lighting', 'binding']" :bordered="false">
      <!-- 基础设置 -->
      <a-collapse-panel key="basic" :header="$t('基础设置')">
        <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
          <a-form-item :label="$t('模型路径')">
            <a-input v-model="localData.modelUrl" :placeholder="$t('输入模型URL或上传')" @change="onDataChange" />
            <a-upload
              :showUploadList="false"
              :beforeUpload="beforeModelUpload"
              :http-request="handleModelUpload"
              accept=".gltf,.glb,.obj,.fbx,.stl"
            >
              <a-button size="small" style="margin-top: 4px;">{{ $t('上传模型') }}</a-button>
            </a-upload>
          </a-form-item>

          <a-form-item :label="$t('宽度')">
            <a-input-number v-model="localData.width" :min="50" :max="2000" @change="onDataChange" />
          </a-form-item>

          <a-form-item :label="$t('高度')">
            <a-input-number v-model="localData.height" :min="50" :max="2000" @change="onDataChange" />
          </a-form-item>

          <a-form-item :label="$t('自动旋转')">
            <a-switch v-model="localData.autoRotate" @change="onDataChange" />
          </a-form-item>

          <a-form-item :label="$t('自动播放动画')">
            <a-switch v-model="localData.autoPlay" @change="onDataChange" />
          </a-form-item>
        </a-form>
      </a-collapse-panel>

      <!-- 变换 -->
      <a-collapse-panel key="transform" :header="$t('变换')">
        <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
          <a-divider orientation="left">{{ $t('位置') }}</a-divider>
          <a-form-item label="X">
            <a-input-number v-model="localData.position.x" :step="0.1" @change="onDataChange" />
          </a-form-item>
          <a-form-item label="Y">
            <a-input-number v-model="localData.position.y" :step="0.1" @change="onDataChange" />
          </a-form-item>
          <a-form-item label="Z">
            <a-input-number v-model="localData.position.z" :step="0.1" @change="onDataChange" />
          </a-form-item>

          <a-divider orientation="left">{{ $t('旋转') }}</a-divider>
          <a-form-item label="X">
            <a-input-number v-model="localData.rotation.x" :step="1" @change="onDataChange" />
          </a-form-item>
          <a-form-item label="Y">
            <a-input-number v-model="localData.rotation.y" :step="1" @change="onDataChange" />
          </a-form-item>
          <a-form-item label="Z">
            <a-input-number v-model="localData.rotation.z" :step="1" @change="onDataChange" />
          </a-form-item>

          <a-divider orientation="left">{{ $t('缩放') }}</a-divider>
          <a-form-item label="X">
            <a-input-number v-model="localData.scale.x" :step="0.1" :min="0.1" @change="onDataChange" />
          </a-form-item>
          <a-form-item label="Y">
            <a-input-number v-model="localData.scale.y" :step="0.1" :min="0.1" @change="onDataChange" />
          </a-form-item>
          <a-form-item label="Z">
            <a-input-number v-model="localData.scale.z" :step="0.1" :min="0.1" @change="onDataChange" />
          </a-form-item>
        </a-form>
      </a-collapse-panel>

      <!-- 相机 -->
      <a-collapse-panel key="camera" :header="$t('相机')">
        <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
          <a-form-item :label="$t('相机位置X')">
            <a-input-number v-model="localData.cameraPosition.x" :step="0.5" @change="onDataChange" />
          </a-form-item>
          <a-form-item :label="$t('相机位置Y')">
            <a-input-number v-model="localData.cameraPosition.y" :step="0.5" @change="onDataChange" />
          </a-form-item>
          <a-form-item :label="$t('相机位置Z')">
            <a-input-number v-model="localData.cameraPosition.z" :step="0.5" @change="onDataChange" />
          </a-form-item>
          <a-form-item :label="$t('背景透明')">
            <a-switch :checked="localData.backgroundAlpha === 0" @change="onBackgroundAlphaChange" />
          </a-form-item>
        </a-form>
      </a-collapse-panel>

      <!-- 光照 -->
      <a-collapse-panel key="lighting" :header="$t('光照')">
        <div v-for="(light, index) in localData.lights" :key="index" class="light-item">
          <a-divider orientation="left">{{ $t('光照') }} {{ index + 1 }}</a-divider>
          <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
            <a-form-item :label="$t('类型')">
              <a-select v-model="light.type" @change="onDataChange">
                <a-select-option value="AmbientLight">{{ $t('环境光') }}</a-select-option>
                <a-select-option value="DirectionalLight">{{ $t('方向光') }}</a-select-option>
                <a-select-option value="PointLight">{{ $t('点光源') }}</a-select-option>
                <a-select-option value="SpotLight">{{ $t('聚光灯') }}</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item :label="$t('颜色')">
              <a-input type="color" v-model="light.colorHex" @change="onColorChange(light, index)" />
            </a-form-item>
            <a-form-item :label="$t('强度')" v-if="light.type !== 'AmbientLight'">
              <a-input-number v-model="light.intensity" :step="0.1" :min="0" :max="10" @change="onDataChange" />
            </a-form-item>
          </a-form>
          <a-button type="danger" size="small" @click="removeLight(index)">{{ $t('删除') }}</a-button>
        </div>
        <a-button type="dashed" block @click="addLight">{{ $t('添加光照') }}</a-button>
      </a-collapse-panel>

      <!-- 数据绑定 -->
      <a-collapse-panel key="binding" :header="$t('数据绑定')">
        <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
          <a-form-item :label="$t('启用绑定')">
            <a-switch v-model="localData.bindData.enabled" @change="onDataChange" />
          </a-form-item>

          <template v-if="localData.bindData.enabled">
            <a-form-item :label="$t('旋转绑定')">
              <a-input v-model="localData.bindData.rotationBind" :placeholder="'device001.rotate_speed'" @change="onDataChange" />
            </a-form-item>
            <a-form-item :label="$t('颜色绑定')">
              <a-input v-model="localData.bindData.colorBind" :placeholder="'device001.temperature'" @change="onDataChange" />
            </a-form-item>
            <a-form-item :label="$t('显隐绑定')">
              <a-input v-model="localData.bindData.visibleBind" :placeholder="'device001.online'" @change="onDataChange" />
            </a-form-item>
          </template>
        </a-form>
      </a-collapse-panel>
    </a-collapse>
  </div>
</template>

<script>
import { METHOD, request } from '@/utils/request'
import { SYSTEMIMAGEUPLOAD } from '@/services/api'

export default {
  name: 'ThreeModelEditor',
  i18n: require('@/i18n/language'),
  props: {
    componentData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      localData: {
        modelPath: '',
        modelUrl: '',
        width: 300,
        height: 200,
        autoRotate: false,
        autoPlay: true,
        backgroundColor: 0x000000,
        backgroundAlpha: 0,
        position: { x: 0, y: 0, z: 0 },
        rotation: { x: 0, y: 0, z: 0 },
        scale: { x: 1, y: 1, z: 1 },
        cameraPosition: { x: 0, y: 0, z: 5 },
        lights: [
          { type: 'AmbientLight', color: 0x404040, colorHex: '#404040' }
        ],
        bindData: {
          enabled: false,
          rotationBind: '',
          colorBind: '',
          visibleBind: ''
        }
      }
    }
  },
  watch: {
    componentData: {
      handler(val) {
        if (val && val.detail) {
          this.localData = { ...this.getDefaultData(), ...val.detail }
          // 确保 colorHex 存在
          if (this.localData.lights) {
            this.localData.lights.forEach(light => {
              if (!light.colorHex) {
                light.colorHex = '#' + (light.color || 0xffffff).toString(16).padStart(6, '0')
              }
            })
          }
        }
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    getDefaultData() {
      return {
        modelPath: '',
        modelUrl: '',
        width: 300,
        height: 200,
        autoRotate: false,
        autoPlay: true,
        backgroundColor: 0x000000,
        backgroundAlpha: 0,
        position: { x: 0, y: 0, z: 0 },
        rotation: { x: 0, y: 0, z: 0 },
        scale: { x: 1, y: 1, z: 1 },
        cameraPosition: { x: 0, y: 0, z: 5 },
        lights: [
          { type: 'AmbientLight', color: 0x404040, colorHex: '#404040' }
        ],
        bindData: {
          enabled: false,
          rotationBind: '',
          colorBind: '',
          visibleBind: ''
        }
      }
    },
    onDataChange() {
      this.$emit('update', {
        detail: { ...this.localData }
      })
    },
    onBackgroundAlphaChange(checked) {
      this.localData.backgroundAlpha = checked ? 0 : 1
      this.onDataChange()
    },
    onColorChange(light, index) {
      const hex = light.colorHex || '#ffffff'
      light.color = parseInt(hex.replace('#', ''), 16)
      this.onDataChange()
    },
    addLight() {
      this.localData.lights.push({
        type: 'DirectionalLight',
        color: 0xffffff,
        colorHex: '#ffffff',
        intensity: 1
      })
      this.onDataChange()
    },
    removeLight(index) {
      this.localData.lights.splice(index, 1)
      this.onDataChange()
    },
    beforeModelUpload(file) {
      const isValid = file.name.match(/\.(gltf|glb|obj|fbx|stl)$/i)
      if (!isValid) {
        this.$message && this.$message.error(this.$t('只支持 glTF/GLB/OBJ/FBX/STL 格式'))
      }
      return isValid
    },
    handleModelUpload({ file }) {
      const formData = new FormData()
      formData.append('file', file)
      request(SYSTEMIMAGEUPLOAD, METHOD.POST, formData).then((res) => {
        const result = res && res.data ? res.data : {}
        if (result.Code !== 200 && result.Code !== 2002) {
          throw new Error(result.Message || 'Model upload failed')
        }
        const path = result.Path || result.path || ''
        this.localData.modelPath = path
        this.localData.modelUrl = path
        this.onDataChange()
        this.$message && this.$message.success('Model uploaded successfully')
      }).catch((error) => {
        this.$message && this.$message.error(error.message || 'Model upload failed')
      })
    }
    }
  }
}
</script>

<style scoped>
.three-model-editor {
  padding: 8px;
}
.light-item {
  margin-bottom: 8px;
  padding: 8px;
  border: 1px solid #eee;
  border-radius: 4px;
}
</style>
