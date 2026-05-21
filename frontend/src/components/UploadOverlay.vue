<template>
  <div class="upload-overlay">
    <div class="upload-center">
      <div class="logo-area">
        <span class="logo-glow">✦</span>
        <h1>AI 修图助手</h1>
        <p>上传一张图片，开始 AI 修图之旅</p>
      </div>

      <div
        class="dropzone"
        :class="{ active: dragging }"
        @dragover.prevent="dragging = true"
        @dragleave.prevent="dragging = false"
        @drop.prevent="onDrop"
      >
        <div class="dropzone-content">
          <span class="drop-icon">📷</span>
          <p>拖拽图片到此处</p>
          <span class="divider">或</span>
          <label class="upload-btn">
            选择文件
            <input type="file" accept="image/*" @change="onFileChange" />
          </label>
        </div>
      </div>

      <p v-if="store.statusMessage" class="upload-status" :class="{ error: isError }">
        <span v-if="loading" class="spinner"></span>
        {{ store.statusMessage }}
      </p>
      <p class="format-hint">支持 JPG、PNG、WebP 格式</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSessionStore } from '@/stores/session'

const store = useSessionStore()
const dragging = ref(false)
const loading = ref(false)
const isError = ref(false)

const emit = defineEmits<{ uploaded: [] }>()

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) handleFile(file)
}

function onDrop(event: DragEvent) {
  dragging.value = false
  const file = event.dataTransfer?.files?.[0]
  if (file) handleFile(file)
}

async function handleFile(file: File) {
  loading.value = true
  isError.value = false
  try {
    await store.uploadAndCreateSession(file)
    emit('uploaded')
  } catch {
    isError.value = true
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.upload-overlay {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
  position: relative;
  overflow: hidden;
}

.upload-overlay::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(ellipse at center, rgba(99,102,241,0.06) 0%, transparent 70%);
  pointer-events: none;
}

.upload-center {
  text-align: center;
  max-width: 480px;
  width: 100%;
  padding: 2rem;
  position: relative;
  z-index: 1;
}

.logo-area {
  margin-bottom: 2.5rem;
}

.logo-glow {
  font-size: 3.5rem;
  display: block;
  margin-bottom: 1rem;
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  filter: drop-shadow(0 0 20px var(--accent-glow));
}

.logo-area h1 {
  font-size: 1.8rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.logo-area p {
  font-size: 0.95rem;
  color: var(--text-secondary);
}

.dropzone {
  border: 2px dashed var(--border);
  border-radius: var(--radius-xl);
  background: var(--bg-card);
  transition: all var(--transition);
  cursor: pointer;
}

.dropzone.active {
  border-color: var(--accent);
  background: var(--bg-card-hover);
  box-shadow: 0 0 0 6px var(--accent-glow);
  transform: scale(1.02);
}

.dropzone-content {
  padding: 3rem 2rem;
}

.drop-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 1rem;
  opacity: 0.7;
}

.dropzone p {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.8rem;
}

.divider {
  font-size: 0.8rem;
  color: var(--text-muted);
  display: block;
  margin-bottom: 1rem;
}

.upload-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.7rem 1.8rem;
  background: var(--accent-gradient);
  color: #fff;
  border-radius: var(--radius);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition);
}

.upload-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px var(--accent-glow);
}

.upload-status {
  margin-top: 1.2rem;
  font-size: 0.85rem;
  color: var(--success);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.upload-status.error {
  color: var(--danger);
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.2);
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.format-hint {
  margin-top: 1rem;
  font-size: 0.75rem;
  color: var(--text-muted);
}
</style>
