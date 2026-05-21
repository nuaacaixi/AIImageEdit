<template>
  <div class="editor-page">
    <header class="top-bar">
      <div class="top-left">
        <span class="logo-icon">✦</span>
        <span class="logo-text">AI Image Editor</span>
        <span class="badge">Beta</span>
      </div>
      <div class="top-right">
        <button class="top-btn" @click="resetAll" title="新建会话">
          ➕ 新图片
        </button>
      </div>
    </header>

    <div class="editor-body">
      <div class="viewer-area">
        <ImmersiveViewer
          :current-url="displayUrl"
          :original-url="originalUrl"
          alt="图片预览"
          :is-generating="store.isStreaming"
          :status-message="store.statusMessage"
        />
        <HistoryStrip @select="onVersionSelect" />
      </div>
    </div>

    <ToolBar />

    <ChatPanel
      @send="onChatSend"
      @view-image="onViewImage"
    />

    <UploadOverlay
      v-if="!store.hasImage"
      @uploaded="onUploaded"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useSessionStore } from '@/stores/session'
import { useViewerStore } from '@/stores/viewer'
import ImmersiveViewer from '@/components/ImmersiveViewer.vue'
import ToolBar from '@/components/ToolBar.vue'
import ChatPanel from '@/components/ChatPanel.vue'
import HistoryStrip from '@/components/HistoryStrip.vue'
import UploadOverlay from '@/components/UploadOverlay.vue'

const store = useSessionStore()
const viewer = useViewerStore()

const displayUrl = computed(() => {
  if (store.currentImage) return store.currentImage.url
  if (store.originalImage) return store.originalImage.url
  return ''
})

const originalUrl = computed(() => {
  return store.originalImage?.url ?? ''
})

defineEmits<{ uploaded: [] }>()

function onUploaded() {
  // Image uploaded, UploadOverlay will hide automatically
}

function onChatSend(message: string) {
  store.sendMessage(message)
}

function onViewImage(url: string) {
  window.open(url, '_blank')
}

function onVersionSelect(index: number) {
  const list = [...store.doneTurns]
  if (index === 0 && store.originalImage) {
    store.currentImage = store.originalImage
  } else {
    const turn = list[index - 1]
    if (turn && turn.resultImageUrl) {
      store.currentImage = {
        id: turn.turnId,
        url: turn.resultImageUrl,
      }
    }
  }
  viewer.resetView()
}

function resetAll() {
  store.resetSession()
  viewer.resetView()
}
</script>

<style scoped>
.editor-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.5rem;
  height: 3.2rem;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.top-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.logo-icon {
  font-size: 1.2rem;
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logo-text {
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.badge {
  font-size: 0.6rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  padding: 0.2rem 0.5rem;
  border-radius: 100px;
  background: var(--accent-gradient);
  color: #fff;
}

.top-btn {
  padding: 0.4rem 0.9rem;
  font-size: 0.78rem;
  font-weight: 600;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition);
}

.top-btn:hover {
  border-color: var(--accent);
  color: var(--accent-hover);
}

.editor-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1.2rem;
  gap: 0.8rem;
}

.viewer-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}
</style>
