<template>
  <div class="viewer">
    <!-- Empty state -->
    <div v-if="!hasImage" class="viewer-empty">
      <div class="empty-content">
        <span class="empty-icon">🖼️</span>
        <p class="empty-title">暂无图片</p>
        <p class="empty-desc">上传图片后在这里预览</p>
      </div>
    </div>

    <!-- Image display -->
    <div v-else class="image-area">
      <Transition name="crossfade" mode="out-in">
        <div :key="currentUrl" class="image-wrapper" :style="imageStyle">
          <img
            :src="currentUrl"
            :alt="alt"
            class="main-image"
            @error="onImageError"
            @load="onImageLoad"
            :class="{ 'image-hidden': imageError || imageLoading }"
          />
          <div v-if="imageLoading && !imageError" class="image-loading">
            <span class="loading-spinner"></span>
            <span class="loading-text">加载中...</span>
          </div>
          <div v-if="imageError" class="image-broken">
            <span class="broken-icon">🖼️</span>
            <p class="broken-title">图片加载失败</p>
            <p class="broken-desc">图片可能已被删除或无法访问</p>
            <button class="retry-btn" @click="retryLoad">重新加载</button>
          </div>
        </div>
      </Transition>
    </div>

    <!-- Controls (only show when image is loaded) -->
    <div v-if="hasImage && !imageError" class="viewer-controls">
      <button class="ctrl-btn" @click="viewer.resetView()" title="重置视图">↺</button>
      <button class="ctrl-btn" @click="viewer.zoom = Math.max(0.25, viewer.zoom - 0.25)" title="缩小">−</button>
      <span class="zoom-label">{{ Math.round(viewer.zoom * 100) }}%</span>
      <button class="ctrl-btn" @click="viewer.zoom = Math.min(3, viewer.zoom + 0.25)" title="放大">+</button>
      <button
        class="ctrl-btn"
        :class="{ active: viewer.compareMode }"
        @click="viewer.toggleCompare()"
        title="对比"
      >⇄</button>
    </div>

    <!-- Compare mode -->
    <div v-if="viewer.compareMode && !imageError" class="compare-panel" @click="viewer.toggleCompare()">
      <div class="compare-images" @click.stop>
        <div class="compare-side">
          <span class="compare-label">原图</span>
          <img :src="originalUrl" alt="原图" @error="e => (e.target as HTMLImageElement).style.display = 'none'" />
        </div>
        <div class="compare-side">
          <span class="compare-label">当前</span>
          <img :src="currentUrl" alt="当前" />
        </div>
      </div>
    </div>

    <!-- Generating overlay -->
    <div v-if="isGenerating" class="generating-overlay">
      <div class="generating-pulse"></div>
      <p>{{ statusMessage || 'AI 正在处理...' }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useViewerStore } from '@/stores/viewer'

const props = defineProps<{
  currentUrl: string
  originalUrl: string
  alt: string
  isGenerating: boolean
  statusMessage: string
}>()

const viewer = useViewerStore()
const imageError = ref(false)
const imageLoading = ref(true)

const hasImage = computed(() => props.currentUrl !== '')

watch(() => props.currentUrl, () => {
  imageError.value = false
  imageLoading.value = true
})

function onImageError() {
  imageError.value = true
  imageLoading.value = false
}

function onImageLoad() {
  imageError.value = false
  imageLoading.value = false
}

function retryLoad() {
  imageError.value = false
  imageLoading.value = true
}

const imageStyle = computed(() => ({
  transform: `scale(${viewer.zoom}) translate(${viewer.panX}px, ${viewer.panY}px)`,
}))
</script>

<style scoped>
.viewer {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: var(--bg-input);
  border-radius: var(--radius-lg);
  overflow: hidden;
  min-height: 400px;
  border: 1px solid var(--border);
}

.viewer-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 3rem;
}

.empty-content {
  text-align: center;
}

.empty-icon {
  font-size: 3.5rem;
  display: block;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-secondary);
  margin-bottom: 0.4rem;
}

.empty-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.image-area {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s ease-out;
  max-width: 100%;
  max-height: 100%;
}

.image-hidden {
  display: none;
}

.main-image {
  max-width: 90%;
  max-height: 70vh;
  object-fit: contain;
  border-radius: var(--radius-lg);
}

.crossfade-enter-active,
.crossfade-leave-active {
  transition: opacity 0.35s ease-out, filter 0.35s ease-out;
}
.crossfade-enter-from { opacity: 0; filter: blur(12px); }
.crossfade-leave-to { opacity: 0; }
.crossfade-enter-to { opacity: 1; filter: blur(0); }

/* Loading state */
.image-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.8rem;
  padding: 3rem;
}

.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 0.85rem;
  color: var(--text-muted);
}

/* Error state */
.image-broken {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 3rem;
  text-align: center;
}

.broken-icon {
  font-size: 3rem;
  opacity: 0.4;
  margin-bottom: 0.4rem;
}

.broken-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.broken-desc {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.retry-btn {
  margin-top: 0.6rem;
  padding: 0.45rem 1.2rem;
  font-size: 0.78rem;
  font-weight: 600;
  background: var(--accent-gradient);
  color: #fff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition);
}

.retry-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px var(--accent-glow);
}

/* Controls */
.viewer-controls {
  position: absolute;
  bottom: 1.2rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(12px);
  padding: 0.4rem 0.8rem;
  border-radius: 100px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  z-index: 10;
}

.ctrl-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: transparent;
  color: #fff;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
  transition: background var(--transition);
}

.ctrl-btn:hover {
  background: rgba(255, 255, 255, 0.15);
}

.ctrl-btn.active {
  background: var(--accent);
}

.zoom-label {
  color: #fff;
  font-size: 0.75rem;
  font-weight: 600;
  min-width: 40px;
  text-align: center;
}

/* Compare panel */
.compare-panel {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.55);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 20;
  cursor: pointer;
}

.compare-images {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  padding: 2rem;
  max-width: 90%;
}

.compare-side {
  text-align: center;
}

.compare-label {
  display: inline-block;
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #fff;
  margin-bottom: 0.6rem;
  padding: 0.3rem 0.8rem;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 100px;
  backdrop-filter: blur(8px);
}

.compare-side img {
  width: 100%;
  border-radius: var(--radius);
  max-height: 400px;
  object-fit: contain;
}

/* Generating overlay */
.generating-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.35);
  gap: 1rem;
  z-index: 15;
}

.generating-pulse {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: var(--accent-gradient);
  animation: pulse 1.5s ease-in-out infinite;
  opacity: 0.6;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 0.4; }
  50% { transform: scale(1.25); opacity: 0.8; }
}

.generating-overlay p {
  color: #fff;
  font-size: 0.95rem;
  font-weight: 600;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
}
</style>
