<template>
  <div class="upload-panel">
    <div class="panel-header">
      <h2 class="panel-title">
        <span class="icon">&#8683;</span>
        上传图片
      </h2>
      <p class="panel-desc">选择一张图片开始 AI 修图</p>
    </div>

    <div
      class="dropzone"
      :class="{ 'dropzone-active': dragging, 'dropzone-has-file': session.originalUrl }"
      @dragover.prevent="dragging = true"
      @dragleave.prevent="dragging = false"
      @drop.prevent="onDrop"
    >
      <div v-if="!session.originalUrl" class="dropzone-empty">
        <div class="upload-icon">
          <span>&#128247;</span>
        </div>
        <p class="upload-hint">拖拽图片到此处</p>
        <p class="upload-sub">或</p>
        <label class="upload-btn">
          选择文件
          <input type="file" accept="image/*" @change="onFileChange" />
        </label>
        <p class="upload-formats">支持 JPG、PNG、WebP</p>
      </div>

      <div v-else class="dropzone-preview">
        <img :src="session.originalUrl" alt="原图" />
        <div class="preview-overlay">
          <label class="change-btn">
            更换图片
            <input type="file" accept="image/*" @change="onFileChange" />
          </label>
        </div>
      </div>
    </div>

    <button
      v-if="!session.originalUrl"
      class="submit-btn"
      :disabled="!selectedFile || loading"
      @click="uploadImage"
    >
      <span v-if="loading" class="spinner"></span>
      <span v-else>&#8593;</span>
      {{ loading ? '上传中...' : '上传到编辑器' }}
    </button>

    <p v-if="statusMessage" class="status" :class="{ 'status-error': isError }">
      {{ statusMessage }}
    </p>

    <div v-if="session.originalUrl" class="image-info">
      <div class="info-row">
        <span class="info-label">状态</span>
        <span class="info-value success">已就绪</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useSessionStore } from '@/stores/session'

const selectedFile = ref<File | null>(null)
const loading = ref(false)
const dragging = ref(false)
const statusMessage = ref('')
const isError = ref(false)
const session = useSessionStore()

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    selectedFile.value = file
    uploadImage()
  }
}

function onDrop(event: DragEvent) {
  dragging.value = false
  const file = event.dataTransfer?.files?.[0]
  if (file) {
    selectedFile.value = file
    uploadImage()
  }
}

async function uploadImage() {
  if (!selectedFile.value) return

  const formData = new FormData()
  formData.append('image', selectedFile.value)
  loading.value = true
  isError.value = false
  statusMessage.value = ''

  try {
    const response = await axios.post('/api/upload', formData)
    const data = response.data as { imageId: string; originalUrl: string }
    session.setImage(data.imageId, data.originalUrl)
    statusMessage.value = ''
  } catch {
    isError.value = true
    statusMessage.value = '上传失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.upload-panel {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.panel-header {
  margin-bottom: 0.2rem;
}

.panel-title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.panel-title .icon {
  font-size: 1.1rem;
  color: var(--accent);
}

.panel-desc {
  font-size: 0.8rem;
  color: var(--text-secondary);
  margin-top: 0.2rem;
}

.dropzone {
  border: 2px dashed var(--border);
  border-radius: var(--radius-lg);
  background: var(--bg-card);
  transition: all var(--transition);
  overflow: hidden;
  min-height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dropzone-active {
  border-color: var(--accent);
  background: var(--bg-card-hover);
  box-shadow: 0 0 0 4px var(--accent-glow);
}

.dropzone-has-file {
  border-style: solid;
  border-color: var(--border);
}

.dropzone-empty {
  text-align: center;
  padding: 2rem 1.5rem;
}

.upload-icon span {
  font-size: 2.8rem;
  display: block;
  margin-bottom: 0.8rem;
  opacity: 0.8;
}

.upload-hint {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.3rem;
}

.upload-sub {
  font-size: 0.8rem;
  color: var(--text-muted);
  margin-bottom: 1rem;
}

.upload-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.6rem 1.4rem;
  background: var(--accent-gradient);
  color: #fff;
  border-radius: var(--radius);
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition);
}

.upload-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px var(--accent-glow);
}

.upload-formats {
  font-size: 0.72rem;
  color: var(--text-muted);
  margin-top: 0.8rem;
}

.dropzone-preview {
  position: relative;
  width: 100%;
}

.dropzone-preview img {
  width: 100%;
  border-radius: var(--radius-lg);
  object-fit: cover;
  max-height: 260px;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.45);
  opacity: 0;
  transition: opacity var(--transition);
  border-radius: var(--radius-lg);
}

.dropzone-preview:hover .preview-overlay {
  opacity: 1;
}

.change-btn {
  padding: 0.5rem 1.2rem;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(8px);
  color: #fff;
  border-radius: var(--radius);
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all var(--transition);
}

.change-btn:hover {
  background: rgba(255, 255, 255, 0.25);
}

.submit-btn {
  width: 100%;
  padding: 0.8rem;
  background: var(--accent-gradient);
  color: #fff;
  font-size: 0.9rem;
  font-weight: 600;
  border-radius: var(--radius);
  box-shadow: 0 2px 12px var(--accent-glow);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 24px var(--accent-glow);
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.status {
  font-size: 0.82rem;
  padding: 0.6rem 0.9rem;
  border-radius: var(--radius-sm);
  background: rgba(34, 197, 94, 0.1);
  color: var(--success);
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.status-error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.image-info {
  padding: 0.8rem;
  background: var(--bg-card);
  border-radius: var(--radius);
  border: 1px solid var(--border);
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 0.82rem;
  color: var(--text-secondary);
}

.info-value {
  font-size: 0.8rem;
  font-weight: 600;
}

.info-value.success {
  color: var(--success);
}
</style>
