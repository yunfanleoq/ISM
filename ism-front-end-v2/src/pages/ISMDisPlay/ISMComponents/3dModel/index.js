export { default as ThreeModel } from './ThreeModel.vue'
export { default as ThreeModelEditor } from './ThreeModelEditor.vue'

// 组件配置
export const ThreeModelConfig = {
  type: '3dModel',
  name: '3D 模型',
  icon: 'el-icon-monitor',
  defaultData: {
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
