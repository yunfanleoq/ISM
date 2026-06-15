<template>
  <div class="local-model-browser">
    <!-- 顶部操作栏 -->
    <div class="lmb-toolbar">
      <a-button size="small" type="primary" @click="onSelectFiles">
        <i class="fas fa-folder-open" style="margin-right:4px;"></i>选择模型文件
      </a-button>
      <a-button size="small" @click="loadManifest">
        <i class="fas fa-sync-alt" style="margin-right:4px;"></i>刷新列表
      </a-button>
      <input
        ref="fileInput"
        type="file"
        accept=".glb,.gltf"
        multiple
        style="display:none"
        @change="onFilesSelected"
      />
    </div>

    <!-- 拖拽上传区 -->
    <div
      class="lmb-dropzone"
      :class="{ dragover: dragOver }"
      @dragover.prevent="onDragOver"
      @dragleave.prevent="onDragLeave"
      @drop.prevent="onDrop"
    >
      <i class="fas fa-cloud-upload-alt" style="font-size:24px;color:#13c2c2;display:block;margin-bottom:6px;"></i>
      <div style="font-size:12px;color:#666;">拖拽 .glb / .gltf 文件到此处</div>
      <div style="font-size:11px;color:#999;margin-top:3px;">直接加载到场景，无需上传服务器</div>
    </div>

    <!-- 本地模型列表（来自 manifest.json） -->
    <div v-if="manifestModels.length > 0" class="lmb-section-title">本地模型库（public/models/）</div>
    <div v-if="manifestModels.length > 0" class="lmb-grid scrollbar">
      <div
        v-for="m in manifestModels"
        :key="'mf_' + m.id"
        class="lmb-card"
        @click="onManifestModelClick(m)"
        :title="m.name"
      >
        <div class="lmb-card-thumb">
          <img v-if="m.thumbnail" :src="m.thumbnail" alt="thumb" />
          <div v-else class="lmb-card-placeholder"><i class="fas fa-cube"></i></div>
        </div>
        <div class="lmb-card-label">{{ m.name }}</div>
      </div>
    </div>

    <!-- 本次会话已加载的模型 -->
    <div v-if="sessionModels.length > 0" class="lmb-section-title">本次已加载</div>
    <div v-if="sessionModels.length > 0" class="lmb-grid scrollbar">
      <div
        v-for="m in sessionModels"
        :key="'sess_' + m.id"
        class="lmb-card"
        @click="$emit('add-object', m.payload)"
        :title="m.name"
      >
        <div class="lmb-card-thumb">
          <div class="lmb-card-placeholder"><i class="fas fa-check-circle" style="color:#52c41a;"></i></div>
        </div>
        <div class="lmb-card-label">{{ m.name }}</div>
      </div>
    </div>

    <div v-if="manifestModels.length === 0 && sessionModels.length === 0 && !loading" class="lmb-empty">
      暂无模型，请拖拽或选择 .glb / .gltf 文件
    </div>
  </div>
</template>

<script>
export default {
  name: 'LocalModelBrowser',
  props: {},
  data() {
    return {
      dragOver: false,
      loading: false,
      manifestModels: [],
      sessionModels: [],
      sessionCounter: 0,
    }
  },
  created() {
    this.loadManifest()
  },
  methods: {
    async loadManifest() {
      this.loading = true
      try {
        const res = await fetch('/models/manifest.json?t=' + Date.now())
        if (!res.ok) {
          console.warn('[LocalModelBrowser] manifest.json not found (' + res.status + ')')
          this.manifestModels = []
          return
        }
        const text = await res.text()
        let json
        try {
          json = JSON.parse(text)
        } catch (e) {
          console.warn('[LocalModelBrowser] manifest.json parse error')
          this.manifestModels = []
          return
        }
        this.manifestModels = (json.models || []).map((m, i) => ({
          id: m.id || 'mf_' + i,
          name: m.name || m.file || '本地模型',
          file: m.file || '',
          path: m.path || '/models/' + m.file,
          thumbnail: m.thumbnail ? (m.thumbnail.startsWith('http') ? m.thumbnail : '/models/' + m.thumbnail) : '',
        }))
      } catch (e) {
        console.warn('[LocalModelBrowser] 加载 manifest.json 失败', e)
        this.manifestModels = []
      } finally {
        this.loading = false
      }
    },
    onSelectFiles() {
      if (this.$refs.fileInput) {
        this.$refs.fileInput.click()
      }
    },
    onFilesSelected(e) {
      const files = e.target && e.target.files
      if (!files || files.length === 0) return
      Array.from(files).forEach(f => {
        if (/\.(glb|gltf)$/i.test(f.name)) {
          this.loadLocalFile(f)
        }
      })
      // 清空 input 以支持重复选择同一文件
      e.target.value = ''
    },
    onDragOver(e) {
      e.preventDefault()
      this.dragOver = true
    },
    onDragLeave(e) {
      e.preventDefault()
      this.dragOver = false
    },
    onDrop(e) {
      e.preventDefault()
      this.dragOver = false
      const files = e.dataTransfer && e.dataTransfer.files
      if (!files || files.length === 0) return
      Array.from(files).forEach(f => {
        if (/\.(glb|gltf)$/i.test(f.name)) {
          this.loadLocalFile(f)
        }
      })
    },
    loadLocalFile(file) {
      const reader = new FileReader()
      reader.onload = (ev) => {
        const arrayBuffer = ev.target.result
        const modelId = 'sess_' + (++this.sessionCounter)
        const name = file.name.replace(/\.(glb|gltf)$/i, '')
        // 保存到 sessionModels 以便重复添加
        this.sessionModels.push({
          id: modelId,
          name,
          payload: {
            type: 'gltf-buffer',
            buffer: arrayBuffer,
            name,
            fitSize: 2,
          }
        })
        // 立即添加到场景
        this.$emit('add-object', {
          type: 'gltf-buffer',
          buffer: arrayBuffer,
          name,
          fitSize: 2,
        })
        this.$message && this.$message.success('模型「' + file.name + '」已添加到场景！')
      }
      reader.onerror = () => {
        this.$message && this.$message.error('读取文件失败：' + file.name)
      }
      reader.readAsArrayBuffer(file)
    },
    onManifestModelClick(m) {
      if (!m || !m.path) return
      this.$emit('add-object', {
        type: 'gltf',
        modelPath: m.path,
        name: m.name || '本地模型',
        thumbnail: m.thumbnail || '',
        fitSize: 2,
      })
      this.$message && this.$message.success('模型「' + m.name + '」已添加到场景！')
    },
  },
}
</script>

<style scoped>
.local-model-browser {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  background: #fafafa;
}
/* 操作栏 */
.lmb-toolbar {
  display: flex;
  gap: 4px;
  padding: 6px 8px;
  border-bottom: 1px solid #eee;
  background: #fff;
  flex-shrink: 0;
}
/* 拖拽区 */
.lmb-dropzone {
  margin: 6px;
  padding: 16px 8px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  text-align: center;
  background: #fafafa;
  transition: all 0.2s;
  cursor: pointer;
  flex-shrink: 0;
}
.lmb-dropzone:hover,
.lmb-dropzone.dragover {
  border-color: #13c2c2;
  background: #e6fffb;
}
/* 分区标题 */
.lmb-section-title {
  font-size: 11px;
  color: #999;
  padding: 6px 8px 2px;
  flex-shrink: 0;
}
/* 模型网格 */
.lmb-grid {
  flex: 1;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(78px, 1fr));
  align-content: flex-start;
  padding: 6px;
  gap: 6px;
}
/* 卡片 */
.lmb-card {
  min-width: 0;
  cursor: pointer;
  border: 1px solid #eee;
  border-radius: 6px;
  padding: 3px;
  background: #fff;
  transition: all 0.2s;
}
.lmb-card:hover {
  border-color: #13c2c2;
  box-shadow: 0 2px 8px rgba(19,194,194,0.18);
  transform: translateY(-2px);
}
.lmb-card-thumb {
  width: 100%;
  height: 74px;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}
.lmb-card-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.lmb-card-placeholder {
  color: #ccc;
  font-size: 28px;
}
.lmb-card-label {
  margin-top: 3px;
  font-size: 11px;
  color: #555;
  text-align: center;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  line-height: 14px;
  height: 14px;
}
/* 空状态 */
.lmb-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #bbb;
  font-size: 12px;
  padding: 24px 0;
}
/* 滚动条 */
.scrollbar::-webkit-scrollbar {
  width: 5px;
}
.scrollbar::-webkit-scrollbar-thumb {
  border-radius: 10px;
  background-color: #91d5ff;
}
.scrollbar::-webkit-scrollbar-track {
  background: #ededed;
  border-radius: 10px;
}
</style>
