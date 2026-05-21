<template>
  <div class="toolbar">
    <button
      v-for="tool in tools"
      :key="tool.name"
      class="tool-btn"
      :class="{ active: store.activeTool === tool.name }"
      @click="store.setActiveTool(store.activeTool === tool.name ? null : tool.name)"
    >
      <span class="tool-icon">{{ tool.icon }}</span>
      <span class="tool-label">{{ tool.label }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { useSessionStore } from '@/stores/session'

const store = useSessionStore()

const tools = [
  { name: 'edit_image', label: 'AI 编辑', icon: '🎨' },
  { name: 'generate_image', label: 'AI 生成', icon: '✨' },
  { name: 'remove_background', label: '去背景', icon: '✂️' },
  { name: 'upscale', label: '超分', icon: '🔍' },
]
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 0.6rem;
  justify-content: center;
  padding: 0.8rem 0;
}

.tool-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
  padding: 0.7rem 1.2rem;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition);
  min-width: 80px;
}

.tool-btn:hover {
  border-color: var(--accent);
  color: var(--text-primary);
  background: var(--bg-card-hover);
  transform: translateY(-1px);
}

.tool-btn.active {
  border-color: var(--accent);
  background: rgba(99, 102, 241, 0.1);
  color: var(--accent-hover);
  box-shadow: 0 0 16px var(--accent-glow);
}

.tool-btn.active .tool-icon {
  animation: glow 1.5s ease-in-out infinite;
}

@keyframes glow {
  0%, 100% { filter: drop-shadow(0 0 2px var(--accent)); }
  50% { filter: drop-shadow(0 0 8px var(--accent-hover)); }
}

.tool-icon {
  font-size: 1.3rem;
}

.tool-label {
  font-size: 0.7rem;
  font-weight: 600;
}
</style>
